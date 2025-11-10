// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vet provides static analysis for Moxie code
package vet

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Severity levels for issues
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue represents a linting issue found in code
type Issue struct {
	File     string   // File path
	Line     int      // Line number
	Column   int      // Column number
	Severity Severity // Issue severity
	Category string   // Check category (memory, channels, etc.)
	Check    string   // Specific check name
	Message  string   // Issue description
	Help     string   // Suggested fix
}

// Config holds vet configuration
type Config struct {
	Checks       []string         // Enabled check categories
	MinSeverity  Severity         // Minimum severity to report
	Format       string           // Output format (text, json, github)
	CheckConfigs map[string]Check // Category-specific configs
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Checks:      []string{"memory", "channels", "types", "const", "errors"},
		MinSeverity: SeverityWarning,
		Format:      "text",
		CheckConfigs: map[string]Check{
			"memory":   &MemoryCheck{},
			"channels": &ChannelCheck{},
			"types":    &TypeCheck{},
		},
	}
}

// Check interface for different check types
type Check interface {
	Name() string
	Category() string
	Analyze(fset *token.FileSet, file *ast.File) []*Issue
}

// Analyzer performs static analysis on Moxie code
type Analyzer struct {
	config *Config
	fset   *token.FileSet
	issues []*Issue
}

// NewAnalyzer creates a new analyzer with the given config
func NewAnalyzer(config *Config) *Analyzer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Analyzer{
		config: config,
		fset:   token.NewFileSet(),
		issues: make([]*Issue, 0),
	}
}

// AnalyzeFile analyzes a single Moxie source file
func (a *Analyzer) AnalyzeFile(filename string) error {
	// Read the file
	src, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filename, err)
	}

	// Preprocess Moxie syntax (similar to transpiler)
	// This handles channel literals and endianness tuples
	source := preprocessMoxieSyntax(string(src))

	// Parse as Go AST
	file, err := parser.ParseFile(a.fset, filename, source, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", filename, err)
	}

	// Run enabled checks
	for _, category := range a.config.Checks {
		if check, ok := a.config.CheckConfigs[category]; ok {
			issues := check.Analyze(a.fset, file)
			a.issues = append(a.issues, issues...)
		}
	}

	return nil
}

// AnalyzeDir recursively analyzes all Moxie files in a directory
func (a *Analyzer) AnalyzeDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories and build artifacts
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || strings.HasPrefix(name, "moxie-build-") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only analyze Moxie files
		if filepath.Ext(path) != ".mx" && filepath.Ext(path) != ".x" {
			return nil
		}

		return a.AnalyzeFile(path)
	})
}

// FilterIssues filters issues based on minimum severity
func (a *Analyzer) FilterIssues() []*Issue {
	filtered := make([]*Issue, 0)
	minLevel := severityLevel(a.config.MinSeverity)

	for _, issue := range a.issues {
		if severityLevel(issue.Severity) >= minLevel {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// severityLevel returns numeric level for comparison
func severityLevel(s Severity) int {
	switch s {
	case SeverityInfo:
		return 0
	case SeverityWarning:
		return 1
	case SeverityError:
		return 2
	default:
		return 0
	}
}

// Issues returns all issues found
func (a *Analyzer) Issues() []*Issue {
	return a.issues
}

// HasErrors returns true if any errors were found
func (a *Analyzer) HasErrors() bool {
	for _, issue := range a.issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Summary returns issue counts by severity
func (a *Analyzer) Summary() (errors, warnings, info int) {
	for _, issue := range a.issues {
		switch issue.Severity {
		case SeverityError:
			errors++
		case SeverityWarning:
			warnings++
		case SeverityInfo:
			info++
		}
	}
	return
}

// preprocessMoxieSyntax performs text-level preprocessing
// Reused from main transpiler logic
func preprocessMoxieSyntax(source string) string {
	// For now, just return as-is
	// In full implementation, this would call the preprocessor from main package
	// or have shared preprocessing logic
	return source
}
