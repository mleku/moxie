package main

import "testing"

func TestFuncMapper_IsBuiltin(t *testing.T) {
	fm := NewFuncMapper()

	builtins := []string{
		"append", "cap", "close", "complex", "copy", "delete",
		"imag", "len", "make", "new", "panic", "print", "println",
		"real", "recover", "clear", "max", "min",
	}

	for _, fn := range builtins {
		if !fm.IsBuiltin(fn) {
			t.Errorf("IsBuiltin(%q) = false; want true", fn)
		}
	}

	notBuiltins := []string{
		"MyFunc", "ProcessData", "getUserID", "parseHTTPRequest",
	}

	for _, fn := range notBuiltins {
		if fm.IsBuiltin(fn) {
			t.Errorf("IsBuiltin(%q) = true; want false", fn)
		}
	}
}

func TestFuncMapper_IsSpecial(t *testing.T) {
	fm := NewFuncMapper()

	special := []string{
		"init", "main", "Error", "String",
	}

	for _, fn := range special {
		if !fm.IsSpecial(fn) {
			t.Errorf("IsSpecial(%q) = false; want true", fn)
		}
	}

	notSpecial := []string{
		"MyFunc", "ProcessData", "append", "len",
	}

	for _, fn := range notSpecial {
		if fm.IsSpecial(fn) {
			t.Errorf("IsSpecial(%q) = true; want false", fn)
		}
	}
}

func TestFuncMapper_ShouldTransform(t *testing.T) {
	fm := NewFuncMapper()

	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		// Builtins - never transform
		{"builtin append", "append", false},
		{"builtin len", "len", false},
		{"builtin make", "make", false},

		// Special - never transform
		{"special init", "init", false},
		{"special main", "main", false},
		{"special Error", "Error", false},

		// Empty/single char - never transform
		{"empty", "", false},
		{"single char", "f", false},

		// User functions - would transform if enabled
		{"user func", "MyFunc", false}, // false because fm.enabled = false
		{"get user", "GetUser", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.ShouldTransform(tt.funcName)
			if got != tt.want {
				t.Errorf("ShouldTransform(%q) = %v; want %v", tt.funcName, got, tt.want)
			}
		})
	}
}

func TestFuncMapper_ShouldTransformWhenEnabled(t *testing.T) {
	fm := NewFuncMapper()
	fm.Enable() // Enable transformation

	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		// Builtins - still never transform
		{"builtin append", "append", false},
		{"builtin len", "len", false},

		// Special - still never transform
		{"special init", "init", false},
		{"special main", "main", false},

		// User functions - now should transform
		{"user func", "MyFunc", true},
		{"get user", "GetUser", true},
		{"parse request", "ParseRequest", true},

		// Edge cases - never transform
		{"empty", "", false},
		{"single char", "f", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.ShouldTransform(tt.funcName)
			if got != tt.want {
				t.Errorf("ShouldTransform(%q) = %v; want %v (enabled=true)", tt.funcName, got, tt.want)
			}
		})
	}
}

func TestFuncMapper_TransformFuncName_Disabled(t *testing.T) {
	fm := NewFuncMapper()
	// Disabled by default

	tests := []string{
		"MyFunc",
		"GetUser",
		"ParseHTTPRequest",
		"processData",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := fm.TransformFuncName(input)
			if result != input {
				t.Errorf("TransformFuncName(%q) = %q; want %q (when disabled)", input, result, input)
			}
		})
	}
}

func TestFuncMapper_TransformFuncName_Enabled(t *testing.T) {
	fm := NewFuncMapper()
	fm.Enable()

	tests := []struct {
		input string
		want  string
	}{
		// snake_case -> PascalCase/camelCase (Moxie -> Go direction)
		// Exported functions (start with uppercase)
		{"Get_user", "GetUser"},                       // Should be GetUser (exported)
		{"Parse_http_request", "ParseHTTPRequest"},    // Exported, with acronym
		{"Process_data", "ProcessData"},               // Exported

		// Unexported functions (start with lowercase)
		{"get_user", "getUser"},                       // Should be getUser (unexported)
		{"parse_http_request", "parseHTTPRequest"},    // Unexported, with acronym
		{"process_data", "processData"},               // Unexported
		{"parse_request", "parseRequest"},             // Unexported

		// Already PascalCase/camelCase - unchanged
		{"MyFunc", "MyFunc"},
		{"GetUser", "GetUser"},
		{"parseData", "parseData"},

		// Builtins - unchanged even when enabled
		{"append", "append"},
		{"len", "len"},

		// Special - unchanged even when enabled
		{"init", "init"},
		{"main", "main"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := fm.TransformFuncName(tt.input)
			if result != tt.want {
				t.Errorf("TransformFuncName(%q) = %q; want %q (when enabled)", tt.input, result, tt.want)
			}
		})
	}
}

func TestFuncMapper_RegisterUserFunc(t *testing.T) {
	fm := NewFuncMapper()

	funcName := "MyCustomFunc"
	fm.RegisterUserFunc(funcName)

	if !fm.userFuncs[funcName] {
		t.Errorf("RegisterUserFunc(%q) did not register function", funcName)
	}
}

func TestFuncMapper_EnableDisable(t *testing.T) {
	fm := NewFuncMapper()

	// Should start disabled
	if fm.IsEnabled() {
		t.Error("FuncMapper should start disabled")
	}

	fm.Enable()
	if !fm.IsEnabled() {
		t.Error("FuncMapper.Enable() did not enable")
	}

	fm.Disable()
	if fm.IsEnabled() {
		t.Error("FuncMapper.Disable() did not disable")
	}
}

func TestFuncMapper_TransformFuncNameReverse(t *testing.T) {
	fm := NewFuncMapper()
	fm.Enable()

	tests := []struct {
		input string
		want  string
	}{
		// PascalCase -> snake_case (Go -> Moxie direction)
		{"GetUser", "get_user"},
		{"ParseHTTPRequest", "parse_http_request"},
		{"ProcessData", "process_data"},

		// camelCase -> snake_case
		{"getUserID", "get_user_id"},
		{"parseRequest", "parse_request"},

		// Builtins - unchanged
		{"append", "append"},
		{"len", "len"},

		// Special - unchanged
		{"init", "init"},
		{"main", "main"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := fm.TransformFuncNameReverse(tt.input)
			if result != tt.want {
				t.Errorf("TransformFuncNameReverse(%q) = %q; want %q", tt.input, result, tt.want)
			}
		})
	}
}
