package ast

// ============================================================================
// File and Package Structure
// ============================================================================

// File represents a Moxie source file.
type File struct {
	Package    *PackageClause  // Package clause
	Imports    []*ImportDecl   // Import declarations
	Decls      []Decl          // Top-level declarations (const, var, type, func)
	Comments   []*CommentGroup // Comments in the file
	StartPos   Position        // Start of file
	EndPos     Position        // End of file
}

func (f *File) Pos() Position { return f.StartPos }
func (f *File) End() Position { return f.EndPos }
func (f *File) node()         {}

// PackageClause represents the package declaration.
type PackageClause struct {
	Package Position // Position of "package" keyword
	Name    *Ident   // Package name
}

func (p *PackageClause) Pos() Position { return p.Package }
func (p *PackageClause) End() Position { return p.Name.End() }
func (p *PackageClause) node()         {}

// Comment represents a single comment (// or /* */).
type Comment struct {
	Slash Position // Position of "/" starting the comment
	Text  string   // Comment text (includes // or /* */)
}

func (c *Comment) Pos() Position { return c.Slash }
func (c *Comment) End() Position { return Position{Line: c.Slash.Line, Column: c.Slash.Column + len(c.Text)} }
func (c *Comment) node()         {}

// CommentGroup represents a sequence of comments with no blank lines between them.
type CommentGroup struct {
	List []*Comment
}

func (g *CommentGroup) Pos() Position {
	if len(g.List) > 0 {
		return g.List[0].Pos()
	}
	return Position{}
}
func (g *CommentGroup) End() Position {
	if n := len(g.List); n > 0 {
		return g.List[n-1].End()
	}
	return Position{}
}
func (g *CommentGroup) node() {}

// Text returns the text of the comment group.
func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	var text string
	for _, c := range g.List {
		text += c.Text + "\n"
	}
	return text
}

// ============================================================================
// Import Declarations
// ============================================================================

// ImportDecl represents an import declaration.
type ImportDecl struct {
	Import Position     // Position of "import" keyword
	Lparen Position     // Position of "(" (invalid if not grouped)
	Specs  []*ImportSpec // Import specs
	Rparen Position     // Position of ")" (invalid if not grouped)
}

func (d *ImportDecl) Pos() Position { return d.Import }
func (d *ImportDecl) End() Position {
	if d.Rparen.IsValid() {
		return d.Rparen
	}
	if len(d.Specs) > 0 {
		return d.Specs[len(d.Specs)-1].End()
	}
	return d.Import
}
func (d *ImportDecl) node() {}
func (d *ImportDecl) decl() {}

// ImportSpec represents a single import specification.
type ImportSpec struct {
	Name   *Ident    // Local name (may be nil for default import, "." for dot import, "_" for side-effect)
	Path   *BasicLit // Import path (string literal)
}

func (s *ImportSpec) Pos() Position {
	if s.Name != nil {
		return s.Name.Pos()
	}
	return s.Path.Pos()
}
func (s *ImportSpec) End() Position { return s.Path.End() }
func (s *ImportSpec) node()         {}
func (s *ImportSpec) spec()         {}

// ============================================================================
// Constant Declarations
// ============================================================================

// ConstDecl represents a const declaration.
type ConstDecl struct {
	Const  Position    // Position of "const" keyword
	Lparen Position    // Position of "(" (invalid if not grouped)
	Specs  []*ConstSpec // Const specs
	Rparen Position    // Position of ")" (invalid if not grouped)
}

func (d *ConstDecl) Pos() Position { return d.Const }
func (d *ConstDecl) End() Position {
	if d.Rparen.IsValid() {
		return d.Rparen
	}
	if len(d.Specs) > 0 {
		return d.Specs[len(d.Specs)-1].End()
	}
	return d.Const
}
func (d *ConstDecl) node() {}
func (d *ConstDecl) decl() {}

// ConstSpec represents a const specification.
type ConstSpec struct {
	Names  []*Ident  // Constant names
	Type   Type      // Type (may be nil)
	Values []Expr    // Values (initializers)
}

func (s *ConstSpec) Pos() Position {
	if len(s.Names) > 0 {
		return s.Names[0].Pos()
	}
	return Position{}
}
func (s *ConstSpec) End() Position {
	if len(s.Values) > 0 {
		return s.Values[len(s.Values)-1].End()
	}
	if s.Type != nil {
		return s.Type.End()
	}
	if len(s.Names) > 0 {
		return s.Names[len(s.Names)-1].End()
	}
	return Position{}
}
func (s *ConstSpec) node() {}
func (s *ConstSpec) spec() {}

// ============================================================================
// Variable Declarations
// ============================================================================

// VarDecl represents a var declaration.
type VarDecl struct {
	Var    Position   // Position of "var" keyword
	Lparen Position   // Position of "(" (invalid if not grouped)
	Specs  []*VarSpec // Var specs
	Rparen Position   // Position of ")" (invalid if not grouped)
}

func (d *VarDecl) Pos() Position { return d.Var }
func (d *VarDecl) End() Position {
	if d.Rparen.IsValid() {
		return d.Rparen
	}
	if len(d.Specs) > 0 {
		return d.Specs[len(d.Specs)-1].End()
	}
	return d.Var
}
func (d *VarDecl) node() {}
func (d *VarDecl) decl() {}

// VarSpec represents a var specification.
type VarSpec struct {
	Names  []*Ident  // Variable names
	Type   Type      // Type (may be nil if values are present)
	Values []Expr    // Values (initializers, may be nil)
}

func (s *VarSpec) Pos() Position {
	if len(s.Names) > 0 {
		return s.Names[0].Pos()
	}
	return Position{}
}
func (s *VarSpec) End() Position {
	if len(s.Values) > 0 {
		return s.Values[len(s.Values)-1].End()
	}
	if s.Type != nil {
		return s.Type.End()
	}
	if len(s.Names) > 0 {
		return s.Names[len(s.Names)-1].End()
	}
	return Position{}
}
func (s *VarSpec) node() {}
func (s *VarSpec) spec() {}

// ============================================================================
// Type Declarations
// ============================================================================

// TypeDecl represents a type declaration.
type TypeDecl struct {
	Type   Position    // Position of "type" keyword
	Lparen Position    // Position of "(" (invalid if not grouped)
	Specs  []*TypeSpec // Type specs
	Rparen Position    // Position of ")" (invalid if not grouped)
}

func (d *TypeDecl) Pos() Position { return d.Type }
func (d *TypeDecl) End() Position {
	if d.Rparen.IsValid() {
		return d.Rparen
	}
	if len(d.Specs) > 0 {
		return d.Specs[len(d.Specs)-1].End()
	}
	return d.Type
}
func (d *TypeDecl) node() {}
func (d *TypeDecl) decl() {}

// TypeSpec represents a type specification (type definition or alias).
type TypeSpec struct {
	Name       *Ident     // Type name
	TypeParams *FieldList // Type parameters (generics), may be nil
	Assign     Position   // Position of "=" (invalid if not an alias)
	Type       Type       // Underlying type
}

func (s *TypeSpec) Pos() Position { return s.Name.Pos() }
func (s *TypeSpec) End() Position { return s.Type.End() }
func (s *TypeSpec) node()         {}
func (s *TypeSpec) spec()         {}

// IsAlias returns true if this is a type alias (type A = B).
func (s *TypeSpec) IsAlias() bool {
	return s.Assign.IsValid()
}

// ============================================================================
// Function Declarations
// ============================================================================

// FuncDecl represents a function declaration.
type FuncDecl struct {
	Recv *FieldList // Receiver (for methods), may be nil
	Name *Ident     // Function name
	Type *FuncType  // Function signature
	Body *BlockStmt // Function body (may be nil for external functions)
}

func (d *FuncDecl) Pos() Position {
	if d.Recv != nil {
		return d.Recv.Pos()
	}
	return d.Type.Pos()
}
func (d *FuncDecl) End() Position {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}
func (d *FuncDecl) node() {}
func (d *FuncDecl) decl() {}

// IsMethod returns true if this is a method (has a receiver).
func (d *FuncDecl) IsMethod() bool {
	return d.Recv != nil
}
