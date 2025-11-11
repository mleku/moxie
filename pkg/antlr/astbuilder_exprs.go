package antlr

import (
	"github.com/mleku/moxie/pkg/ast"
)

// ============================================================================
// Expressions
// ============================================================================

// VisitExpression transforms an expression (handles precedence).
func (b *ASTBuilder) VisitExpression(ctx *ExpressionContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Check for unary expression first
	if unaryCtx := ctx.UnaryExpr(); unaryCtx != nil {
		return b.VisitUnaryExpr(unaryCtx)
	}

	// Check for primary expression
	if primaryCtx := ctx.PrimaryExpr(); primaryCtx != nil {
		return b.VisitPrimaryExpr(primaryCtx)
	}

	// Binary expression (left op right)
	exprs := ctx.AllExpression()
	if len(exprs) >= 2 {
		left := b.VisitExpression(exprs[0])
		right := b.VisitExpression(exprs[1])

		if left != nil && right != nil {
			binary := &ast.BinaryExpr{
				X:     left.(ast.Expr),
				OpPos: b.pos(ctx),
				Y:     right.(ast.Expr),
			}

			// Determine operator from context
			if mulOp := ctx.Mul_op(); mulOp != nil {
				binary.Op = b.VisitMul_op(mulOp).(ast.Token)
			} else if addOp := ctx.Add_op(); addOp != nil {
				binary.Op = b.VisitAdd_op(addOp).(ast.Token)
			} else if relOp := ctx.Rel_op(); relOp != nil {
				binary.Op = b.VisitRel_op(relOp).(ast.Token)
			}

			return binary
		}
	}

	// Fallback to first expression
	if len(exprs) > 0 {
		return b.VisitExpression(exprs[0])
	}

	return nil
}

// VisitPrimaryExpr transforms a primary expression.
func (b *ASTBuilder) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Operand (literal, identifier, etc.)
	if operandCtx := ctx.Operand(); operandCtx != nil {
		return b.VisitOperand(operandCtx)
	}

	// Conversion
	if convCtx := ctx.Conversion(); convCtx != nil {
		return b.VisitConversion(convCtx)
	}

	// Selector (x.y)
	if selCtx := ctx.Selector(); selCtx != nil {
		base := b.VisitPrimaryExpr(ctx.PrimaryExpr())
		sel := b.VisitSelector(selCtx)
		if base != nil && sel != nil {
			return &ast.SelectorExpr{
				X:   base.(ast.Expr),
				Sel: sel.(*ast.Ident),
			}
		}
	}

	// Index (x[i])
	if idxCtx := ctx.Index(); idxCtx != nil {
		base := b.VisitPrimaryExpr(ctx.PrimaryExpr())
		idx := b.VisitIndex(idxCtx)
		if base != nil && idx != nil {
			return &ast.IndexExpr{
				X:      base.(ast.Expr),
				Lbrack: b.pos(ctx),
				Index:  idx.(ast.Expr),
				Rbrack: b.endPos(ctx),
			}
		}
	}

	// Slice (x[i:j] or x[i:j:k])
	if sliceCtx := ctx.Slice_(); sliceCtx != nil {
		base := b.VisitPrimaryExpr(ctx.PrimaryExpr())
		slice := b.VisitSlice_(sliceCtx)
		if base != nil && slice != nil {
			sliceExpr := slice.(*ast.SliceExpr)
			sliceExpr.X = base.(ast.Expr)
			return sliceExpr
		}
	}

	// Type assertion (x.(T))
	if assertCtx := ctx.TypeAssertion(); assertCtx != nil {
		base := b.VisitPrimaryExpr(ctx.PrimaryExpr())
		assert := b.VisitTypeAssertion(assertCtx)
		if base != nil && assert != nil {
			assertExpr := assert.(*ast.TypeAssertExpr)
			assertExpr.X = base.(ast.Expr)
			return assertExpr
		}
	}

	// Arguments (function call)
	if argsCtx := ctx.Arguments(); argsCtx != nil {
		base := b.VisitPrimaryExpr(ctx.PrimaryExpr())
		args := b.VisitArguments(argsCtx)
		if base != nil {
			call := &ast.CallExpr{
				Fun:    base.(ast.Expr),
				Lparen: b.pos(ctx),
				Rparen: b.endPos(ctx),
			}
			if args != nil {
				if argList, ok := args.([]ast.Expr); ok {
					call.Args = argList
				}
			}
			return call
		}
	}

	return nil
}

// VisitUnaryExpr transforms a unary expression.
func (b *ASTBuilder) VisitUnaryExpr(ctx *UnaryExprContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Primary expression (base case)
	if primaryCtx := ctx.PrimaryExpr(); primaryCtx != nil {
		return b.VisitPrimaryExpr(primaryCtx)
	}

	// Unary operator + expression
	if unaryOpCtx := ctx.Unary_op(); unaryOpCtx != nil {
		unary := &ast.UnaryExpr{
			OpPos: b.pos(ctx),
		}

		if op := b.VisitUnary_op(unaryOpCtx); op != nil {
			unary.Op = op.(ast.Token)
		}

		if exprCtx := ctx.Expression(); exprCtx != nil {
			if expr := b.VisitExpression(exprCtx); expr != nil {
				unary.X = expr.(ast.Expr)
			}
		}

		return unary
	}

	return nil
}

// VisitOperand transforms an operand.
func (b *ASTBuilder) VisitOperand(ctx *OperandContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Literal
	if litCtx := ctx.Literal(); litCtx != nil {
		return b.VisitLiteral(litCtx)
	}

	// Operand name (identifier)
	if nameCtx := ctx.OperandName(); nameCtx != nil {
		return b.VisitOperandName(nameCtx)
	}

	// Parenthesized expression
	if exprCtx := ctx.Expression(); exprCtx != nil {
		expr := b.VisitExpression(exprCtx)
		if expr != nil {
			return &ast.ParenExpr{
				Lparen: b.pos(ctx),
				X:      expr.(ast.Expr),
				Rparen: b.endPos(ctx),
			}
		}
	}

	return nil
}

// VisitOperandName transforms an operand name (identifier or qualified).
func (b *ASTBuilder) VisitOperandName(ctx *OperandNameContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Qualified identifier
	if qualCtx := ctx.QualifiedIdent(); qualCtx != nil {
		return b.VisitQualifiedIdent(qualCtx)
	}

	// Simple identifier
	if ident := ctx.IDENTIFIER(); ident != nil {
		return b.visitIdentifier(ident)
	}

	return nil
}

// VisitSelector transforms a selector (.field).
func (b *ASTBuilder) VisitSelector(ctx *SelectorContext) interface{} {
	if ctx == nil {
		return nil
	}

	if ident := ctx.IDENTIFIER(); ident != nil {
		return b.visitIdentifier(ident)
	}

	return nil
}

// VisitIndex transforms an index expression.
func (b *ASTBuilder) VisitIndex(ctx *IndexContext) interface{} {
	if ctx == nil {
		return nil
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		return b.VisitExpression(exprCtx)
	}

	return nil
}

// VisitSlice_ transforms a slice expression.
func (b *ASTBuilder) VisitSlice_(ctx *Slice_Context) interface{} {
	if ctx == nil {
		return nil
	}

	slice := &ast.SliceExpr{
		Lbrack: b.pos(ctx),
		Rbrack: b.endPos(ctx),
	}

	exprs := ctx.AllExpression()
	if len(exprs) >= 1 && exprs[0] != nil {
		if expr := b.VisitExpression(exprs[0]); expr != nil {
			slice.Low = expr.(ast.Expr)
		}
	}
	if len(exprs) >= 2 && exprs[1] != nil {
		if expr := b.VisitExpression(exprs[1]); expr != nil {
			slice.High = expr.(ast.Expr)
		}
	}
	if len(exprs) >= 3 && exprs[2] != nil {
		if expr := b.VisitExpression(exprs[2]); expr != nil {
			slice.Max = expr.(ast.Expr)
			slice.Slice3 = true
		}
	}

	return slice
}

// VisitTypeAssertion transforms a type assertion.
func (b *ASTBuilder) VisitTypeAssertion(ctx *TypeAssertionContext) interface{} {
	if ctx == nil {
		return nil
	}

	assert := &ast.TypeAssertExpr{
		Lparen: b.pos(ctx),
		Rparen: b.endPos(ctx),
	}

	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			assert.Type = typ.(ast.Type)
		}
	}

	return assert
}

// VisitArguments transforms function arguments.
func (b *ASTBuilder) VisitArguments(ctx *ArgumentsContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Expression list or type with expression list
	if exprListCtx := ctx.ExpressionList(); exprListCtx != nil {
		return b.VisitExpressionList(exprListCtx)
	}

	return []ast.Expr{}
}

// VisitConversion transforms a type conversion.
func (b *ASTBuilder) VisitConversion(ctx *ConversionContext) interface{} {
	if ctx == nil {
		return nil
	}

	call := &ast.CallExpr{
		Lparen: b.pos(ctx),
		Rparen: b.endPos(ctx),
	}

	// Type (used as function)
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			call.Fun = typ.(ast.Expr)
		}
	}

	// Expression to convert
	if exprCtx := ctx.Expression(); exprCtx != nil {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			call.Args = []ast.Expr{expr.(ast.Expr)}
		}
	}

	return call
}

// VisitExpressionList transforms an expression list.
func (b *ASTBuilder) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	if ctx == nil {
		return nil
	}

	var exprs []ast.Expr
	for _, exprCtx := range ctx.AllExpression() {
		if expr := b.VisitExpression(exprCtx); expr != nil {
			exprs = append(exprs, expr.(ast.Expr))
		}
	}

	return exprs
}

// ============================================================================
// Operators
// ============================================================================

// VisitMul_op transforms a multiplication operator.
func (b *ASTBuilder) VisitMul_op(ctx *Mul_opContext) interface{} {
	if ctx == nil {
		return ast.MUL
	}

	text := ctx.GetText()
	switch text {
	case "*":
		return ast.MUL
	case "/":
		return ast.QUO
	case "%":
		return ast.REM
	case "<<":
		return ast.SHL
	case ">>":
		return ast.SHR
	case "&":
		return ast.AND
	case "&^":
		return ast.AND_NOT
	default:
		return ast.MUL
	}
}

// VisitAdd_op transforms an addition operator.
func (b *ASTBuilder) VisitAdd_op(ctx *Add_opContext) interface{} {
	if ctx == nil {
		return ast.ADD
	}

	text := ctx.GetText()
	switch text {
	case "+":
		return ast.ADD
	case "-":
		return ast.SUB
	case "|":
		return ast.OR
	case "^":
		return ast.XOR
	default:
		return ast.ADD
	}
}

// VisitRel_op transforms a relational operator.
func (b *ASTBuilder) VisitRel_op(ctx *Rel_opContext) interface{} {
	if ctx == nil {
		return ast.EQL
	}

	text := ctx.GetText()
	switch text {
	case "==":
		return ast.EQL
	case "!=":
		return ast.NEQ
	case "<":
		return ast.LSS
	case "<=":
		return ast.LEQ
	case ">":
		return ast.GTR
	case ">=":
		return ast.GEQ
	default:
		return ast.EQL
	}
}

// VisitUnary_op transforms a unary operator.
func (b *ASTBuilder) VisitUnary_op(ctx *Unary_opContext) interface{} {
	if ctx == nil {
		return ast.ADD
	}

	text := ctx.GetText()
	switch text {
	case "+":
		return ast.ADD
	case "-":
		return ast.SUB
	case "!":
		return ast.NOT
	case "^":
		return ast.XOR
	case "*":
		return ast.MUL // Pointer dereference
	case "&":
		return ast.AND // Address-of
	case "<-":
		return ast.ARROW // Channel receive
	default:
		return ast.ADD
	}
}

// ============================================================================
// Literals
// ============================================================================

// VisitLiteral transforms a literal.
func (b *ASTBuilder) VisitLiteral(ctx *LiteralContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Basic literal
	if basicCtx := ctx.BasicLit(); basicCtx != nil {
		return b.VisitBasicLit(basicCtx)
	}

	// Composite literal
	if compCtx := ctx.CompositeLit(); compCtx != nil {
		return b.VisitCompositeLit(compCtx)
	}

	// Function literal
	if funcCtx := ctx.FunctionLit(); funcCtx != nil {
		return b.VisitFunctionLit(funcCtx)
	}

	return nil
}

// VisitBasicLit transforms a basic literal.
func (b *ASTBuilder) VisitBasicLit(ctx *BasicLitContext) interface{} {
	if ctx == nil {
		return nil
	}

	lit := &ast.BasicLit{
		ValuePos: b.pos(ctx),
	}

	// Determine literal kind and value
	if ctx.INT_LIT() != nil {
		lit.Kind = ast.IntLit
		lit.Value = ctx.INT_LIT().GetText()
	} else if ctx.FLOAT_LIT() != nil {
		lit.Kind = ast.FloatLit
		lit.Value = ctx.FLOAT_LIT().GetText()
	} else if ctx.IMAGINARY_LIT() != nil {
		lit.Kind = ast.ImagLit
		lit.Value = ctx.IMAGINARY_LIT().GetText()
	} else if ctx.RUNE_LIT() != nil {
		lit.Kind = ast.RuneLit
		lit.Value = ctx.RUNE_LIT().GetText()
	} else if strCtx := ctx.String_(); strCtx != nil {
		if str := b.VisitString_(strCtx); str != nil {
			return str
		}
	}

	return lit
}

// VisitString_ transforms a string literal.
func (b *ASTBuilder) VisitString_(ctx *String_Context) interface{} {
	if ctx == nil {
		return nil
	}

	lit := &ast.BasicLit{
		ValuePos: b.pos(ctx),
		Kind:     ast.StringLit,
	}

	if ctx.RAW_STRING_LIT() != nil {
		lit.Value = ctx.RAW_STRING_LIT().GetText()
	} else if ctx.INTERPRETED_STRING_LIT() != nil {
		lit.Value = ctx.INTERPRETED_STRING_LIT().GetText()
	}

	return lit
}

// VisitCompositeLit transforms a composite literal.
func (b *ASTBuilder) VisitCompositeLit(ctx *CompositeLitContext) interface{} {
	if ctx == nil {
		return nil
	}

	comp := &ast.CompositeLit{
		Lbrace: b.pos(ctx),
		Rbrace: b.endPos(ctx),
	}

	// Literal type
	if litTypeCtx := ctx.LiteralType(); litTypeCtx != nil {
		if typ := b.VisitLiteralType(litTypeCtx); typ != nil {
			comp.Type = typ.(ast.Type)
		}
	}

	// Literal value (elements)
	if litValCtx := ctx.LiteralValue(); litValCtx != nil {
		if val := b.VisitLiteralValue(litValCtx); val != nil {
			if elts, ok := val.([]ast.Expr); ok {
				comp.Elts = elts
			}
		}
	}

	return comp
}

// VisitLiteralType transforms a literal type.
func (b *ASTBuilder) VisitLiteralType(ctx *LiteralTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Struct, array, slice, map, type name, etc.
	if typeCtx := ctx.Type_(); typeCtx != nil {
		return b.VisitType_(typeCtx)
	}

	return nil
}

// VisitLiteralValue transforms a literal value (element list).
func (b *ASTBuilder) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	if ctx == nil {
		return nil
	}

	if elemListCtx := ctx.ElementList(); elemListCtx != nil {
		return b.VisitElementList(elemListCtx)
	}

	return []ast.Expr{}
}

// VisitElementList transforms an element list.
func (b *ASTBuilder) VisitElementList(ctx *ElementListContext) interface{} {
	if ctx == nil {
		return nil
	}

	var elts []ast.Expr
	for _, keyedElemCtx := range ctx.AllKeyedElement() {
		if elem := b.VisitKeyedElement(keyedElemCtx); elem != nil {
			elts = append(elts, elem.(ast.Expr))
		}
	}

	return elts
}

// VisitKeyedElement transforms a keyed element.
func (b *ASTBuilder) VisitKeyedElement(ctx *KeyedElementContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Check if it's a key:value pair
	if keyCtx := ctx.Key(); keyCtx != nil {
		kv := &ast.KeyValueExpr{
			Colon: b.pos(ctx),
		}

		if key := b.VisitKey(keyCtx); key != nil {
			kv.Key = key.(ast.Expr)
		}

		if elemCtx := ctx.Element(); elemCtx != nil {
			if val := b.VisitElement(elemCtx); val != nil {
				kv.Value = val.(ast.Expr)
			}
		}

		return kv
	}

	// Just an element (no key)
	if elemCtx := ctx.Element(); elemCtx != nil {
		return b.VisitElement(elemCtx)
	}

	return nil
}

// VisitKey transforms a key in a keyed element.
func (b *ASTBuilder) VisitKey(ctx *KeyContext) interface{} {
	if ctx == nil {
		return nil
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		return b.VisitExpression(exprCtx)
	}

	return nil
}

// VisitElement transforms an element value.
func (b *ASTBuilder) VisitElement(ctx *ElementContext) interface{} {
	if ctx == nil {
		return nil
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		return b.VisitExpression(exprCtx)
	}

	if litValCtx := ctx.LiteralValue(); litValCtx != nil {
		return b.VisitLiteralValue(litValCtx)
	}

	return nil
}

// VisitFunctionLit transforms a function literal.
func (b *ASTBuilder) VisitFunctionLit(ctx *FunctionLitContext) interface{} {
	if ctx == nil {
		return nil
	}

	funcLit := &ast.FuncLit{}

	// Function type
	if funcTypeCtx := ctx.FunctionType(); funcTypeCtx != nil {
		if funcType := b.VisitFunctionType(funcTypeCtx); funcType != nil {
			funcLit.Type = funcType.(*ast.FuncType)
		}
	}

	// Function body
	if blockCtx := ctx.Block(); blockCtx != nil {
		if block := b.VisitBlock(blockCtx); block != nil {
			funcLit.Body = block.(*ast.BlockStmt)
		}
	}

	return funcLit
}
