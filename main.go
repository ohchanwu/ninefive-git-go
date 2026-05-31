package main

import (
	"log"
	"os"
)

func main() {
	switch os.Args[1] {
	case "init":
		initRepo()
	}
}

func initRepo() {
	if err := os.Mkdir(".git/", os.FileMode(0o644)); err != nil {
		log.Fatal("error creating repository root")
	}
	if err := os.Mkdir(".git/objects/", os.FileMode(0o644)); err != nil {
		log.Fatal("error creating \".git/objects/\"")
	}
	if err := os.Mkdir(".git/refs/", os.FileMode(0o644)); err != nil {
		log.Fatal("error creating \".git/refs/\"")
	}
	if err := os.WriteFile(".git/HEAD", []byte("ref: refs/heads/main\n"), os.FileMode(0o644)); err != nil {
		log.Fatal("error creating \".git/HEAD\"")
	}
}
