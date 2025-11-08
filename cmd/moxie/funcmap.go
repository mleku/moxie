// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
)

// Function and method name mapping and transformation

// FuncMapper handles function and method name transformations
type FuncMapper struct {
	// Whether to actually transform function names
	// Set to false to keep PascalCase/camelCase (current default)
	enabled bool

	// Tracks user-defined functions we've seen
	userFuncs map[string]bool

	// Built-in functions that should never be transformed
	builtinFuncs map[string]bool

	// Standard library functions that should never be transformed
	// These are typically from imported packages
	stdlibFuncs map[string]bool

	// Special functions that should never be transformed
	specialFuncs map[string]bool
}

// NewFuncMapper creates a new function mapper
func NewFuncMapper() *FuncMapper {
	fm := &FuncMapper{
		enabled:      false, // Keep PascalCase/camelCase for now
		userFuncs:    make(map[string]bool),
		builtinFuncs: make(map[string]bool),
		stdlibFuncs:  make(map[string]bool),
		specialFuncs: make(map[string]bool),
	}

	// Initialize builtin functions
	builtins := []string{
		"append", "cap", "close", "complex", "copy", "delete",
		"imag", "len", "make", "new", "panic", "print", "println",
		"real", "recover",
		"clear", // Go 1.21+
		"max", "min", // Go 1.21+
	}

	for _, f := range builtins {
		fm.builtinFuncs[f] = true
	}

	// Special functions (Go keywords/special names)
	special := []string{
		"init",   // Package initialization
		"main",   // Program entry point
		"Error",  // error interface method
		"String", // Stringer interface method
	}

	for _, f := range special {
		fm.specialFuncs[f] = true
	}

	return fm
}

// IsBuiltin checks if a function name is a built-in function
func (fm *FuncMapper) IsBuiltin(name string) bool {
	return fm.builtinFuncs[name]
}

// IsSpecial checks if a function name is a special function
func (fm *FuncMapper) IsSpecial(name string) bool {
	return fm.specialFuncs[name]
}

// ShouldTransform determines if a function name should be transformed
func (fm *FuncMapper) ShouldTransform(name string) bool {
	if !fm.enabled {
		return false // Transformation disabled, keep current naming
	}

	// Don't transform builtins
	if fm.IsBuiltin(name) {
		return false
	}

	// Don't transform special functions
	if fm.IsSpecial(name) {
		return false
	}

	// Don't transform empty names
	if name == "" {
		return false
	}

	// Don't transform single letters (rare for functions)
	if len(name) == 1 {
		return false
	}

	return true
}

// RegisterUserFunc registers a user-defined function
func (fm *FuncMapper) RegisterUserFunc(name string) {
	fm.userFuncs[name] = true
}

// TransformFuncName transforms a function name from Moxie to Go
// When enabled=false (current), this is a no-op and keeps PascalCase/camelCase
func (fm *FuncMapper) TransformFuncName(moxieName string) string {
	if !fm.ShouldTransform(moxieName) {
		return moxieName
	}

	// When enabled, would convert snake_case → camelCase/PascalCase
	// Preserve export status
	return preserveExportStatus(moxieName, toPascalCase)
}

// TransformFuncNameReverse transforms a function name from Go to Moxie
// When enabled=false (current), this is a no-op
func (fm *FuncMapper) TransformFuncNameReverse(goName string) string {
	if !fm.ShouldTransform(goName) {
		return goName
	}

	// When enabled, would convert PascalCase/camelCase → snake_case
	return toSnakeCase(goName)
}

// Enable enables function name transformation
// Currently kept disabled to maintain PascalCase/camelCase
func (fm *FuncMapper) Enable() {
	fm.enabled = true
}

// Disable disables function name transformation (default)
func (fm *FuncMapper) Disable() {
	fm.enabled = false
}

// IsEnabled returns whether function transformation is enabled
func (fm *FuncMapper) IsEnabled() bool {
	return fm.enabled
}

// transformFuncDecl transforms a function declaration
func (fm *FuncMapper) transformFuncDecl(decl *ast.FuncDecl) {
	if decl == nil || decl.Name == nil {
		return
	}

	// Transform the function name
	if fm.ShouldTransform(decl.Name.Name) {
		fm.RegisterUserFunc(decl.Name.Name)
		decl.Name.Name = fm.TransformFuncName(decl.Name.Name)
	}

	// Note: Function parameters and returns are handled by typeMap
}

// transformCallExpr transforms a function call expression
func (fm *FuncMapper) transformCallExpr(call *ast.CallExpr) {
	if call == nil || call.Fun == nil {
		return
	}

	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Simple function call: myFunc()
		if fm.ShouldTransform(fun.Name) {
			fun.Name = fm.TransformFuncName(fun.Name)
		}

	case *ast.SelectorExpr:
		// Method call or qualified function: obj.Method(), pkg.Func()
		// Only transform the selector (method/function name), not the object/package
		if fm.ShouldTransform(fun.Sel.Name) {
			fun.Sel.Name = fm.TransformFuncName(fun.Sel.Name)
		}

	// Other cases (type conversions, etc.) don't need transformation
	}
}

// transformFuncLit transforms a function literal (anonymous function)
func (fm *FuncMapper) transformFuncLit(lit *ast.FuncLit) {
	// Function literals don't have names, but their parameters and
	// return types are handled by typeMap
	// Nothing to do here for name transformation
}

// Global function mapper instance
var funcMap = NewFuncMapper()
