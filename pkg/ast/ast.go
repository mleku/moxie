// Package ast defines the abstract syntax tree for the Moxie programming language.
package ast

// Node is the base interface for all AST nodes.
type Node interface {
	Pos() Position // Starting position of the node
	End() Position // Ending position of the node
	node()         // Marker method to ensure only AST nodes implement this interface
}

// Expr represents an expression node.
type Expr interface {
	Node
	expr()
}

// Stmt represents a statement node.
type Stmt interface {
	Node
	stmt()
}

// Decl represents a declaration node.
type Decl interface {
	Node
	decl()
}

// Spec represents a specification node (used in grouped declarations).
type Spec interface {
	Node
	spec()
}

// Position represents a source position with line and column information.
type Position struct {
	Filename string // Source file name
	Offset   int    // Byte offset in source (0-based)
	Line     int    // Line number (1-based)
	Column   int    // Column number (1-based)
}

// IsValid returns true if the position is valid (non-zero).
func (p Position) IsValid() bool {
	return p.Line > 0
}

// String returns a string representation of the position.
func (p Position) String() string {
	if !p.IsValid() {
		return "-"
	}
	s := p.Filename
	if s == "" {
		s = "<input>"
	}
	s += ":" + itoa(p.Line)
	if p.Column > 0 {
		s += ":" + itoa(p.Column)
	}
	return s
}

// itoa converts an integer to a string (simple implementation).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var buf [20]byte
	i := len(buf) - 1
	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i--
	}

	if negative {
		buf[i] = '-'
		i--
	}

	return string(buf[i+1:])
}
