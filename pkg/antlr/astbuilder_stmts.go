package antlr

import (
	"github.com/mleku/moxie/pkg/ast"
)

// ============================================================================
// Statements
// ============================================================================

// VisitBlock transforms a block statement.
func (b *ASTBuilder) VisitBlock(ctx *BlockContext) interface{} {
	if ctx == nil {
		return nil
	}

	block := &ast.BlockStmt{
		Lbrace: b.pos(ctx),
		Rbrace: b.endPos(ctx),
	}

	// Statement list
	if stmtListCtx := ctx.StatementList(); stmtListCtx != nil {
		if stmts := b.VisitStatementList(stmtListCtx); stmts != nil {
			block.List = stmts.([]ast.Stmt)
		}
	}

	return block
}

// VisitStatementList transforms a statement list.
func (b *ASTBuilder) VisitStatementList(ctx *StatementListContext) interface{} {
	if ctx == nil {
		return nil
	}

	var stmts []ast.Stmt
	for _, stmtCtx := range ctx.AllStatement() {
		if stmt := b.VisitStatement(stmtCtx); stmt != nil {
			stmts = append(stmts, stmt.(ast.Stmt))
		}
	}

	return stmts
}

// VisitStatement transforms a statement.
func (b *ASTBuilder) VisitStatement(ctx *StatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Declaration statement
	if declCtx := ctx.Declaration(); declCtx != nil {
		if decl := b.VisitDeclaration(declCtx); decl != nil {
			return &ast.DeclStmt{Decl: decl.(ast.Decl)}
		}
	}

	// Simple statement
	if simpleCtx := ctx.SimpleStmt(); simpleCtx != nil {
		return b.VisitSimpleStmt(simpleCtx)
	}

	// Return statement
	if retCtx := ctx.ReturnStmt(); retCtx != nil {
		return b.VisitReturnStmt(retCtx)
	}

	// Break statement
	if breakCtx := ctx.BreakStmt(); breakCtx != nil {
		return b.VisitBreakStmt(breakCtx)
	}

	// Continue statement
	if contCtx := ctx.ContinueStmt(); contCtx != nil {
		return b.VisitContinueStmt(contCtx)
	}

	// Goto statement
	if gotoCtx := ctx.GotoStmt(); gotoCtx != nil {
		return b.VisitGotoStmt(gotoCtx)
	}

	// Fallthrough statement
	if fallthroughCtx := ctx.FallthroughStmt(); fallthroughCtx != nil {
		return b.VisitFallthroughStmt(fallthroughCtx)
	}

	// Block statement
	if blockCtx := ctx.Block(); blockCtx != nil {
		return b.VisitBlock(blockCtx)
	}

	// If statement
	if ifCtx := ctx.IfStmt(); ifCtx != nil {
		return b.VisitIfStmt(ifCtx)
	}

	// Switch statement
	if switchCtx := ctx.SwitchStmt(); switchCtx != nil {
		return b.VisitSwitchStmt(switchCtx)
	}

	// Select statement
	if selectCtx := ctx.SelectStmt(); selectCtx != nil {
		return b.VisitSelectStmt(selectCtx)
	}

	// For statement
	if forCtx := ctx.ForStmt(); forCtx != nil {
		return b.VisitForStmt(forCtx)
	}

	// Defer statement
	if deferCtx := ctx.DeferStmt(); deferCtx != nil {
		return b.VisitDeferStmt(deferCtx)
	}

	// Go statement
	if goCtx := ctx.GoStmt(); goCtx != nil {
		return b.VisitGoStmt(goCtx)
	}

	// Labeled statement
	if labeledCtx := ctx.LabeledStmt(); labeledCtx != nil {
		return b.VisitLabeledStmt(labeledCtx)
	}

	return &ast.EmptyStmt{Semicolon: b.pos(ctx)}
}

// VisitSimpleStmt transforms a simple statement.
func (b *ASTBuilder) VisitSimpleStmt(ctx *SimpleStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Expression statement
	if exprCtx := ctx.ExpressionStmt(); exprCtx != nil {
		return b.VisitExpressionStmt(exprCtx)
	}

	// Send statement
	if sendCtx := ctx.SendStmt(); sendCtx != nil {
		return b.VisitSendStmt(sendCtx)
	}

	// Inc/Dec statement
	if incDecCtx := ctx.IncDecStmt(); incDecCtx != nil {
		return b.VisitIncDecStmt(incDecCtx)
	}

	// Assignment
	if assignCtx := ctx.Assignment(); assignCtx != nil {
		return b.VisitAssignment(assignCtx)
	}

	// Short var declaration
	if shortVarCtx := ctx.ShortVarDecl(); shortVarCtx != nil {
		return b.VisitShortVarDecl(shortVarCtx)
	}

	return nil
}

// VisitExpressionStmt transforms an expression statement.
func (b *ASTBuilder) VisitExpressionStmt(ctx *ExpressionStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			return &ast.ExprStmt{X: expr.(ast.Expr)}
		}
	}

	return nil
}

// VisitSendStmt transforms a send statement.
func (b *ASTBuilder) VisitSendStmt(ctx *SendStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	send := &ast.SendStmt{
		Arrow: b.pos(ctx),
	}

	// Channel expression
	exprs := ctx.AllExpression()
	if len(exprs) >= 2 {
		if ch := b.VisitExpression(exprs[0]); ch != nil {
			send.Chan = ch.(ast.Expr)
		}
		if val := b.VisitExpression(exprs[1]); val != nil {
			send.Value = val.(ast.Expr)
		}
	}

	return send
}

// VisitIncDecStmt transforms an increment/decrement statement.
func (b *ASTBuilder) VisitIncDecStmt(ctx *IncDecStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	incDec := &ast.IncDecStmt{
		TokPos: b.pos(ctx),
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			incDec.X = expr.(ast.Expr)
		}
	}

	// Determine if ++ or --
	text := ctx.GetText()
	if len(text) >= 2 {
		if text[len(text)-2:] == "++" {
			incDec.Tok = ast.INC
		} else {
			incDec.Tok = ast.DEC
		}
	}

	return incDec
}

// VisitAssignment transforms an assignment statement.
func (b *ASTBuilder) VisitAssignment(ctx *AssignmentContext) interface{} {
	if ctx == nil {
		return nil
	}

	assign := &ast.AssignStmt{
		TokPos: b.pos(ctx),
		Tok:    ast.ASSIGN,
	}

	// Left-hand side
	if lhsCtx := ctx.ExpressionList(0); lhsCtx != nil {
		if lhs := b.VisitExpressionList(lhsCtx); lhs != nil {
			assign.Lhs = lhs.([]ast.Expr)
		}
	}

	// Assignment operator
	if opCtx := ctx.Assign_op(); opCtx != nil {
		if op := b.VisitAssign_op(opCtx); op != nil {
			assign.Tok = op.(ast.Token)
		}
	}

	// Right-hand side
	if rhsCtx := ctx.ExpressionList(1); rhsCtx != nil {
		if rhs := b.VisitExpressionList(rhsCtx); rhs != nil {
			assign.Rhs = rhs.([]ast.Expr)
		}
	}

	return assign
}

// VisitAssign_op transforms an assignment operator.
func (b *ASTBuilder) VisitAssign_op(ctx *Assign_opContext) interface{} {
	if ctx == nil {
		return ast.ASSIGN
	}

	text := ctx.GetText()
	switch text {
	case "=":
		return ast.ASSIGN
	case "+=":
		return ast.ADD_ASSIGN
	case "-=":
		return ast.SUB_ASSIGN
	case "*=":
		return ast.MUL_ASSIGN
	case "/=":
		return ast.QUO_ASSIGN
	case "%=":
		return ast.REM_ASSIGN
	case "&=":
		return ast.AND_ASSIGN
	case "|=":
		return ast.OR_ASSIGN
	case "^=":
		return ast.XOR_ASSIGN
	case "<<=":
		return ast.SHL_ASSIGN
	case ">>=":
		return ast.SHR_ASSIGN
	case "&^=":
		return ast.AND_NOT_ASSIGN
	default:
		return ast.ASSIGN
	}
}

// VisitShortVarDecl transforms a short variable declaration.
func (b *ASTBuilder) VisitShortVarDecl(ctx *ShortVarDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	assign := &ast.AssignStmt{
		TokPos: b.pos(ctx),
		Tok:    ast.DEFINE,
	}

	// Left-hand side (identifiers)
	if idListCtx := ctx.IdentifierList(); idListCtx != nil {
		idents := b.visitIdentifierList(idListCtx)
		for _, id := range idents {
			assign.Lhs = append(assign.Lhs, id)
		}
	}

	// Right-hand side (expressions)
	if exprListCtx := ctx.ExpressionList(); exprListCtx != nil {
		if exprs := b.VisitExpressionList(exprListCtx); exprs != nil {
			assign.Rhs = exprs.([]ast.Expr)
		}
	}

	return assign
}

// VisitReturnStmt transforms a return statement.
func (b *ASTBuilder) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	ret := &ast.ReturnStmt{
		Return: b.tokenPos(ctx.RETURN().GetSymbol()),
	}

	// Return values
	if exprListCtx := ctx.ExpressionList(); exprListCtx != nil {
		if exprs := b.VisitExpressionList(exprListCtx); exprs != nil {
			ret.Results = exprs.([]ast.Expr)
		}
	}

	return ret
}

// VisitBreakStmt transforms a break statement.
func (b *ASTBuilder) VisitBreakStmt(ctx *BreakStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	branch := &ast.BranchStmt{
		TokPos: b.tokenPos(ctx.BREAK().GetSymbol()),
		Tok:    ast.BREAK,
	}

	// Label (optional)
	if ident := ctx.IDENTIFIER(); ident != nil {
		branch.Label = b.visitIdentifier(ident)
	}

	return branch
}

// VisitContinueStmt transforms a continue statement.
func (b *ASTBuilder) VisitContinueStmt(ctx *ContinueStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	branch := &ast.BranchStmt{
		TokPos: b.tokenPos(ctx.CONTINUE().GetSymbol()),
		Tok:    ast.CONTINUE,
	}

	// Label (optional)
	if ident := ctx.IDENTIFIER(); ident != nil {
		branch.Label = b.visitIdentifier(ident)
	}

	return branch
}

// VisitGotoStmt transforms a goto statement.
func (b *ASTBuilder) VisitGotoStmt(ctx *GotoStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	branch := &ast.BranchStmt{
		TokPos: b.tokenPos(ctx.GOTO().GetSymbol()),
		Tok:    ast.GOTO,
	}

	// Label (required)
	if ident := ctx.IDENTIFIER(); ident != nil {
		branch.Label = b.visitIdentifier(ident)
	}

	return branch
}

// VisitFallthroughStmt transforms a fallthrough statement.
func (b *ASTBuilder) VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	return &ast.BranchStmt{
		TokPos: b.tokenPos(ctx.FALLTHROUGH().GetSymbol()),
		Tok:    ast.FALLTHROUGH,
	}
}

// VisitDeferStmt transforms a defer statement.
func (b *ASTBuilder) VisitDeferStmt(ctx *DeferStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	deferStmt := &ast.DeferStmt{
		Defer: b.tokenPos(ctx.DEFER().GetSymbol()),
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			if call, ok := expr.(*ast.CallExpr); ok {
				deferStmt.Call = call
			}
		}
	}

	return deferStmt
}

// VisitGoStmt transforms a go statement.
func (b *ASTBuilder) VisitGoStmt(ctx *GoStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	goStmt := &ast.GoStmt{
		Go: b.tokenPos(ctx.GO().GetSymbol()),
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			if call, ok := expr.(*ast.CallExpr); ok {
				goStmt.Call = call
			}
		}
	}

	return goStmt
}

// VisitLabeledStmt transforms a labeled statement.
func (b *ASTBuilder) VisitLabeledStmt(ctx *LabeledStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	labeled := &ast.LabeledStmt{
		Colon: b.pos(ctx),
	}

	// Label
	if ident := ctx.IDENTIFIER(); ident != nil {
		labeled.Label = b.visitIdentifier(ident)
	}

	// Statement
	if stmtCtx := ctx.Statement(); stmtCtx != nil {
		if stmt := b.VisitStatement(stmtCtx); stmt != nil {
			labeled.Stmt = stmt.(ast.Stmt)
		}
	}

	return labeled
}

// VisitIfStmt transforms an if statement.
func (b *ASTBuilder) VisitIfStmt(ctx *IfStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	ifStmt := &ast.IfStmt{
		If: b.tokenPos(ctx.IF().GetSymbol()),
	}

	// Initialization statement (optional)
	if simpleCtx := ctx.SimpleStmt(); simpleCtx != nil {
		if stmt := b.VisitSimpleStmt(simpleCtx); stmt != nil {
			ifStmt.Init = stmt.(ast.Stmt)
		}
	}

	// Condition
	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			ifStmt.Cond = expr.(ast.Expr)
		}
	}

	// Body
	blocks := ctx.AllBlock()
	if len(blocks) >= 1 {
		if block := b.VisitBlock(blocks[0]); block != nil {
			ifStmt.Body = block.(*ast.BlockStmt)
		}
	}

	// Else branch
	if len(blocks) >= 2 {
		if block := b.VisitBlock(blocks[1]); block != nil {
			ifStmt.Else = block.(*ast.BlockStmt)
		}
	} else if elseIfCtx := ctx.IfStmt(); elseIfCtx != nil {
		if elseIf := b.VisitIfStmt(elseIfCtx); elseIf != nil {
			ifStmt.Else = elseIf.(ast.Stmt)
		}
	}

	return ifStmt
}

// VisitForStmt transforms a for statement.
func (b *ASTBuilder) VisitForStmt(ctx *ForStmtContext) interface{} {
	if ctx == nil {
		return nil
	}

	forStmt := &ast.ForStmt{
		For: b.tokenPos(ctx.FOR().GetSymbol()),
	}

	// For clause (init; cond; post)
	if clauseCtx := ctx.ForClause(); clauseCtx != nil {
		if clause := b.VisitForClause(clauseCtx); clause != nil {
			if fs, ok := clause.(*ast.ForStmt); ok {
				forStmt.Init = fs.Init
				forStmt.Cond = fs.Cond
				forStmt.Post = fs.Post
			}
		}
	}

	// Range clause
	if rangeCtx := ctx.RangeClause(); rangeCtx != nil {
		if rangeStmt := b.VisitRangeClause(rangeCtx); rangeStmt != nil {
			// Return range statement instead
			if rs, ok := rangeStmt.(*ast.RangeStmt); ok {
				rs.For = forStmt.For
				if blockCtx := ctx.Block(); blockCtx != nil {
					if block := b.VisitBlock(blockCtx); block != nil {
						rs.Body = block.(*ast.BlockStmt)
					}
				}
				return rs
			}
		}
	}

	// Simple condition (infinite loop if nil)
	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			forStmt.Cond = expr.(ast.Expr)
		}
	}

	// Body
	if blockCtx := ctx.Block(); blockCtx != nil {
		if block := b.VisitBlock(blockCtx); block != nil {
			forStmt.Body = block.(*ast.BlockStmt)
		}
	}

	return forStmt
}

// VisitForClause transforms a for clause.
func (b *ASTBuilder) VisitForClause(ctx *ForClauseContext) interface{} {
	if ctx == nil {
		return nil
	}

	forStmt := &ast.ForStmt{}

	// Init, cond, post
	stmts := ctx.AllSimpleStmt()
	if len(stmts) >= 1 && stmts[0] != nil {
		if stmt := b.VisitSimpleStmt(stmts[0]); stmt != nil {
			forStmt.Init = stmt.(ast.Stmt)
		}
	}
	if len(stmts) >= 2 && stmts[1] != nil {
		if stmt := b.VisitSimpleStmt(stmts[1]); stmt != nil {
			forStmt.Post = stmt.(ast.Stmt)
		}
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			forStmt.Cond = expr.(ast.Expr)
		}
	}

	return forStmt
}

// VisitRangeClause transforms a range clause.
func (b *ASTBuilder) VisitRangeClause(ctx *RangeClauseContext) interface{} {
	if ctx == nil {
		return nil
	}

	rangeStmt := &ast.RangeStmt{
		TokPos: b.pos(ctx),
		Tok:    ast.ASSIGN,
	}

	// Check if it's a short var decl (:=)
	if ctx.GetDefine() != nil {
		rangeStmt.Tok = ast.DEFINE
	}

	// Key and value
	exprs := ctx.AllExpression()
	if len(exprs) >= 1 {
		if expr := b.VisitExpression(exprs[0]); expr != nil {
			rangeStmt.Key = expr.(ast.Expr)
		}
	}
	if len(exprs) >= 2 {
		if expr := b.VisitExpression(exprs[1]); expr != nil {
			rangeStmt.Value = expr.(ast.Expr)
		}
	}

	// Range expression
	if len(exprs) >= 3 {
		if expr := b.VisitExpression(exprs[2]); expr != nil {
			rangeStmt.X = expr.(ast.Expr)
		}
	} else if exprListCtx := ctx.ExpressionList(); exprListCtx != nil {
		if exprs := b.VisitExpressionList(exprListCtx); exprs != nil {
			exprList := exprs.([]ast.Expr)
			if len(exprList) > 0 {
				rangeStmt.X = exprList[len(exprList)-1]
			}
		}
	}

	return rangeStmt
}

// Placeholder stubs for switch and select (can be expanded later)

func (b *ASTBuilder) VisitSwitchStmt(ctx *SwitchStmtContext) interface{} {
	// Simplified: return a basic switch statement
	return &ast.SwitchStmt{
		Switch: b.pos(ctx),
	}
}

func (b *ASTBuilder) VisitSelectStmt(ctx *SelectStmtContext) interface{} {
	return &ast.SelectStmt{
		Select: b.pos(ctx),
	}
}
