package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var markdownHeader = "## "

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: concat <file1> <file2> ...")
		os.Exit(1)
	}

	files := os.Args[1:]
	sort.Strings(files)

	for i, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			os.Exit(1)
		}

		baseName := filepath.Base(file)
		ext := filepath.Ext(baseName)
		nameWithoutExt := strings.TrimSuffix(baseName, ext)

		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s%s\n\n", markdownHeader, nameWithoutExt)
		fmt.Print(string(content))
	}
}
