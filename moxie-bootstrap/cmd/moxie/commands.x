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

// installCommand transpiles and installs a Moxie package
func installCommand(args *[]*[]byte) error {
	tmpDir, err := os.MkdirTemp("", "moxie-install-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := "."
	filteredArgs := &[]*[]byte{}
	if len(*args) > 0 && !strings.HasPrefix(string(*(*args)[0]), "-") {
		srcDir = string(*(*args)[0])
		filteredArgs = &(*args)[1:]
	} else {
		filteredArgs = args
	}

	if err := transpileTree(srcDir, tmpDir); err != nil {
		return fmt.Errorf("transpiling: %w", err)
	}

	if err := copyGoMod(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.mod: %w", err)
	}

	if err := copyGoSum(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.sum: %w", err)
	}

	if err := copyRuntimeDir(tmpDir); err != nil {
		return fmt.Errorf("copying runtime: %w", err)
	}

	installArgs := &[]*[]byte{&[]byte("install")}
	for _, arg := range *filteredArgs {
		installArgs = grow(installArgs, 1)
		(*installArgs)[len(*installArgs)-1] = arg
	}

	cmd := exec.Command("go", installArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go install failed: %w", err)
	}

	return nil
}

// runCommand transpiles and runs a Moxie program
func runCommand(args *[]*[]byte) error {
	tmpDir, err := os.MkdirTemp("", "moxie-run-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := "."
	progArgs := &[]*[]byte{}

	// Find first non-flag argument as source directory/file
	for i := 0; i < len(*args); i++ {
		arg := string(*(*args)[i])
		if !strings.HasPrefix(arg, "-") {
			srcDir = arg
			if i+1 < len(*args) {
				progArgs = &(*args)[i+1:]
			}
			break
		}
	}

	// Check if srcDir is a file or directory
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("checking source: %w", err)
	}

	if !srcInfo.IsDir() {
		// Single file - transpile just that file
		srcFile := srcDir
		srcDir = filepath.Dir(srcFile)

		baseName := filepath.Base(srcFile)
		dstPath := filepath.Join(tmpDir, baseName)
		if filepath.Ext(dstPath) == ".x" {
			dstPath = dstPath[:len(dstPath)-2] + ".go"
		}

		if err := transpileFile(srcFile, dstPath); err != nil {
			return fmt.Errorf("transpiling %s: %w", srcFile, err)
		}
	} else {
		// Directory - transpile entire tree
		if err := transpileTree(srcDir, tmpDir); err != nil {
			return fmt.Errorf("transpiling: %w", err)
		}
	}

	if err := copyGoMod(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.mod: %w", err)
	}

	if err := copyGoSum(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.sum: %w", err)
	}

	if err := copyRuntimeDir(tmpDir); err != nil {
		return fmt.Errorf("copying runtime: %w", err)
	}

	runArgs := &[]*[]byte{&[]byte("run"), &[]byte(".")}
	for _, arg := range *progArgs {
		runArgs = grow(runArgs, 1)
		(*runArgs)[len(*runArgs)-1] = arg
	}

	cmd := exec.Command("go", runArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go run failed: %w", err)
	}

	return nil
}

// testCommand transpiles and tests a Moxie package
func testCommand(args *[]*[]byte) error {
	tmpDir, err := os.MkdirTemp("", "moxie-test-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := "."
	filteredArgs := &[]*[]byte{}
	if len(*args) > 0 && !strings.HasPrefix(string(*(*args)[0]), "-") {
		srcDir = string(*(*args)[0])
		filteredArgs = &(*args)[1:]
	} else {
		filteredArgs = args
	}

	if err := transpileTree(srcDir, tmpDir); err != nil {
		return fmt.Errorf("transpiling: %w", err)
	}

	if err := copyGoMod(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.mod: %w", err)
	}

	if err := copyGoSum(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.sum: %w", err)
	}

	if err := copyRuntimeDir(tmpDir); err != nil {
		return fmt.Errorf("copying runtime: %w", err)
	}

	testArgs := &[]*[]byte{&[]byte("test")}
	for _, arg := range *filteredArgs {
		testArgs = grow(testArgs, 1)
		(*testArgs)[len(*testArgs)-1] = arg
	}

	cmd := exec.Command("go", testArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go test failed: %w", err)
	}

	return nil
}

// Placeholder commands - to be implemented
func formatCommand(args *[]*[]byte) error {
	return fmt.Errorf("format command not yet implemented in bootstrap")
}

func watchCommand(args *[]*[]byte) error {
	return fmt.Errorf("watch command not yet implemented in bootstrap")
}

func vetCommand(args *[]*[]byte) error {
	return fmt.Errorf("vet command not yet implemented in bootstrap")
}

func cleanCommand(args *[]*[]byte) error {
	return fmt.Errorf("clean command not yet implemented in bootstrap")
}

func lspCommand(args *[]*[]byte) error {
	return fmt.Errorf("lsp command not yet implemented in bootstrap")
}
