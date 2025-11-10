// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mleku/moxie/cmd/moxie/lsp"
)

// lspCommand starts the LSP server
func lspCommand(args []string) error {
	// Parse flags
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			lspUsage()
			return nil
		default:
			if arg[0] == '-' {
				return fmt.Errorf("unknown flag: %s", arg)
			}
		}
	}

	// Create and start LSP server
	server := lsp.NewServer()
	ctx := context.Background()

	// Use stdin/stdout for LSP communication
	if err := server.Start(ctx, os.Stdin, os.Stdout); err != nil {
		return fmt.Errorf("LSP server error: %w", err)
	}

	return nil
}

func lspUsage() {
	fmt.Print(`usage: moxie lsp

Start the Moxie Language Server Protocol (LSP) server.

The LSP server provides IDE features such as:
  - Symbol indexing and navigation
  - Go to definition
  - Find references
  - Hover information
  - Code completion
  - Real-time diagnostics
  - Code formatting

The server communicates over stdin/stdout using JSON-RPC 2.0.

This command is typically invoked automatically by IDE extensions.
For VS Code, install the "Moxie" extension from the marketplace.

Examples:

	# Start LSP server (typically called by IDE)
	moxie lsp

Configuration:

The LSP server respects the following settings:
  - Workspace root: Set by the client during initialization
  - File patterns: Automatically discovers .mx and .x files
  - Diagnostics: Integrates with 'moxie vet' for static analysis

Logging:

LSP protocol messages use stdin/stdout. Debug logs are written to stderr.
Set the MOXIE_LSP_DEBUG environment variable for verbose logging.
`)
}
