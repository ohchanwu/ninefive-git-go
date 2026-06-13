package main

import (
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	switch os.Args[1] {
	case "init":
		if err := initRepo(); err != nil {
			log.Fatal(err)
		}
	case "cat-file":
		if err := catFile(os.Args[len(os.Args)-1]); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "hash-object":
		if err := hashObject(os.Args[len(os.Args)-1]); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "ls-tree":
		if err := lsTree(os.Args[len(os.Args)-1]); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func lsTree(sha string) error {
	path := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])
	f, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(0o0644))
	if err != nil {
		return fmt.Errorf("error opening tree at %s: %w", path, err)
	}
	defer f.Close()

	reader, err := zlib.NewReader(f)
	if err != nil {
		return fmt.Errorf("error creating zlib reader: %w", err)
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("error reading decompressed content: %w", err)
	}

	// strip tree header
	content := bytes[slices.Index(bytes, 0x00)+1:]

	for entry := range strings.SplitAfterSeq(string(content), " ") {
		if name, _, found := strings.Cut(entry, string(byte(0))); found {
			fmt.Println(name)
		}
	}

	if slices.Contains(os.Args[2:len(os.Args)-1], "--name-only") {
		return nil
	}

	return nil
}

func hashObject(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file at %s: %w", path, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading file at %s: %w", path, err)
	}

	header := "blob " + strconv.Itoa(len(content)) + string(byte(0))
	store := []byte(header + string(content))
	h := sha1.New()
	h.Write(store)
	sha := fmt.Sprintf("%x", string(h.Sum(nil)))

	fmt.Println(sha)

	// If no "-w" flag is present, do not create a blob object
	// and write the file's compressed contents to it.
	if len(os.Args) == 3 {
		return nil
	}

	objDir := fmt.Sprintf(".git/objects/%s/", sha[:2])
	objPath := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])
	if err := os.Mkdir(objDir, os.FileMode(0o644)); err != nil {
		return fmt.Errorf("error creating blob dir at %s, %w", objDir, err)
	}
	blobFile, err := os.OpenFile(objPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("error creating blob file at %s, %w", objPath, err)
	}
	defer blobFile.Close()

	zlibWriter := zlib.NewWriter(blobFile)
	defer zlibWriter.Close()
	if _, err := zlibWriter.Write(store); err != nil {
		return fmt.Errorf("error writing to blob at %s: %w", objPath, err)
	}

	return nil
}

func catFile(sha string) error {
	if len(sha) != 40 {
		return errors.New("SHA must be 40 characters long")
	}

	path := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])
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

	// load the object's contents into memory
	// not the best idea if the object is very large
	// so may have to refactor this part
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
