package antlr

import (
	"github.com/mleku/moxie/pkg/ast"
)

// ============================================================================
// Constant Declarations
// ============================================================================

// VisitConstDecl transforms a const declaration.
func (b *ASTBuilder) VisitConstDecl(ctx *ConstDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	decl := &ast.ConstDecl{
		Const: b.tokenPos(ctx.CONST().GetSymbol()),
	}

	// Get all const specs
	for _, specCtx := range ctx.AllConstSpec() {
		if sCtx, ok := specCtx.(*ConstSpecContext); ok {
			if spec := b.VisitConstSpec(sCtx); spec != nil {
				decl.Specs = append(decl.Specs, spec.(*ast.ConstSpec))
			}
		}
	}

	return decl
}

// VisitConstSpec transforms a const specification.
func (b *ASTBuilder) VisitConstSpec(ctx *ConstSpecContext) interface{} {
	if ctx == nil {
		return nil
	}

	spec := &ast.ConstSpec{}

	// Constant names
	if idListCtx := ctx.IdentifierList(); idListCtx != nil {
		if idList, ok := idListCtx.(*IdentifierListContext); ok {
			spec.Names = b.visitIdentifierList(idList)
		}
	}

	// Type (optional)
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*Type_Context); ok {
			if typ := b.VisitType_(tCtx); typ != nil {
				spec.Type = typ.(ast.Type)
			}
		}
	}

	// Values
	if exprListCtx := ctx.ExpressionList(); exprListCtx != nil {
		if eCtx, ok := exprListCtx.(*ExpressionListContext); ok {
			if exprs := b.VisitExpressionList(eCtx); exprs != nil {
				spec.Values = exprs.([]ast.Expr)
			}
		}
	}

	return spec
}

// ============================================================================
// Variable Declarations
// ============================================================================

// VisitVarDecl transforms a var declaration.
func (b *ASTBuilder) VisitVarDecl(ctx *VarDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	decl := &ast.VarDecl{
		Var: b.tokenPos(ctx.VAR().GetSymbol()),
	}

	// Get all var specs
	for _, specCtx := range ctx.AllVarSpec() {
		if sCtx, ok := specCtx.(*VarSpecContext); ok {
			if spec := b.VisitVarSpec(sCtx); spec != nil {
				decl.Specs = append(decl.Specs, spec.(*ast.VarSpec))
			}
		}
	}

	return decl
}

// VisitVarSpec transforms a var specification.
func (b *ASTBuilder) VisitVarSpec(ctx *VarSpecContext) interface{} {
	if ctx == nil {
		return nil
	}

	spec := &ast.VarSpec{}

	// Variable names
	if idListCtx := ctx.IdentifierList(); idListCtx != nil {
		if idList, ok := idListCtx.(*IdentifierListContext); ok {
			spec.Names = b.visitIdentifierList(idList)
		}
	}

	// Type (optional if values are present)
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*Type_Context); ok {
			if typ := b.VisitType_(tCtx); typ != nil {
				spec.Type = typ.(ast.Type)
			}
		}
	}

	// Values (optional)
	if exprListCtx := ctx.ExpressionList(); exprListCtx != nil {
		if eCtx, ok := exprListCtx.(*ExpressionListContext); ok {
			if exprs := b.VisitExpressionList(eCtx); exprs != nil {
				spec.Values = exprs.([]ast.Expr)
			}
		}
	}

	return spec
}

// ============================================================================
// Type Declarations
// ============================================================================

// VisitTypeDecl transforms a type declaration.
func (b *ASTBuilder) VisitTypeDecl(ctx *TypeDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	decl := &ast.TypeDecl{
		Type: b.tokenPos(ctx.TYPE().GetSymbol()),
	}

	// Get all type specs
	for _, specCtx := range ctx.AllTypeSpec() {
		// TypeSpec is an interface - check for concrete types
		if aliasCtx, ok := specCtx.(*TypeAliasContext); ok {
			if spec := b.VisitTypeAlias(aliasCtx); spec != nil {
				decl.Specs = append(decl.Specs, spec.(*ast.TypeSpec))
			}
		} else if defCtx, ok := specCtx.(*TypeDefContext); ok {
			if spec := b.VisitTypeDef(defCtx); spec != nil {
				decl.Specs = append(decl.Specs, spec.(*ast.TypeSpec))
			}
		}
	}

	return decl
}

// VisitTypeAlias transforms a type alias (type A = B).
func (b *ASTBuilder) VisitTypeAlias(ctx *TypeAliasContext) interface{} {
	if ctx == nil {
		return nil
	}

	spec := &ast.TypeSpec{}

	// Type name
	if ident := ctx.IDENTIFIER(); ident != nil {
		spec.Name = b.visitIdentifier(ident)
	}

	// Mark as alias - type aliases have "=" in the grammar
	// We'll use a position to mark it (the actual "=" token position would need grammar analysis)
	spec.Assign = b.pos(ctx)

	// Underlying type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*Type_Context); ok {
			if typ := b.VisitType_(tCtx); typ != nil {
				spec.Type = typ.(ast.Type)
			}
		}
	}

	return spec
}

// VisitTypeDef transforms a type definition (type A B).
func (b *ASTBuilder) VisitTypeDef(ctx *TypeDefContext) interface{} {
	if ctx == nil {
		return nil
	}

	spec := &ast.TypeSpec{}

	// Type name
	if ident := ctx.IDENTIFIER(); ident != nil {
		spec.Name = b.visitIdentifier(ident)
	}

	// Type parameters (generics)
	if typeParamsCtx := ctx.TypeParameters(); typeParamsCtx != nil {
		if tpCtx, ok := typeParamsCtx.(*TypeParametersContext); ok {
			if typeParams := b.VisitTypeParameters(tpCtx); typeParams != nil {
				spec.TypeParams = typeParams.(*ast.FieldList)
			}
		}
	}

	// Underlying type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*Type_Context); ok {
			if typ := b.VisitType_(tCtx); typ != nil {
				spec.Type = typ.(ast.Type)
			}
		}
	}

	return spec
}

// VisitTypeParameters transforms type parameters (generics).
func (b *ASTBuilder) VisitTypeParameters(ctx *TypeParametersContext) interface{} {
	if ctx == nil {
		return nil
	}

	fieldList := &ast.FieldList{
		Opening: b.pos(ctx),
		Closing: b.endPos(ctx),
	}

	// Add type parameter declarations
	for _, paramCtx := range ctx.AllTypeParameterDecl() {
		if pCtx, ok := paramCtx.(*TypeParameterDeclContext); ok {
			if param := b.VisitTypeParameterDecl(pCtx); param != nil {
				fieldList.List = append(fieldList.List, param.(*ast.Field))
			}
		}
	}

	return fieldList
}

// VisitTypeParameterDecl transforms a type parameter declaration.
func (b *ASTBuilder) VisitTypeParameterDecl(ctx *TypeParameterDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	field := &ast.Field{}

	// Type parameter names
	if idListCtx := ctx.IdentifierList(); idListCtx != nil {
		if idList, ok := idListCtx.(*IdentifierListContext); ok {
			field.Names = b.visitIdentifierList(idList)
		}
	}

	// Type constraint
	if constraintCtx := ctx.TypeConstraint(); constraintCtx != nil {
		if cCtx, ok := constraintCtx.(*TypeConstraintContext); ok {
			if constraint := b.VisitTypeConstraint(cCtx); constraint != nil {
				field.Type = constraint.(ast.Type)
			}
		}
	}

	return field
}

// VisitTypeConstraint transforms a type constraint.
func (b *ASTBuilder) VisitTypeConstraint(ctx *TypeConstraintContext) interface{} {
	if ctx == nil {
		return nil
	}

	if typeCtx := ctx.Type_(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*Type_Context); ok {
			return b.VisitType_(tCtx)
		}
	}

	return nil
}

// ============================================================================
// Function Declarations
// ============================================================================

// VisitFunctionDecl transforms a function declaration.
func (b *ASTBuilder) VisitFunctionDecl(ctx *FunctionDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	decl := &ast.FuncDecl{}

	// Function name
	if ident := ctx.IDENTIFIER(); ident != nil {
		decl.Name = b.visitIdentifier(ident)
	}

	// Function signature
	if sigCtx := ctx.Signature(); sigCtx != nil {
		if sCtx, ok := sigCtx.(*SignatureContext); ok {
			if sig := b.VisitSignature(sCtx); sig != nil {
				decl.Type = sig.(*ast.FuncType)
			}
		}
	}

	// Function body (may be nil for external/FFI functions)
	if blockCtx := ctx.Block(); blockCtx != nil {
		if bCtx, ok := blockCtx.(*BlockContext); ok {
			if block := b.VisitBlock(bCtx); block != nil {
				decl.Body = block.(*ast.BlockStmt)
			}
		}
	}

	return decl
}

// VisitMethodDecl transforms a method declaration.
func (b *ASTBuilder) VisitMethodDecl(ctx *MethodDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	decl := &ast.FuncDecl{}

	// Receiver
	if recvCtx := ctx.Receiver(); recvCtx != nil {
		if rCtx, ok := recvCtx.(*ReceiverContext); ok {
			if recv := b.VisitReceiver(rCtx); recv != nil {
				decl.Recv = recv.(*ast.FieldList)
			}
		}
	}

	// Method name
	if ident := ctx.IDENTIFIER(); ident != nil {
		decl.Name = b.visitIdentifier(ident)
	}

	// Method signature
	if sigCtx := ctx.Signature(); sigCtx != nil {
		if sCtx, ok := sigCtx.(*SignatureContext); ok {
			if sig := b.VisitSignature(sCtx); sig != nil {
				decl.Type = sig.(*ast.FuncType)
			}
		}
	}

	// Method body
	if blockCtx := ctx.Block(); blockCtx != nil {
		if bCtx, ok := blockCtx.(*BlockContext); ok {
			if block := b.VisitBlock(bCtx); block != nil {
				decl.Body = block.(*ast.BlockStmt)
			}
		}
	}

	return decl
}

// VisitReceiver transforms a method receiver.
func (b *ASTBuilder) VisitReceiver(ctx *ReceiverContext) interface{} {
	if ctx == nil {
		return nil
	}

	fieldList := &ast.FieldList{
		Opening: b.pos(ctx),
		Closing: b.endPos(ctx),
	}

	// Receiver is a single parameter
	field := &ast.Field{}

	// Receiver name (optional)
	if ident := ctx.IDENTIFIER(); ident != nil {
		field.Names = []*ast.Ident{b.visitIdentifier(ident)}
	}

	// Receiver type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if tCtx, ok := typeCtx.(*Type_Context); ok {
			if typ := b.VisitType_(tCtx); typ != nil {
				field.Type = typ.(ast.Type)
			}
		}
	}

	fieldList.List = []*ast.Field{field}

	return fieldList
}
