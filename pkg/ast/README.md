# Moxie Abstract Syntax Tree (AST)

This package defines the abstract syntax tree for the Moxie programming language.

## Overview

The AST provides a structured representation of Moxie source code that is:
- **Type-safe**: All nodes implement specific interfaces
- **Position-aware**: Every node tracks its source location
- **Go-compatible**: Based on Go's AST design but extended for Moxie features
- **Ready for analysis**: Supports semantic analysis and code generation

## Package Structure

### Core Files

- **ast.go** - Base interfaces and position tracking
- **types.go** - Type nodes and type expressions
- **decls.go** - Declarations (package, import, const, var, type, func)
- **stmts.go** - Statements (if, for, switch, return, etc.)
- **exprs.go** - Expressions (binary, unary, call, index, etc.)
- **literals.go** - Literals and tokens

## Node Hierarchy

```
Node (base interface)
├── Expr (expressions)
│   ├── Type (type expressions)
│   │   ├── Ident
│   │   ├── BasicType
│   │   ├── PointerType
│   │   ├── SliceType (Moxie: *[]T)
│   │   ├── ArrayType
│   │   ├── MapType (Moxie: *map[K]V)
│   │   ├── ChanType (Moxie: *chan T)
│   │   ├── StructType
│   │   ├── InterfaceType
│   │   ├── FuncType (with generics support)
│   │   └── ParenType
│   ├── BasicLit (literals)
│   ├── CompositeLit
│   ├── FuncLit
│   ├── UnaryExpr
│   ├── BinaryExpr
│   ├── CallExpr
│   ├── IndexExpr
│   ├── SliceExpr
│   ├── SelectorExpr
│   ├── TypeAssertExpr
│   ├── ChanLit (Moxie: &chan T{cap: 10})
│   ├── SliceLit (Moxie: &[]T{...})
│   ├── MapLit (Moxie: &map[K]V{...})
│   ├── TypeCoercion (Moxie FFI)
│   └── FFICall (Moxie: dlsym, dlopen, etc.)
├── Stmt (statements)
│   ├── DeclStmt
│   ├── EmptyStmt
│   ├── LabeledStmt
│   ├── ExprStmt
│   ├── SendStmt
│   ├── IncDecStmt
│   ├── AssignStmt
│   ├── GoStmt
│   ├── DeferStmt
│   ├── ReturnStmt
│   ├── BranchStmt
│   ├── BlockStmt
│   ├── IfStmt
│   ├── SwitchStmt
│   ├── TypeSwitchStmt
│   ├── SelectStmt
│   ├── ForStmt
│   └── RangeStmt
└── Decl (declarations)
    ├── ImportDecl
    ├── ConstDecl
    ├── VarDecl
    ├── TypeDecl
    └── FuncDecl
```

## Moxie-Specific Features

### Explicit Pointer Types

Moxie makes reference types explicit with pointer syntax:

```go
// SliceType with Pointer=true
&[]int{1, 2, 3}  → SliceLit with Ampersand

// MapType with Pointer=true
&map[string]int{} → MapLit with Ampersand

// ChanType with Pointer=true
&chan int{cap: 10} → ChanLit with Ampersand
```

### Built-in Functions

Moxie adds new built-in functions:

```go
clone(x)  // Deep copy (Token: CLONE)
free(x)   // Explicit memory release (Token: FREE)
grow(x,n) // Pre-allocate capacity (Token: GROW)
clear(x)  // Reset container (Token: CLEAR)
```

### FFI Support

Native FFI with special expression types:

```go
// FFICall node
dlopen("libc.so", RTLD_LAZY)
dlsym[func(*byte) int64](lib, "strlen")
dlclose(lib)

// TypeCoercion node
(*[]uint32)(bytes)  // Zero-copy type coercion
```

### Mutable Strings

In Moxie, `string` is an alias for `*[]byte`:

```go
s := "hello"  // Type: string (= *[]byte)
s[0] = 'H'    // Mutable
```

### Generics

Full support for type parameters:

```go
// FuncType.TypeParams field
func filter[T any](items *[]T, pred func(T) bool) *[]T

// TypeSpec.TypeParams field
type Stack[T any] struct { ... }
```

## Usage Example

### Building an AST

```go
// Create a simple function declaration
funcDecl := &ast.FuncDecl{
    Name: &ast.Ident{Name: "add"},
    Type: &ast.FuncType{
        Params: &ast.FieldList{
            List: []*ast.Field{
                {
                    Names: []*ast.Ident{{Name: "a"}, {Name: "b"}},
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
                        X: &ast.Ident{Name: "a"},
                        Op: ast.ADD,
                        Y: &ast.Ident{Name: "b"},
                    },
                },
            },
        },
    },
}
```

### Walking the AST

```go
// Visitor pattern
type Visitor struct{}

func (v *Visitor) Visit(node ast.Node) {
    switch n := node.(type) {
    case *ast.FuncDecl:
        fmt.Printf("Function: %s\n", n.Name.Name)
    case *ast.BinaryExpr:
        fmt.Printf("Binary op: %s\n", n.Op)
    }
}
```

## Position Tracking

Every node implements position tracking:

```go
type Position struct {
    Filename string // Source file name
    Offset   int    // Byte offset (0-based)
    Line     int    // Line number (1-based)
    Column   int    // Column number (1-based)
}

// Usage
node.Pos()  // Start position
node.End()  // End position
```

## Next Steps

1. **AST Builder** - Transform ANTLR parse tree → AST
2. **Symbol Table** - Track declarations and scopes
3. **Type Checker** - Semantic analysis and type checking
4. **Code Generator** - Walk AST to generate target code

## Design Principles

1. **Explicit over implicit** - All pointer types are explicit
2. **Position-aware** - Every node tracks source location for error reporting
3. **Type-safe** - Strong typing with Go interfaces
4. **Extensible** - Easy to add new node types for Moxie features
5. **Compatible** - Similar to Go's AST for familiarity

## Comparison with Go AST

### Similar to Go
- Overall structure and node hierarchy
- Statement and expression types
- Position tracking approach
- Visitor pattern support

### Different from Go
- `SliceType.Pointer` field for `*[]T`
- `MapType.Pointer` field for `*map[K]V`
- `ChanType.Pointer` field for `*chan T`
- New literal types: `ChanLit`, `SliceLit`, `MapLit`
- `FFICall` and `TypeCoercion` nodes
- New tokens: `CLONE`, `FREE`, `GROW`, `CLEAR`, `DLOPEN`, `DLSYM`, `DLCLOSE`
- String type is mutable (= `*[]byte`)

## References

- [Moxie Language Specification](../../go-language-revision.md)
- [ANTLR Parser](../antlr/README_PARSER.md)
- [Go AST Package](https://pkg.go.dev/go/ast)
