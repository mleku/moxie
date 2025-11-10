// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildCommand transpiles and builds a Moxie project
func buildCommand(args *[]*[]byte) error {
	// Create temporary directory for transpiled code
	tmpDir, err := os.MkdirTemp("", "moxie-build-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	// DEBUG: Don't remove tmpDir so we can inspect it
	fmt.Fprintf(os.Stderr, "DEBUG: Build directory: %s\n", tmpDir)
	// defer os.RemoveAll(tmpDir)

	// Determine source directory and output binary name
	srcDir := "."
	srcArg := "." // Track the original source argument to skip it in build args
	outputName := ""
	skipNext := false

	// Parse arguments to find source dir and -o flag
	for i := 0; i < len(*args); i++ {
		if skipNext {
			skipNext = false
			continue
		}
		arg := string(*(*args)[i])
		if arg == "-o" && i+1 < len(*args) {
			outputName = string(*(*args)[i+1])
			skipNext = true
		} else if !strings.HasPrefix(arg, "-") && srcDir == "." {
			srcArg = arg
			srcDir = arg
		}
	}

	// Check if srcDir is actually a file
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("checking source: %w", err)
	}

	if !srcInfo.IsDir() {
		// srcDir is a file, extract directory and use file-specific logic
		srcFile := srcDir
		srcDir = filepath.Dir(srcFile)

		// If no -o flag, use file base name without extension
		if outputName == "" {
			baseName := filepath.Base(srcFile)
			ext := filepath.Ext(baseName)
			outputName = baseName[:len(baseName)-len(ext)]
		}

		// Transpile single file directly to tmpDir (flatten structure)
		baseName := filepath.Base(srcFile)
		dstPath := filepath.Join(tmpDir, baseName)
		if filepath.Ext(dstPath) == ".x" {
			dstPath = dstPath[:len(dstPath)-2] + ".go"
		}

		// Transpile the file
		if err := transpileFile(srcFile, dstPath); err != nil {
			return fmt.Errorf("transpiling %s: %w", srcFile, err)
		}
	} else {
		// If no -o flag, determine output name from source directory
		if outputName == "" {
			if srcDir == "." {
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}
				outputName = filepath.Base(cwd)
			} else {
				outputName = filepath.Base(srcDir)
			}
		}

		// Transpile source code tree
		if err := transpileTree(srcDir, tmpDir); err != nil {
			return fmt.Errorf("transpiling: %w", err)
		}
	}

	// Copy go.mod if it exists
	if err := copyGoMod(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.mod: %w", err)
	}

	// Copy go.sum to ensure dependencies are resolved
	if err := copyGoSum(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.sum: %w", err)
	}

	// Copy runtime directory
	if err := copyRuntimeDir(tmpDir); err != nil {
		return fmt.Errorf("copying runtime: %w", err)
	}

	// Build with output name in temp dir
	tmpBinary := filepath.Join(tmpDir, outputName)

	// Filter out the package path from args since we transpiled to tmpDir root
	buildArgs := &[]*[]byte{
		&[]byte("build"),
		&[]byte("-o"),
		&[]byte(tmpBinary),
	}
	for _, arg := range *args {
		argStr := string(*arg)
		// Skip the source argument (file or directory)
		if argStr == srcArg {
			continue
		}
		buildArgs = grow(buildArgs, 1)
		(*buildArgs)[len(*buildArgs)-1] = arg
	}

	cmd := exec.Command("go", buildArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	// Copy binary back to original directory
	if err := copyFileBinary(tmpBinary, outputName); err != nil {
		return fmt.Errorf("copying binary: %w", err)
	}

	return nil
}
