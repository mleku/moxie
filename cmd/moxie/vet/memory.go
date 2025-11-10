// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vet

import (
	"go/ast"
	"go/token"
)

// MemoryCheck performs memory management analysis
type MemoryCheck struct{}

// Name returns the check name
func (m *MemoryCheck) Name() string {
	return "memory"
}

// Category returns the check category
func (m *MemoryCheck) Category() string {
	return "memory"
}

// Analyze performs memory safety checks
func (m *MemoryCheck) Analyze(fset *token.FileSet, file *ast.File) []*Issue {
	issues := make([]*Issue, 0)

	// Track clone() calls and their usage
	cloneCalls := make(map[*ast.CallExpr]bool)
	freeCalls := make(map[string]int) // variable name -> free count

	// First pass: find all clone() and free() calls
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			// Check for clone() calls
			if ident, ok := node.Fun.(*ast.Ident); ok && ident.Name == "clone" {
				cloneCalls[node] = false // Mark as not used yet
			}

			// Check for free() calls
			if ident, ok := node.Fun.(*ast.Ident); ok && ident.Name == "free" {
				if len(node.Args) > 0 {
					if argIdent, ok := node.Args[0].(*ast.Ident); ok {
						freeCalls[argIdent.Name]++

						// Check for double free
						if freeCalls[argIdent.Name] > 1 {
							pos := fset.Position(node.Pos())
							issues = append(issues, &Issue{
								File:     pos.Filename,
								Line:     pos.Line,
								Column:   pos.Column,
								Severity: SeverityError,
								Category: "memory",
								Check:    "double_free",
								Message:  "multiple free() calls on same variable",
								Help:     "variable '" + argIdent.Name + "' is freed " +
									"more than once; set to nil after first free()",
							})
						}
					}
				}
			}

		case *ast.AssignStmt:
			// Check if clone() result is assigned
			for i, rhs := range node.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok {
					if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "clone" {
						if i < len(node.Lhs) {
							// Mark as used
							cloneCalls[call] = true

							// Track the variable for free() check
							if lhsIdent, ok := node.Lhs[i].(*ast.Ident); ok {
								// Check if there's a corresponding free() in the function
								if freeCalls[lhsIdent.Name] == 0 {
									// This will be checked in second pass
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	// Check for unused clone() calls
	for call, used := range cloneCalls {
		if !used {
			pos := fset.Position(call.Pos())
			issues = append(issues, &Issue{
				File:     pos.Filename,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: SeverityWarning,
				Category: "memory",
				Check:    "unused_clone",
				Message:  "clone() result is not used",
				Help:     "assign result to a variable or remove the call",
			})
		}
	}

	// Second pass: check for missing free() calls
	// This is a simplified check - a full implementation would do control flow analysis
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Body != nil {
				allocations := make(map[string]bool)

				// Find allocations (clone calls assigned to variables)
				ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
					if assign, ok := n.(*ast.AssignStmt); ok {
						for i, rhs := range assign.Rhs {
							if call, ok := rhs.(*ast.CallExpr); ok {
								if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "clone" {
									if i < len(assign.Lhs) {
										if lhsIdent, ok := assign.Lhs[i].(*ast.Ident); ok {
											allocations[lhsIdent.Name] = true
										}
									}
								}
							}
						}
					}
					return true
				})

				// Check if each allocation has a free()
				for varName := range allocations {
					if freeCalls[varName] == 0 {
						// Find the allocation position for error reporting
						var allocPos token.Position
						ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
							if assign, ok := n.(*ast.AssignStmt); ok {
								for _, lhs := range assign.Lhs {
									if ident, ok := lhs.(*ast.Ident); ok && ident.Name == varName {
										allocPos = fset.Position(assign.Pos())
										return false
									}
								}
							}
							return true
						})

						issues = append(issues, &Issue{
							File:     allocPos.Filename,
							Line:     allocPos.Line,
							Column:   allocPos.Column,
							Severity: SeverityWarning,
							Category: "memory",
							Check:    "missing_free",
							Message:  "potential memory leak: allocated variable '" + varName + "' is not freed",
							Help:     "add 'defer free(" + varName + ")' after allocation",
						})
					}
				}
			}
		}
		return true
	})

	return issues
}
