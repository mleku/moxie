// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ConstChecker tracks const declarations and enforces compile-time immutability
type ConstChecker struct {
	// Map of const identifier names to their declaration positions
	constDecls map[string]token.Pos
	errors     []error
}

// NewConstChecker creates a new const checker
func NewConstChecker() *ConstChecker {
	return &ConstChecker{
		constDecls: make(map[string]token.Pos),
		errors:     make([]error, 0),
	}
}

// Check performs const enforcement checking on an AST
func (cc *ConstChecker) Check(file *ast.File) []error {
	cc.errors = nil

	// First pass: collect all const declarations
	ast.Inspect(file, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.CONST {
				cc.collectConstDecls(decl)
			}
		}
		return true
	})

	// Second pass: check for mutations
	ast.Inspect(file, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			cc.checkAssignment(stmt)
		case *ast.IncDecStmt:
			cc.checkIncDec(stmt)
		}
		return true
	})

	return cc.errors
}

// collectConstDecls collects all const declarations
func (cc *ConstChecker) collectConstDecls(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				cc.constDecls[name.Name] = name.Pos()
			}
		}
	}
}

// checkAssignment checks if an assignment mutates a const
func (cc *ConstChecker) checkAssignment(stmt *ast.AssignStmt) {
	for _, lhs := range stmt.Lhs {
		if ident := cc.extractIdentifier(lhs); ident != nil {
			if pos, isConst := cc.constDecls[ident.Name]; isConst {
				cc.errors = append(cc.errors, fmt.Errorf(
					"cannot assign to const %s (declared at %v)",
					ident.Name,
					pos,
				))
			}
		}
	}
}

// checkIncDec checks if increment/decrement mutates a const
func (cc *ConstChecker) checkIncDec(stmt *ast.IncDecStmt) {
	if ident := cc.extractIdentifier(stmt.X); ident != nil {
		if pos, isConst := cc.constDecls[ident.Name]; isConst {
			op := "++"
			if stmt.Tok == token.DEC {
				op = "--"
			}
			cc.errors = append(cc.errors, fmt.Errorf(
				"cannot %s const %s (declared at %v)",
				op,
				ident.Name,
				pos,
			))
		}
	}
}

// extractIdentifier extracts the identifier from an expression
// Handles cases like:
//   - ident (simple identifier)
//   - *ident (dereference)
//   - ident.field (selector - checks base ident)
//   - ident[index] (index - checks base ident)
func (cc *ConstChecker) extractIdentifier(expr ast.Expr) *ast.Ident {
	switch e := expr.(type) {
	case *ast.Ident:
		return e
	case *ast.StarExpr:
		// For *ident, check if ident itself is const
		return cc.extractIdentifier(e.X)
	case *ast.SelectorExpr:
		// For ident.field, check if ident is const
		// Note: This doesn't prevent mutation of fields of a const struct pointer
		// For full protection, we'd need type information
		return cc.extractIdentifier(e.X)
	case *ast.IndexExpr:
		// For ident[index], check if ident is const
		return cc.extractIdentifier(e.X)
	case *ast.ParenExpr:
		return cc.extractIdentifier(e.X)
	}
	return nil
}

// IsConst checks if an identifier is declared as const
func (cc *ConstChecker) IsConst(name string) bool {
	_, isConst := cc.constDecls[name]
	return isConst
}
