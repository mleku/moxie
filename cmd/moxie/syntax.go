// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

// SyntaxTransformer handles Moxie→Go syntax transformations
type SyntaxTransformer struct {
	errors           []error
	needsRuntimeImport bool
}

// NewSyntaxTransformer creates a new syntax transformer
func NewSyntaxTransformer() *SyntaxTransformer {
	return &SyntaxTransformer{
		errors:           make([]error, 0),
		needsRuntimeImport: false,
	}
}

// Transform applies all Moxie→Go syntax transformations to an AST
func (st *SyntaxTransformer) Transform(file *ast.File) error {
	st.errors = nil
	st.needsRuntimeImport = false

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			st.transformAssignStmt(node)
		case *ast.CallExpr:
			st.transformCallExpr(node)
		case *ast.UnaryExpr:
			st.transformUnaryExpr(node)
		case *ast.CompositeLit:
			st.transformCompositeLit(node)
		}
		return true
	})

	// Add runtime import if needed
	if st.needsRuntimeImport {
		st.addRuntimeImport(file)
	}

	if len(st.errors) > 0 {
		return st.errors[0] // Return first error
	}
	return nil
}

// transformAssignStmt handles assignment statement transformations
// Specifically for append() assignments: s = append(s, x) -> *s = append(*s, x)
func (st *SyntaxTransformer) transformAssignStmt(node *ast.AssignStmt) {
	// Check if RHS contains append() call
	for i, rhs := range node.Rhs {
		if call, ok := rhs.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "append" {
				// This is an append() call
				// Transform: s = append(s, items...) -> *s = append(*s, items...)
				if len(call.Args) > 0 {
					// Dereference first argument: append(s, ...) -> append(*s, ...)
					call.Args[0] = &ast.StarExpr{X: call.Args[0]}
				}

				// Dereference LHS: s -> *s
				if i < len(node.Lhs) {
					node.Lhs[i] = &ast.StarExpr{X: node.Lhs[i]}
				}
			}
		}
	}
}

// transformCallExpr handles function call transformations
func (st *SyntaxTransformer) transformCallExpr(node *ast.CallExpr) {
	// Check if this is a built-in function call
	if ident, ok := node.Fun.(*ast.Ident); ok {
		switch ident.Name {
		case "make":
			// make() is not allowed in Moxie
			st.errors = append(st.errors, fmt.Errorf(
				"make() is not available in Moxie; use &[]T{}, &map[K]V{}, or &chan T{} instead",
			))

		case "append":
			// append() in Moxie takes *[]T and returns *[]T
			// In Go, append takes []T and returns []T
			// We need a runtime wrapper that handles the pointer semantics
			// For now, leave as-is and implement runtime.Append later
			// TODO: Transform to moxie.Append() or implement proper wrapper

		case "grow":
			// grow() is a Moxie built-in
			// Transform: grow(s, n) -> moxie.Grow(s, n)
			st.transformToRuntimeCall(node, "Grow")

		case "clone":
			// clone() is a Moxie built-in
			// Transform: clone(v) -> moxie.CloneSlice(v) or moxie.CloneMap(v)
			// For now, use CloneSlice - type-specific versions can be optimized later
			st.transformToRuntimeCall(node, "CloneSlice") // TODO: Detect type and use appropriate function

		case "clear":
			// clear() exists in Go 1.21+
			// In Moxie, slices/maps are pointers, so clear(*map[K]V) needs to become clear((*map[K]V))
			// We need to dereference pointer types
			st.transformClearCall(node)

		case "free":
			// free() is a Moxie built-in that provides GC hints
			// Transform: free(v) -> moxie.FreeSlice(v)
			st.transformToRuntimeCall(node, "FreeSlice") // TODO: Detect type and use appropriate function
		}
	}
}

// transformUnaryExpr handles unary expression transformations
// This handles address-of operator for composite literals: &[]T{}, &map[K]V{}, &chan T{}
func (st *SyntaxTransformer) transformUnaryExpr(node *ast.UnaryExpr) {
	if node.Op != token.AND {
		return
	}

	// Check if this is &compositeLit
	compLit, ok := node.X.(*ast.CompositeLit)
	if !ok {
		return
	}

	// Check the type
	switch typ := compLit.Type.(type) {
	case *ast.ArrayType:
		// &[]T{...} - slice literal with address-of
		// This is already valid Go syntax, no transformation needed
		if typ.Len == nil {
			// This is a slice type (len == nil means []T not [N]T)
			// &[]T{1, 2, 3} is valid in Go 1.18+
		}

	case *ast.MapType:
		// &map[K]V{...} - map literal with address-of
		// This is already valid Go syntax, no transformation needed

	case *ast.ChanType:
		// &chan T{...} - channel literal with address-of
		// This is NOT valid Go syntax and needs transformation
		// In Moxie: &chan int{cap: 10}
		// In Go: we need make(chan int, 10)
		// However, we're removing make()...
		// Solution: Transform to a helper function call
		st.transformChannelLiteral(node, typ, compLit)
	}
}

// transformChannelLiteral transforms &chan T{...} to valid Go
func (st *SyntaxTransformer) transformChannelLiteral(node *ast.UnaryExpr, chanType *ast.ChanType, compLit *ast.CompositeLit) {
	// Extract capacity from composite literal
	// &chan int{cap: 10} has Elts containing KeyValueExpr with Key="cap" and Value=10
	capacity := &ast.BasicLit{
		Kind:  token.INT,
		Value: "0", // Default unbuffered
	}

	for _, elt := range compLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "cap" {
				capacity = kv.Value.(*ast.BasicLit)
				break
			}
		}
	}

	// Transform to make(chan T, capacity)
	// But wait - we're removing make()!
	// We need a runtime helper: runtime.MakeChan[T](capacity)
	// For now, generate make() call and we'll add a warning
	makeCall := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			chanType, // chan T
		},
	}

	// Add capacity if non-zero
	if capacity.Value != "0" {
		makeCall.Args = append(makeCall.Args, capacity)
	}

	// Replace the &chan T{} node with make() call
	*node = ast.UnaryExpr{
		Op: token.ILLEGAL, // Mark for removal
		X:  makeCall,
	}

	// Actually, we can't modify the parent node this way
	// We need to return the replacement node
	// This requires restructuring the visitor pattern
	// For now, add an error
	st.errors = append(st.errors, fmt.Errorf(
		"channel literal syntax &chan T{cap: N} not yet fully implemented; use manual make() for now",
	))
}

// transformToRuntimeCall transforms a built-in function call to a runtime package call
// Example: grow(s, n) -> moxie.Grow(s, n)
func (st *SyntaxTransformer) transformToRuntimeCall(node *ast.CallExpr, runtimeFunc string) {
	// Replace the function identifier with a selector expression: moxie.Function
	node.Fun = &ast.SelectorExpr{
		X: &ast.Ident{
			Name: "moxie",
		},
		Sel: &ast.Ident{
			Name: runtimeFunc,
		},
	}

	// Mark that we need to add the runtime import
	st.needsRuntimeImport = true
}

// addRuntimeImport adds the moxie runtime import to the file
func (st *SyntaxTransformer) addRuntimeImport(file *ast.File) {
	// Check if import already exists
	for _, imp := range file.Imports {
		if imp.Path != nil && imp.Path.Value == `"github.com/mleku/moxie/runtime"` {
			// Check if it has the alias "moxie"
			if imp.Name != nil && imp.Name.Name == "moxie" {
				return // Already exists with correct alias
			}
		}
	}

	// Add import: moxie "github.com/mleku/moxie/runtime"
	newImport := &ast.ImportSpec{
		Name: &ast.Ident{Name: "moxie"},
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"github.com/mleku/moxie/runtime"`,
		},
	}

	// Find or create an import declaration
	var importDecl *ast.GenDecl
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			importDecl = genDecl
			break
		}
	}

	if importDecl == nil {
		// Create new import declaration
		importDecl = &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{newImport},
		}
		// Insert at beginning of declarations
		file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
	} else {
		// Add to existing import declaration
		importDecl.Specs = append(importDecl.Specs, newImport)
	}

	// Also add to file.Imports
	file.Imports = append(file.Imports, newImport)
}

// transformAppendCall transforms append() calls to handle pointer types
// In Moxie: s = append(s, items...) where s is *[]T
// In Go: s = append(s, items...) where s is []T
// BUT we need to handle the pointer nature
// Actually, this is tricky - we can't just wrap it, we need to think about this differently
func (st *SyntaxTransformer) transformAppendCall(node *ast.CallExpr) {
	if len(node.Args) == 0 {
		return
	}

	// First argument is the slice
	sliceArg := node.Args[0]

	// Dereference the slice: append(s, ...) -> append(*s, ...)
	node.Args[0] = &ast.StarExpr{
		X: sliceArg,
	}

	// Note: The caller needs to handle re-addressing the result
	// In Moxie: s = append(s, items)  (s is *[]T)
	// In Go: *s = append(*s, items)   (s is *[]T, append returns []T)
	// But this transformation happens at the assignment level, not here
	// We just transform the call itself
}

// transformClearCall transforms clear() calls to handle pointer types
// In Moxie: clear(m) where m is *map[K]V
// In Go: clear(*m) - need to dereference
func (st *SyntaxTransformer) transformClearCall(node *ast.CallExpr) {
	if len(node.Args) != 1 {
		return
	}

	arg := node.Args[0]

	// Wrap argument in dereference: clear(m) -> clear(*m)
	// This is safe because in Moxie all maps/slices are pointers
	node.Args[0] = &ast.StarExpr{
		X: arg,
	}
}

// transformCompositeLit handles composite literal transformations
func (st *SyntaxTransformer) transformCompositeLit(node *ast.CompositeLit) {
	// Check for channel composite literals that aren't behind &
	if _, ok := node.Type.(*ast.ChanType); ok {
		st.errors = append(st.errors, fmt.Errorf(
			"chan T{} is not valid; channels must use &chan T{} syntax",
		))
	}
}
