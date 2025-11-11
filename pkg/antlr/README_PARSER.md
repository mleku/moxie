# Moxie Parser Test

This directory contains the ANTLR4-generated parser for the Moxie programming language.

## Files Generated

- `moxie_lexer.go` - Lexical analyzer (tokenizer)
- `moxie_parser.go` - Syntax analyzer (parser)
- `moxie_listener.go` - Listener interface for tree walking
- `moxie_base_listener.go` - Base listener implementation

## Test Files

### `print_ast_test.go`

Contains three test functions that demonstrate parsing and AST analysis:

1. **TestPrintAST** - Parses `example.x` and prints:
   - Parse errors (if any)
   - Compact parse tree representation
   - Readable parse tree (depth-limited to 5 levels)

2. **TestPrintASTWithListener** - Uses the listener pattern to extract:
   - Package name
   - Imports
   - Type declarations
   - Constants
   - Functions
   - Variables

3. **CustomErrorListener** - Collects parsing errors during analysis

## Running the Tests

```bash
# Run all tests
GOROOT=/home/mleku/go GOTOOLCHAIN=local /home/mleku/go/bin/go test -v

# Run specific test
GOROOT=/home/mleku/go GOTOOLCHAIN=local /home/mleku/go/bin/go test -v -run TestPrintASTWithListener
```

## Example Output

```
=== AST SUMMARY ===
Package: main

Imports (1):
  - "fmt"

Types (1):
  - Server

Constants (1):
  - MaxConnections

Functions (3):
  - main
  - divide
  - filter

Variables (0):
```

## Parse Tree Structure

ANTLR generates a **parse tree** (concrete syntax tree) that includes:
- **Rule nodes** - Grammar rules like `sourceFile`, `functionDecl`, `expression`
- **Terminal nodes** - Tokens like `IDENTIFIER`, `INT_LIT`, keywords

### Walking the Tree

Two patterns are supported:

1. **Listener Pattern** (push-based)
   - Implement `MoxieListener` interface
   - ANTLR walks the tree and calls your methods
   - Good for: Collecting information, simple transformations

2. **Visitor Pattern** (pull-based)
   - Would need to regenerate with `-visitor` flag
   - You control the tree traversal
   - Good for: Complex transformations, evaluation

## Next Steps

To build a proper AST for code generation:

1. **Define AST node types** - Create Go structs for each language construct
2. **Transform parse tree → AST** - Use listener/visitor to build simplified AST
3. **Semantic analysis** - Type checking, symbol resolution, scope analysis
4. **Code generation** - Walk the AST and emit target code

### Current Status

✅ Lexer - Tokenizes Moxie source code
✅ Parser - Builds parse tree from tokens
✅ Tree walking - Can extract declarations and structure
⚠️ Grammar issues - Some constructs in `example.x` don't parse correctly yet
⏳ AST - Need to define proper AST node structures
⏳ Semantic analysis - Type checking, symbol tables
⏳ Code generation - Target code emission

## Grammar Issues Found

The test reveals several grammar issues that need fixing:

- Line 3: Missing semicolon/EOS handling after package clause
- Line 6: Const expression parsing (numeric literals)
- Line 9: Struct field type parsing (`*[]byte`)
- Line 16-19: Statement parsing in function body
- Line 47: FFI function type syntax
- Line 59: Multiple return types syntax
- Line 67: Generic function parameter syntax

These issues indicate the grammar needs refinement to fully support Moxie's extended syntax.
