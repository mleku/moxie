package ast

// Type represents a type expression.
type Type interface {
	Expr
	typeNode()
}

// ============================================================================
// Type Nodes
// ============================================================================

// Ident represents an identifier (used for type names, variable names, etc.).
type Ident struct {
	NamePos Position // Position of the identifier
	Name    string   // Identifier name
}

func (i *Ident) Pos() Position { return i.NamePos }
func (i *Ident) End() Position { return Position{Line: i.NamePos.Line, Column: i.NamePos.Column + len(i.Name)} }
func (i *Ident) node()         {}
func (i *Ident) expr()         {}
func (i *Ident) typeNode()     {}

// BasicType represents a built-in type (int, float64, bool, byte, rune, etc.).
type BasicType struct {
	NamePos Position // Position of the type keyword
	Kind    BasicKind
}

type BasicKind int

const (
	Invalid BasicKind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String // In Moxie, string = *[]byte
	Byte   // alias for uint8
	Rune   // alias for int32
)

func (t *BasicType) Pos() Position { return t.NamePos }
func (t *BasicType) End() Position { return t.NamePos }
func (t *BasicType) node()         {}
func (t *BasicType) expr()         {}
func (t *BasicType) typeNode()     {}

// PointerType represents a pointer type: *T
type PointerType struct {
	Star Position // Position of "*"
	Base Type     // Base type
}

func (t *PointerType) Pos() Position { return t.Star }
func (t *PointerType) End() Position { return t.Base.End() }
func (t *PointerType) node()         {}
func (t *PointerType) expr()         {}
func (t *PointerType) typeNode()     {}

// SliceType represents a slice type: []T or *[]T (explicit pointer in Moxie)
type SliceType struct {
	Lbrack  Position // Position of "["
	Pointer bool     // true if *[]T (explicit pointer)
	Elem    Type     // Element type
}

func (t *SliceType) Pos() Position { return t.Lbrack }
func (t *SliceType) End() Position { return t.Elem.End() }
func (t *SliceType) node()         {}
func (t *SliceType) expr()         {}
func (t *SliceType) typeNode()     {}

// ArrayType represents an array type: [N]T
type ArrayType struct {
	Lbrack Position // Position of "["
	Len    Expr     // Length expression (constant)
	Elem   Type     // Element type
}

func (t *ArrayType) Pos() Position { return t.Lbrack }
func (t *ArrayType) End() Position { return t.Elem.End() }
func (t *ArrayType) node()         {}
func (t *ArrayType) expr()         {}
func (t *ArrayType) typeNode()     {}

// MapType represents a map type: map[K]V or *map[K]V (explicit pointer in Moxie)
type MapType struct {
	Map     Position // Position of "map" keyword
	Lbrack  Position // Position of "["
	Pointer bool     // true if *map[K]V (explicit pointer)
	Key     Type     // Key type
	Value   Type     // Value type
}

func (t *MapType) Pos() Position { return t.Map }
func (t *MapType) End() Position { return t.Value.End() }
func (t *MapType) node()         {}
func (t *MapType) expr()         {}
func (t *MapType) typeNode()     {}

// ChanType represents a channel type: chan T, chan<- T, <-chan T, or *chan T (explicit pointer in Moxie)
type ChanType struct {
	Begin   Position // Position of "chan" keyword or "<-"
	Arrow   Position // Position of "<-" (invalid if no arrow)
	Dir     ChanDir  // Channel direction
	Pointer bool     // true if *chan T (explicit pointer)
	Value   Type     // Value type
}

type ChanDir int

const (
	ChanBoth ChanDir = iota // chan T (send and receive)
	ChanSend                // chan<- T (send only)
	ChanRecv                // <-chan T (receive only)
)

func (t *ChanType) Pos() Position { return t.Begin }
func (t *ChanType) End() Position { return t.Value.End() }
func (t *ChanType) node()         {}
func (t *ChanType) expr()         {}
func (t *ChanType) typeNode()     {}

// StructType represents a struct type.
type StructType struct {
	Struct Position     // Position of "struct" keyword
	Fields *FieldList   // List of fields
	Lbrace Position     // Position of "{"
	Rbrace Position     // Position of "}"
}

func (t *StructType) Pos() Position { return t.Struct }
func (t *StructType) End() Position { return t.Rbrace }
func (t *StructType) node()         {}
func (t *StructType) expr()         {}
func (t *StructType) typeNode()     {}

// InterfaceType represents an interface type.
type InterfaceType struct {
	Interface Position    // Position of "interface" keyword
	Methods   *FieldList  // List of methods
	Lbrace    Position    // Position of "{"
	Rbrace    Position    // Position of "}"
}

func (t *InterfaceType) Pos() Position { return t.Interface }
func (t *InterfaceType) End() Position { return t.Rbrace }
func (t *InterfaceType) node()         {}
func (t *InterfaceType) expr()         {}
func (t *InterfaceType) typeNode()     {}

// FuncType represents a function type.
type FuncType struct {
	Func       Position    // Position of "func" keyword (may be invalid)
	TypeParams *FieldList  // Type parameters (generics) [T any, U comparable]
	Params     *FieldList  // Function parameters
	Results    *FieldList  // Function results (return values)
}

func (t *FuncType) Pos() Position {
	if t.Func.IsValid() {
		return t.Func
	}
	if t.Params != nil {
		return t.Params.Pos()
	}
	return Position{}
}
func (t *FuncType) End() Position {
	if t.Results != nil {
		return t.Results.End()
	}
	if t.Params != nil {
		return t.Params.End()
	}
	return Position{}
}
func (t *FuncType) node()     {}
func (t *FuncType) expr()     {}
func (t *FuncType) typeNode() {}

// FieldList represents a list of fields (struct fields, function parameters, etc.).
type FieldList struct {
	Opening Position  // Position of opening delimiter "(" or "{"
	List    []*Field  // List of fields
	Closing Position  // Position of closing delimiter ")" or "}"
}

func (f *FieldList) Pos() Position {
	if f.Opening.IsValid() {
		return f.Opening
	}
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}
	return Position{}
}
func (f *FieldList) End() Position {
	if f.Closing.IsValid() {
		return f.Closing
	}
	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}
	return Position{}
}
func (f *FieldList) node() {}

// Field represents a field in a struct, interface, or function parameter/result list.
type Field struct {
	Names []*Ident // Field names (may be empty for anonymous fields or unnamed parameters)
	Type  Type     // Field type
	Tag   *BasicLit // Field tag (for struct fields only, may be nil)
}

func (f *Field) Pos() Position {
	if len(f.Names) > 0 {
		return f.Names[0].Pos()
	}
	if f.Type != nil {
		return f.Type.Pos()
	}
	return Position{}
}
func (f *Field) End() Position {
	if f.Tag != nil {
		return f.Tag.End()
	}
	if f.Type != nil {
		return f.Type.End()
	}
	if len(f.Names) > 0 {
		return f.Names[len(f.Names)-1].End()
	}
	return Position{}
}
func (f *Field) node() {}

// ParenType represents a parenthesized type: (T)
type ParenType struct {
	Lparen Position // Position of "("
	X      Type     // Type inside parentheses
	Rparen Position // Position of ")"
}

func (t *ParenType) Pos() Position { return t.Lparen }
func (t *ParenType) End() Position { return t.Rparen }
func (t *ParenType) node()         {}
func (t *ParenType) expr()         {}
func (t *ParenType) typeNode()     {}

// TypeAssertExpr represents a type assertion: x.(T)
type TypeAssertExpr struct {
	X      Expr     // Expression being asserted
	Lparen Position // Position of "("
	Type   Type     // Asserted type (nil for type switch x.(type))
	Rparen Position // Position of ")"
}

func (e *TypeAssertExpr) Pos() Position { return e.X.Pos() }
func (e *TypeAssertExpr) End() Position { return e.Rparen }
func (e *TypeAssertExpr) node()         {}
func (e *TypeAssertExpr) expr()         {}
