# Moxie ANTLR Parser

This package contains the ANTLR4-generated lexer and parser for the Moxie programming language.

## Generated Files

The following files are **auto-generated** from `grammar/Moxie.g4`. Do not edit them directly:

- `moxie_lexer.go` - Lexical analyzer (tokenizer)
- `moxie_parser.go` - Syntax parser
- `moxie_listener.go` - Listener interface for tree walking
- `moxie_base_listener.go` - Base listener implementation
- `Moxie.tokens` - Token definitions
- `MoxieLexer.tokens` - Lexer token definitions
- `Moxie.interp` - Parser interpreter data
- `MoxieLexer.interp` - Lexer interpreter data

## Usage

### Basic Parsing

```go
package main

import (
    "fmt"

    "github.com/antlr4-go/antlr/v4"
    "github.com/mleku/moxie/pkg/antlr"
)

func main() {
    // Create input from source code
    input := antlr.NewInputStream(`package main

func main() {
    s := &[]int{1, 2, 3}
    fmt.Println("Hello, Moxie!")
}
`)

    // Create lexer
    lexer := antlr.NewMoxieLexer(input)

    // Create token stream
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

    // Create parser
    parser := antlr.NewMoxieParser(stream)

    // Parse the source file
    tree := parser.SourceFile()

    // tree is now the root of the parse tree
    fmt.Printf("Package: %s\n", tree.PackageClause().GetText())
}
```

### Parsing from File

```go
func ParseFile(filename string) (*antlr.SourceFileContext, error) {
    input, err := antlr.NewFileStream(filename)
    if err != nil {
        return nil, err
    }

    lexer := antlr.NewMoxieLexer(input)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    parser := antlr.NewMoxieParser(stream)

    return parser.SourceFile(), nil
}
```

### Error Handling

```go
type ErrorListener struct {
    *antlr.DefaultErrorListener
    Errors []string
}

func (l *ErrorListener) SyntaxError(
    recognizer antlr.Recognizer,
    offendingSymbol interface{},
    line, column int,
    msg string,
    e antlr.RecognitionException,
) {
    l.Errors = append(l.Errors,
        fmt.Sprintf("line %d:%d %s", line, column, msg))
}

func ParseWithErrors(input string) (tree *antlr.SourceFileContext, errors []string) {
    is := antlr.NewInputStream(input)
    lexer := antlr.NewMoxieLexer(is)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    parser := antlr.NewMoxieParser(stream)

    // Add custom error listener
    errorListener := &ErrorListener{}
    parser.RemoveErrorListeners()
    parser.AddErrorListener(errorListener)

    tree = parser.SourceFile()
    return tree, errorListener.Errors
}
```

### Walking the Parse Tree

Use the listener pattern to traverse the parse tree:

```go
type MyListener struct {
    *antlr.BaseMoxieListener
}

func (l *MyListener) EnterFunctionDecl(ctx *antlr.FunctionDeclContext) {
    fmt.Printf("Found function: %s\n", ctx.IDENTIFIER().GetText())
}

func (l *MyListener) EnterVarSpec(ctx *antlr.VarSpecContext) {
    fmt.Println("Found variable declaration")
}

func WalkTree(tree antlr.Tree) {
    listener := &MyListener{}
    antlr.ParseTreeWalkerDefault.Walk(listener, tree)
}
```

### Visitor Pattern (Custom)

For more control, implement a custom visitor:

```go
type Visitor struct {
    functions []string
}

func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
    switch t := tree.(type) {
    case *antlr.FunctionDeclContext:
        v.functions = append(v.functions, t.IDENTIFIER().GetText())
    }

    // Visit children
    for i := 0; i < tree.GetChildCount(); i++ {
        child := tree.GetChild(i)
        if child != nil {
            child.Accept(v)
        }
    }

    return nil
}

func (v *Visitor) VisitChildren(node antlr.RuleNode) interface{} {
    for i := 0; i < node.GetChildCount(); i++ {
        node.GetChild(i).(antlr.ParseTree).Accept(v)
    }
    return nil
}

func (v *Visitor) VisitTerminal(node antlr.TerminalNode) interface{} {
    return nil
}

func (v *Visitor) VisitErrorNode(node antlr.ErrorNode) interface{} {
    return nil
}
```

## Parser Entry Points

### Main Entry Point

- `SourceFile()` - Parses a complete Moxie source file

### Sub-parsers

You can also parse specific language constructs:

- `PackageClause()` - Parse package declaration
- `ImportDecl()` - Parse import declaration
- `TopLevelDecl()` - Parse top-level declaration
- `FunctionDecl()` - Parse function declaration
- `Type_()` - Parse type expression
- `Expression()` - Parse expression
- `Statement()` - Parse statement
- `Block()` - Parse block

Example:

```go
// Parse just an expression
func ParseExpression(input string) *antlr.ExpressionContext {
    is := antlr.NewInputStream(input)
    lexer := antlr.NewMoxieLexer(is)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    parser := antlr.NewMoxieParser(stream)

    return parser.Expression()
}

// Usage
expr := ParseExpression(`s1 | s2 | "suffix"`)
```

## Common Parse Tree Nodes

### SourceFileContext

Root node of the parse tree.

Methods:
- `PackageClause()` - Get package declaration
- `AllImportDecl()` - Get all imports
- `AllTopLevelDecl()` - Get all top-level declarations

### FunctionDeclContext

Function declaration node.

Methods:
- `IDENTIFIER()` - Function name
- `Signature()` - Function signature
- `Block()` - Function body
- `TypeParameters()` - Generic type parameters (if any)

### ExpressionContext

Expression node. Has various subtypes:
- `UnaryExpression`
- `MultiplicativeExpr`
- `AdditiveExpr`
- `ConcatenationExpr` - `a | b` concatenation
- `RelationalExpr`
- `LogicalAndExpr`
- `LogicalOrExpr`

### TypeContext

Type expression node. Subtypes:
- `NamedType` - Simple type name
- `TypeLiteral` - Array, slice, map, struct, etc.
- `ParenType` - Parenthesized type
- `ConstType` - `const T`

## Moxie-Specific Features

### Concatenation Operator

The `|` operator is parsed as `ConcatenationExpr`:

```go
// Input: result := s1 | s2
expr := parser.Expression()
// expr is a ConcatenationExpr with left=s1, right=s2
```

### Explicit Pointer Types

Slices, maps, and channels use explicit pointer syntax:

```go
// Input: s := &[]int{1, 2, 3}
// Parse tree has pointer type with slice type
```

### Type Casting with Endianness

Zero-copy casts include optional endianness:

```go
// Input: (*[]uint32, LittleEndian)(bytes)
// Parsed as SliceCastEndianExpr
```

### Built-in Functions

Built-ins are parsed as regular function calls:
- `clone()`
- `copy()`
- `grow()`
- `clear()`
- `free()`
- `dlopen()`, `dlsym()`, etc.

## Testing

Run the test suite:

```bash
go test -v
```

See `example_test.go` for comprehensive examples of parsing various Moxie constructs.

## Regenerating

If you modify `grammar/Moxie.g4`, regenerate the parser:

```bash
# Using Docker (if available)
docker run --rm -v $(pwd)/grammar:/grammar -v $(pwd)/pkg/antlr:/output \
    antlr/antlr4:latest -Dlanguage=Go -o /output -package antlr /grammar/Moxie.g4

# Using ANTLR jar directly (requires Java)
java -jar /tmp/antlr-4.13.1-complete.jar \
    -Dlanguage=Go -o pkg/antlr -package antlr grammar/Moxie.g4

# Move generated files if needed
mv pkg/antlr/grammar/* pkg/antlr/
rmdir pkg/antlr/grammar
```

## Dependencies

- `github.com/antlr4-go/antlr/v4` - ANTLR Go runtime

## License

Generated code inherits the license of the Moxie project (BSD-style).

## References

- [ANTLR Documentation](https://github.com/antlr/antlr4/blob/master/doc/index.md)
- [Moxie Grammar](../../grammar/Moxie.g4)
- [Moxie Language Specification](../../go-language-revision.md)
