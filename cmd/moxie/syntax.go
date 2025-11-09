// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// SyntaxTransformer handles Moxie→Go syntax transformations
type SyntaxTransformer struct {
	errors             []error
	needsRuntimeImport bool
	needsBytesImport   bool
	typeTracker        *TypeTracker
}

// NewSyntaxTransformer creates a new syntax transformer
func NewSyntaxTransformer() *SyntaxTransformer {
	return &SyntaxTransformer{
		errors:             make([]error, 0),
		needsRuntimeImport: false,
		needsBytesImport:   false,
		typeTracker:        NewTypeTracker(),
	}
}

// Transform applies all Moxie→Go syntax transformations to an AST
func (st *SyntaxTransformer) Transform(file *ast.File) error {
	st.errors = nil
	st.needsRuntimeImport = false

	// First pass: transform types, literals, comparisons
	// Multiple passes needed for chained string concatenation
	maxPasses := 10
	for pass := 0; pass < maxPasses; pass++ {
		changed := false

		// Use astutil.Apply for transformations that need node replacement
		astutil.Apply(file, func(cursor *astutil.Cursor) bool {
			node := cursor.Node()
			switch n := node.(type) {
			case *ast.GenDecl:
				// Record type information from declarations (first pass only)
				if pass == 0 {
					st.typeTracker.RecordDecl(n)
				}
			case *ast.Ident:
				// Transform string type to *[]byte (only in first pass)
				if pass == 0 {
					if replacement := st.tryTransformStringType(n, cursor); replacement != nil {
						cursor.Replace(replacement)
						changed = true
						return true
					}
					// Transform FFI constants
					if replacement := st.tryTransformFFIConstant(n); replacement != nil {
						cursor.Replace(replacement)
						changed = true
						return true
					}
					// Transform endianness constants
					if replacement := st.tryTransformEndiannessConstant(n); replacement != nil {
						cursor.Replace(replacement)
						changed = true
						return true
					}
				}
			case *ast.BasicLit:
				// Transform string literals to byte slice literals
				// NOTE: We can't replace BasicLit directly because it might be in a type-strict context
				// Instead, we'll handle this in a post-processing step or parent node transformation
				// For now, skip direct BasicLit transformation
				// TODO: Handle in BinaryExpr, AssignStmt, etc.
			case *ast.AssignStmt:
				if pass == 0 {
					// Record type information before transforming
					st.typeTracker.RecordAssign(n)
					st.transformAssignStmt(n)
					st.transformStringLiteralsInAssign(n)
				}
			case *ast.ReturnStmt:
				if pass == 0 {
					st.transformStringLiteralsInReturn(n)
				}
			case *ast.BinaryExpr:
				// Handle string concatenation and comparison
				if pass == 0 {
					st.transformBinaryExpr(n)
				}
				// Try to replace with Concat() for concatenation (all passes)
				if replacement := st.tryTransformStringConcat(n); replacement != nil {
					cursor.Replace(replacement)
					changed = true
					return true
				}
				// Try to replace with bytes.Equal/Compare for comparison (first pass only)
				if pass == 0 {
					if replacement := st.tryTransformStringComparison(n); replacement != nil {
						cursor.Replace(replacement)
						changed = true
						return true
					}
				}
			case *ast.CallExpr:
				if pass == 0 {
					st.transformCallExpr(n)
					st.transformStringLiteralsInCall(n)
					// Try to transform type coercion (e.g., (*[]uint32)(bytes))
					if replacement := st.tryTransformTypeCoercion(n); replacement != nil {
						cursor.Replace(replacement)
						changed = true
						return true
					}
				}
			case *ast.UnaryExpr:
				// Check for channel literal that needs replacement (first pass only)
				if pass == 0 {
					if replacement := st.tryTransformChannelLiteral(n); replacement != nil {
						cursor.Replace(replacement)
						changed = true
						return true
					}
				}
			case *ast.CompositeLit:
				if pass == 0 {
					st.transformCompositeLit(n)
				}
			}
			return true
		}, nil)

		// If no changes in this pass, we're done
		if !changed {
			break
		}
	}

	// Add runtime import if needed
	if st.needsRuntimeImport {
		st.addRuntimeImport(file)
	}

	// Add bytes import if needed
	if st.needsBytesImport {
		st.addBytesImport(file)
	}

	if len(st.errors) > 0 {
		return st.errors[0] // Return first error
	}
	return nil
}

// transformStringLiteralsInAssign transforms string literals in assignment RHS
func (st *SyntaxTransformer) transformStringLiteralsInAssign(node *ast.AssignStmt) {
	for i, rhs := range node.Rhs {
		if lit, ok := rhs.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if replacement := st.tryTransformStringLiteral(lit); replacement != nil {
				node.Rhs[i] = replacement
			}
		}
	}
}

// transformStringLiteralsInReturn transforms string literals in return statements
func (st *SyntaxTransformer) transformStringLiteralsInReturn(node *ast.ReturnStmt) {
	for i, result := range node.Results {
		if lit, ok := result.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if replacement := st.tryTransformStringLiteral(lit); replacement != nil {
				node.Results[i] = replacement
			}
		}
	}
}

// transformStringLiteralsInCall transforms string literals in function call arguments
func (st *SyntaxTransformer) transformStringLiteralsInCall(node *ast.CallExpr) {
	// Skip string transformations for fmt package functions
	// fmt functions expect Go strings, not *[]byte
	if st.isFmtFunction(node) {
		return
	}

	for i, arg := range node.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if replacement := st.tryTransformStringLiteral(lit); replacement != nil {
				node.Args[i] = replacement
			}
		}
	}
}

// isFmtFunction checks if a call expression is calling a fmt package function
func (st *SyntaxTransformer) isFmtFunction(node *ast.CallExpr) bool {
	if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "fmt"
		}
	}
	return false
}

// transformBinaryExpr handles binary expression transformations
// Specifically for string concatenation: s1 + s2 -> append(*s1, *s2...)
func (st *SyntaxTransformer) transformBinaryExpr(node *ast.BinaryExpr) {
	// First, transform string literals in operands
	if lit, ok := node.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if replacement := st.tryTransformStringLiteral(lit); replacement != nil {
			node.X = replacement
		}
	}
	if lit, ok := node.Y.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if replacement := st.tryTransformStringLiteral(lit); replacement != nil {
			node.Y = replacement
		}
	}
}

// tryTransformStringConcat transforms string/array concatenation
// s1 + s2 -> moxie.Concat(s1, s2) for strings (*[]byte)
// a1 + a2 -> moxie.ConcatSlice[T](a1, a2) for other slices
// Returns nil if not a slice concatenation
func (st *SyntaxTransformer) tryTransformStringConcat(node *ast.BinaryExpr) ast.Expr {
	// Only transform + operator
	if node.Op != token.ADD {
		return nil
	}

	// Check if this looks like slice concatenation
	// We assume if one operand is &[]T{...} or a variable, it might be a concat
	isSliceConcat := false

	// Check left operand
	if _, ok := node.X.(*ast.UnaryExpr); ok {
		// Likely &[]T{...}
		isSliceConcat = true
	}
	if starExpr, ok := node.X.(*ast.StarExpr); ok {
		// Could be *s (pointer dereference)
		if _, ok := starExpr.X.(*ast.Ident); ok {
			// Likely a variable
			isSliceConcat = true
		}
	}
	if _, ok := node.X.(*ast.Ident); ok {
		// Variable that might be *[]T
		isSliceConcat = true
	}
	if callExpr, ok := node.X.(*ast.CallExpr); ok {
		// Check if it's moxie.Concat or moxie.ConcatSlice (from previous transformation)
		if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "moxie" {
				if sel.Sel.Name == "Concat" || sel.Sel.Name == "ConcatSlice" {
					isSliceConcat = true
				}
			}
		}
	}

	// Check right operand too
	if _, ok := node.Y.(*ast.UnaryExpr); ok {
		isSliceConcat = true
	}
	if _, ok := node.Y.(*ast.Ident); ok {
		isSliceConcat = true
	}

	if !isSliceConcat {
		return nil
	}

	// Try to extract element type from operands
	elemType := st.extractSliceElementType(node.X)
	if elemType == nil {
		elemType = st.extractSliceElementType(node.Y)
	}

	// Prepare arguments
	var arg1, arg2 ast.Expr
	arg1 = node.X
	arg2 = node.Y

	// Mark that we need runtime import
	st.needsRuntimeImport = true

	// Check if this is a string ([]byte) concatenation
	// Priority: use Concat for byte slices (strings)
	isString := false
	if elemType != nil {
		if ident, ok := elemType.(*ast.Ident); ok && ident.Name == "byte" {
			isString = true
		}
	}

	// Check if either operand is a moxie.Concat call (string concat)
	if callExpr, ok := node.X.(*ast.CallExpr); ok {
		if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "moxie" {
				if sel.Sel.Name == "Concat" {
					isString = true
				}
			}
		}
	}

	if isString {
		// Use Concat for strings (byte slices)
		return &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "Concat"},
			},
			Args: []ast.Expr{arg1, arg2},
		}
	}

	// Use ConcatSlice[T] for other types
	// Create type parameter index expression
	var funcExpr ast.Expr
	if elemType != nil {
		// moxie.ConcatSlice[T]
		funcExpr = &ast.IndexExpr{
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "ConcatSlice"},
			},
			Index: elemType,
		}
	} else {
		// Can't determine type - just use ConcatSlice without type param
		// (will cause compile error, but user can fix)
		funcExpr = &ast.SelectorExpr{
			X:   &ast.Ident{Name: "moxie"},
			Sel: &ast.Ident{Name: "ConcatSlice"},
		}
	}

	return &ast.CallExpr{
		Fun:  funcExpr,
		Args: []ast.Expr{arg1, arg2},
	}
}

// extractSliceElementType tries to extract the element type from a slice expression
// Returns the element type AST node, or nil if it cannot be determined
func (st *SyntaxTransformer) extractSliceElementType(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		// &[]T{...} - extract T from composite literal
		if e.Op == token.AND {
			if compLit, ok := e.X.(*ast.CompositeLit); ok {
				if arrType, ok := compLit.Type.(*ast.ArrayType); ok {
					return arrType.Elt
				}
			}
		}
	case *ast.CallExpr:
		// Check for moxie.ConcatSlice[T](...) - extract T
		if indexExpr, ok := e.Fun.(*ast.IndexExpr); ok {
			return indexExpr.Index
		}
	}
	return nil
}

// tryTransformStringComparison transforms string comparison to bytes functions
// s1 == s2 -> bytes.Equal(*s1, *s2)
// s1 != s2 -> !bytes.Equal(*s1, *s2)
// s1 < s2  -> bytes.Compare(*s1, *s2) < 0
// s1 <= s2 -> bytes.Compare(*s1, *s2) <= 0
// s1 > s2  -> bytes.Compare(*s1, *s2) > 0
// s1 >= s2 -> bytes.Compare(*s1, *s2) >= 0
func (st *SyntaxTransformer) tryTransformStringComparison(node *ast.BinaryExpr) ast.Expr {
	// Only handle comparison operators
	switch node.Op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		// Potentially a string comparison
	default:
		return nil
	}

	// Check if this looks like a string/byte slice comparison
	// We assume if one operand is &[]byte{...} or *[]byte type, it's a comparison
	isStringCompare := false

	// Check left operand
	if _, ok := node.X.(*ast.UnaryExpr); ok {
		isStringCompare = true
	}
	if _, ok := node.X.(*ast.Ident); ok {
		isStringCompare = true
	}

	if !isStringCompare {
		return nil
	}

	// Mark that we need bytes import
	st.needsBytesImport = true

	// Prepare dereferenced arguments
	var arg1, arg2 ast.Expr

	// Left operand
	if unary, ok := node.X.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		// &[]byte{...} -> []byte{...} (remove &)
		arg1 = unary.X
	} else {
		// Variable -> *variable
		arg1 = &ast.StarExpr{X: node.X}
	}

	// Right operand
	if unary, ok := node.Y.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		// &[]byte{...} -> []byte{...}
		arg2 = unary.X
	} else {
		// Variable -> *variable
		arg2 = &ast.StarExpr{X: node.Y}
	}

	// Handle equality/inequality with bytes.Equal
	if node.Op == token.EQL || node.Op == token.NEQ {
		equalCall := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "bytes"},
				Sel: &ast.Ident{Name: "Equal"},
			},
			Args: []ast.Expr{arg1, arg2},
		}

		if node.Op == token.EQL {
			// s1 == s2 -> bytes.Equal(s1, s2)
			return equalCall
		} else {
			// s1 != s2 -> !bytes.Equal(s1, s2)
			return &ast.UnaryExpr{
				Op: token.NOT,
				X:  equalCall,
			}
		}
	}

	// Handle ordering comparisons with bytes.Compare
	compareCall := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "bytes"},
			Sel: &ast.Ident{Name: "Compare"},
		},
		Args: []ast.Expr{arg1, arg2},
	}

	// Return: bytes.Compare(s1, s2) <op> 0
	return &ast.BinaryExpr{
		X:  compareCall,
		Op: node.Op, // Keep the same operator
		Y: &ast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		},
	}
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
				// But only if s is NOT already dereferenced
				if len(call.Args) > 0 {
					// Check if first argument is already a dereference
					if _, alreadyDeref := call.Args[0].(*ast.StarExpr); !alreadyDeref {
						// Dereference first argument: append(s, ...) -> append(*s, ...)
						call.Args[0] = &ast.StarExpr{X: call.Args[0]}
					}
				}

				// Dereference LHS: s -> *s
				// But only if LHS is NOT already dereferenced
				if i < len(node.Lhs) {
					if _, alreadyDeref := node.Lhs[i].(*ast.StarExpr); !alreadyDeref {
						node.Lhs[i] = &ast.StarExpr{X: node.Lhs[i]}
					}
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
			// make() is generally not allowed in Moxie, except for channels
			// Check if the first argument is a channel type
			if len(node.Args) > 0 {
				if _, isChan := node.Args[0].(*ast.ChanType); !isChan {
					// Not a channel - error
					st.errors = append(st.errors, fmt.Errorf(
						"make() is not available in Moxie; use &[]T{} or &map[K]V{} instead",
					))
				}
				// Channels are allowed to use make() until parser supports &chan T{}
			}

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
			// Transform based on type:
			//   clone(slice) -> moxie.CloneSlice[T](slice)
			//   clone(map)   -> moxie.CloneMap[K, V](map)
			//   clone(struct) -> moxie.DeepCopy[T](struct)
			st.transformCloneCall(node)

		case "clear":
			// clear() exists in Go 1.21+
			// In Moxie, slices/maps are pointers, so clear(*map[K]V) needs to become clear((*map[K]V))
			// We need to dereference pointer types
			st.transformClearCall(node)

		case "free":
			// free() is a Moxie built-in that provides GC hints
			// Transform: free(v) -> moxie.Free[T](v) or moxie.FreeSlice[T](v) or moxie.FreeMap[K,V](v)
			st.transformFreeCall(node)

		case "dlopen":
			// dlopen() is a Moxie FFI built-in
			// Transform: dlopen(filename, flags) -> moxie.Dlopen(filename, flags)
			st.transformToRuntimeCall(node, "Dlopen")

		case "dlsym":
			// dlsym() is a Moxie FFI built-in
			// Transform: dlsym[T](lib, name) -> moxie.Dlsym[T](lib, name)
			st.transformToRuntimeCall(node, "Dlsym")

		case "dlclose":
			// dlclose() is a Moxie FFI built-in
			// Transform: dlclose(lib) -> moxie.Dlclose(lib)
			st.transformToRuntimeCall(node, "Dlclose")

		case "dlerror":
			// dlerror() is a Moxie FFI built-in
			// Transform: dlerror() -> moxie.Dlerror()
			st.transformToRuntimeCall(node, "Dlerror")

		case "string":
			// string() type conversion
			// Transform based on argument type:
			// string(int) -> moxie.IntToString(int)
			// string(rune) -> moxie.RuneToString(rune)
			// string(*[]rune) -> moxie.RunesToString(*[]rune)
			if len(node.Args) > 0 {
				st.transformStringConversion(node)
			}
		}
	}

	// Check for []rune(string) or *[]rune(string) conversions
	if st.isRuneSliceConversion(node) {
		st.transformRuneSliceConversion(node)
	}
}

// tryTransformChannelLiteral checks if a UnaryExpr is &__MoxieChan[T]{} and returns replacement node
// Returns nil if no transformation needed
func (st *SyntaxTransformer) tryTransformChannelLiteral(node *ast.UnaryExpr) ast.Expr {
	if node.Op != token.AND {
		return nil
	}

	// Check if this is &compositeLit
	compLit, ok := node.X.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	// Check if the composite literal type is __MoxieChan[T], __MoxieChanSend[T], or __MoxieChanRecv[T]
	// This will be an IndexExpr: __MoxieChan[T] or IndexListExpr for multiple type params
	var chanType *ast.ChanType
	var markerName string

	switch typ := compLit.Type.(type) {
	case *ast.IndexExpr:
		// __MoxieChan[T] or similar
		if ident, ok := typ.X.(*ast.Ident); ok {
			markerName = ident.Name
			switch markerName {
			case "__MoxieChan":
				// Bidirectional channel: chan T
				chanType = &ast.ChanType{
					Dir:   ast.SEND | ast.RECV,
					Value: typ.Index,
				}
			case "__MoxieChanSend":
				// Send-only channel: chan<- T
				chanType = &ast.ChanType{
					Dir:   ast.SEND,
					Value: typ.Index,
				}
			case "__MoxieChanRecv":
				// Receive-only channel: <-chan T
				chanType = &ast.ChanType{
					Dir:   ast.RECV,
					Value: typ.Index,
				}
			default:
				return nil // Not a Moxie channel marker
			}
		}
	case *ast.ChanType:
		// Old-style direct channel type (shouldn't happen with preprocessor, but handle it)
		chanType = typ
	default:
		return nil
	}

	if chanType == nil {
		return nil
	}

	// This is &__MoxieChan[T]{...} which needs transformation
	// Channel literal has one anonymous int64 field for buffer count
	// &chan int{} → make(chan int) - unbuffered
	// &chan int{10} → make(chan int, 10) - buffered with capacity 10

	// Transform to make(chan T) or make(chan T, capacity)
	makeCall := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			chanType, // chan T
		},
	}

	// Check if there's an anonymous element (buffer count)
	if len(compLit.Elts) > 0 {
		// First element is the buffer count (anonymous int64 field)
		capacity := compLit.Elts[0]

		// Only add capacity if it's not literally 0
		if basicLit, ok := capacity.(*ast.BasicLit); ok && basicLit.Value == "0" {
			// Skip - unbuffered channel
		} else {
			// Add capacity argument (could be literal or expression)
			makeCall.Args = append(makeCall.Args, capacity)
		}
	}

	// Return the make() call to replace &__MoxieChan[T]{}
	return makeCall
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

// addBytesImport adds "bytes" import to the file if not already present
func (st *SyntaxTransformer) addBytesImport(file *ast.File) {
	// Check if import already exists
	for _, imp := range file.Imports {
		if imp.Path != nil && imp.Path.Value == `"bytes"` {
			return // Already exists
		}
	}

	// Add import: "bytes"
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"bytes"`,
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
			Tok:   token.IMPORT,
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

	// Transform string literals in struct composite literal fields
	// This fixes the bug where string literals in struct fields cause type errors
	// Only transform if the composite literal type is a struct (not an array of strings)
	isStringArray := false
	if arrType, ok := node.Type.(*ast.ArrayType); ok {
		if ident, ok := arrType.Elt.(*ast.Ident); ok && ident.Name == "string" {
			isStringArray = true
		}
	}

	if !isStringArray {
		// Only transform string literals in non-string-array contexts
		for _, elt := range node.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				// Struct field: {Name: "value"}
				if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					// Transform string literal to &[]byte{...}
					if transformed := st.tryTransformStringLiteral(lit); transformed != nil {
						kv.Value = transformed
					}
				}
			}
		}
	}
}

// tryTransformStringLiteral transforms string literals to byte slice literals
// "hello" -> &[]byte{'h', 'e', 'l', 'l', 'o'}
func (st *SyntaxTransformer) tryTransformStringLiteral(lit *ast.BasicLit) ast.Expr {
	// Only transform string literals
	if lit.Kind != token.STRING {
		return nil
	}

	// Parse the string literal value
	// lit.Value includes the quotes, e.g., "hello" or `hello`
	value := lit.Value

	// Handle raw string literals (backticks)
	isRaw := value[0] == '`'

	// Remove quotes
	var str string
	if isRaw {
		str = value[1 : len(value)-1]
	} else {
		// For quoted strings, we need to handle escape sequences
		// Use strconv.Unquote to properly parse the string
		var err error
		str, err = parseStringLiteral(value)
		if err != nil {
			// If we can't parse it, leave it as-is
			return nil
		}
	}

	// Convert string to []byte literal elements
	elts := make([]ast.Expr, len(str))
	for i := 0; i < len(str); i++ {
		elts[i] = &ast.BasicLit{
			Kind:  token.CHAR,
			Value: charToLiteral(str[i]),
		}
	}

	// Create &[]byte{...} composite literal
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.CompositeLit{
			Type: &ast.ArrayType{
				Elt: &ast.Ident{Name: "byte"},
			},
			Elts: elts,
		},
	}
}

// parseStringLiteral parses a Go string literal (with quotes) and returns the actual string
func parseStringLiteral(s string) (string, error) {
	// Simple parser for common cases
	// For full correctness, we'd use strconv.Unquote, but that requires importing strconv
	// For now, handle the common escape sequences manually

	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return "", fmt.Errorf("invalid string literal")
	}

	s = s[1 : len(s)-1] // Remove quotes
	result := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			// Escape sequence
			switch s[i+1] {
			case 'n':
				result = append(result, '\n')
				i++
			case 't':
				result = append(result, '\t')
				i++
			case 'r':
				result = append(result, '\r')
				i++
			case '\\':
				result = append(result, '\\')
				i++
			case '"':
				result = append(result, '"')
				i++
			case '\'':
				result = append(result, '\'')
				i++
			default:
				// Unknown escape, keep as-is
				result = append(result, s[i])
			}
		} else {
			result = append(result, s[i])
		}
	}

	return string(result), nil
}

// charToLiteral converts a byte to a character literal string
func charToLiteral(b byte) string {
	switch b {
	case '\n':
		return `'\n'`
	case '\t':
		return `'\t'`
	case '\r':
		return `'\r'`
	case '\\':
		return `'\\'`
	case '\'':
		return `'\''`
	case '"':
		return `'"'`
	default:
		if b >= 32 && b <= 126 {
			// Printable ASCII
			return fmt.Sprintf("'%c'", b)
		}
		// Non-printable, use hex
		return fmt.Sprintf("'\\x%02x'", b)
	}
}

// tryTransformStringType transforms 'string' type identifier to '*[]byte'
// Returns nil if not a string type identifier or if in wrong context
func (st *SyntaxTransformer) tryTransformStringType(ident *ast.Ident, cursor *astutil.Cursor) ast.Expr {
	// Only transform if this is the identifier "string"
	if ident.Name != "string" {
		return nil
	}

	// Check the parent node to determine if this is a type context
	parent := cursor.Parent()

	// We want to transform string in type contexts:
	// - Field types in structs
	// - Function parameter types
	// - Function return types
	// - Variable declaration types
	// - Type assertions

	// But NOT transform string in:
	// - Package names
	// - Variable names
	// - Function names
	// - String literal values

	switch p := parent.(type) {
	case *ast.Field:
		// This is a parameter or struct field type
		// string -> *[]byte
		return &ast.StarExpr{
			X: &ast.ArrayType{
				Elt: &ast.Ident{Name: "byte"},
			},
		}
	case *ast.ValueSpec:
		// Variable declaration type
		return &ast.StarExpr{
			X: &ast.ArrayType{
				Elt: &ast.Ident{Name: "byte"},
			},
		}
	case *ast.TypeAssertExpr:
		// Type assertion
		if p.Type == ident {
			return &ast.StarExpr{
				X: &ast.ArrayType{
					Elt: &ast.Ident{Name: "byte"},
				},
			}
		}
	}

	// Don't transform in other contexts
	return nil
}

// tryTransformFFIConstant transforms FFI constants to moxie.CONSTANT
func (st *SyntaxTransformer) tryTransformFFIConstant(ident *ast.Ident) ast.Expr {
	switch ident.Name {
	case "RTLD_LAZY", "RTLD_NOW", "RTLD_GLOBAL", "RTLD_LOCAL":
		st.needsRuntimeImport = true
		return &ast.SelectorExpr{
			X:   &ast.Ident{Name: "moxie"},
			Sel: &ast.Ident{Name: ident.Name},
		}
	}
	return nil
}

// tryTransformEndiannessConstant transforms endianness constants to moxie.CONSTANT
func (st *SyntaxTransformer) tryTransformEndiannessConstant(ident *ast.Ident) ast.Expr {
	switch ident.Name {
	case "NativeEndian", "LittleEndian", "BigEndian":
		st.needsRuntimeImport = true
		return &ast.SelectorExpr{
			X:   &ast.Ident{Name: "moxie"},
			Sel: &ast.Ident{Name: ident.Name},
		}
	}
	return nil
}

// tryTransformTypeCoercion transforms type coercion expressions
// Detects: (*[]TargetType)(sourceSlice) or (*[]TargetType, Endian)(sourceSlice)
// Transforms to: moxie.Coerce[SourceType, TargetType](sourceSlice) or
//                moxie.Coerce[SourceType, TargetType](sourceSlice, moxie.Endian)
func (st *SyntaxTransformer) tryTransformTypeCoercion(call *ast.CallExpr) ast.Expr {
	// Check if this is a type conversion/cast (Fun is a type expression)
	// Pattern: (Type)(expr) where Type is (*[]T)

	// The Fun should be a ParenExpr containing a StarExpr containing an ArrayType
	// Or it could be a direct StarExpr (Go parser handles this differently)

	var targetType ast.Expr
	var endianExpr ast.Expr

	// Check for pattern: (*[]T, Endian)(expr) - this would be parsed differently
	// For now, handle the simpler case: (*[]T)(expr)

	// Try to extract the star expression
	switch fun := call.Fun.(type) {
	case *ast.ParenExpr:
		targetType = fun.X
	case *ast.StarExpr:
		targetType = fun
	default:
		return nil
	}

	// Now targetType should be *ast.StarExpr with X being *ast.ArrayType
	starExpr, ok := targetType.(*ast.StarExpr)
	if !ok {
		return nil
	}

	// Check if X is an array type (slice)
	arrayType, ok := starExpr.X.(*ast.ArrayType)
	if !ok || arrayType.Len != nil {
		// Not a slice type, or it's an array with fixed length
		return nil
	}

	// We have a slice coercion: (*[]T)(expr)
	// Need to extract source type from expr
	if len(call.Args) == 0 {
		return nil
	}

	sourceExpr := call.Args[0]

	// Try to determine source element type
	// For now, we'll use a generic "byte" as source if we can't determine it
	// TODO: Implement proper type inference
	sourceElemType := &ast.Ident{Name: "byte"}
	targetElemType := arrayType.Elt

	// Check if there's an endianness argument (would be the second arg)
	if len(call.Args) > 1 {
		endianExpr = call.Args[1]
	}

	// Build the Coerce call: moxie.Coerce[SourceType, TargetType](expr, [endian])
	// For generics with multiple type parameters, use IndexListExpr
	coerceCall := &ast.CallExpr{
		Fun: &ast.IndexListExpr{
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "Coerce"},
			},
			Indices: []ast.Expr{
				sourceElemType,
				targetElemType,
			},
		},
		Args: []ast.Expr{sourceExpr},
	}

	// Add endianness argument if present
	if endianExpr != nil {
		coerceCall.Args = append(coerceCall.Args, endianExpr)
	}

	st.needsRuntimeImport = true
	return coerceCall
}

// transformStringConversion transforms string(x) to appropriate runtime function
// based on the type of x
func (st *SyntaxTransformer) transformStringConversion(node *ast.CallExpr) {
	if len(node.Args) == 0 {
		return
	}

	arg := node.Args[0]

	// Try to infer the type from the expression
	// This is a best-effort heuristic without full type checking

	// Check for numeric literals
	if lit, ok := arg.(*ast.BasicLit); ok {
		switch lit.Kind {
		case token.INT:
			// string(42) -> moxie.IntToString(42)
			st.transformToRuntimeCallWithArgs(node, "IntToString", node.Args)
			return
		case token.CHAR:
			// string('A') -> moxie.RuneToString('A')
			st.transformToRuntimeCallWithArgs(node, "RuneToString", node.Args)
			return
		}
	}

	// Check for star expressions (pointer dereference or address)
	if _, ok := arg.(*ast.StarExpr); ok {
		// Likely *[]rune -> string
		// string(*runeSlice) -> moxie.RunesToString(runeSlice)
		st.transformToRuntimeCallWithArgs(node, "RunesToString", node.Args)
		return
	}

	// Check for identifiers - try to guess based on name patterns
	if ident, ok := arg.(*ast.Ident); ok {
		name := ident.Name
		// Common rune/int32 variable names
		if name == "r" || name == "ch" || name == "c" {
			st.transformToRuntimeCallWithArgs(node, "RuneToString", node.Args)
			return
		}
		// If the variable ends with "Runes" or "runes", it's likely a rune slice
		if len(name) > 5 && (name[len(name)-5:] == "Runes" || name[len(name)-5:] == "runes") {
			st.transformToRuntimeCallWithArgs(node, "RunesToString", node.Args)
			return
		}
	}

	// Default: assume int conversion (most common case)
	// string(n) -> moxie.IntToString(n)
	st.transformToRuntimeCallWithArgs(node, "IntToString", node.Args)
}

// isRuneSliceConversion checks if a CallExpr is []rune(x) or *[]rune(x)
func (st *SyntaxTransformer) isRuneSliceConversion(node *ast.CallExpr) bool {
	// Check for []rune(x)
	if arrType, ok := node.Fun.(*ast.ArrayType); ok {
		if ident, ok := arrType.Elt.(*ast.Ident); ok && ident.Name == "rune" {
			return true
		}
	}

	// Check for *[]rune(x)
	if starExpr, ok := node.Fun.(*ast.StarExpr); ok {
		if arrType, ok := starExpr.X.(*ast.ArrayType); ok {
			if ident, ok := arrType.Elt.(*ast.Ident); ok && ident.Name == "rune" {
				return true
			}
		}
	}

	return false
}

// transformRuneSliceConversion transforms []rune(s) or *[]rune(s) to moxie.StringToRunes(s)
func (st *SyntaxTransformer) transformRuneSliceConversion(node *ast.CallExpr) {
	if len(node.Args) == 0 {
		return
	}

	// Transform: []rune(s) -> moxie.StringToRunes(s)
	// Transform: *[]rune(s) -> moxie.StringToRunes(s)
	st.transformToRuntimeCallWithArgs(node, "StringToRunes", node.Args)
}

// transformToRuntimeCallWithArgs is a helper that transforms a CallExpr to call a moxie runtime function
// with explicit arguments (used for conversions)
func (st *SyntaxTransformer) transformToRuntimeCallWithArgs(node *ast.CallExpr, runtimeFunc string, args []ast.Expr) {
	node.Fun = &ast.SelectorExpr{
		X:   &ast.Ident{Name: "moxie"},
		Sel: &ast.Ident{Name: runtimeFunc},
	}
	node.Args = args
	st.needsRuntimeImport = true
}

// transformCloneCall transforms clone() calls based on argument type
func (st *SyntaxTransformer) transformCloneCall(node *ast.CallExpr) {
	if len(node.Args) == 0 {
		st.errors = append(st.errors, fmt.Errorf("clone() requires an argument"))
		return
	}

	arg := node.Args[0]
	var argType ast.Expr

	// Try to infer type from the argument
	switch a := arg.(type) {
	case *ast.Ident:
		// Look up identifier in type tracker
		argType = st.typeTracker.GetType(a.Name)
	case *ast.UnaryExpr:
		// &T{...}
		if a.Op == token.AND {
			if compLit, ok := a.X.(*ast.CompositeLit); ok {
				argType = compLit.Type
			}
		}
	default:
		// Try to infer from the expression
		argType = st.typeTracker.inferTypeFromExpr(arg)
	}

	// If we couldn't determine type, default to DeepCopy (safest option)
	if argType == nil {
		st.transformToRuntimeCall(node, "DeepCopy")
		return
	}

	// Determine which clone function to use based on type
	if st.typeTracker.IsSliceType(argType) {
		// Transform to CloneSlice[T]
		elemType := st.typeTracker.extractElementType(argType)
		if elemType != nil {
			// moxie.CloneSlice[T](arg)
			node.Fun = &ast.IndexExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "moxie"},
					Sel: &ast.Ident{Name: "CloneSlice"},
				},
				Index: elemType,
			}
		} else {
			// Can't determine element type, use DeepCopy
			st.transformToRuntimeCall(node, "DeepCopy")
		}
	} else if st.typeTracker.IsMapType(argType) {
		// Transform to CloneMap[K, V]
		keyType, valueType := st.typeTracker.GetMapKeyValueTypes(argType)
		if keyType != nil && valueType != nil {
			// moxie.CloneMap[K, V](arg)
			node.Fun = &ast.IndexListExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "moxie"},
					Sel: &ast.Ident{Name: "CloneMap"},
				},
				Indices: []ast.Expr{keyType, valueType},
			}
		} else {
			// Can't determine key/value types, use DeepCopy
			st.transformToRuntimeCall(node, "DeepCopy")
		}
	} else {
		// For structs and other types, use DeepCopy
		// Extract the base type (remove pointer if present)
		baseType := argType
		if starExpr, ok := argType.(*ast.StarExpr); ok {
			baseType = starExpr.X
		}

		// moxie.DeepCopy[T](arg)
		node.Fun = &ast.IndexExpr{
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "DeepCopy"},
			},
			Index: baseType,
		}
	}

	st.needsRuntimeImport = true
}

// transformFreeCall transforms free() calls based on argument type
func (st *SyntaxTransformer) transformFreeCall(node *ast.CallExpr) {
	if len(node.Args) == 0 {
		st.errors = append(st.errors, fmt.Errorf("free() requires an argument"))
		return
	}

	arg := node.Args[0]
	var argType ast.Expr

	// Try to infer type from the argument
	switch a := arg.(type) {
	case *ast.Ident:
		// Look up identifier in type tracker
		argType = st.typeTracker.GetType(a.Name)
		if argType == nil {
			// Can't find type, use generic Free without type params
			node.Fun = &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "Free"},
			}
			st.needsRuntimeImport = true
			return
		}
	case *ast.UnaryExpr:
		// &T{...}
		if a.Op == token.AND {
			if compLit, ok := a.X.(*ast.CompositeLit); ok {
				argType = compLit.Type
			}
		}
	default:
		// Try to infer from the expression
		argType = st.typeTracker.inferTypeFromExpr(arg)
	}

	// If we couldn't determine type, use generic Free
	if argType == nil {
		node.Fun = &ast.SelectorExpr{
			X:   &ast.Ident{Name: "moxie"},
			Sel: &ast.Ident{Name: "Free"},
		}
		st.needsRuntimeImport = true
		return
	}

	// Determine which free function to use based on type
	if st.typeTracker.IsSliceType(argType) {
		// Transform to FreeSlice[T]
		elemType := st.typeTracker.extractElementType(argType)
		if elemType != nil {
			// moxie.FreeSlice[T](arg)
			node.Fun = &ast.IndexExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "moxie"},
					Sel: &ast.Ident{Name: "FreeSlice"},
				},
				Index: elemType,
			}
		} else {
			// Can't determine element type, use generic Free
			node.Fun = &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "Free"},
			}
		}
	} else if st.typeTracker.IsMapType(argType) {
		// Transform to FreeMap[K, V]
		keyType, valueType := st.typeTracker.GetMapKeyValueTypes(argType)
		if keyType != nil && valueType != nil {
			// moxie.FreeMap[K, V](arg)
			node.Fun = &ast.IndexListExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "moxie"},
					Sel: &ast.Ident{Name: "FreeMap"},
				},
				Indices: []ast.Expr{keyType, valueType},
			}
		} else {
			// Can't determine key/value types, use generic Free
			node.Fun = &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "Free"},
			}
		}
	} else {
		// For structs and other types, use Free[T]
		// Extract the base type (remove pointer if present)
		baseType := argType
		if starExpr, ok := argType.(*ast.StarExpr); ok {
			baseType = starExpr.X
		}

		// moxie.Free[T](arg)
		node.Fun = &ast.IndexExpr{
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "moxie"},
				Sel: &ast.Ident{Name: "Free"},
			},
			Index: baseType,
		}
	}

	st.needsRuntimeImport = true
}
