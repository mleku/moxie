package main

import "testing"

func TestTypeMapper_IsBuiltin(t *testing.T) {
	tm := NewTypeMapper()

	builtins := []string{
		"bool", "byte", "string", "int", "int32", "int64",
		"uint", "uint32", "uint64", "float32", "float64",
		"complex64", "complex128", "error", "rune", "uintptr",
		"any", "comparable",
	}

	for _, typ := range builtins {
		if !tm.IsBuiltin(typ) {
			t.Errorf("IsBuiltin(%q) = false; want true", typ)
		}
	}

	notBuiltins := []string{
		"MyType", "HTTPServer", "User", "Request",
	}

	for _, typ := range notBuiltins {
		if tm.IsBuiltin(typ) {
			t.Errorf("IsBuiltin(%q) = true; want false", typ)
		}
	}
}

func TestTypeMapper_ShouldTransform(t *testing.T) {
	tm := NewTypeMapper()

	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		// Builtins - never transform
		{"builtin int", "int", false},
		{"builtin string", "string", false},
		{"builtin error", "error", false},

		// Empty/single char - never transform
		{"empty", "", false},
		{"single char", "T", false},

		// User types - would transform if enabled
		{"user type", "MyType", false}, // false because tm.enabled = false
		{"http server", "HTTPServer", false},

		// Stdlib types - never transform
		{"stdlib", "Request", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.ShouldTransform(tt.typeName)
			if got != tt.want {
				t.Errorf("ShouldTransform(%q) = %v; want %v", tt.typeName, got, tt.want)
			}
		})
	}
}

func TestTypeMapper_ShouldTransformWhenEnabled(t *testing.T) {
	tm := NewTypeMapper()
	tm.Enable() // Enable transformation

	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		// Builtins - still never transform
		{"builtin int", "int", false},
		{"builtin string", "string", false},

		// User types - now should transform
		{"user type", "MyType", true},
		{"http server", "HTTPServer", true},
		{"user id", "UserID", true},

		// Stdlib types - never transform
		{"stdlib", "Request", false},
		{"stdlib", "Response", false},

		// Edge cases - never transform
		{"empty", "", false},
		{"single char", "T", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.ShouldTransform(tt.typeName)
			if got != tt.want {
				t.Errorf("ShouldTransform(%q) = %v; want %v (enabled=true)", tt.typeName, got, tt.want)
			}
		})
	}
}

func TestTypeMapper_TransformTypeName_Disabled(t *testing.T) {
	tm := NewTypeMapper()
	// Disabled by default

	tests := []string{
		"MyType",
		"HTTPServer",
		"UserID",
		"XMLParser",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := tm.TransformTypeName(input)
			if result != input {
				t.Errorf("TransformTypeName(%q) = %q; want %q (when disabled)", input, result, input)
			}
		})
	}
}

func TestTypeMapper_TransformTypeName_Enabled(t *testing.T) {
	tm := NewTypeMapper()
	tm.Enable()

	tests := []struct {
		input string
		want  string
	}{
		// snake_case -> PascalCase (Moxie -> Go direction)
		{"my_type", "MyType"},
		{"http_server", "HTTPServer"},
		{"user_id", "UserID"},
		{"xml_parser", "XMLParser"},

		// Already PascalCase - unchanged
		{"MyType", "MyType"},
		{"HTTPServer", "HTTPServer"},

		// Builtins - unchanged even when enabled
		{"int", "int"},
		{"string", "string"},
		{"error", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tm.TransformTypeName(tt.input)
			if result != tt.want {
				t.Errorf("TransformTypeName(%q) = %q; want %q (when enabled)", tt.input, result, tt.want)
			}
		})
	}
}

func TestTypeMapper_RegisterUserType(t *testing.T) {
	tm := NewTypeMapper()

	typeName := "MyCustomType"
	tm.RegisterUserType(typeName)

	if !tm.userTypes[typeName] {
		t.Errorf("RegisterUserType(%q) did not register type", typeName)
	}
}

func TestTypeMapper_EnableDisable(t *testing.T) {
	tm := NewTypeMapper()

	// Should start disabled
	if tm.IsEnabled() {
		t.Error("TypeMapper should start disabled")
	}

	tm.Enable()
	if !tm.IsEnabled() {
		t.Error("TypeMapper.Enable() did not enable")
	}

	tm.Disable()
	if tm.IsEnabled() {
		t.Error("TypeMapper.Disable() did not disable")
	}
}
