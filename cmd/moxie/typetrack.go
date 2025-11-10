// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/token"
)

// TypeTracker tracks variable types throughout AST traversal
type TypeTracker struct {
	types     map[string]ast.Expr      // Maps variable names to their type expressions
	functions map[string]*ast.FuncType // Maps function names to their signatures
}

// NewTypeTracker creates a new type tracker
func NewTypeTracker() *TypeTracker {
	return &TypeTracker{
		types:     make(map[string]ast.Expr),
		functions: make(map[string]*ast.FuncType),
	}
}

// RecordDecl records type information from a declaration
func (tt *TypeTracker) RecordDecl(decl *ast.GenDecl) {
	if decl.Tok != token.VAR && decl.Tok != token.CONST {
		return
	}

	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Record type for each name in the spec
		for _, name := range valueSpec.Names {
			if name.Name == "_" {
				continue // Skip blank identifier
			}

			// Use explicit type if available
			if valueSpec.Type != nil {
				tt.types[name.Name] = valueSpec.Type
			} else if len(valueSpec.Values) > 0 {
				// Try to infer type from value
				if inferredType := tt.inferTypeFromExpr(valueSpec.Values[0]); inferredType != nil {
					tt.types[name.Name] = inferredType
				}
			}
		}
	}
}

// RecordFunc records function signature information
func (tt *TypeTracker) RecordFunc(funcDecl *ast.FuncDecl) {
	if funcDecl == nil || funcDecl.Name == nil {
		return
	}

	// Store the function type (signature)
	tt.functions[funcDecl.Name.Name] = funcDecl.Type

	// Also record parameter types in the types map
	// This helps when we're inside a function and need to look up parameter types
	if funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			if field.Type != nil {
				for _, name := range field.Names {
					if name.Name != "_" {
						tt.types[name.Name] = field.Type
					}
				}
			}
		}
	}
}

// RecordAssign records type information from an assignment
func (tt *TypeTracker) RecordAssign(assign *ast.AssignStmt) {
	// Only handle short variable declarations (:=) and regular assignments
	if assign.Tok != token.DEFINE && assign.Tok != token.ASSIGN {
		return
	}

	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}

		// For :=, always record the type
		// For =, only record if we don't already have type info
		if assign.Tok == token.DEFINE || tt.types[ident.Name] == nil {
			if i < len(assign.Rhs) {
				if inferredType := tt.inferTypeFromExpr(assign.Rhs[i]); inferredType != nil {
					tt.types[ident.Name] = inferredType
				}
			}
		}
	}
}

// GetType returns the type expression for a given identifier
func (tt *TypeTracker) GetType(name string) ast.Expr {
	return tt.types[name]
}

// inferTypeFromExpr attempts to infer the type from an expression
func (tt *TypeTracker) inferTypeFromExpr(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.CompositeLit:
		// Composite literal has explicit type
		return e.Type

	case *ast.UnaryExpr:
		if e.Op == token.AND {
			// &T{...} → type is T
			if compLit, ok := e.X.(*ast.CompositeLit); ok {
				return compLit.Type
			}
		}

	case *ast.CallExpr:
		// Check for Moxie built-in functions (before transformation)
		if ident, ok := e.Fun.(*ast.Ident); ok {
			switch ident.Name {
			case "make":
				if len(e.Args) > 0 {
					return e.Args[0] // make(T, ...) → type is T
				}
			case "clone", "grow":
				// clone(x) and grow(x, n) return the same type as their first argument
				if len(e.Args) > 0 {
					if argType := tt.inferTypeFromExpr(e.Args[0]); argType != nil {
						return argType
					}
				}
			default:
				// Check if this is a user-defined function call
				// Look up the function signature and return the first return type
				if funcType, ok := tt.functions[ident.Name]; ok {
					if funcType.Results != nil && len(funcType.Results.List) > 0 {
						return funcType.Results.List[0].Type
					}
				}
			}
		}

		// Check for moxie runtime functions that return typed values (after transformation)
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "moxie" {
				// For functions like Concat, ConcatSlice, etc., try to infer from arguments
				switch sel.Sel.Name {
				case "Concat":
					// Returns *[]byte (Moxie string)
					return &ast.StarExpr{
						X: &ast.ArrayType{
							Elt: &ast.Ident{Name: "byte"},
						},
					}
				case "ConcatSlice":
					// Returns *[]T - extract type parameter
					if indexExpr, ok := e.Fun.(*ast.SelectorExpr); ok {
						if idx, ok := indexExpr.X.(*ast.IndexExpr); ok {
							return &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: idx.Index,
								},
							}
						}
					}
				case "CloneSlice", "Grow":
					// These return the same type as their input
					if len(e.Args) > 0 {
						if argType := tt.inferTypeFromExpr(e.Args[0]); argType != nil {
							return argType
						}
					}
				}
			}
		}

		// Check for type conversions
		if len(e.Args) == 1 {
			// T(x) might be a type conversion
			// The Fun itself might be a type
			return e.Fun
		}

	case *ast.Ident:
		// Look up identifier in our type map
		return tt.types[e.Name]

	case *ast.IndexExpr:
		// Array/slice/map indexing - need to extract element type
		if baseType := tt.inferTypeFromExpr(e.X); baseType != nil {
			return tt.extractElementType(baseType)
		}

	case *ast.StarExpr:
		// Pointer dereference - return the dereferenced type
		if innerType := tt.inferTypeFromExpr(e.X); innerType != nil {
			// If innerType is *T, return T
			if starType, ok := innerType.(*ast.StarExpr); ok {
				return starType.X
			}
		}
	}

	return nil
}

// extractElementType extracts the element type from a container type
func (tt *TypeTracker) extractElementType(typeExpr ast.Expr) ast.Expr {
	switch t := typeExpr.(type) {
	case *ast.ArrayType:
		return t.Elt

	case *ast.StarExpr:
		// *[]T or *map[K]V
		if arrType, ok := t.X.(*ast.ArrayType); ok {
			return arrType.Elt
		}
		if mapType, ok := t.X.(*ast.MapType); ok {
			return mapType.Value
		}

	case *ast.MapType:
		return t.Value
	}

	return nil
}

// IsSliceType checks if a type expression represents a slice
func (tt *TypeTracker) IsSliceType(typeExpr ast.Expr) bool {
	switch t := typeExpr.(type) {
	case *ast.ArrayType:
		return t.Len == nil // Slice has no length

	case *ast.StarExpr:
		// *[]T
		if arrType, ok := t.X.(*ast.ArrayType); ok {
			return arrType.Len == nil
		}
	}

	return false
}

// IsMapType checks if a type expression represents a map
func (tt *TypeTracker) IsMapType(typeExpr ast.Expr) bool {
	switch t := typeExpr.(type) {
	case *ast.MapType:
		return true

	case *ast.StarExpr:
		// *map[K]V
		_, ok := t.X.(*ast.MapType)
		return ok
	}

	return false
}

// IsStructType checks if a type expression represents a struct
func (tt *TypeTracker) IsStructType(typeExpr ast.Expr) bool {
	switch t := typeExpr.(type) {
	case *ast.StructType:
		return true

	case *ast.StarExpr:
		// *struct{...}
		_, ok := t.X.(*ast.StructType)
		return ok

	case *ast.Ident:
		// Could be a named struct type, but we can't determine without full type info
		// For now, assume it might be a struct
		return true

	case *ast.SelectorExpr:
		// pkg.Type - could be a struct
		return true
	}

	return false
}

// GetMapKeyValueTypes extracts key and value types from a map type
func (tt *TypeTracker) GetMapKeyValueTypes(typeExpr ast.Expr) (key, value ast.Expr) {
	switch t := typeExpr.(type) {
	case *ast.MapType:
		return t.Key, t.Value

	case *ast.StarExpr:
		// *map[K]V
		if mapType, ok := t.X.(*ast.MapType); ok {
			return mapType.Key, mapType.Value
		}
	}

	return nil, nil
}
