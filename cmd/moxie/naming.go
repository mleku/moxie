// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"unicode"
)

// Name conversion utilities for transforming between Moxie and Go naming conventions

// toSnakeCase converts PascalCase or camelCase to snake_case
// Examples:
//   - MyStruct -> my_struct
//   - HTTPServer -> http_server
//   - UserID -> user_id
//   - XMLParser -> xml_parser
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) + 5) // Estimate extra space for underscores

	// Track if we're in a sequence of uppercase letters (acronym)
	prevUpper := false
	prevPrevUpper := false

	runes := []rune(s)
	for i, r := range runes {
		isUpper := unicode.IsUpper(r)
		isDigit := unicode.IsDigit(r)

		// Add underscore before uppercase letter if:
		// 1. Not at start
		// 2. Previous char was lowercase or digit
		// 3. OR: This is lowercase and previous TWO were uppercase (end of acronym)
		if i > 0 {
			if isUpper && (!prevUpper || (i > 1 && prevPrevUpper && i < len(runes)-1 && unicode.IsLower(runes[i+1]))) {
				result.WriteRune('_')
			} else if isDigit && !unicode.IsDigit(runes[i-1]) && !prevUpper {
				result.WriteRune('_')
			}
		}

		result.WriteRune(unicode.ToLower(r))

		prevPrevUpper = prevUpper
		prevUpper = isUpper
	}

	return result.String()
}

// toPascalCase converts snake_case to PascalCase
// Examples:
//   - my_struct -> MyStruct
//   - http_server -> HTTPServer (if acronym)
//   - user_id -> UserID (if acronym)
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	parts := strings.Split(s, "_")
	var result strings.Builder

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Check if this part is a known acronym
		if isAcronym(part) {
			result.WriteString(strings.ToUpper(part))
		} else {
			// Capitalize first letter, keep rest as-is
			runes := []rune(part)
			result.WriteRune(unicode.ToUpper(runes[0]))
			if len(runes) > 1 {
				result.WriteString(string(runes[1:]))
			}
		}
	}

	return result.String()
}

// isAcronym checks if a lowercase word should be treated as an acronym
var commonAcronyms = map[string]bool{
	"http":  true,
	"https": true,
	"url":   true,
	"uri":   true,
	"id":    true,
	"uuid":  true,
	"api":   true,
	"rest":  true,
	"json":  true,
	"xml":   true,
	"html":  true,
	"css":   true,
	"sql":   true,
	"tcp":   true,
	"udp":   true,
	"ip":    true,
	"dns":   true,
	"tls":   true,
	"ssl":   true,
	"ssh":   true,
	"ftp":   true,
	"smtp":  true,
	"imap":  true,
	"pop":   true,
	"rpc":   true,
	"grpc":  true,
	"cpu":   true,
	"gpu":   true,
	"ram":   true,
	"io":    true,
	"os":    true,
	"ui":    true,
	"gui":   true,
	"cli":   true,
	"ascii": true,
	"utf":   true,
	"utc":   true,
}

func isAcronym(s string) bool {
	return commonAcronyms[strings.ToLower(s)]
}

// isExported checks if a name is exported (starts with uppercase)
func isExported(name string) bool {
	if name == "" {
		return false
	}
	r := []rune(name)[0]
	return unicode.IsUpper(r)
}

// preserveExportStatus converts a name while preserving its export status
// If input is exported, output is exported; if unexported, output is unexported
func preserveExportStatus(name string, converter func(string) string) string {
	if name == "" {
		return name
	}

	wasExported := isExported(name)
	converted := converter(name)

	if converted == "" {
		return converted
	}

	// Make sure export status matches
	if wasExported && !isExported(converted) {
		// Need to export it
		runes := []rune(converted)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	} else if !wasExported && isExported(converted) {
		// Need to unexport it - handle leading acronyms
		// If starts with an acronym like "HTTPServer", make it "httpServer"
		runes := []rune(converted)

		// Find where the leading uppercase sequence ends
		uppercaseCount := 0
		for i, r := range runes {
			if unicode.IsUpper(r) {
				uppercaseCount = i + 1
			} else {
				break
			}
		}

		// If we have multiple uppercase letters followed by lowercase,
		// lowercase all but the last uppercase letter
		// HTTPServer -> httpServer (lowercase HTTP, keep S)
		// ID -> id (lowercase entire acronym at start)
		if uppercaseCount > 1 && uppercaseCount < len(runes) {
			// Multiple uppercase followed by more content
			for i := 0; i < uppercaseCount-1; i++ {
				runes[i] = unicode.ToLower(runes[i])
			}
		} else {
			// Just one uppercase or all uppercase - lowercase first char
			runes[0] = unicode.ToLower(runes[0])
		}

		return string(runes)
	}

	return converted
}
