// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
)

// Type name mapping and transformation for Moxie types

// TypeMapper handles type name transformations
type TypeMapper struct {
	// Whether to actually transform type names
	// Set to false to keep PascalCase (current default)
	enabled bool

	// Tracks user-defined types we've seen
	userTypes map[string]bool

	// Built-in types that should never be transformed
	builtinTypes map[string]bool

	// Standard library types that should never be transformed
	stdlibTypes map[string]bool
}

// NewTypeMapper creates a new type mapper
func NewTypeMapper() *TypeMapper {
	tm := &TypeMapper{
		enabled:      false, // Keep PascalCase for now
		userTypes:    make(map[string]bool),
		builtinTypes: make(map[string]bool),
		stdlibTypes:  make(map[string]bool),
	}

	// Initialize builtin types
	builtins := []string{
		"bool",
		"byte",
		"complex64", "complex128",
		"error",
		"float32", "float64",
		"int", "int8", "int16", "int32", "int64",
		"rune",
		"string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"any", // Go 1.18+
		"comparable", // Go 1.18+
	}

	for _, t := range builtins {
		tm.builtinTypes[t] = true
	}

	// Common stdlib types to never transform
	stdlib := []string{
		// These are qualified (pkg.Type) so just tracking for reference
		"Request", "Response", "Client", "Server", // http
		"Encoder", "Decoder", "Marshaler", "Unmarshaler", // json/xml
		"Context", // context
		"Time", "Duration", // time
		"File", // os
		"Reader", "Writer", "ReadWriter", // io
		"Error", // errors (though error is builtin)
	}

	for _, t := range stdlib {
		tm.stdlibTypes[t] = true
	}

	return tm
}

// IsBuiltin checks if a type name is a built-in type
func (tm *TypeMapper) IsBuiltin(name string) bool {
	return tm.builtinTypes[name]
}

// IsStdlib checks if a type name is a stdlib type
func (tm *TypeMapper) IsStdlib(name string) bool {
	return tm.stdlibTypes[name]
}

// ShouldTransform determines if a type name should be transformed
func (tm *TypeMapper) ShouldTransform(name string) bool {
	if !tm.enabled {
		return false // Transformation disabled, keep PascalCase
	}

	// Don't transform builtins
	if tm.IsBuiltin(name) {
		return false
	}

	// Don't transform stdlib types
	if tm.IsStdlib(name) {
		return false
	}

	// Don't transform empty names
	if name == "" {
		return false
	}

	// Don't transform single letters (usually type parameters)
	if len(name) == 1 {
		return false
	}

	return true
}

// RegisterUserType registers a user-defined type
func (tm *TypeMapper) RegisterUserType(name string) {
	tm.userTypes[name] = true
}

// TransformTypeName transforms a type name from Moxie to Go
// When enabled=false (current), this is a no-op and keeps PascalCase
func (tm *TypeMapper) TransformTypeName(moxieName string) string {
	if !tm.ShouldTransform(moxieName) {
		return moxieName
	}

	// When enabled, would convert snake_case → PascalCase
	// For now, just return unchanged
	return toPascalCase(moxieName)
}

// TransformTypeNameReverse transforms a type name from Go to Moxie
// When enabled=false (current), this is a no-op
func (tm *TypeMapper) TransformTypeNameReverse(goName string) string {
	if !tm.ShouldTransform(goName) {
		return goName
	}

	// When enabled, would convert PascalCase → snake_case
	// For now, just return unchanged
	return toSnakeCase(goName)
}

// Enable enables type name transformation
// Currently kept disabled to maintain PascalCase
func (tm *TypeMapper) Enable() {
	tm.enabled = true
}

// Disable disables type name transformation (default)
func (tm *TypeMapper) Disable() {
	tm.enabled = false
}

// IsEnabled returns whether type transformation is enabled
func (tm *TypeMapper) IsEnabled() bool {
	return tm.enabled
}

// transformTypeExpr transforms type expressions in the AST
// This recursively handles all type expression variants
func (tm *TypeMapper) transformTypeExpr(expr ast.Expr) {
	if expr == nil {
		return
	}

	switch t := expr.(type) {
	case *ast.Ident:
		// Simple type name: MyType
		if tm.ShouldTransform(t.Name) {
			t.Name = tm.TransformTypeName(t.Name)
		}

	case *ast.StarExpr:
		// Pointer type: *MyType
		tm.transformTypeExpr(t.X)

	case *ast.ArrayType:
		// Array/slice type: []MyType, [10]MyType
		tm.transformTypeExpr(t.Elt)

	case *ast.MapType:
		// Map type: map[KeyType]ValueType
		tm.transformTypeExpr(t.Key)
		tm.transformTypeExpr(t.Value)

	case *ast.ChanType:
		// Channel type: chan MyType, <-chan MyType
		tm.transformTypeExpr(t.Value)

	case *ast.FuncType:
		// Function type: func(MyType) MyType
		if t.Params != nil {
			tm.transformFieldList(t.Params)
		}
		if t.Results != nil {
			tm.transformFieldList(t.Results)
		}

	case *ast.StructType:
		// Struct type: struct { Field MyType }
		if t.Fields != nil {
			tm.transformFieldList(t.Fields)
		}

	case *ast.InterfaceType:
		// Interface type: interface { Method() MyType }
		if t.Methods != nil {
			tm.transformFieldList(t.Methods)
		}

	case *ast.SelectorExpr:
		// Qualified type: pkg.MyType
		// Don't transform the selector (package name)
		// Only transform if it's a type (not common for qualified types)
		// Most qualified types are from stdlib, so we skip them

	case *ast.ParenExpr:
		// Parenthesized type: (MyType)
		tm.transformTypeExpr(t.X)

	case *ast.IndexExpr:
		// Generic type: MyType[T]
		tm.transformTypeExpr(t.X)
		tm.transformTypeExpr(t.Index)

	case *ast.IndexListExpr:
		// Generic type with multiple params: MyType[T, U]
		tm.transformTypeExpr(t.X)
		for _, index := range t.Indices {
			tm.transformTypeExpr(index)
		}

	// Other cases (Ellipsis, etc.) are rare for type expressions
	}
}

// transformFieldList transforms a field list (parameters, results, struct fields)
func (tm *TypeMapper) transformFieldList(fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		tm.transformTypeExpr(field.Type)
	}
}

// Global type mapper instance
var typeMap = NewTypeMapper()
