// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Moxie transpiler - Self-hosted implementation written in Moxie
package main

import (
	"fmt"
	"os"
)

const (
	VERSION = "1.0.0-bootstrap"
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
	case "fmt":
		if err := formatCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie fmt: %v\n", err)
			os.Exit(1)
		}
	case "watch":
		if err := watchCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie watch: %v\n", err)
			os.Exit(1)
		}
	case "vet":
		if err := vetCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie vet: %v\n", err)
			os.Exit(1)
		}
	case "clean":
		if err := cleanCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie clean: %v\n", err)
			os.Exit(1)
		}
	case "lsp":
		if err := lspCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "moxie lsp: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("moxie version %s (self-hosted)\n", VERSION)
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
	fmt         format Moxie source files
	watch       watch for changes and rebuild automatically
	vet         run static analysis on Moxie code
	clean       remove cached build artifacts
	lsp         start Language Server Protocol server (for IDEs)
	version     print Moxie version

Use "moxie help <command>" for more information about a command.

This is the self-hosted Moxie compiler, written in Moxie itself!
`)
}
