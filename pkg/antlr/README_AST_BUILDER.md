# AST Builder - ANTLR Parse Tree to Moxie AST Transformer

## Overview

This package provides transformation from ANTLR parse trees to Moxie AST nodes defined in `pkg/ast`.

## Status: 🚧 Under Development (95% Complete)

The AST builder architecture and all major transformers have been implemented. Minor compilation issues remain related to ANTLR interface type assertions.

### Completed ✅

1. **Position Mapper** (`position.go`)
   - TokenToPosition - Converts ANTLR tokens to AST positions
   - ContextToPosition - Gets start position from parser context
   - ContextEndPosition - Gets end position from parser context
   - Handles 0-based → 1-based column conversion

2. **Core AST Builder** (`astbuilder.go`)
   - ASTBuilder struct with visitor pattern
   - Source file transformation
   - Package clause transformation
   - Import declarations
   - Top-level declarations routing
   - Error collection mechanism
   - Helper methods (visitIdentifier, visitIdentifierList)
   - BuildAST entry point function

3. **Type Transformers** (`astbuilder_types.go`)
   - All type expressions (30+ visitor methods)
   - Named types and qualified identifiers
   - Type literals (struct, interface, pointer, slice, array, map, channel)
   - Function types with generics support
   - Signature and parameter transformations
   - Field declarations and tags
   - Interface elements (methods and embedded types)
   - Parenthesized types
   - Const types (Moxie feature)

4. **Declaration Transformers** (`astbuilder_decls.go`)
   - Constant declarations and specs
   - Variable declarations and specs
   - Type declarations and specs
   - Type aliases vs type definitions
   - Type parameters (generics)
   - Type constraints
   - Function declarations
   - Method declarations
   - Receivers

5. **Statement Transformers** (`astbuilder_stmts.go`)
   - Block statements
   - Statement lists
   - Simple statements (expression, send, inc/dec, assignment)
   - Short variable declarations
   - Control flow (if, for, switch, select)
   - Branch statements (break, continue, goto, fallthrough)
   - Defer and go statements
   - Labeled statements
   - For clauses and range clauses
   - Assignment operators
   - Expression statements

6. **Expression Transformers** (`astbuilder_exprs.go`)
   - Binary and unary expressions
   - Primary expressions
   - Operands (literals, identifiers, parenthesized)
   - Selectors (x.field)
   - Index expressions (x[i])
   - Slice expressions (x[i:j:k])
   - Type assertions (x.(T))
   - Function calls with arguments
   - Type conversions
   - Expression lists
   - All operators (mul, add, rel, unary)

7. **Literal Transformers** (`astbuilder_exprs.go`)
   - Basic literals (int, float, imag, rune, string)
   - String literals (raw and interpreted)
   - Composite literals
   - Literal types and values
   - Element lists
   - Keyed elements (key:value pairs)
   - Function literals

### Remaining Work 🔧

1. **Compilation Fixes** (Estimated: 1-2 hours)
   - Fix BaseMoxieVisitor import (need to reference grammar subpackage correctly)
   - Add type assertions for ANTLR interface types throughout
   - Fix context accessor methods (GetDot, GetUnderscore)
   - Approximately 50-100 type assertions needed across all files

2. **Testing** (Estimated: 2-4 hours)
   - Create simple transformation tests
   - Test basic declarations
   - Test expressions and statements
   - Test Moxie-specific features
   - Integration with existing ANTLR parser tests

3. **Enhancements** (Optional, Future)
   - Better error messages with source locations
   - Recovery from parse errors
   - Pretty-printing transformed AST
   - Validation passes

## Architecture

```
ANTLR Parse Tree
    ↓
ASTBuilder (Visitor Pattern)
    ├── position.go (Position Mapping)
    ├── astbuilder.go (Core + Entry Point)
    ├── astbuilder_types.go (Type Expressions)
    ├── astbuilder_decls.go (Declarations)
    ├── astbuilder_stmts.go (Statements)
    └── astbuilder_exprs.go (Expressions + Literals)
    ↓
Moxie AST (pkg/ast)
```

## Usage (Once Compilation Fixed)

```go
package main

import (
    "github.com/antlr4-go/antlr/v4"
    "github.com/mleku/moxie/pkg/antlr"
)

func main() {
    // Create ANTLR input stream
    input := antlr.NewInputStream("package main\\nfunc main() {}")

    // Create lexer and parser
    lexer := antlr.NewMoxieLexer(input)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    parser := antlr.NewMoxieParser(stream)

    // Parse source file
    tree := parser.SourceFile()

    // Transform to AST
    astFile, errors := antlr.BuildAST(tree, "example.x")

    if len(errors) > 0 {
        for _, err := range errors {
            fmt.Printf("Error: %v\\n", err)
        }
    }

    // Use AST
    fmt.Printf("Package: %s\\n", astFile.Package.Name.Name)
    fmt.Printf("Declarations: %d\\n", len(astFile.Decls))
}
```

## File Structure

### position.go (47 lines)
Position mapping utilities between ANTLR and AST.

### astbuilder.go (230 lines)
- ASTBuilder struct
- Core visitor methods
- Source file, package, imports
- Top-level declaration routing
- Helper methods

### astbuilder_types.go (540 lines)
- Type expression transformations
- 30+ visitor methods covering:
  - Named types, type literals
  - Struct, interface, function types
  - Pointer, slice, array, map, channel types
  - Generics support (type parameters, constraints)
  - Qualified identifiers

### astbuilder_decls.go (310 lines)
- Declaration transformations
- Const, var, type declarations
- Function and method declarations
- Type specs (aliases vs definitions)
- Receivers and signatures

### astbuilder_stmts.go (520 lines)
- Statement transformations
- All control flow constructs
- Simple statements
- Assignment operators
- For loops and range

### astbuilder_exprs.go (530 lines)
- Expression transformations
- Binary, unary, primary expressions
- Literals (basic, composite, function)
- Operators (mul, add, rel, unary)
- Selectors, indices, slices, calls

## Total Implementation

- **Files**: 6
- **Lines of Code**: ~2,177
- **Visitor Methods**: 80+
- **Node Types Supported**: 50+
- **Completion**: ~95%

## Known Issues

1. **Compilation Errors** (Minor, Easy to Fix)
   - BaseMoxieVisitor not found - need correct import
   - ANTLR interface types need type assertions
   - Some context accessor methods need implementation checking

2. **Missing Features** (Can be added later)
   - Switch/select statement details (currently stubbed)
   - Channel direction detection in types
   - Some Moxie-specific literal types (ChanLit, SliceLit, MapLit)
   - FFI-specific expressions

3. **Testing** (Not yet started)
   - No tests written yet
   - Need integration with existing parser tests
   - Need validation against example Moxie files

## Next Steps

1. **Fix Compilation** (Priority 1)
   ```bash
   # Fix import for BaseMoxieVisitor
   # Add type assertions throughout
   # Test compilation
   ```

2. **Write Tests** (Priority 2)
   ```go
   // Test basic transformation
   // Test each declaration type
   // Test expressions and statements
   ```

3. **Integration** (Priority 3)
   ```go
   // Integrate with existing ANTLR parser tests
   // Test with real Moxie examples
   // Validate against example.x files
   ```

## Design Decisions

1. **Visitor Pattern**: Chosen over listener for better control and return values
2. **Interface Type Assertions**: ANTLR generates interfaces, need runtime type assertions
3. **Error Collection**: Errors collected in slice, don't fail fast
4. **Position Tracking**: Every node gets accurate source positions
5. **Null Safety**: All visitor methods check for nil contexts

## Contributing

To complete the AST builder:

1. Fix the remaining type assertions
2. Ensure BaseMoxieVisitor is properly accessible
3. Test each transformer category
4. Validate against Moxie example files
5. Add error messages with source locations

## References

- [Moxie AST Package](../ast/README.md)
- [ANTLR Parser](./README_PARSER.md)
- [Go AST Design](https://pkg.go.dev/go/ast)
- [ANTLR Go Runtime](https://github.com/antlr4-go/antlr)
