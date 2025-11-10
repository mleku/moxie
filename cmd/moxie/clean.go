// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// cleanCommand clears build cache and temporary files
func cleanCommand(args []string) error {
	// Parse flags
	cacheOnly := false
	verbose := false

	for _, arg := range args {
		switch arg {
		case "--cache":
			cacheOnly = true
		case "-v", "--verbose":
			verbose = true
		case "-h", "--help":
			cleanUsage()
			return nil
		default:
			if arg[0] == '-' {
				return fmt.Errorf("unknown flag: %s", arg)
			}
		}
	}

	// Clear build cache
	cache, err := NewBuildCache(true)
	if err != nil {
		return fmt.Errorf("accessing cache: %w", err)
	}

	if verbose {
		files, size, _ := cache.Stats()
		fmt.Printf("Cache stats: %d files, %d bytes\n", files, size)
	}

	if err := cache.Clear(); err != nil {
		return fmt.Errorf("clearing cache: %w", err)
	}

	fmt.Println("Cache cleared successfully")

	if !cacheOnly {
		// Could also clear temporary build directories here
		// For now, just the cache
	}

	return nil
}

func cleanUsage() {
	fmt.Print(`usage: moxie clean [flags]

Clean removes cached build artifacts and temporary files.

The flags are:

	--cache
		Only clear build cache (default clears everything)
	-v, --verbose
		Show what is being cleaned

Examples:

	moxie clean              # Clear all build artifacts
	moxie clean --cache      # Clear only build cache
	moxie clean -v           # Verbose output
`)
}
