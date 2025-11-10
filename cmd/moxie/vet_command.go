// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mleku/moxie/cmd/moxie/vet"
)

// vetCommand runs static analysis on Moxie code
func vetCommand(args []string) error {
	// Parse flags
	config := vet.DefaultConfig()
	var paths []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "--checks":
			if i+1 < len(args) {
				config.Checks = strings.Split(args[i+1], ",")
				i++
			} else {
				return fmt.Errorf("--checks requires argument")
			}

		case arg == "--format":
			if i+1 < len(args) {
				config.Format = args[i+1]
				i++
			} else {
				return fmt.Errorf("--format requires argument")
			}

		case arg == "--min-severity":
			if i+1 < len(args) {
				switch args[i+1] {
				case "info":
					config.MinSeverity = vet.SeverityInfo
				case "warning":
					config.MinSeverity = vet.SeverityWarning
				case "error":
					config.MinSeverity = vet.SeverityError
				default:
					return fmt.Errorf("invalid severity: %s", args[i+1])
				}
				i++
			} else {
				return fmt.Errorf("--min-severity requires argument")
			}

		case arg == "-h", arg == "--help":
			vetUsage()
			return nil

		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown flag: %s", arg)

		default:
			paths = append(paths, arg)
		}
	}

	// Default to current directory
	if len(paths) == 0 {
		paths = []string{"."}
	}

	// Create analyzer
	analyzer := vet.NewAnalyzer(config)

	// Analyze each path
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("accessing %s: %w", path, err)
		}

		if info.IsDir() {
			if err := analyzer.AnalyzeDir(path); err != nil {
				return fmt.Errorf("analyzing %s: %w", path, err)
			}
		} else {
			if err := analyzer.AnalyzeFile(path); err != nil {
				return fmt.Errorf("analyzing %s: %w", path, err)
			}
		}
	}

	// Filter and report issues
	issues := analyzer.FilterIssues()
	reporter := vet.NewReporter(config.Format)

	if err := reporter.Report(issues); err != nil {
		return fmt.Errorf("reporting issues: %w", err)
	}

	// Print summary for text format
	if config.Format == "text" {
		errors, warnings, info := analyzer.Summary()
		reporter.PrintSummary(errors, warnings, info)
	}

	// Exit with error code if issues found
	if analyzer.HasErrors() {
		os.Exit(1)
	}

	return nil
}

func vetUsage() {
	fmt.Print(`usage: moxie vet [flags] [path ...]

Vet runs static analysis on Moxie code to detect common errors and anti-patterns.

The flags are:

	--checks <categories>
		Comma-separated list of check categories to run.
		Available: memory, channels, types, const, errors
		Default: all checks

	--format <format>
		Output format: text, json, github
		Default: text

	--min-severity <severity>
		Minimum severity to report: info, warning, error
		Default: warning

Examples:

	moxie vet file.x                    # Vet single file
	moxie vet ./...                     # Vet all files recursively
	moxie vet --checks=memory ./...     # Only memory checks
	moxie vet --format=json ./...       # JSON output
	moxie vet --min-severity=error ./...  # Only errors

Check categories:

	memory     Memory management (clone/free usage)
	channels   Channel safety (deadlocks, leaks)
	types      Type safety (coercions, endianness)
	const      Const correctness
	errors     Error handling
`)
}
