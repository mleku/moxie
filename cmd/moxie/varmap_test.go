package main

import "testing"

func TestVarMapper_IsBuiltin(t *testing.T) {
	vm := NewVarMapper()

	builtins := []string{
		"nil", "true", "false", "iota",
	}

	for _, v := range builtins {
		if !vm.IsBuiltin(v) {
			t.Errorf("IsBuiltin(%q) = false; want true", v)
		}
	}

	notBuiltins := []string{
		"myVar", "userName", "count", "data",
	}

	for _, v := range notBuiltins {
		if vm.IsBuiltin(v) {
			t.Errorf("IsBuiltin(%q) = true; want false", v)
		}
	}
}

func TestVarMapper_IsSpecial(t *testing.T) {
	vm := NewVarMapper()

	special := []string{
		"_", // Blank identifier
	}

	for _, v := range special {
		if !vm.IsSpecial(v) {
			t.Errorf("IsSpecial(%q) = false; want true", v)
		}
	}

	notSpecial := []string{
		"myVar", "userName", "_private", "__unused",
	}

	for _, v := range notSpecial {
		if vm.IsSpecial(v) {
			t.Errorf("IsSpecial(%q) = true; want false", v)
		}
	}
}

func TestVarMapper_ShouldTransform(t *testing.T) {
	vm := NewVarMapper()

	tests := []struct {
		name    string
		varName string
		want    bool
	}{
		// Builtins - never transform
		{"builtin nil", "nil", false},
		{"builtin true", "true", false},
		{"builtin false", "false", false},
		{"builtin iota", "iota", false},

		// Special - never transform
		{"special blank", "_", false},

		// Empty/single char - never transform
		{"empty", "", false},
		{"single char", "x", false},
		{"single char i", "i", false},

		// User variables - would transform if enabled
		{"user var", "myVar", false}, // false because vm.enabled = false
		{"user name", "userName", false},
		{"count", "count", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vm.ShouldTransform(tt.varName)
			if got != tt.want {
				t.Errorf("ShouldTransform(%q) = %v; want %v", tt.varName, got, tt.want)
			}
		})
	}
}

func TestVarMapper_ShouldTransformWhenEnabled(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable() // Enable transformation

	tests := []struct {
		name    string
		varName string
		want    bool
	}{
		// Builtins - still never transform
		{"builtin nil", "nil", false},
		{"builtin true", "true", false},
		{"builtin false", "false", false},

		// Special - still never transform
		{"special blank", "_", false},

		// User variables - now should transform
		{"user var", "myVar", true},
		{"user name", "userName", true},
		{"count", "count", true},
		{"data ptr", "dataPtr", true},

		// Edge cases - never transform
		{"empty", "", false},
		{"single char", "x", false},
		{"single char i", "i", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vm.ShouldTransform(tt.varName)
			if got != tt.want {
				t.Errorf("ShouldTransform(%q) = %v; want %v (enabled=true)", tt.varName, got, tt.want)
			}
		})
	}
}

func TestVarMapper_TransformVarName_Disabled(t *testing.T) {
	vm := NewVarMapper()
	// Disabled by default

	tests := []string{
		"myVar",
		"userName",
		"count",
		"dataPtr",
		"httpServer",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := vm.TransformVarName(input)
			if result != input {
				t.Errorf("TransformVarName(%q) = %q; want %q (when disabled)", input, result, input)
			}
		})
	}
}

func TestVarMapper_TransformVarName_Enabled(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable()

	tests := []struct {
		input string
		want  string
	}{
		// snake_case -> camelCase (Moxie -> Go direction)
		// Variables are typically unexported (start with lowercase)
		{"my_var", "myVar"},
		{"user_name", "userName"},
		{"user_id", "userID"},               // ID is an acronym
		{"http_server", "httpServer"},       // HTTP is an acronym
		{"max_count", "maxCount"},
		{"data_ptr", "dataPtr"},

		// Exported variables (rare but possible)
		{"My_var", "MyVar"},
		{"User_name", "UserName"},

		// Already camelCase - unchanged
		{"myVar", "myVar"},
		{"userName", "userName"},
		{"count", "count"},

		// Builtins - unchanged even when enabled
		{"nil", "nil"},
		{"true", "true"},
		{"false", "false"},
		{"iota", "iota"},

		// Special - unchanged even when enabled
		{"_", "_"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := vm.TransformVarName(tt.input)
			if result != tt.want {
				t.Errorf("TransformVarName(%q) = %q; want %q (when enabled)", tt.input, result, tt.want)
			}
		})
	}
}

func TestVarMapper_RegisterUserVar(t *testing.T) {
	vm := NewVarMapper()

	varName := "myCustomVar"
	vm.RegisterUserVar(varName)

	if !vm.userVars[varName] {
		t.Errorf("RegisterUserVar(%q) did not register variable", varName)
	}
}

func TestVarMapper_EnableDisable(t *testing.T) {
	vm := NewVarMapper()

	// Should start disabled
	if vm.IsEnabled() {
		t.Error("VarMapper should start disabled")
	}

	vm.Enable()
	if !vm.IsEnabled() {
		t.Error("VarMapper.Enable() did not enable")
	}

	vm.Disable()
	if vm.IsEnabled() {
		t.Error("VarMapper.Disable() did not disable")
	}
}

func TestVarMapper_TransformVarNameReverse(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable()

	tests := []struct {
		input string
		want  string
	}{
		// camelCase -> snake_case (Go -> Moxie direction)
		{"myVar", "my_var"},
		{"userName", "user_name"},
		{"userID", "user_id"},
		{"httpServer", "http_server"},
		{"maxCount", "max_count"},

		// PascalCase -> snake_case (exported variables)
		{"MyVar", "my_var"},
		{"UserName", "user_name"},

		// Builtins - unchanged
		{"nil", "nil"},
		{"true", "true"},
		{"false", "false"},

		// Special - unchanged
		{"_", "_"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := vm.TransformVarNameReverse(tt.input)
			if result != tt.want {
				t.Errorf("TransformVarNameReverse(%q) = %q; want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestVarMapper_LoopVariables(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable()

	// Common loop variable names should not be transformed
	loopVars := []string{
		"i", "j", "k", "x", "y", "z",
	}

	for _, v := range loopVars {
		t.Run(v, func(t *testing.T) {
			result := vm.TransformVarName(v)
			if result != v {
				t.Errorf("TransformVarName(%q) = %q; want %q (single letter should not transform)", v, result, v)
			}
		})
	}
}

func TestVarMapper_CommonVariablePatterns(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Common Go variable patterns
		{"error variable", "err", "err"},
		{"context", "ctx", "ctx"},
		{"ok variable", "ok", "ok"},

		// Multi-word variables
		{"user data", "user_data", "userData"},
		{"request body", "request_body", "requestBody"},
		{"response code", "response_code", "responseCode"},

		// With acronyms
		{"user ID", "user_id", "userID"},
		{"http client", "http_client", "httpClient"},
		{"json data", "json_data", "jsonData"},
		{"xml parser", "xml_parser", "xmlParser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vm.TransformVarName(tt.input)
			if result != tt.want {
				t.Errorf("TransformVarName(%q) = %q; want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestVarMapper_ExportedVsUnexported(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Unexported (lowercase first letter)
		{"unexported var", "my_var", "myVar"},
		{"unexported data", "user_data", "userData"},

		// Exported (uppercase first letter) - rare for variables but possible
		{"exported var", "My_var", "MyVar"},
		{"exported data", "User_data", "UserData"},

		// Already correct case
		{"already unexported", "myVar", "myVar"},
		{"already exported", "MyVar", "MyVar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vm.TransformVarName(tt.input)
			if result != tt.want {
				t.Errorf("TransformVarName(%q) = %q; want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestVarMapper_ConstantNames(t *testing.T) {
	vm := NewVarMapper()
	vm.Enable()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Constants are often exported (PascalCase)
		{"max size", "Max_size", "MaxSize"},
		{"default timeout", "Default_timeout", "DefaultTimeout"},
		{"buffer size", "Buffer_size", "BufferSize"},

		// Unexported constants
		{"max retry", "max_retry", "maxRetry"},
		{"default port", "default_port", "defaultPort"},

		// With acronyms
		{"max http retries", "max_http_retries", "maxHTTPRetries"},
		{"default api url", "default_api_url", "defaultAPIURL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vm.TransformVarName(tt.input)
			if result != tt.want {
				t.Errorf("TransformVarName(%q) = %q; want %q", tt.input, result, tt.want)
			}
		})
	}
}
