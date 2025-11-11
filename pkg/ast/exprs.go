package ast

// ============================================================================
// Expression Nodes
// ============================================================================

// BadExpr represents an expression containing syntax errors.
type BadExpr struct {
	From Position // Start of bad expression
	To   Position // End of bad expression
}

func (e *BadExpr) Pos() Position { return e.From }
func (e *BadExpr) End() Position { return e.To }
func (e *BadExpr) node()         {}
func (e *BadExpr) expr()         {}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Lparen Position // Position of "("
	X      Expr     // Expression inside parentheses
	Rparen Position // Position of ")"
}

func (e *ParenExpr) Pos() Position { return e.Lparen }
func (e *ParenExpr) End() Position { return e.Rparen }
func (e *ParenExpr) node()         {}
func (e *ParenExpr) expr()         {}

// SelectorExpr represents a selector expression: x.Sel
type SelectorExpr struct {
	X   Expr   // Expression
	Sel *Ident // Selector
}

func (e *SelectorExpr) Pos() Position { return e.X.Pos() }
func (e *SelectorExpr) End() Position { return e.Sel.End() }
func (e *SelectorExpr) node()         {}
func (e *SelectorExpr) expr()         {}

// IndexExpr represents an index expression: x[i]
type IndexExpr struct {
	X      Expr     // Expression
	Lbrack Position // Position of "["
	Index  Expr     // Index expression
	Rbrack Position // Position of "]"
}

func (e *IndexExpr) Pos() Position { return e.X.Pos() }
func (e *IndexExpr) End() Position { return e.Rbrack }
func (e *IndexExpr) node()         {}
func (e *IndexExpr) expr()         {}

// SliceExpr represents a slice expression: x[low:high] or x[low:high:max]
type SliceExpr struct {
	X      Expr     // Expression
	Lbrack Position // Position of "["
	Low    Expr     // Low bound (may be nil)
	High   Expr     // High bound (may be nil)
	Max    Expr     // Maximum capacity (may be nil for 2-index slices)
	Slice3 bool     // true for 3-index slice (x[i:j:k])
	Rbrack Position // Position of "]"
}

func (e *SliceExpr) Pos() Position { return e.X.Pos() }
func (e *SliceExpr) End() Position { return e.Rbrack }
func (e *SliceExpr) node()         {}
func (e *SliceExpr) expr()         {}

// CallExpr represents a function call or type conversion.
type CallExpr struct {
	Fun      Expr     // Function or type
	Lparen   Position // Position of "("
	Args     []Expr   // Arguments
	Ellipsis Position // Position of "..." (invalid if not variadic)
	Rparen   Position // Position of ")"
}

func (e *CallExpr) Pos() Position { return e.Fun.Pos() }
func (e *CallExpr) End() Position { return e.Rparen }
func (e *CallExpr) node()         {}
func (e *CallExpr) expr()         {}

// StarExpr represents a pointer dereference or pointer type: *x
type StarExpr struct {
	Star Position // Position of "*"
	X    Expr     // Operand
}

func (e *StarExpr) Pos() Position { return e.Star }
func (e *StarExpr) End() Position { return e.X.End() }
func (e *StarExpr) node()         {}
func (e *StarExpr) expr()         {}

// UnaryExpr represents a unary expression.
type UnaryExpr struct {
	OpPos Position // Position of operator
	Op    Token    // Operator (ADD, SUB, NOT, XOR, MUL for pointer, AND for address-of)
	X     Expr     // Operand
}

func (e *UnaryExpr) Pos() Position { return e.OpPos }
func (e *UnaryExpr) End() Position { return e.X.End() }
func (e *UnaryExpr) node()         {}
func (e *UnaryExpr) expr()         {}

// BinaryExpr represents a binary expression.
type BinaryExpr struct {
	X     Expr     // Left operand
	OpPos Position // Position of operator
	Op    Token    // Operator
	Y     Expr     // Right operand
}

func (e *BinaryExpr) Pos() Position { return e.X.Pos() }
func (e *BinaryExpr) End() Position { return e.Y.End() }
func (e *BinaryExpr) node()         {}
func (e *BinaryExpr) expr()         {}

// KeyValueExpr represents a key-value pair in a composite literal.
type KeyValueExpr struct {
	Key   Expr     // Key
	Colon Position // Position of ":"
	Value Expr     // Value
}

func (e *KeyValueExpr) Pos() Position { return e.Key.Pos() }
func (e *KeyValueExpr) End() Position { return e.Value.End() }
func (e *KeyValueExpr) node()         {}
func (e *KeyValueExpr) expr()         {}

// CompositeLit represents a composite literal: T{...}
type CompositeLit struct {
	Type       Type     // Literal type (may be nil)
	Lbrace     Position // Position of "{"
	Elts       []Expr   // Elements (expressions or KeyValueExpr)
	Rbrace     Position // Position of "}"
	Incomplete bool     // true if "}}" is missing (for error recovery)
}

func (e *CompositeLit) Pos() Position {
	if e.Type != nil {
		return e.Type.Pos()
	}
	return e.Lbrace
}
func (e *CompositeLit) End() Position { return e.Rbrace }
func (e *CompositeLit) node()         {}
func (e *CompositeLit) expr()         {}

// FuncLit represents a function literal (anonymous function).
type FuncLit struct {
	Type *FuncType  // Function type
	Body *BlockStmt // Function body
}

func (e *FuncLit) Pos() Position { return e.Type.Pos() }
func (e *FuncLit) End() Position { return e.Body.End() }
func (e *FuncLit) node()         {}
func (e *FuncLit) expr()         {}

// Ellipsis represents the "..." in parameter lists or array types.
type Ellipsis struct {
	Ellipsis Position // Position of "..."
	Elt      Type     // Element type (may be nil for variadic parameters without type)
}

func (e *Ellipsis) Pos() Position { return e.Ellipsis }
func (e *Ellipsis) End() Position {
	if e.Elt != nil {
		return e.Elt.End()
	}
	return e.Ellipsis
}
func (e *Ellipsis) node() {}
func (e *Ellipsis) expr() {}

// IndexListExpr represents an index expression with multiple indices (for generics).
// Example: F[T1, T2, T3]
type IndexListExpr struct {
	X       Expr     // Expression
	Lbrack  Position // Position of "["
	Indices []Expr   // Index expressions
	Rbrack  Position // Position of "]"
}

func (e *IndexListExpr) Pos() Position { return e.X.Pos() }
func (e *IndexListExpr) End() Position { return e.Rbrack }
func (e *IndexListExpr) node()         {}
func (e *IndexListExpr) expr()         {}

// ============================================================================
// Moxie-specific Expression Nodes
// ============================================================================

// ChanLit represents a channel literal (Moxie syntax): &chan T{cap: 10}
type ChanLit struct {
	Ampersand Position     // Position of "&" (explicit pointer)
	Chan      Position     // Position of "chan" keyword
	Dir       ChanDir      // Channel direction
	Type      Type         // Element type
	Lbrace    Position     // Position of "{"
	Cap       Expr         // Capacity expression (in cap: expr)
	Rbrace    Position     // Position of "}"
}

func (e *ChanLit) Pos() Position { return e.Ampersand }
func (e *ChanLit) End() Position { return e.Rbrace }
func (e *ChanLit) node()         {}
func (e *ChanLit) expr()         {}

// SliceLit represents an explicit slice literal (Moxie syntax): &[]T{...}
type SliceLit struct {
	Ampersand Position     // Position of "&" (explicit pointer)
	Lbrack    Position     // Position of "["
	Type      Type         // Element type
	Lbrace    Position     // Position of "{"
	Elts      []Expr       // Elements
	Rbrace    Position     // Position of "}"
}

func (e *SliceLit) Pos() Position { return e.Ampersand }
func (e *SliceLit) End() Position { return e.Rbrace }
func (e *SliceLit) node()         {}
func (e *SliceLit) expr()         {}

// MapLit represents an explicit map literal (Moxie syntax): &map[K]V{...}
type MapLit struct {
	Ampersand Position     // Position of "&" (explicit pointer)
	Map       Position     // Position of "map" keyword
	Lbrack    Position     // Position of "["
	Key       Type         // Key type
	Value     Type         // Value type
	Lbrace    Position     // Position of "{"
	Elts      []Expr       // Elements (KeyValueExpr)
	Rbrace    Position     // Position of "}"
}

func (e *MapLit) Pos() Position { return e.Ampersand }
func (e *MapLit) End() Position { return e.Rbrace }
func (e *MapLit) node()         {}
func (e *MapLit) expr()         {}

// TypeCoercion represents a type coercion (Moxie FFI feature): (*[]uint32)(bytes)
type TypeCoercion struct {
	Lparen Position // Position of "("
	Target Type     // Target type
	Rparen Position // Position of ")"
	Expr   Expr     // Expression to coerce
}

func (e *TypeCoercion) Pos() Position { return e.Lparen }
func (e *TypeCoercion) End() Position { return e.Expr.End() }
func (e *TypeCoercion) node()         {}
func (e *TypeCoercion) expr()         {}

// FFICall represents an FFI call using dlsym (Moxie feature).
// Example: dlsym[func(*byte) int64](lib, "strlen")
type FFICall struct {
	Name   *Ident     // Function name (dlsym, dlopen, dlclose, etc.)
	Lbrack Position   // Position of "[" (type parameter start)
	Type   Type       // Function type
	Rbrack Position   // Position of "]" (type parameter end)
	Args   []Expr     // Arguments to dlsym
}

func (e *FFICall) Pos() Position { return e.Name.Pos() }
func (e *FFICall) End() Position {
	if len(e.Args) > 0 {
		return e.Args[len(e.Args)-1].End()
	}
	return e.Rbrack
}
func (e *FFICall) node() {}
func (e *FFICall) expr() {}
