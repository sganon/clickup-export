package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func rmSpace(s string) string {
	s = strings.ReplaceAll(s, "/", "")
	return strings.ReplaceAll(s, " ", "")
}

func mustCreateDir(dirPath string) {
	fmt.Println("Create DIR: ", dirPath)
	if err := os.MkdirAll(dirPath, 0770); err != nil {
		var pe *os.PathError
		if !errors.As(err, &pe) {
			panic(err)
		}
	}
}

func SaveSpace(outDir string, s *Space) error {
	spaceDir := filepath.Join(outDir, rmSpace(s.Name))
	mustCreateDir(spaceDir)
	for _, doc := range s.Docs {
		docDir := filepath.Join(spaceDir, rmSpace(doc.Name))
		if err := WriteDocPages(docDir, doc.Pages); err != nil {
			return fmt.Errorf("error writing page: %w", err)
		}
	}

	return nil
}

func WriteDocPages(baseDir string, pages []Page) error {
	mustCreateDir(baseDir)

	for _, p := range pages {
		name := rmSpace(p.Name)
		pageFile := filepath.Join(baseDir, name+".md")
		if err := os.WriteFile(pageFile, []byte(p.Content), 0770); err != nil {
			return fmt.Errorf("error writing file %s: %w", pageFile, err)
		}

		if len(p.Pages) > 0 {
			pageBaseDir := filepath.Join(baseDir, name)
			if err := WriteDocPages(pageBaseDir, p.Pages); err != nil {
				return fmt.Errorf("error writing sub page: %w", err)
			}
		}
	}

	return nil
}
