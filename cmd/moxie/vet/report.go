// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vet

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Reporter formats and outputs issues
type Reporter struct {
	format string
	output io.Writer
}

// NewReporter creates a new reporter
func NewReporter(format string) *Reporter {
	return &Reporter{
		format: format,
		output: os.Stdout,
	}
}

// Report outputs issues in the specified format
func (r *Reporter) Report(issues []*Issue) error {
	switch r.format {
	case "json":
		return r.reportJSON(issues)
	case "github":
		return r.reportGitHub(issues)
	default:
		return r.reportText(issues)
	}
}

// reportText outputs issues in human-readable format
func (r *Reporter) reportText(issues []*Issue) error {
	if len(issues) == 0 {
		fmt.Fprintln(r.output, "No issues found.")
		return nil
	}

	for _, issue := range issues {
		// Format: file:line:column: [category/check] message
		fmt.Fprintf(r.output, "%s:%d:%d: %s: [%s/%s] %s\n",
			issue.File,
			issue.Line,
			issue.Column,
			issue.Severity,
			issue.Category,
			issue.Check,
			issue.Message,
		)

		if issue.Help != "" {
			fmt.Fprintf(r.output, "  help: %s\n", issue.Help)
		}
		fmt.Fprintln(r.output)
	}

	return nil
}

// jsonIssue represents an issue in JSON format
type jsonIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Category string `json:"category"`
	Check    string `json:"check"`
	Message  string `json:"message"`
	Help     string `json:"help,omitempty"`
}

// jsonSummary represents issue summary in JSON format
type jsonSummary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
}

// jsonOutput represents the complete JSON output
type jsonOutput struct {
	Issues  []jsonIssue `json:"issues"`
	Summary jsonSummary `json:"summary"`
}

// reportJSON outputs issues in JSON format
func (r *Reporter) reportJSON(issues []*Issue) error {
	// Convert issues
	jsonIssues := make([]jsonIssue, len(issues))
	summary := jsonSummary{}

	for i, issue := range issues {
		jsonIssues[i] = jsonIssue{
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Severity: string(issue.Severity),
			Category: issue.Category,
			Check:    issue.Check,
			Message:  issue.Message,
			Help:     issue.Help,
		}

		switch issue.Severity {
		case SeverityError:
			summary.Errors++
		case SeverityWarning:
			summary.Warnings++
		case SeverityInfo:
			summary.Info++
		}
	}

	output := jsonOutput{
		Issues:  jsonIssues,
		Summary: summary,
	}

	encoder := json.NewEncoder(r.output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// reportGitHub outputs issues in GitHub Actions format
func (r *Reporter) reportGitHub(issues []*Issue) error {
	for _, issue := range issues {
		// GitHub format: ::warning file=file,line=line,col=col::message
		level := "warning"
		if issue.Severity == SeverityError {
			level = "error"
		}

		fmt.Fprintf(r.output, "::%s file=%s,line=%d,col=%d::%s [%s/%s]\n",
			level,
			issue.File,
			issue.Line,
			issue.Column,
			issue.Message,
			issue.Category,
			issue.Check,
		)
	}

	return nil
}

// PrintSummary prints a summary of issues
func (r *Reporter) PrintSummary(errors, warnings, info int) {
	if errors == 0 && warnings == 0 && info == 0 {
		return
	}

	fmt.Fprintf(r.output, "\nSummary:")
	if errors > 0 {
		fmt.Fprintf(r.output, " %d error(s)", errors)
	}
	if warnings > 0 {
		fmt.Fprintf(r.output, " %d warning(s)", warnings)
	}
	if info > 0 {
		fmt.Fprintf(r.output, " %d info", info)
	}
	fmt.Fprintln(r.output)
}
