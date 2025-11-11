package antlr

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/mleku/moxie/pkg/ast"
)

// TokenToPosition converts an ANTLR token to an AST position.
func TokenToPosition(token antlr.Token, filename string) ast.Position {
	if token == nil {
		return ast.Position{}
	}
	return ast.Position{
		Filename: filename,
		Offset:   token.GetStart(),
		Line:     token.GetLine(),
		Column:   token.GetColumn() + 1, // ANTLR columns are 0-based, AST are 1-based
	}
}

// ContextToPosition returns the starting position of a parser context.
func ContextToPosition(ctx antlr.ParserRuleContext, filename string) ast.Position {
	if ctx == nil {
		return ast.Position{}
	}
	token := ctx.GetStart()
	return TokenToPosition(token, filename)
}

// ContextEndPosition returns the ending position of a parser context.
func ContextEndPosition(ctx antlr.ParserRuleContext, filename string) ast.Position {
	if ctx == nil {
		return ast.Position{}
	}
	token := ctx.GetStop()
	if token == nil {
		return ContextToPosition(ctx, filename)
	}

	pos := TokenToPosition(token, filename)
	// Add the token length to get the true end position
	text := token.GetText()
	pos.Column += len(text)
	pos.Offset += len(text)
	return pos
}
