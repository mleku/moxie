// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
)

// ConstChecker checks for const immutability violations
type ConstChecker struct {
	consts *map[string]bool
}

// NewConstChecker creates a new const checker
func NewConstChecker() *ConstChecker {
	return &ConstChecker{
		consts: &map[string]bool{},
	}
}

// Check checks a file for const violations
func (cc *ConstChecker) Check(file *ast.File) *[]error {
	errors := &[]error{}

	// First pass: collect all const declarations
	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			if genDecl.Tok.String() == "const" {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							(*cc.consts)[name.Name] = true
						}
					}
				}
			}
		}
		return true
	})

	// Second pass: check for assignments to consts
	ast.Inspect(file, func(n ast.Node) bool {
		// Check assignments
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			for _, lhs := range assignStmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if (*cc.consts)[ident.Name] {
						*errors = append(*errors, fmt.Errorf("cannot assign to const '%s'", ident.Name))
					}
				}
			}
		}

		// Check increment/decrement
		if incDecStmt, ok := n.(*ast.IncDecStmt); ok {
			if ident, ok := incDecStmt.X.(*ast.Ident); ok {
				if (*cc.consts)[ident.Name] {
					*errors = append(*errors, fmt.Errorf("cannot modify const '%s'", ident.Name))
				}
			}
		}

		return true
	})

	return errors
}
