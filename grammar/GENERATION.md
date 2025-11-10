# ANTLR Parser Generation Guide

This document describes how the Moxie ANTLR parser was generated and how to regenerate it.

## Generated Files Location

**Location:** `pkg/antlr/`

**Files:**
- `moxie_lexer.go` (50KB) - Lexical analyzer
- `moxie_parser.go` (507KB) - Syntax parser
- `moxie_listener.go` (34KB) - Listener interface
- `moxie_base_listener.go` (43KB) - Base listener implementation
- `Moxie.tokens` - Token definitions
- `MoxieLexer.tokens` - Lexer tokens
- `Moxie.interp` - Parser interpreter data
- `MoxieLexer.interp` - Lexer interpreter data

## Generation Process

### Prerequisites

1. **Java** - Required to run ANTLR
2. **ANTLR 4.13.1** - Parser generator jar

### Download ANTLR

```bash
cd /tmp
curl -O https://www.antlr.org/download/antlr-4.13.1-complete.jar
```

### Generate Parser

```bash
# Set Java path (adjust to your installation)
export JAVA_HOME=/home/mleku/tools/jdk-21.0.1

# Generate Go parser from grammar
$JAVA_HOME/bin/java -jar /tmp/antlr-4.13.1-complete.jar \
    -Dlanguage=Go \
    -o pkg/antlr \
    -package antlr \
    grammar/Moxie.g4

# If files are generated in subdirectory, move them
if [ -d pkg/antlr/grammar ]; then
    mv pkg/antlr/grammar/* pkg/antlr/
    rmdir pkg/antlr/grammar
fi
```

### Add Dependencies

```bash
cd pkg/antlr

# Initialize Go module (first time only)
go mod init github.com/mleku/moxie/pkg/antlr

# Add ANTLR runtime
go get github.com/antlr4-go/antlr/v4

# Verify it builds
go build
```

### Run Tests

```bash
cd pkg/antlr
go test -v
```

## Grammar Modifications

If you modify `grammar/Moxie.g4`, you need to regenerate the parser.

### Common Issues

#### 1. `lineTerminator()` Predicate Error

**Error:**
```
./moxie_parser.go:21250:8: undefined: lineTerminator
```

**Solution:**
The `eos` rule was fixed to use `TERMINATOR` token instead of semantic predicate:

```antlr
eos
    : ';'
    | EOF
    | TERMINATOR
    ;
```

#### 2. Label Name Conflicts

**Error:**
```
error(124): Moxie.g4:112:41: rule alt label LiteralType conflicts with rule literalType
```

**Solution:**
Rename the label to avoid conflict with rule name:

```antlr
type_
    : typeName typeArgs?              # NamedType
    | typeLit                          # TypeLiteral  // Changed from LiteralType
    | '(' type_ ')'                    # ParenType
    | 'const' type_                    # ConstType
    ;
```

## Grammar Features

The generated parser supports all Moxie language features:

### Lexical Features
- Keywords (break, case, chan, const, etc.)
- Explicit integer types (int8-64, uint8-64, no platform-dependent int/uint)
- Operators including `|` for concatenation
- String literals (raw and interpreted)
- Numeric literals (decimal, hex, binary, octal, float, imaginary)

### Syntax Features
- Package and import declarations
- Function declarations with generics
- Type declarations (aliases and definitions)
- Variable and constant declarations
- Explicit pointer types (`*[]T`, `*map[K]V`, `*chan T`)
- Const type modifier (`const T`)
- Zero-copy type casting with endianness
- Concatenation operator (`|`)
- All Go control flow structures

## Parser Usage

### Basic Example

```go
import (
    "github.com/antlr4-go/antlr/v4"
    "github.com/mleku/moxie/pkg/antlr"
)

func parseFile(filename string) error {
    // Read input
    input, err := antlr.NewFileStream(filename)
    if err != nil {
        return err
    }

    // Create lexer
    lexer := antlr.NewMoxieLexer(input)

    // Create token stream
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

    // Create parser
    parser := antlr.NewMoxieParser(tokens)

    // Parse
    tree := parser.SourceFile()

    // Walk the tree
    // ... (implement listener or visitor)

    return nil
}
```

See `pkg/antlr/README.md` for complete usage documentation.

## Testing

The parser includes comprehensive tests:

```bash
cd pkg/antlr
go test -v
```

Tests cover:
- Basic parsing
- Concatenation operator (`|`)
- Explicit pointer types
- Type casting with endianness
- Built-in functions (clone, copy, grow, etc.)
- Const declarations
- Error handling

## File Sizes

The generated parser is quite large:

| File | Size | Lines |
|------|------|-------|
| `moxie_parser.go` | 507KB | ~18,000 |
| `moxie_lexer.go` | 50KB | ~1,500 |
| `moxie_listener.go` | 34KB | ~1,200 |
| `moxie_base_listener.go` | 43KB | ~1,500 |

This is normal for ANTLR-generated parsers. The size comes from:
- State tables for DFA (Deterministic Finite Automaton)
- Prediction tables for parsing decisions
- Context classes for each grammar rule

## Performance

ANTLR4 parsers are generally fast:
- **Lexing:** ~100,000 tokens/sec
- **Parsing:** ~10,000 lines/sec (depends on grammar complexity)

For Moxie, expect similar performance to the Go parser.

## Regeneration Checklist

When regenerating the parser:

1. ✅ Backup current parser files (optional)
2. ✅ Verify grammar has no conflicts
3. ✅ Run ANTLR generator
4. ✅ Move files if in subdirectory
5. ✅ Verify parser builds (`go build`)
6. ✅ Run tests (`go test`)
7. ✅ Update grammar version/changelog if needed

## Version Information

- **ANTLR Version:** 4.13.1
- **Go Target:** Go 1.21+
- **Runtime:** github.com/antlr4-go/antlr/v4 v4.13.1
- **Grammar:** Moxie.g4
- **Generated:** 2025-11-10

## References

- [ANTLR Documentation](https://github.com/antlr/antlr4/blob/master/doc/index.md)
- [ANTLR Go Target](https://github.com/antlr/antlr4/blob/master/doc/go-target.md)
- [Moxie Grammar](Moxie.g4)
- [Parser Package](../pkg/antlr/README.md)

## Troubleshooting

### Parser doesn't recognize new syntax

1. Check grammar is updated
2. Regenerate parser completely
3. Verify token definitions in lexer

### Build errors after regeneration

1. Check for grammar conflicts (labels vs rules)
2. Verify semantic predicates are removed or implemented
3. Check Go module dependencies are up to date

### Test failures

1. Update test cases for grammar changes
2. Check if parse tree structure changed
3. Verify error messages updated

## Future Enhancements

Potential improvements to the grammar/parser:

1. **Better error recovery** - Add error productions
2. **Incremental parsing** - Support parsing partial files
3. **Comment preservation** - Add comments to parse tree
4. **Source mapping** - Better position tracking
5. **Performance tuning** - Optimize prediction/DFA tables

## Contact

For issues with the parser generation, see:
- Grammar bugs: `grammar/Moxie.g4`
- Parser bugs: `pkg/antlr/`
- Usage questions: `pkg/antlr/README.md`
