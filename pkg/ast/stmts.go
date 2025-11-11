package ast

// ============================================================================
// Statement Nodes
// ============================================================================

// BadStmt represents a statement containing syntax errors.
type BadStmt struct {
	From Position // Start of bad statement
	To   Position // End of bad statement
}

func (s *BadStmt) Pos() Position { return s.From }
func (s *BadStmt) End() Position { return s.To }
func (s *BadStmt) node()         {}
func (s *BadStmt) stmt()         {}

// DeclStmt represents a declaration in a statement list.
type DeclStmt struct {
	Decl Decl // Declaration (const, var, or type)
}

func (s *DeclStmt) Pos() Position { return s.Decl.Pos() }
func (s *DeclStmt) End() Position { return s.Decl.End() }
func (s *DeclStmt) node()         {}
func (s *DeclStmt) stmt()         {}

// EmptyStmt represents an empty statement.
type EmptyStmt struct {
	Semicolon Position // Position of ";"
	Implicit  bool     // true if semicolon was implicit (inserted by scanner)
}

func (s *EmptyStmt) Pos() Position { return s.Semicolon }
func (s *EmptyStmt) End() Position { return s.Semicolon }
func (s *EmptyStmt) node()         {}
func (s *EmptyStmt) stmt()         {}

// LabeledStmt represents a labeled statement.
type LabeledStmt struct {
	Label *Ident   // Label
	Colon Position // Position of ":"
	Stmt  Stmt     // Statement following the label
}

func (s *LabeledStmt) Pos() Position { return s.Label.Pos() }
func (s *LabeledStmt) End() Position { return s.Stmt.End() }
func (s *LabeledStmt) node()         {}
func (s *LabeledStmt) stmt()         {}

// ExprStmt represents an expression used as a statement.
type ExprStmt struct {
	X Expr // Expression
}

func (s *ExprStmt) Pos() Position { return s.X.Pos() }
func (s *ExprStmt) End() Position { return s.X.End() }
func (s *ExprStmt) node()         {}
func (s *ExprStmt) stmt()         {}

// SendStmt represents a send statement: ch <- x
type SendStmt struct {
	Chan  Expr     // Channel expression
	Arrow Position // Position of "<-"
	Value Expr     // Value to send
}

func (s *SendStmt) Pos() Position { return s.Chan.Pos() }
func (s *SendStmt) End() Position { return s.Value.End() }
func (s *SendStmt) node()         {}
func (s *SendStmt) stmt()         {}

// IncDecStmt represents an increment or decrement statement: x++ or x--
type IncDecStmt struct {
	X      Expr     // Expression
	TokPos Position // Position of "++" or "--"
	Tok    Token    // INC or DEC
}

func (s *IncDecStmt) Pos() Position { return s.X.Pos() }
func (s *IncDecStmt) End() Position { return s.TokPos }
func (s *IncDecStmt) node()         {}
func (s *IncDecStmt) stmt()         {}

// AssignStmt represents an assignment or short variable declaration.
type AssignStmt struct {
	Lhs    []Expr   // Left-hand side
	TokPos Position // Position of assignment token
	Tok    Token    // Assignment token (ASSIGN, DEFINE, ADD_ASSIGN, etc.)
	Rhs    []Expr   // Right-hand side
}

func (s *AssignStmt) Pos() Position { return s.Lhs[0].Pos() }
func (s *AssignStmt) End() Position { return s.Rhs[len(s.Rhs)-1].End() }
func (s *AssignStmt) node()         {}
func (s *AssignStmt) stmt()         {}

// GoStmt represents a go statement: go f(x)
type GoStmt struct {
	Go   Position  // Position of "go" keyword
	Call *CallExpr // Function call
}

func (s *GoStmt) Pos() Position { return s.Go }
func (s *GoStmt) End() Position { return s.Call.End() }
func (s *GoStmt) node()         {}
func (s *GoStmt) stmt()         {}

// DeferStmt represents a defer statement: defer f(x)
type DeferStmt struct {
	Defer Position  // Position of "defer" keyword
	Call  *CallExpr // Function call
}

func (s *DeferStmt) Pos() Position { return s.Defer }
func (s *DeferStmt) End() Position { return s.Call.End() }
func (s *DeferStmt) node()         {}
func (s *DeferStmt) stmt()         {}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Return  Position // Position of "return" keyword
	Results []Expr   // Result expressions (may be nil)
}

func (s *ReturnStmt) Pos() Position { return s.Return }
func (s *ReturnStmt) End() Position {
	if len(s.Results) > 0 {
		return s.Results[len(s.Results)-1].End()
	}
	return s.Return
}
func (s *ReturnStmt) node() {}
func (s *ReturnStmt) stmt() {}

// BranchStmt represents a break, continue, goto, or fallthrough statement.
type BranchStmt struct {
	TokPos Position // Position of branch keyword
	Tok    Token    // BREAK, CONTINUE, GOTO, or FALLTHROUGH
	Label  *Ident   // Label (may be nil except for goto)
}

func (s *BranchStmt) Pos() Position { return s.TokPos }
func (s *BranchStmt) End() Position {
	if s.Label != nil {
		return s.Label.End()
	}
	return s.TokPos
}
func (s *BranchStmt) node() {}
func (s *BranchStmt) stmt() {}

// BlockStmt represents a block statement.
type BlockStmt struct {
	Lbrace Position // Position of "{"
	List   []Stmt   // Statements in the block
	Rbrace Position // Position of "}"
}

func (s *BlockStmt) Pos() Position { return s.Lbrace }
func (s *BlockStmt) End() Position { return s.Rbrace }
func (s *BlockStmt) node()         {}
func (s *BlockStmt) stmt()         {}

// IfStmt represents an if statement.
type IfStmt struct {
	If   Position  // Position of "if" keyword
	Init Stmt      // Initialization statement (may be nil)
	Cond Expr      // Condition
	Body *BlockStmt // Body
	Else Stmt      // Else branch (IfStmt or BlockStmt, may be nil)
}

func (s *IfStmt) Pos() Position { return s.If }
func (s *IfStmt) End() Position {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}
func (s *IfStmt) node() {}
func (s *IfStmt) stmt() {}

// CaseClause represents a case or default clause in a switch or select statement.
type CaseClause struct {
	Case  Position  // Position of "case" or "default" keyword
	List  []Expr    // List of expressions (nil for default case)
	Colon Position  // Position of ":"
	Body  []Stmt    // Statements in the case
}

func (s *CaseClause) Pos() Position { return s.Case }
func (s *CaseClause) End() Position {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon
}
func (s *CaseClause) node() {}
func (s *CaseClause) stmt() {}

// SwitchStmt represents an expression switch statement.
type SwitchStmt struct {
	Switch Position      // Position of "switch" keyword
	Init   Stmt          // Initialization statement (may be nil)
	Tag    Expr          // Tag expression (may be nil)
	Body   *BlockStmt    // Body (contains case clauses)
}

func (s *SwitchStmt) Pos() Position { return s.Switch }
func (s *SwitchStmt) End() Position { return s.Body.End() }
func (s *SwitchStmt) node()         {}
func (s *SwitchStmt) stmt()         {}

// TypeSwitchStmt represents a type switch statement.
type TypeSwitchStmt struct {
	Switch Position   // Position of "switch" keyword
	Init   Stmt       // Initialization statement (may be nil)
	Assign Stmt       // Type assertion (x := y.(type))
	Body   *BlockStmt // Body (contains case clauses)
}

func (s *TypeSwitchStmt) Pos() Position { return s.Switch }
func (s *TypeSwitchStmt) End() Position { return s.Body.End() }
func (s *TypeSwitchStmt) node()         {}
func (s *TypeSwitchStmt) stmt()         {}

// CommClause represents a case clause in a select statement.
type CommClause struct {
	Case  Position  // Position of "case" or "default" keyword
	Comm  Stmt      // Send or receive statement (nil for default)
	Colon Position  // Position of ":"
	Body  []Stmt    // Statements in the case
}

func (s *CommClause) Pos() Position { return s.Case }
func (s *CommClause) End() Position {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon
}
func (s *CommClause) node() {}
func (s *CommClause) stmt() {}

// SelectStmt represents a select statement.
type SelectStmt struct {
	Select Position   // Position of "select" keyword
	Body   *BlockStmt // Body (contains comm clauses)
}

func (s *SelectStmt) Pos() Position { return s.Select }
func (s *SelectStmt) End() Position { return s.Body.End() }
func (s *SelectStmt) node()         {}
func (s *SelectStmt) stmt()         {}

// ForStmt represents a for loop.
type ForStmt struct {
	For  Position   // Position of "for" keyword
	Init Stmt       // Initialization statement (may be nil)
	Cond Expr       // Condition (may be nil for infinite loop)
	Post Stmt       // Post iteration statement (may be nil)
	Body *BlockStmt // Body
}

func (s *ForStmt) Pos() Position { return s.For }
func (s *ForStmt) End() Position { return s.Body.End() }
func (s *ForStmt) node()         {}
func (s *ForStmt) stmt()         {}

// RangeStmt represents a for...range statement.
type RangeStmt struct {
	For        Position   // Position of "for" keyword
	Key        Expr       // Key variable (or index for slices/arrays)
	Value      Expr       // Value variable (may be nil)
	TokPos     Position   // Position of assignment token
	Tok        Token      // ASSIGN or DEFINE
	X          Expr       // Value to range over
	Body       *BlockStmt // Body
}

func (s *RangeStmt) Pos() Position { return s.For }
func (s *RangeStmt) End() Position { return s.Body.End() }
func (s *RangeStmt) node()         {}
func (s *RangeStmt) stmt()         {}
