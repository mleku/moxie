// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// formatCommand formats Moxie source files
func formatCommand(args []string) error {
	// Parse flags
	write := false      // -w flag: write result to source file instead of stdout
	list := false       // -l flag: list files that need formatting
	diff := false       // -d flag: display diffs instead of rewriting files

	var paths []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-w":
			write = true
		case "-l":
			list = true
		case "-d":
			diff = true
		case "-h", "--help":
			formatUsage()
			return nil
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			paths = append(paths, args[i])
		}
	}

	// Default to current directory if no paths specified
	if len(paths) == 0 {
		paths = []string{"."}
	}

	// Validate flag combinations
	if list && (write || diff) {
		return fmt.Errorf("cannot use -l with -w or -d")
	}
	if write && diff {
		return fmt.Errorf("cannot use -w with -d")
	}

	// Process each path
	hadError := false
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moxie fmt: %v\n", err)
			hadError = true
			continue
		}

		if info.IsDir() {
			// Format all .mx files in directory
			err = formatDir(path, write, list, diff)
		} else {
			// Format single file
			err = formatFile(path, write, list, diff)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "moxie fmt: %v\n", err)
			hadError = true
		}
	}

	if hadError {
		return fmt.Errorf("formatting errors encountered")
	}

	return nil
}

func formatUsage() {
	fmt.Print(`usage: moxie fmt [flags] [path ...]

Fmt formats Moxie source code files (.mx).

The flags are:

	-w
		Write result to source file instead of stdout.
	-l
		List files whose formatting differs from moxie fmt's.
	-d
		Display diffs instead of rewriting files.

When formatting a directory, moxie fmt processes all .mx files recursively.

Examples:

	moxie fmt file.mx           Format file.mx and print to stdout
	moxie fmt -w file.mx        Format file.mx and write back to file
	moxie fmt -l .              List all .mx files that need formatting
	moxie fmt -w ./...          Format all .mx files in current directory and subdirs
`)
}

// formatDir recursively formats all .mx files in a directory
func formatDir(dir string, write, list, diff bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and build artifacts
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || strings.HasPrefix(name, "moxie-build-") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .mx files
		if filepath.Ext(path) != ".mx" && filepath.Ext(path) != ".x" {
			return nil
		}

		return formatFile(path, write, list, diff)
	})
}

// formatFile formats a single Moxie source file
func formatFile(filename string, write, list, diff bool) error {
	// Read original source
	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Format the source
	formatted, err := formatMoxieSource(src)
	if err != nil {
		return fmt.Errorf("%s: %w", filename, err)
	}

	// Check if formatting changed anything
	if bytes.Equal(src, formatted) {
		// No changes needed
		return nil
	}

	// Handle different output modes
	if list {
		// Just list the filename
		fmt.Println(filename)
		return nil
	}

	if diff {
		// Show diff (simplified - just show that file differs)
		fmt.Printf("--- %s (original)\n", filename)
		fmt.Printf("+++ %s (formatted)\n", filename)
		// In a real implementation, we'd show a proper diff here
		// For now, just indicate the file differs
		fmt.Printf("File would be reformatted\n\n")
		return nil
	}

	if write {
		// Write back to source file
		// Preserve original file permissions
		info, err := os.Stat(filename)
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, formatted, info.Mode())
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "formatted %s\n", filename)
		return nil
	}

	// Default: write to stdout
	_, err = os.Stdout.Write(formatted)
	return err
}

// formatMoxieSource formats Moxie source code
func formatMoxieSource(src []byte) ([]byte, error) {
	// Step 1: Preprocess Moxie syntax to make it parseable by Go parser
	preprocessed := preprocessMoxieSyntax(string(src))

	// Step 2: Parse as Go code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", preprocessed, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	// Step 3: Format using Go's formatter
	var buf bytes.Buffer
	err = format.Node(&buf, fset, file)
	if err != nil {
		return nil, fmt.Errorf("formatting: %w", err)
	}

	// Step 4: Reverse preprocessing to restore Moxie syntax
	formatted := reversePreprocessMoxieSyntax(buf.String())

	// Step 5: Apply Moxie-specific formatting rules
	formatted = applyMoxieFormatting(formatted)

	return []byte(formatted), nil
}

// reversePreprocessMoxieSyntax reverses the preprocessing to restore Moxie syntax
func reversePreprocessMoxieSyntax(source string) string {
	// Reverse channel literal markers
	source = strings.ReplaceAll(source, "__MoxieChan[", "&chan ")
	source = strings.ReplaceAll(source, "__MoxieChanSend[", "&chan<- ")
	source = strings.ReplaceAll(source, "__MoxieChanRecv[", "&<-chan ")

	// Clean up remaining brackets from channel markers
	// This is a bit tricky - we need to only remove the ] that was part of the marker
	// For now, use a simple approach: replace "]{ with " {
	source = strings.ReplaceAll(source, "]{", " {")

	// Reverse endianness markers
	// __MoxieCoerceLE[T](...) → (*[]T, LittleEndian)(...)
	source = replaceEndiannessMarkers(source, "__MoxieCoerceLE", "LittleEndian")
	source = replaceEndiannessMarkers(source, "__MoxieCoerceBE", "BigEndian")

	return source
}

// replaceEndiannessMarkers replaces endianness function markers with tuple syntax
func replaceEndiannessMarkers(source, marker, endianness string) string {
	// This is a simplified replacement
	// In a real implementation, we'd use proper AST-based transformation
	// Pattern: __MoxieCoerceLE[Type](expr) → (*[]Type, LittleEndian)(expr)

	for {
		start := strings.Index(source, marker+"[")
		if start == -1 {
			break
		}

		// Find the matching ]
		bracketCount := 1
		i := start + len(marker) + 1
		typeEnd := -1

		for i < len(source) && bracketCount > 0 {
			if source[i] == '[' {
				bracketCount++
			} else if source[i] == ']' {
				bracketCount--
				if bracketCount == 0 {
					typeEnd = i
					break
				}
			}
			i++
		}

		if typeEnd == -1 {
			break // Malformed, stop processing
		}

		// Extract the type
		typeName := source[start+len(marker)+1 : typeEnd]

		// Replace with tuple syntax
		replacement := fmt.Sprintf("(*[]%s, %s)", typeName, endianness)
		source = source[:start] + replacement + source[typeEnd+1:]
	}

	return source
}

// applyMoxieFormatting applies Moxie-specific formatting rules
func applyMoxieFormatting(source string) string {
	// Apply any Moxie-specific formatting preferences
	// For now, this is a placeholder for future enhancements

	// Ensure consistent spacing around channel literals
	source = strings.ReplaceAll(source, "&chan  ", "&chan ")
	source = strings.ReplaceAll(source, "&chan<-  ", "&chan<- ")
	source = strings.ReplaceAll(source, "&<-chan  ", "&<-chan ")

	// Ensure spacing in endianness tuples
	source = strings.ReplaceAll(source, ",LittleEndian)", ", LittleEndian)")
	source = strings.ReplaceAll(source, ",BigEndian)", ", BigEndian)")
	source = strings.ReplaceAll(source, ",NativeEndian)", ", NativeEndian)")

	return source
}

// formatBytes is a helper that formats []byte from an io.Reader
func formatBytes(r io.Reader) ([]byte, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return formatMoxieSource(src)
}
