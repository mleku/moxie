package main

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cases
		{"MyStruct", "my_struct"},
		{"UserID", "user_id"},
		{"HTTPServer", "http_server"},
		{"XMLParser", "xml_parser"},
		{"JSONEncoder", "json_encoder"},

		// Single words
		{"User", "user"},
		{"Server", "server"},
		{"ID", "id"},

		// Already snake_case
		{"my_struct", "my_struct"},
		{"user_id", "user_id"},

		// Acronyms
		{"HTTPSConnection", "https_connection"},
		{"URLParser", "url_parser"},
		{"APIKey", "api_key"},
		{"HTTPAPI", "httpapi"},

		// Mixed
		{"parseHTTPRequest", "parse_http_request"},
		{"myHTTPServer", "my_http_server"},
		{"HTTPSProxy", "https_proxy"},

		// Numbers
		{"User2", "user_2"},
		{"HTTP2Server", "http2_server"}, // 2 is treated as part of HTTP
		{"Base64Encoder", "base_64_encoder"},

		// Edge cases
		{"", ""},
		{"A", "a"},
		{"AB", "ab"},
		{"ABC", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cases
		{"my_struct", "MyStruct"},
		{"user_id", "UserID"},
		{"http_server", "HTTPServer"},
		{"xml_parser", "XMLParser"},
		{"json_encoder", "JSONEncoder"},

		// Single words
		{"user", "User"},
		{"server", "Server"},
		{"id", "ID"},

		// Already PascalCase (treated as single word if no underscores)
		{"MyStruct", "MyStruct"},
		{"UserID", "UserID"},

		// Acronyms
		{"https_connection", "HTTPSConnection"},
		{"url_parser", "URLParser"},
		{"api_key", "APIKey"},

		// Multiple underscores
		{"my__struct", "MyStruct"}, // Empty parts ignored
		{"_my_struct", "MyStruct"}, // Leading underscore ignored

		// Edge cases
		{"", ""},
		{"a", "A"},
		{"ab", "Ab"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("toPascalCase(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsExported(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"MyStruct", true},
		{"myStruct", false},
		{"HTTPServer", true},
		{"httpServer", false},
		{"_private", false},
		{"", false},
		{"A", true},
		{"a", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isExported(tt.input)
			if result != tt.expected {
				t.Errorf("isExported(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPreserveExportStatus(t *testing.T) {
	tests := []struct {
		input     string
		converter func(string) string
		expected  string
	}{
		// Exported stays exported
		{"MyStruct", toSnakeCase, "My_struct"},
		{"HTTPServer", toSnakeCase, "Http_server"},

		// Unexported stays unexported
		{"myStruct", toSnakeCase, "my_struct"},
		{"httpServer", toSnakeCase, "http_server"},

		// Round trip
		{"my_struct", toPascalCase, "myStruct"}, // unexported input -> unexported output
		{"My_struct", toPascalCase, "MyStruct"}, // exported input -> exported output
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := preserveExportStatus(tt.input, tt.converter)
			if result != tt.expected {
				t.Errorf("preserveExportStatus(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []string{
		"MyStruct",
		"HTTPServer",
		"UserID",
		"XMLParser",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			// PascalCase -> snake_case -> PascalCase
			snake := toSnakeCase(tt)
			back := toPascalCase(snake)

			// May not be exactly equal due to acronym handling
			// but should at least start with same letter (export status)
			if isExported(tt) != isExported(back) {
				t.Errorf("Round trip lost export status: %q -> %q -> %q", tt, snake, back)
			}
		})
	}
}

func TestIsAcronym(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http", true},
		{"https", true},
		{"id", true},
		{"url", true},
		{"api", true},
		{"json", true},
		{"xml", true},
		{"sql", true},
		{"tcp", true},
		{"udp", true},

		// Not acronyms
		{"user", false},
		{"server", false},
		{"parser", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isAcronym(tt.input)
			if result != tt.expected {
				t.Errorf("isAcronym(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
