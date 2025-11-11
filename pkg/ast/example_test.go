package ast_test

import (
	"fmt"
	"github.com/mleku/moxie/pkg/ast"
)

// Example demonstrates building a simple function AST.
func Example_buildAST() {
	// Build AST for: func add(a, b int) int { return a + b }
	funcDecl := &ast.FuncDecl{
		Name: &ast.Ident{
			NamePos: ast.Position{Line: 1, Column: 6},
			Name:    "add",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "a"},
							{Name: "b"},
						},
						Type: &ast.Ident{Name: "int"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "int"}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BinaryExpr{
							X:  &ast.Ident{Name: "a"},
							Op: ast.ADD,
							Y:  &ast.Ident{Name: "b"},
						},
					},
				},
			},
		},
	}

	fmt.Printf("Function: %s\n", funcDecl.Name.Name)
	fmt.Printf("Is method: %v\n", funcDecl.IsMethod())

	// Output:
	// Function: add
	// Is method: false
}

// Example demonstrates Moxie-specific slice literal.
func Example_moxieSliceLit() {
	// Build AST for: &[]int{1, 2, 3}
	sliceLit := &ast.SliceLit{
		Type: &ast.Ident{Name: "int"},
		Elts: []ast.Expr{
			&ast.BasicLit{Kind: ast.IntLit, Value: "1"},
			&ast.BasicLit{Kind: ast.IntLit, Value: "2"},
			&ast.BasicLit{Kind: ast.IntLit, Value: "3"},
		},
	}

	fmt.Printf("Slice element type: %s\n", sliceLit.Type.(*ast.Ident).Name)
	fmt.Printf("Number of elements: %d\n", len(sliceLit.Elts))

	// Output:
	// Slice element type: int
	// Number of elements: 3
}

// Example demonstrates Moxie FFI call.
func Example_moxieFFI() {
	// Build AST for: dlsym[func(*byte) int64](lib, "strlen")
	ffiCall := &ast.FFICall{
		Name: &ast.Ident{Name: "dlsym"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.PointerType{
						Base: &ast.Ident{Name: "byte"},
					}},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "int64"}},
				},
			},
		},
		Args: []ast.Expr{
			&ast.Ident{Name: "lib"},
			&ast.BasicLit{Kind: ast.StringLit, Value: `"strlen"`},
		},
	}

	fmt.Printf("FFI function: %s\n", ffiCall.Name.Name)
	fmt.Printf("Arguments: %d\n", len(ffiCall.Args))

	// Output:
	// FFI function: dlsym
	// Arguments: 2
}

// Example demonstrates token operations.
func Example_tokens() {
	// Token information
	fmt.Printf("ADD token: %s\n", ast.ADD)
	fmt.Printf("Is operator: %v\n", ast.ADD.IsOperator())
	fmt.Printf("Precedence: %d\n", ast.ADD.Precedence())

	fmt.Printf("FUNC token: %s\n", ast.FUNC)
	fmt.Printf("Is keyword: %v\n", ast.FUNC.IsKeyword())

	fmt.Printf("CLONE token: %s\n", ast.CLONE)
	fmt.Printf("Is keyword: %v\n", ast.CLONE.IsKeyword())

	// Output:
	// ADD token: +
	// Is operator: true
	// Precedence: 4
	// FUNC token: func
	// Is keyword: true
	// CLONE token: clone
	// Is keyword: true
}
