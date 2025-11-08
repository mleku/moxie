package main

import (
	"testing"
)

func TestTransformImportPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "github.com/mleku/moxie/src/fmt",
			expected: "fmt",
		},
		{
			input:    "github.com/mleku/moxie/internal/os",
			expected: "os",
		},
		{
			input:    "github.com/mleku/moxie/src/net/http",
			expected: "net/http",
		},
		{
			input:    "github.com/mleku/moxie/internal/encoding/json",
			expected: "encoding/json",
		},
		{
			input:    "github.com/other/package",
			expected: "github.com/other/package",
		},
		{
			input:    "fmt",
			expected: "fmt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := transformImportPath(tt.input)
			if result != tt.expected {
				t.Errorf("transformImportPath(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
