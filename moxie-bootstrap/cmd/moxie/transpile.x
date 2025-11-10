// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// transpileTree walks the source tree and transpiles all .x and .go files
func transpileTree(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories we don't want to process
		if info.IsDir() {
			name := info.Name()
			// Only skip hidden directories, not the root
			if path != srcDir && (name == ".git" || name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".")) {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .x and .go files
		ext := filepath.Ext(path)
		if ext != ".x" && ext != ".go" {
			return nil
		}

		// Compute relative path and destination
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		// Convert .x extension to .go for output
		if filepath.Ext(dstPath) == ".x" {
			dstPath = dstPath[:len(dstPath)-2] + ".go"
		}

		// Create destination directory
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		// Transpile the file
		if err := transpileFile(path, dstPath); err != nil {
			return fmt.Errorf("transpiling %s: %w", path, err)
		}
		return nil
	})
}

// transpileFile transpiles a single Moxie file to standard Go
func transpileFile(src, dst string) error {
	// Read source file
	sourceBytes, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading %s: %w", src, err)
	}

	// Preprocess Moxie syntax (e.g., &chan T{} literals)
	preprocessed := preprocessMoxieSyntax(string(sourceBytes))

	// Parse the preprocessed source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, src, preprocessed, parser.ParseComments)
	if err != nil {
		// Post-process error message to show original Moxie syntax
		errMsg := postprocessMoxieSyntax(err.Error())
		return fmt.Errorf("parsing %s: %s", src, errMsg)
	}

	// Transform the AST
	transformAST(file)

	// Write the transformed code
	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	cfg := printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}
	if err := cfg.Fprint(outFile, fset, file); err != nil {
		return fmt.Errorf("writing %s: %w", dst, err)
	}

	return nil
}

// transformAST transforms the AST from Moxie to standard Go
func transformAST(file *ast.File) {
	// Phase 0: Check const enforcement (compile-time immutability)
	constChecker := NewConstChecker()
	if constErrors := constChecker.Check(file); len(constErrors) > 0 {
		fmt.Fprintf(os.Stderr, "const enforcement errors:\n")
		for _, err := range constErrors {
			fmt.Fprintf(os.Stderr, "  %v\n", err)
		}
		// For now, print warnings but continue
		// TODO: Make this a hard error in the future
	}

	// Phase 2: Apply syntax transformations (Moxie-specific syntax to Go)
	// This is where most of the work happens - see syntax.x
	applySyntaxTransformations(file)
}
