package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	var err error

	// Check cli arg format
	if len(os.Args) > 2 {
		failAndExit("Usage: undup [directory]")
	}
	path := "."
	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	// make path absolute
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			failAndExit(fmt.Sprintf("Error getting absolute path: %v", err))
		}
	}

	fmt.Printf("Will undup directory structure at %s\n", path)

	base := filepath.Base(path)
	duplicatePath := filepath.Join(path, base)

	// Check if the structure is actually duplicated
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		failAndExit("Not a directory")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		failAndExit(fmt.Sprintf("Error reading directory: %v", err))
	}
	if len(entries) != 1 {
		failAndExit(fmt.Sprintf("Have %d entries, expected 1", len(entries)))
	}

	if entries[0].Name() != base || !entries[0].IsDir() {
		fmt.Println("Not a duplicated directory structure")
		os.Exit(1)
	}

	// Move contents up one level
	contents, err := os.ReadDir(duplicatePath)
	if err != nil {
		fmt.Printf("Error reading duplicate directory: %v\n", err)
		os.Exit(1)
	}

	for _, item := range contents {
		oldPath := filepath.Join(duplicatePath, item.Name())
		newPath := filepath.Join(path, item.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("Error moving %s: %v\n", item.Name(), err)
			continue
		}
	}

	// Remove empty directory
	if err := os.Remove(duplicatePath); err != nil {
		fmt.Printf("Error removing duplicate directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully undup directory structure")
}

func failAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
