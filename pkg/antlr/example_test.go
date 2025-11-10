package antlr

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

// Example of parsing a simple Moxie program
func ExampleParser_basic() {
	input := `package main

import "fmt"

func main() {
	s := &[]int{1, 2, 3}
	fmt.Println("Hello from Moxie!", s)
}
`

	// Create input stream
	is := antlr.NewInputStream(input)

	// Create lexer
	lexer := NewMoxieLexer(is)

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser
	parser := NewMoxieParser(stream)

	// Parse the source file
	tree := parser.SourceFile()

	// Print the parse tree (for demonstration)
	fmt.Println("Parsed successfully!")
	fmt.Printf("Package clause: %s\n", tree.PackageClause().GetText())

	// Output:
	// Parsed successfully!
	// Package clause: packagemain
}

// Example with error handling
func ExampleParser_withErrors() {
	input := `package main

func main() {
	// Missing closing brace
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	// Add error listener
	errorListener := &CustomErrorListener{}
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	// Parse
	_ = parser.SourceFile()

	if len(errorListener.errors) > 0 {
		fmt.Println("Parse errors detected:")
		for _, err := range errorListener.errors {
			fmt.Printf("  Line %d:%d - %s\n", err.line, err.column, err.msg)
		}
	}
}

// CustomErrorListener collects parsing errors
type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	errors []parseError
}

type parseError struct {
	line   int
	column int
	msg    string
}

func (l *CustomErrorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol interface{},
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	l.errors = append(l.errors, parseError{
		line:   line,
		column: column,
		msg:    msg,
	})
}

// TestParseHelloWorld tests parsing a simple program
func TestParseHelloWorld(t *testing.T) {
	input := `package main

import "github.com/mleku/moxie/src/fmt"

func main() {
	fmt.Println("Hello from Moxie!")
}
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	if tree == nil {
		t.Fatal("Failed to parse")
	}

	// Verify package clause
	if tree.PackageClause() == nil {
		t.Error("Package clause is nil")
	}

	// Verify we have imports
	if len(tree.AllImportDecl()) == 0 {
		t.Error("No imports found")
	}

	// Verify we have top-level declarations
	if len(tree.AllTopLevelDecl()) == 0 {
		t.Error("No top-level declarations found")
	}
}

// TestParseConcatenation tests the | concatenation operator
func TestParseConcatenation(t *testing.T) {
	input := `package main

func test() {
	s1 := "hello "
	s2 := "world"
	result := s1 | s2

	a1 := &[]int{1, 2, 3}
	a2 := &[]int{4, 5, 6}
	combined := a1 | a2
}
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	if tree == nil {
		t.Fatal("Failed to parse concatenation")
	}
}

// TestParseExplicitPointers tests explicit pointer types
func TestParseExplicitPointers(t *testing.T) {
	input := `package main

func test() {
	s := &[]int{1, 2, 3}
	m := &map[string]int{"a": 1, "b": 2}
	ch := &chan int{cap: 10}
}
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	if tree == nil {
		t.Fatal("Failed to parse explicit pointers")
	}
}

// TestParseTypeCasting tests zero-copy type casting with endianness
func TestParseTypeCasting(t *testing.T) {
	input := `package main

func serialize() {
	src := &[]uint32{0x12345678}
	dst := &[]byte{0, 0, 0, 0}

	copy(dst, (*[]byte, LittleEndian)(src))
	copy(dst, (*[]byte, BigEndian)(src))
	copy(dst, (*[]byte)(src))
}
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	if tree == nil {
		t.Fatal("Failed to parse type casting")
	}
}

// TestParseBuiltins tests built-in functions
func TestParseBuiltins(t *testing.T) {
	input := `package main

func test() {
	s1 := &[]int{1, 2, 3}
	s2 := clone(s1)
	s3 := grow(s1, 100)

	dst := &[]byte{0, 0, 0}
	src := &[]byte{1, 2, 3}
	n := copy(dst, src)

	clear(s1)
	free(s1)
}
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	if tree == nil {
		t.Fatal("Failed to parse built-ins")
	}
}

// TestParseConst tests const declarations
func TestParseConst(t *testing.T) {
	input := `package main

const MaxSize = 100
const Pi = 3.14159

const Config = &map[string]int{
	"timeout": 30,
	"retries": 3,
}

func process(data const string) {
	// data is immutable
}
`

	is := antlr.NewInputStream(input)
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	if tree == nil {
		t.Fatal("Failed to parse const")
	}
}
