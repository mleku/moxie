// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyGoMod copies go.mod from src to dst directory if it exists
func copyGoMod(srcDir, dstDir string) error {
	srcPath := filepath.Join(srcDir, "go.mod")
	dstPath := filepath.Join(dstDir, "go.mod")

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		// go.mod doesn't exist, that's OK
		return nil
	}

	return copyFileBinary(srcPath, dstPath)
}

// copyGoSum copies go.sum from src to dst directory if it exists
func copyGoSum(srcDir, dstDir string) error {
	srcPath := filepath.Join(srcDir, "go.sum")
	dstPath := filepath.Join(dstDir, "go.sum")

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		// go.sum doesn't exist, that's OK
		return nil
	}

	return copyFileBinary(srcPath, dstPath)
}

// copyRuntimeDir copies the runtime directory to the build directory
func copyRuntimeDir(dstDir string) error {
	// Find runtime directory (relative to moxie executable)
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable: %w", err)
	}

	exeDir := filepath.Dir(exe)
	runtimeSrc := filepath.Join(exeDir, "..", "..", "runtime")

	// Check if runtime exists at this path
	if _, err := os.Stat(runtimeSrc); os.IsNotExist(err) {
		// Try relative to current directory (development mode)
		runtimeSrc = "runtime"
		if _, err := os.Stat(runtimeSrc); os.IsNotExist(err) {
			return fmt.Errorf("runtime directory not found")
		}
	}

	runtimeDst := filepath.Join(dstDir, "github.com", "mleku", "moxie", "runtime")
	if err := os.MkdirAll(runtimeDst, 0755); err != nil {
		return fmt.Errorf("creating runtime directory: %w", err)
	}

	// Copy all .go files from runtime
	return filepath.Walk(runtimeSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only copy .go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		relPath, err := filepath.Rel(runtimeSrc, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(runtimeDst, relPath)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		return copyFileBinary(path, dstPath)
	})
}

// copyFileBinary copies a file from src to dst
func copyFileBinary(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
