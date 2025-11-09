// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Moxie transpiler - converts Moxie source to Go and invokes the Go toolchain
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	moxieModule = "github.com/mleku/moxie"
	goModule    = "std"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "build":
		if err := buildCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie build: %v\n", err)
			os.Exit(1)
		}
	case "install":
		if err := installCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie install: %v\n", err)
			os.Exit(1)
		}
	case "run":
		if err := runCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie run: %v\n", err)
			os.Exit(1)
		}
	case "test":
		if err := testCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie test: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("moxie version 0.1.0 (transpiler mode)")
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Print(`Moxie is a tool for managing Moxie source code.

Usage:

	moxie <command> [arguments]

The commands are:

	build       transpile and compile packages and dependencies
	install     transpile and compile and install packages and dependencies
	run         transpile and run Moxie program
	test        transpile and test packages
	version     print Moxie version

Use "moxie help <command>" for more information about a command.
`)
}

func buildCommand(args []string) error {
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
	outputName := ""
	skipNext := false

	// Parse arguments to find source dir and -o flag
	for i := 0; i < len(args); i++ {
		if skipNext {
			skipNext = false
			continue
		}
		if args[i] == "-o" && i+1 < len(args) {
			outputName = args[i+1]
			skipNext = true
		} else if !strings.HasPrefix(args[i], "-") && srcDir == "." {
			srcDir = args[i]
		}
	}

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

	// Transpile source code
	if err := transpileTree(srcDir, tmpDir); err != nil {
		return fmt.Errorf("transpiling: %w", err)
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
	buildArgs := []string{"build", "-o", tmpBinary}
	for _, arg := range args {
		// Skip the source directory argument
		if arg == srcDir {
			continue
		}
		buildArgs = append(buildArgs, arg)
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
	if err := copyFile(tmpBinary, outputName); err != nil {
		return fmt.Errorf("copying binary: %w", err)
	}

	return nil
}

func installCommand(args []string) error {
	tmpDir, err := os.MkdirTemp("", "moxie-install-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := "."
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		srcDir = args[0]
		args = args[1:]
	}

	if err := transpileTree(srcDir, tmpDir); err != nil {
		return fmt.Errorf("transpiling: %w", err)
	}

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

	cmd := exec.Command("go", append([]string{"install"}, args...)...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go install failed: %w", err)
	}

	return nil
}

func runCommand(args []string) error {
	tmpDir, err := os.MkdirTemp("", "moxie-run-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := "."
	var sourceFile string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		ext := filepath.Ext(args[0])
		if ext == ".go" || ext == ".mx" {
			srcDir = filepath.Dir(args[0])
			sourceFile = filepath.Base(args[0])
		}
	}

	if err := transpileTree(srcDir, tmpDir); err != nil {
		return fmt.Errorf("transpiling: %w", err)
	}

	if err := copyGoMod(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.mod: %w", err)
	}

	// Copy runtime directory
	if err := copyRuntimeDir(tmpDir); err != nil {
		return fmt.Errorf("copying runtime: %w", err)
	}

	// Convert .mx file reference to .go for go run
	runArgs := make([]string, len(args))
	copy(runArgs, args)
	if sourceFile != "" && filepath.Ext(sourceFile) == ".mx" {
		runArgs[0] = sourceFile[:len(sourceFile)-3] + ".go"
	}

	cmd := exec.Command("go", append([]string{"run"}, runArgs...)...)
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

func testCommand(args []string) error {
	tmpDir, err := os.MkdirTemp("", "moxie-test-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := "."
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		srcDir = args[0]
		args = args[1:]
	}

	if err := transpileTree(srcDir, tmpDir); err != nil {
		return fmt.Errorf("transpiling: %w", err)
	}

	if err := copyGoMod(srcDir, tmpDir); err != nil {
		return fmt.Errorf("copying go.mod: %w", err)
	}

	// Copy runtime directory
	if err := copyRuntimeDir(tmpDir); err != nil {
		return fmt.Errorf("copying runtime: %w", err)
	}

	cmd := exec.Command("go", append([]string{"test"}, args...)...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go test failed: %w", err)
	}

	return nil
}

// transpileTree walks the source tree and transpiles all .go files
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

		// Only process .mx and .go files
		ext := filepath.Ext(path)
		if ext != ".mx" && ext != ".go" {
			return nil
		}

		// Compute relative path and destination
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		// Convert .mx extension to .go for output
		if filepath.Ext(dstPath) == ".mx" {
			dstPath = dstPath[:len(dstPath)-3] + ".go"
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

// transpileFile transpiles a single Go file to standard Go
func transpileFile(src, dst string) error {
	// Parse the source file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", src, err)
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
	syntaxTx := NewSyntaxTransformer()
	if err := syntaxTx.Transform(file); err != nil {
		// For now, print warning but continue
		// TODO: Make this a hard error once channel literals are fully implemented
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	// Transform package name (if needed)
	// For now, Moxie and Go use the same package names (lowercase)
	// This allows for future divergence if needed
	if file.Name != nil {
		moxiePkg := file.Name.Name
		goPkg := pkgMap.MoxieToGo(moxiePkg)
		if goPkg != moxiePkg {
			file.Name.Name = goPkg
		}
	}

	// Transform import paths
	for _, imp := range file.Imports {
		if imp.Path != nil {
			oldPath := strings.Trim(imp.Path.Value, `"`)
			newPath := transformImportPath(oldPath)
			if newPath != oldPath {
				imp.Path.Value = `"` + newPath + `"`
			}
		}
	}

	// Transform type, function, and variable names throughout the AST
	// Note: Currently disabled (typeMap.enabled = false, funcMap.enabled = false, varMap.enabled = false)
	// to maintain PascalCase/camelCase. Can be enabled in the future if desired
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			// Handle type, var, const declarations
			for _, spec := range node.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					// Transform type declaration: type MyType struct {}
					if typeMap.ShouldTransform(s.Name.Name) {
						typeMap.RegisterUserType(s.Name.Name)
						s.Name.Name = typeMap.TransformTypeName(s.Name.Name)
					}
					// Transform the type definition itself
					typeMap.transformTypeExpr(s.Type)

					// Transform struct fields if this is a struct type
					if structType, ok := s.Type.(*ast.StructType); ok {
						varMap.transformFieldList(structType.Fields)
					}

				case *ast.ValueSpec:
					// Transform variable/constant names and types
					varMap.transformValueSpec(s)
					typeMap.transformTypeExpr(s.Type)
				}
			}

		case *ast.FuncDecl:
			// Transform function/method declaration
			funcMap.transformFuncDecl(node)

			// Transform function receiver, parameters, and results
			if node.Recv != nil {
				typeMap.transformFieldList(node.Recv)
				varMap.transformFieldList(node.Recv)
			}
			if node.Type != nil {
				if node.Type.Params != nil {
					typeMap.transformFieldList(node.Type.Params)
					varMap.transformFieldList(node.Type.Params)
				}
				if node.Type.Results != nil {
					typeMap.transformFieldList(node.Type.Results)
					varMap.transformFieldList(node.Type.Results)
				}
			}

			// Transform function body
			if node.Body != nil {
				varMap.transformBlockStmt(node.Body)
			}

		case *ast.CallExpr:
			// Transform function calls
			funcMap.transformCallExpr(node)

		case *ast.FuncLit:
			// Transform function literals (anonymous functions)
			funcMap.transformFuncLit(node)

			// Transform parameters and results
			if node.Type != nil {
				if node.Type.Params != nil {
					varMap.transformFieldList(node.Type.Params)
				}
				if node.Type.Results != nil {
					varMap.transformFieldList(node.Type.Results)
				}
			}

			// Transform function body
			if node.Body != nil {
				varMap.transformBlockStmt(node.Body)
			}
		}
		return true
	})
}

// transformImportPath converts Moxie import paths to standard Go paths
func transformImportPath(path string) string {
	// Special case: Don't transform the runtime package import
	if path == "github.com/mleku/moxie/runtime" {
		return path
	}

	// Handle internal Moxie standard library paths
	if strings.HasPrefix(path, moxieModule+"/") {
		// github.com/mleku/moxie/internal/fmt -> fmt
		// github.com/mleku/moxie/src/fmt -> fmt
		// github.com/mleku/moxie/internal/net/http -> net/http
		remainder := strings.TrimPrefix(path, moxieModule+"/")

		// Strip internal/ or src/ prefix
		remainder = strings.TrimPrefix(remainder, "internal/")
		remainder = strings.TrimPrefix(remainder, "src/")

		// Return the full package path (not just the first part)
		return remainder
	}

	// Leave other imports unchanged
	return path
}

func copyGoMod(srcDir, dstDir string) error {
	srcMod := filepath.Join(srcDir, "go.mod")
	dstMod := filepath.Join(dstDir, "go.mod")

	if _, err := os.Stat(srcMod); os.IsNotExist(err) {
		// No go.mod, create a basic one with runtime module support
		content := `module moxie-build

go 1.24

require github.com/mleku/moxie/runtime v0.0.0

replace github.com/mleku/moxie/runtime => ./runtime
`
		return os.WriteFile(dstMod, []byte(content), 0644)
	}

	// Copy existing go.mod and add runtime module if needed
	content, err := os.ReadFile(srcMod)
	if err != nil {
		return err
	}

	// Check if runtime module is already present
	contentStr := string(content)
	if !strings.Contains(contentStr, "github.com/mleku/moxie/runtime") {
		// Add runtime module (use relative path since we copy runtime to build dir)
		contentStr += `

require github.com/mleku/moxie/runtime v0.0.0

replace github.com/mleku/moxie/runtime => ./runtime
`
	}

	return os.WriteFile(dstMod, []byte(contentStr), 0644)
}

// copyGoSum copies go.sum from source to destination directory
// This ensures transitive dependencies (like purego) are properly resolved
func copyGoSum(srcDir, dstDir string) error {
	// First, try to copy from the source directory
	srcSum := filepath.Join(srcDir, "go.sum")
	dstSum := filepath.Join(dstDir, "go.sum")

	if _, err := os.Stat(srcSum); err == nil {
		// Source has go.sum, copy it
		return copyFile(srcSum, dstSum)
	}

	// If source doesn't have go.sum, try to get it from the moxie project root
	// Get moxie executable path to find project root
	exePath, err := os.Executable()
	if err != nil {
		return nil // Don't fail if we can't find it
	}

	moxieRoot := filepath.Dir(exePath)
	projectSum := filepath.Join(moxieRoot, "go.sum")

	if _, err := os.Stat(projectSum); err == nil {
		// Project root has go.sum, copy it
		return copyFile(projectSum, dstSum)
	}

	// No go.sum found, that's okay - go will generate one
	return nil
}

// copyRuntimeDir copies the runtime directory to the build directory
func copyRuntimeDir(dstDir string) error {
	// Get moxie executable path to find runtime directory
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Runtime is in the same directory as the moxie executable
	// e.g., /path/to/moxie/moxie -> /path/to/moxie/runtime
	moxieRoot := filepath.Dir(exePath)
	runtimeSrc := filepath.Join(moxieRoot, "runtime")
	runtimeDst := filepath.Join(dstDir, "runtime")

	// Check if runtime directory exists
	if _, err := os.Stat(runtimeSrc); os.IsNotExist(err) {
		return fmt.Errorf("runtime directory not found at %s", runtimeSrc)
	}

	// Create destination runtime directory
	if err := os.MkdirAll(runtimeDst, 0755); err != nil {
		return err
	}

	// Copy all files from runtime directory
	entries, err := os.ReadDir(runtimeSrc)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip subdirectories for now
		}

		srcFile := filepath.Join(runtimeSrc, entry.Name())
		dstFile := filepath.Join(runtimeDst, entry.Name())

		if err := copyFile(srcFile, dstFile); err != nil {
			return fmt.Errorf("copying %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}
