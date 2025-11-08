// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
)

// Variable and constant name mapping and transformation

// VarMapper handles variable and constant name transformations
type VarMapper struct {
	// Whether to actually transform variable/constant names
	// Set to false to keep camelCase (current default)
	enabled bool

	// Tracks user-defined variables/constants we've seen
	userVars map[string]bool

	// Built-in identifiers that should never be transformed
	builtinVars map[string]bool

	// Standard library constants that should never be transformed
	stdlibConsts map[string]bool

	// Special identifiers that should never be transformed
	specialVars map[string]bool
}

// NewVarMapper creates a new variable/constant mapper
func NewVarMapper() *VarMapper {
	vm := &VarMapper{
		enabled:      false, // Keep camelCase for now
		userVars:     make(map[string]bool),
		builtinVars:  make(map[string]bool),
		stdlibConsts: make(map[string]bool),
		specialVars:  make(map[string]bool),
	}

	// Initialize builtin identifiers
	builtins := []string{
		"nil",   // Nil value
		"true",  // Boolean true
		"false", // Boolean false
		"iota",  // Constant generator
	}

	for _, v := range builtins {
		vm.builtinVars[v] = true
	}

	// Special identifiers (blank identifier, etc.)
	special := []string{
		"_", // Blank identifier
	}

	for _, v := range special {
		vm.specialVars[v] = true
	}

	return vm
}

// IsBuiltin checks if a variable name is a built-in identifier
func (vm *VarMapper) IsBuiltin(name string) bool {
	return vm.builtinVars[name]
}

// IsSpecial checks if a variable name is a special identifier
func (vm *VarMapper) IsSpecial(name string) bool {
	return vm.specialVars[name]
}

// ShouldTransform determines if a variable/constant name should be transformed
func (vm *VarMapper) ShouldTransform(name string) bool {
	if !vm.enabled {
		return false // Transformation disabled, keep current naming
	}

	// Don't transform builtins
	if vm.IsBuiltin(name) {
		return false
	}

	// Don't transform special identifiers
	if vm.IsSpecial(name) {
		return false
	}

	// Don't transform empty names
	if name == "" {
		return false
	}

	// Don't transform single letters (common for loop variables)
	if len(name) == 1 {
		return false
	}

	return true
}

// RegisterUserVar registers a user-defined variable/constant
func (vm *VarMapper) RegisterUserVar(name string) {
	vm.userVars[name] = true
}

// TransformVarName transforms a variable/constant name from Moxie to Go
// When enabled=false (current), this is a no-op and keeps camelCase
func (vm *VarMapper) TransformVarName(moxieName string) string {
	if !vm.ShouldTransform(moxieName) {
		return moxieName
	}

	// When enabled, would convert snake_case → camelCase
	// Variables are typically unexported (start lowercase)
	return preserveExportStatus(moxieName, toPascalCase)
}

// TransformVarNameReverse transforms a variable/constant name from Go to Moxie
// When enabled=false (current), this is a no-op
func (vm *VarMapper) TransformVarNameReverse(goName string) string {
	if !vm.ShouldTransform(goName) {
		return goName
	}

	// When enabled, would convert camelCase → snake_case
	return toSnakeCase(goName)
}

// Enable enables variable/constant name transformation
// Currently kept disabled to maintain camelCase
func (vm *VarMapper) Enable() {
	vm.enabled = true
}

// Disable disables variable/constant name transformation (default)
func (vm *VarMapper) Disable() {
	vm.enabled = false
}

// IsEnabled returns whether variable/constant transformation is enabled
func (vm *VarMapper) IsEnabled() bool {
	return vm.enabled
}

// transformValueSpec transforms a variable or constant declaration
func (vm *VarMapper) transformValueSpec(spec *ast.ValueSpec) {
	if spec == nil {
		return
	}

	// Transform each name in the declaration
	for _, name := range spec.Names {
		if name != nil && vm.ShouldTransform(name.Name) {
			vm.RegisterUserVar(name.Name)
			name.Name = vm.TransformVarName(name.Name)
		}
	}

	// Note: Type transformation is handled by typeMap
}

// transformField transforms a struct field or function parameter
func (vm *VarMapper) transformField(field *ast.Field) {
	if field == nil {
		return
	}

	// Transform field names
	for _, name := range field.Names {
		if name != nil && vm.ShouldTransform(name.Name) {
			vm.RegisterUserVar(name.Name)
			name.Name = vm.TransformVarName(name.Name)
		}
	}

	// Note: Type transformation is handled by typeMap
}

// transformFieldList transforms a list of fields (parameters, results, struct fields)
func (vm *VarMapper) transformFieldList(fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		vm.transformField(field)
	}
}

// transformAssignStmt transforms an assignment statement
func (vm *VarMapper) transformAssignStmt(stmt *ast.AssignStmt) {
	if stmt == nil {
		return
	}

	// For := (short variable declaration), transform LHS identifiers
	if stmt.Tok.String() == ":=" {
		for _, lhs := range stmt.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				if vm.ShouldTransform(ident.Name) {
					vm.RegisterUserVar(ident.Name)
					ident.Name = vm.TransformVarName(ident.Name)
				}
			}
		}
	}

	// Transform identifiers in RHS (references to variables)
	for _, rhs := range stmt.Rhs {
		vm.transformExpr(rhs)
	}
}

// transformExpr transforms identifiers in expressions
func (vm *VarMapper) transformExpr(expr ast.Expr) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.Ident:
		// Transform variable references
		if vm.ShouldTransform(e.Name) {
			e.Name = vm.TransformVarName(e.Name)
		}

	case *ast.SelectorExpr:
		// Transform the X part (object), but not Sel (field/method)
		vm.transformExpr(e.X)

	case *ast.IndexExpr:
		// Transform array/slice/map index expressions
		vm.transformExpr(e.X)
		vm.transformExpr(e.Index)

	case *ast.CallExpr:
		// Transform function call arguments
		vm.transformExpr(e.Fun)
		for _, arg := range e.Args {
			vm.transformExpr(arg)
		}

	case *ast.UnaryExpr:
		// Transform unary expressions
		vm.transformExpr(e.X)

	case *ast.BinaryExpr:
		// Transform binary expressions
		vm.transformExpr(e.X)
		vm.transformExpr(e.Y)

	case *ast.ParenExpr:
		// Transform parenthesized expressions
		vm.transformExpr(e.X)

	case *ast.CompositeLit:
		// Transform composite literals
		vm.transformExpr(e.Type)
		for _, elt := range e.Elts {
			vm.transformExpr(elt)
		}

	case *ast.KeyValueExpr:
		// Transform key-value expressions
		vm.transformExpr(e.Key)
		vm.transformExpr(e.Value)

	// Other expression types don't need variable transformation
	}
}

// transformStmt transforms identifiers in statements
func (vm *VarMapper) transformStmt(stmt ast.Stmt) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		vm.transformAssignStmt(s)

	case *ast.DeclStmt:
		// Variable/constant declarations in statements
		if decl, ok := s.Decl.(*ast.GenDecl); ok {
			for _, spec := range decl.Specs {
				if vspec, ok := spec.(*ast.ValueSpec); ok {
					vm.transformValueSpec(vspec)
				}
			}
		}

	case *ast.ExprStmt:
		// Expression statements
		vm.transformExpr(s.X)

	case *ast.IfStmt:
		// If statements
		if s.Init != nil {
			vm.transformStmt(s.Init)
		}
		vm.transformExpr(s.Cond)
		vm.transformBlockStmt(s.Body)
		if s.Else != nil {
			vm.transformStmt(s.Else)
		}

	case *ast.ForStmt:
		// For loops
		if s.Init != nil {
			vm.transformStmt(s.Init)
		}
		if s.Cond != nil {
			vm.transformExpr(s.Cond)
		}
		if s.Post != nil {
			vm.transformStmt(s.Post)
		}
		vm.transformBlockStmt(s.Body)

	case *ast.RangeStmt:
		// Range loops
		if s.Key != nil {
			if ident, ok := s.Key.(*ast.Ident); ok && vm.ShouldTransform(ident.Name) {
				vm.RegisterUserVar(ident.Name)
				ident.Name = vm.TransformVarName(ident.Name)
			}
		}
		if s.Value != nil {
			if ident, ok := s.Value.(*ast.Ident); ok && vm.ShouldTransform(ident.Name) {
				vm.RegisterUserVar(ident.Name)
				ident.Name = vm.TransformVarName(ident.Name)
			}
		}
		vm.transformExpr(s.X)
		vm.transformBlockStmt(s.Body)

	case *ast.ReturnStmt:
		// Return statements
		for _, result := range s.Results {
			vm.transformExpr(result)
		}

	case *ast.BlockStmt:
		vm.transformBlockStmt(s)

	// Other statement types handled as needed
	}
}

// transformBlockStmt transforms a block of statements
func (vm *VarMapper) transformBlockStmt(block *ast.BlockStmt) {
	if block == nil {
		return
	}

	for _, stmt := range block.List {
		vm.transformStmt(stmt)
	}
}

// Global variable mapper instance
var varMap = NewVarMapper()
