// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// SourceMapper tracks line mappings from .go to .mx files
type SourceMapper struct {
	mappings map[string]map[int]int // goFile -> (goLine -> mxLine)
}

// NewSourceMapper creates a new source mapper
func NewSourceMapper() *SourceMapper {
	return &SourceMapper{
		mappings: make(map[string]map[int]int),
	}
}

// AddMapping adds a line mapping from Go to Moxie source
func (sm *SourceMapper) AddMapping(goFile string, goLine, mxLine int) {
	if sm.mappings[goFile] == nil {
		sm.mappings[goFile] = make(map[int]int)
	}
	sm.mappings[goFile][goLine] = mxLine
}

// MapLine maps a Go line number to Moxie source
func (sm *SourceMapper) MapLine(goFile string, goLine int) (mxFile string, mxLine int, ok bool) {
	// Convert .go to .mx/.x
	mxFile = strings.TrimSuffix(goFile, ".go")
	if !strings.HasSuffix(mxFile, ".mx") && !strings.HasSuffix(mxFile, ".x") {
		mxFile = mxFile + ".x"
	}

	// For now, simple 1:1 mapping since preprocessing doesn't change line numbers significantly
	// In a full implementation, we'd track actual transformations
	mxLine = goLine
	ok = true

	return
}

// TranslateError translates a Go compiler error to Moxie source
func (sm *SourceMapper) TranslateError(errorMsg string) string {
	// Pattern: ./file.go:line:col: error message
	re := regexp.MustCompile(`(\.?[^:]+\.go):(\d+):(\d+): (.+)`)
	matches := re.FindStringSubmatch(errorMsg)

	if len(matches) != 5 {
		// Not a standard error format, return as-is
		return errorMsg
	}

	goFile := matches[1]
	goLine := matches[2]
	goCol := matches[3]
	message := matches[4]

	// Map to Moxie source
	mxFile := strings.TrimSuffix(goFile, ".go")
	if !strings.HasSuffix(mxFile, ".mx") && !strings.HasSuffix(mxFile, ".x") {
		mxFile = mxFile + ".x"
	}

	// Reconstruct error message with .mx file
	return fmt.Sprintf("%s:%s:%s: %s", mxFile, goLine, goCol, message)
}

// TranslateGoOutput translates Go compiler output to Moxie source references
func (sm *SourceMapper) TranslateGoOutput(output string) string {
	lines := strings.Split(output, "\n")
	translated := make([]string, 0, len(lines))

	for _, line := range lines {
		translated = append(translated, sm.TranslateError(line))
	}

	return strings.Join(translated, "\n")
}

// ProcessGoErrors processes Go command output and translates errors
func ProcessGoErrors(output []byte) []byte {
	sm := NewSourceMapper()
	translated := sm.TranslateGoOutput(string(output))
	return []byte(translated)
}

// AddContextToError adds source code context to error messages
func AddContextToError(filename string, line int, errorMsg string) string {
	// Read the source file
	file, err := os.Open(filename)
	if err != nil {
		return errorMsg
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	contextLines := make([]string, 0)

	// Collect lines around the error
	startLine := line - 2
	if startLine < 1 {
		startLine = 1
	}
	endLine := line + 2

	for scanner.Scan() {
		if lineNum >= startLine && lineNum <= endLine {
			prefix := fmt.Sprintf("%4d | ", lineNum)
			if lineNum == line {
				prefix = fmt.Sprintf("%4d > ", lineNum) // Mark error line
			}
			contextLines = append(contextLines, prefix+scanner.Text())
		}
		lineNum++
		if lineNum > endLine {
			break
		}
	}

	if len(contextLines) == 0 {
		return errorMsg
	}

	// Build enhanced error message
	var result strings.Builder
	result.WriteString(errorMsg)
	result.WriteString("\n\n")
	result.WriteString(strings.Join(contextLines, "\n"))
	result.WriteString("\n")

	return result.String()
}
