package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	headerFile := "scripts/header/header.txt"
	header, err := os.ReadFile(headerFile)
	if err != nil {
		log.Fatalf("Error reading header: %v", err)
	}

	err = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			injectHeader(path, header)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through files: %v", err)
	}
}

func injectHeader(path string, header []byte) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", path, err)
		return
	}

	if bytes.HasPrefix(content, header) {
		return
	}

	if len(content) > 0 && content[0] != '\n' {
		content = append([]byte{'\n'}, content...)
	}

	newContent := append(header, content...)
	err = os.WriteFile(path, newContent, 0644)
	if err != nil {
		log.Printf("Error writing file %s: %v", path, err)
	} else {
		log.Printf("Injected header into: %s\n", path)
	}
}
