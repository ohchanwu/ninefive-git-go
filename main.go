package main

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

func main() {
	switch os.Args[1] {
	case "init":
		if err := initRepo(); err != nil {
			log.Fatal(err)
		}
	case "cat-file":
		if err := catFile(os.Args[len(os.Args)-1], os.Args[2:len(os.Args)]...); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func catFile(sha string, flags ...string) error {
	if len(sha) != 40 {
		return errors.New("SHA must be 40 characters long")
	}

	path := fmt.Sprintf(".git/objects/%s/%s", sha[0:2], sha[2:])
	f, err := os.OpenFile(path, os.O_RDWR, os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("error opening file at %s: %w", path, err)
	}
	defer f.Close()

	zlibReader, err := zlib.NewReader(f)
	if err != nil {
		return fmt.Errorf("error creating zlib reader: %w", err)
	}
	defer zlibReader.Close()

	bytes, err := io.ReadAll(zlibReader)
	if err != nil {
		return fmt.Errorf("error reading decompressed content: %w", err)
	}
	nullByteIndex := slices.Index(bytes, byte(0))
	bytes = bytes[nullByteIndex+1:]

	if _, err := fmt.Print(string(bytes)); err != nil {
		return fmt.Errorf("error writing decompressed content to stdout: %w", err)
	}

	return nil
}

func initRepo() error {
	if err := os.Mkdir(".git/", os.FileMode(0o644)); err != nil {
		return errors.New("error creating repository root")
	}
	if err := os.Mkdir(".git/objects/", os.FileMode(0o644)); err != nil {
		return errors.New("error creating \".git/objects/\"")
	}
	if err := os.Mkdir(".git/refs/", os.FileMode(0o644)); err != nil {
		return errors.New("error creating \".git/refs/\"")
	}
	if err := os.WriteFile(".git/HEAD", []byte("ref: refs/heads/main\n"), os.FileMode(0o644)); err != nil {
		return errors.New("error creating \".git/HEAD\"")
	}
	return nil
}
