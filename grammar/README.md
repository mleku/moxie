# Moxie ANTLR4 Grammar

This directory contains the ANTLR4 grammar specification for the Moxie programming language.

## Files

- `Moxie.g4` - Complete ANTLR4 grammar for lexical analysis and parsing

## Overview

The Moxie grammar is based on the Go language specification with modifications defined in `go-language-revision.md`. It uses ANTLR4 to provide both lexical analysis (tokenization) and syntactic analysis (parsing).

## Key Differences from Go

### 1. Explicit Pointer Types for Reference Types

**Go:**
```go
s := []int{1, 2, 3}        // Implicit reference
m := make(map[string]int)  // Implicit reference
ch := make(chan int, 10)   // Implicit reference
```

**Moxie:**
```go
s := &[]int{1, 2, 3}         // Explicit pointer
m := &map[string]int{}       // Explicit pointer
ch := &chan int{cap: 10}     // Explicit pointer
```

The grammar supports both forms for parsing, but semantically Moxie requires the explicit `*` prefix.

### 2. No `make()` Function

The `make()` built-in is eliminated. Use composite literals with pointer syntax instead.

### 3. Mutable Strings

Strings are mutable in Moxie (aliased as `*[]byte`). The grammar treats strings the same syntactically but they have different semantics.

### 4. Concatenation Operator: | (Vertical Bar)

Moxie uses `|` for concatenation of strings and slices, following standard cryptographic notation where `a | b` means concatenation:

```go
s1 := "hello "
s2 := "world!"
s3 := s1 | s2  // "hello world!"

a1 := &[]int{1, 2, 3}
a2 := &[]int{4, 5, 6}
a3 := a1 | a2  // &[]int{1, 2, 3, 4, 5, 6}
```

The `append()` built-in is removed - use `|` instead:

```go
// Before (Go):
s = append(s, 4, 5, 6)

// After (Moxie):
s = s | &[]int{4, 5, 6}
```

The grammar includes a `ConcatenationExpr` production: `expression '|' expression`

### 5. Const for All Types

The `const` keyword works for all types (not just primitives):

```go
const Config = &map[string]int{"timeout": 30}
const Message = "immutable"
```

The grammar includes a `ConstType` production: `'const' type_`

### 6. Zero-Copy Type Coercion

New casting syntax for zero-copy slice reinterpretation with endianness control:

```go
// Zero-copy cast with native endian
u64s := (*[]uint64)(bytes)

// Zero-copy cast with specific endianness
u64s := (*[]uint64, LittleEndian)(bytes)
u32s := (*[]uint32, BigEndian)(bytes)

// Explicit copy
u64s := &(*[]uint64)(bytes)
```

Grammar productions:
- `SliceCastExpr` - Zero-copy cast
- `SliceCastEndianExpr` - Zero-copy cast with endianness
- `SliceCastCopyExpr` - Explicit copy cast
- `SliceCastCopyEndianExpr` - Copy cast with endianness

### 7. No Platform-Dependent Integer Types

Moxie removes `int` and `uint` types. You must use explicit sizes:

- `int8`, `int16`, `int32`, `int64`
- `uint8`, `uint16`, `uint32`, `uint64`

The grammar defines these as keywords but does NOT include `int` or `uint` tokens.

### 8. Built-in Functions

These are treated as identifiers (not keywords) but have special semantic meaning:

**Memory and Data Operations:**
- `clone(v)` - Deep copy (allocates new backing array)
- `copy(dst, src)` - Copy elements from src to dst (same type or compatible cast)
- `grow(s, n)` - Pre-allocate capacity
- `clear(v)` - Reset slice/map
- `free(v)` - Explicit memory release

**FFI Operations:**
- `dlopen(file, flags)` - Load dynamic library
- `dlsym[T](lib, name)` - Type-safe symbol lookup
- `dlclose(lib)` - Close library
- `dlerror()` - Get error string
- `dlopen_mem(data, flags)` - Load library from memory

**Standard Go Built-ins (unchanged):**
- `len(v)` - Length of slice/map/string/array/chan
- `cap(v)` - Capacity of slice/chan
- `delete(m, key)` - Remove map key
- `close(ch)` - Close channel
- `panic(v)` - Trigger panic
- `recover()` - Recover from panic
- `new(T)` - Allocate zero value
- `print(...)` - Print to stderr
- `println(...)` - Print line to stderr

**Removed built-ins:**
- `append()` - Use `|` concatenation operator instead
- `make()` - Use composite literals with explicit pointers

## Using the Grammar

### Prerequisites

Install ANTLR4:

```bash
# On macOS
brew install antlr

# On Ubuntu/Debian
sudo apt-get install antlr4

# Or download directly
cd /usr/local/lib
sudo curl -O https://www.antlr.org/download/antlr-4.13.1-complete.jar
export CLASSPATH=".:/usr/local/lib/antlr-4.13.1-complete.jar:$CLASSPATH"
alias antlr4='java -jar /usr/local/lib/antlr-4.13.1-complete.jar'
alias grun='java org.antlr.v4.gui.TestRig'
```

### Generate Go Parser

```bash
# Generate Go target code
antlr4 -Dlanguage=Go -o ../parser Moxie.g4

# This creates:
# - moxie_lexer.go       (lexical analyzer)
# - moxie_parser.go      (parser)
# - moxie_listener.go    (listener interface)
# - moxie_base_listener.go
```

### Generate Java Parser (for testing)

```bash
# Generate Java target
antlr4 Moxie.g4

# Compile
javac Moxie*.java

# Test parsing a file
grun Moxie sourceFile ../examples/hello/main.x

# View parse tree GUI
grun Moxie sourceFile -gui ../examples/hello/main.x

# View tokens
grun Moxie sourceFile -tokens ../examples/hello/main.x
```

### Integration with Moxie Compiler

The generated parser can be integrated into the Moxie compiler:

```go
package main

import (
    "github.com/antlr/antlr4/runtime/Go/antlr"
    "github.com/mleku/moxie/parser"
)

func parseFile(filename string) *parser.SourceFileContext {
    // Read input
    input, _ := antlr.NewFileStream(filename)

    // Create lexer
    lexer := parser.NewMoxieLexer(input)

    // Create token stream
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

    // Create parser
    p := parser.NewMoxieParser(tokens)

    // Parse
    tree := p.SourceFile()

    return tree
}
```

## Grammar Structure

### Source File

```
sourceFile
    : packageClause importDecl* topLevelDecl*
```

Every Moxie source file starts with a package declaration, followed by optional imports and top-level declarations.

### Declarations

- **Const** - Constant declarations with enhanced type support
- **Type** - Type definitions and aliases with generics
- **Var** - Variable declarations
- **Func** - Function declarations with generics
- **Method** - Method declarations on receivers

### Types

- **Basic Types** - `bool`, `int8`-`int64`, `uint8`-`uint64`, `float32`, `float64`, etc.
- **Array** - `[N]T`
- **Slice** - `*[]T` (explicit pointer)
- **Struct** - `struct { fields }`
- **Pointer** - `*T`
- **Function** - `func(params) result`
- **Interface** - `interface { methods }`
- **Map** - `*map[K]V` (explicit pointer)
- **Channel** - `*chan T` (explicit pointer)
- **Const** - `const T` (immutable type modifier)

### Statements

Standard Go control flow with same syntax:
- `if`, `for`, `switch`, `select`
- `go`, `defer`, `return`, `break`, `continue`, `fallthrough`, `goto`

### Expressions

Standard Go operators with modifications:
- `|` for slice/string concatenation (replaces `append()` and `+` for concatenation)
- Bitwise OR: Use `||` for logical OR (unchanged)
- Note: Bitwise OR (`|`) is now concatenation - this is a breaking change from Go
- Zero-copy type coercion syntax

## Semantic Analysis

The ANTLR grammar handles **syntax only**. Semantic analysis (type checking, const enforcement, etc.) must be performed in a separate phase using a visitor or listener pattern.

### Example Listener

```go
type MoxieSemanticAnalyzer struct {
    *parser.BaseMoxieListener

    errors []error
    symtab *SymbolTable
}

func (l *MoxieSemanticAnalyzer) EnterVarSpec(ctx *parser.VarSpecContext) {
    // Check variable declaration semantics
}

func (l *MoxieSemanticAnalyzer) EnterAssignment(ctx *parser.AssignmentContext) {
    // Check const mutation, type compatibility, etc.
}
```

## Testing the Grammar

Test files are available in `../examples/`:

```bash
# Test basic program
grun Moxie sourceFile -tree ../examples/hello/main.x

# Test slice pointers
grun Moxie sourceFile -tree ../examples/phase2/test_append.x

# Test string concatenation
grun Moxie sourceFile -tree ../examples/phase3/test_string_concat.x

# Test const enforcement
grun Moxie sourceFile -tree ../examples/phase6/test_const_enforcement.x
```

## Future Enhancements

Potential grammar additions:

1. **Inline Assembly** - For low-level operations
2. **Attribute Syntax** - `@inline`, `@noescape` annotations
3. **Pattern Matching** - Enhanced switch statements
4. **Macro System** - Compile-time code generation

## References

- [ANTLR4 Documentation](https://github.com/antlr/antlr4/blob/master/doc/index.md)
- [Go Language Specification](https://go.dev/ref/spec)
- [Moxie Language Revision](../go-language-revision.md)
- [ANTLR4 Grammar for Go](https://github.com/antlr/grammars-v4/tree/master/golang)

## Contributing

When modifying the grammar:

1. Test with existing `.x` files in `examples/`
2. Regenerate parser: `antlr4 -Dlanguage=Go Moxie.g4`
3. Update this documentation
4. Add test cases for new syntax
