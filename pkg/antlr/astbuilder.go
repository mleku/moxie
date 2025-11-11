package antlr

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/mleku/moxie/pkg/ast"
)

// ASTBuilder transforms ANTLR parse trees into Moxie AST nodes.
// It embeds BaseMoxieVisitor to implement the MoxieVisitor interface.
type ASTBuilder struct {
	BaseMoxieVisitor
	filename string
	errors   []error
}

// NewASTBuilder creates a new AST builder for the given filename.
func NewASTBuilder(filename string) *ASTBuilder {
	return &ASTBuilder{
		filename: filename,
		errors:   []error{},
	}
}

// Errors returns the list of errors encountered during AST building.
func (b *ASTBuilder) Errors() []error {
	return b.errors
}

// addError adds an error to the error list.
func (b *ASTBuilder) addError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// pos returns the starting position of a context.
func (b *ASTBuilder) pos(ctx antlr.ParserRuleContext) ast.Position {
	return ContextToPosition(ctx, b.filename)
}

// endPos returns the ending position of a context.
func (b *ASTBuilder) endPos(ctx antlr.ParserRuleContext) ast.Position {
	return ContextEndPosition(ctx, b.filename)
}

// tokenPos returns the position of a token.
func (b *ASTBuilder) tokenPos(token antlr.Token) ast.Position {
	return TokenToPosition(token, b.filename)
}

// ============================================================================
// Top-level: Source File
// ============================================================================

// VisitSourceFile transforms the top-level source file.
func (b *ASTBuilder) VisitSourceFile(ctx *SourceFileContext) interface{} {
	file := &ast.File{
		StartPos: b.pos(ctx),
		EndPos:   b.endPos(ctx),
	}

	// Package clause
	if pkgCtx := ctx.PackageClause(); pkgCtx != nil {
		if pkg, ok := pkgCtx.(*PackageClauseContext); ok {
			file.Package = b.VisitPackageClause(pkg).(*ast.PackageClause)
		}
	}

	// Imports
	for _, importCtx := range ctx.AllImportDecl() {
		if iCtx, ok := importCtx.(*ImportDeclContext); ok {
			if decl := b.VisitImportDecl(iCtx); decl != nil {
				file.Imports = append(file.Imports, decl.(*ast.ImportDecl))
			}
		}
	}

	// Top-level declarations
	for _, declCtx := range ctx.AllTopLevelDecl() {
		if dCtx, ok := declCtx.(*TopLevelDeclContext); ok {
			if decl := b.VisitTopLevelDecl(dCtx); decl != nil {
				file.Decls = append(file.Decls, decl.(ast.Decl))
			}
		}
	}

	return file
}

// VisitPackageClause transforms a package clause.
func (b *ASTBuilder) VisitPackageClause(ctx *PackageClauseContext) interface{} {
	if ctx == nil {
		return nil
	}

	pkg := &ast.PackageClause{
		Package: b.tokenPos(ctx.PACKAGE().GetSymbol()),
	}

	if ident := ctx.IDENTIFIER(); ident != nil {
		pkg.Name = &ast.Ident{
			NamePos: b.tokenPos(ident.GetSymbol()),
			Name:    ident.GetText(),
		}
	}

	return pkg
}

// VisitTopLevelDecl transforms a top-level declaration.
func (b *ASTBuilder) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Function or method declaration
	if funcCtx := ctx.FunctionDecl(); funcCtx != nil {
		if fCtx, ok := funcCtx.(*FunctionDeclContext); ok {
			return b.VisitFunctionDecl(fCtx)
		}
	}

	if methCtx := ctx.MethodDecl(); methCtx != nil {
		if mCtx, ok := methCtx.(*MethodDeclContext); ok {
			return b.VisitMethodDecl(mCtx)
		}
	}

	// Other declarations
	if declCtx := ctx.Declaration(); declCtx != nil {
		if dCtx, ok := declCtx.(*DeclarationContext); ok {
			return b.VisitDeclaration(dCtx)
		}
	}

	return nil
}

// VisitDeclaration transforms a declaration (const, type, var).
func (b *ASTBuilder) VisitDeclaration(ctx *DeclarationContext) interface{} {
	if ctx == nil {
		return nil
	}

	if constCtx := ctx.ConstDecl(); constCtx != nil {
		if cCtx, ok := constCtx.(*ConstDeclContext); ok {
			return b.VisitConstDecl(cCtx)
		}
	}

	if typeCtx := ctx.TypeDecl(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*TypeDeclContext); ok {
			return b.VisitTypeDecl(tCtx)
		}
	}

	if varCtx := ctx.VarDecl(); varCtx != nil {
		if vCtx, ok := varCtx.(*VarDeclContext); ok {
			return b.VisitVarDecl(vCtx)
		}
	}

	return nil
}

// ============================================================================
// Imports
// ============================================================================

// VisitImportDecl transforms an import declaration.
func (b *ASTBuilder) VisitImportDecl(ctx *ImportDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	decl := &ast.ImportDecl{
		Import: b.tokenPos(ctx.IMPORT().GetSymbol()),
	}

	// Get all import specs
	for _, specCtx := range ctx.AllImportSpec() {
		if sCtx, ok := specCtx.(*ImportSpecContext); ok {
			if spec := b.VisitImportSpec(sCtx); spec != nil {
				decl.Specs = append(decl.Specs, spec.(*ast.ImportSpec))
			}
		}
	}

	return decl
}

// VisitImportSpec transforms an import specification.
func (b *ASTBuilder) VisitImportSpec(ctx *ImportSpecContext) interface{} {
	if ctx == nil {
		return nil
	}

	spec := &ast.ImportSpec{}

	// Import alias (., _, or identifier)
	if ident := ctx.IDENTIFIER(); ident != nil {
		spec.Name = &ast.Ident{
			NamePos: b.tokenPos(ident.GetSymbol()),
			Name:    ident.GetText(),
		}
	}

	// Import path (string literal)
	if str := ctx.String_(); str != nil {
		if sCtx, ok := str.(*String_Context); ok {
			if lit := b.VisitString_(sCtx); lit != nil {
				spec.Path = lit.(*ast.BasicLit)
			}
		}
	}

	return spec
}

// ============================================================================
// Helper Methods
// ============================================================================

// visitIdentifier creates an identifier from a token.
func (b *ASTBuilder) visitIdentifier(token antlr.TerminalNode) *ast.Ident {
	if token == nil {
		return nil
	}
	return &ast.Ident{
		NamePos: b.tokenPos(token.GetSymbol()),
		Name:    token.GetText(),
	}
}

// visitIdentifierList creates a list of identifiers.
func (b *ASTBuilder) visitIdentifierList(ctx IIdentifierListContext) []*ast.Ident {
	if ctx == nil {
		return nil
	}

	var idents []*ast.Ident
	for _, idToken := range ctx.AllIDENTIFIER() {
		idents = append(idents, b.visitIdentifier(idToken))
	}
	return idents
}

// BuildAST is the main entry point for building an AST from a parse tree.
func BuildAST(tree *SourceFileContext, filename string) (*ast.File, []error) {
	builder := NewASTBuilder(filename)
	file := builder.VisitSourceFile(tree).(*ast.File)
	return file, builder.Errors()
}
