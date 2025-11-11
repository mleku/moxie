package ast

// ============================================================================
// Literal Nodes
// ============================================================================

// BasicLit represents a literal of basic type (int, float, string, char, etc.).
type BasicLit struct {
	ValuePos Position  // Position of the literal
	Kind     LitKind   // Literal kind
	Value    string    // Literal value as a string
}

type LitKind int

const (
	IntLit LitKind = iota    // 123, 0x1A, 0o777, 0b1010
	FloatLit                 // 1.23, 1.23e10, 0x1.Fp-2
	ImagLit                  // 1.23i
	RuneLit                  // 'a', '\n', '\x00'
	StringLit                // "hello", `raw string`
)

func (l *BasicLit) Pos() Position { return l.ValuePos }
func (l *BasicLit) End() Position {
	return Position{Line: l.ValuePos.Line, Column: l.ValuePos.Column + len(l.Value)}
}
func (l *BasicLit) node() {}
func (l *BasicLit) expr() {}

// ============================================================================
// Token Type (for operators and keywords)
// ============================================================================

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Literals
	IDENT  // main, foo, x, y
	INT    // 123, 0x1A
	FLOAT  // 1.23
	IMAG   // 1.23i
	CHAR   // 'a'
	STRING // "hello"
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operator_end

	keyword_beg
	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR

	FUNC
	GO
	GOTO
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR

	// Moxie-specific keywords/built-ins
	CLONE  // clone() built-in
	FREE   // free() built-in
	GROW   // grow() built-in
	CLEAR  // clear() built-in

	DLOPEN  // dlopen() FFI function
	DLSYM   // dlsym() FFI function
	DLCLOSE // dlclose() FFI function

	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK:    "break",
	CASE:     "case",
	CHAN:     "chan",
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FALLTHROUGH: "fallthrough",
	FOR:         "for",

	FUNC:   "func",
	GO:     "go",
	GOTO:   "goto",
	IF:     "if",
	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	SELECT: "select",
	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",

	CLONE: "clone",
	FREE:  "free",
	GROW:  "grow",
	CLEAR: "clear",

	DLOPEN:  "dlopen",
	DLSYM:   "dlsym",
	DLCLOSE: "dlclose",
}

// String returns the string representation of the token.
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + itoa(int(tok)) + ")"
	}
	return s
}

// IsLiteral returns true if the token is a literal.
func (tok Token) IsLiteral() bool {
	return literal_beg < tok && tok < literal_end
}

// IsOperator returns true if the token is an operator.
func (tok Token) IsOperator() bool {
	return operator_beg < tok && tok < operator_end
}

// IsKeyword returns true if the token is a keyword.
func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}

// Precedence returns the operator precedence of the binary operator.
func (tok Token) Precedence() int {
	switch tok {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return 0
}
