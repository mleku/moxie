package antlr

import (
	"github.com/mleku/moxie/pkg/ast"
)

// ============================================================================
// Type Expressions
// ============================================================================

// VisitType_ transforms a type expression.
func (b *ASTBuilder) VisitType_(ctx *Type_Context) interface{} {
	if ctx == nil {
		return nil
	}

	// Named type (identifier or qualified identifier)
	if namedCtx, ok := ctx.(*NamedTypeContext); ok {
		return b.VisitNamedType(namedCtx)
	}

	// Type literal (struct, interface, array, slice, map, chan, func, pointer)
	if litCtx, ok := ctx.(*TypeLiteralContext); ok {
		return b.VisitTypeLiteral(litCtx)
	}

	// Parenthesized type
	if parenCtx, ok := ctx.(*ParenTypeContext); ok {
		return b.VisitParenType(parenCtx)
	}

	// Const type (Moxie feature)
	if constCtx, ok := ctx.(*ConstTypeContext); ok {
		return b.VisitConstType(constCtx)
	}

	return nil
}

// VisitNamedType transforms a named type (identifier or qualified).
func (b *ASTBuilder) VisitNamedType(ctx *NamedTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	if typeNameCtx := ctx.TypeName(); typeNameCtx != nil {
		return b.VisitTypeName(typeNameCtx)
	}

	return nil
}

// VisitTypeName transforms a type name.
func (b *ASTBuilder) VisitTypeName(ctx *TypeNameContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Check for qualified identifier (package.Type)
	if qualCtx := ctx.QualifiedIdent(); qualCtx != nil {
		return b.VisitQualifiedIdent(qualCtx)
	}

	// Simple identifier
	if ident := ctx.IDENTIFIER(); ident != nil {
		return b.visitIdentifier(ident)
	}

	return nil
}

// VisitTypeLiteral transforms a type literal.
func (b *ASTBuilder) VisitTypeLiteral(ctx *TypeLiteralContext) interface{} {
	if ctx == nil {
		return nil
	}

	if litCtx := ctx.TypeLit(); litCtx != nil {
		return b.VisitTypeLit(litCtx)
	}

	return nil
}

// VisitTypeLit transforms a type literal (struct, interface, array, etc.).
func (b *ASTBuilder) VisitTypeLit(ctx *TypeLitContext) interface{} {
	if ctx == nil {
		return nil
	}

	if arrayCtx := ctx.ArrayType(); arrayCtx != nil {
		return b.VisitArrayType(arrayCtx)
	}

	if structCtx := ctx.StructType(); structCtx != nil {
		return b.VisitStructType(structCtx)
	}

	if ptrCtx := ctx.PointerType(); ptrCtx != nil {
		return b.VisitPointerType(ptrCtx)
	}

	if funcCtx := ctx.FunctionType(); funcCtx != nil {
		return b.VisitFunctionType(funcCtx)
	}

	if ifaceCtx := ctx.InterfaceType(); ifaceCtx != nil {
		return b.VisitInterfaceType(ifaceCtx)
	}

	if sliceCtx := ctx.SliceType(); sliceCtx != nil {
		return b.VisitSliceType(sliceCtx)
	}

	if mapCtx := ctx.MapType(); mapCtx != nil {
		return b.VisitMapType(mapCtx)
	}

	if chanCtx := ctx.ChannelType(); chanCtx != nil {
		return b.VisitChannelType(chanCtx)
	}

	return nil
}

// VisitParenType transforms a parenthesized type.
func (b *ASTBuilder) VisitParenType(ctx *ParenTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	paren := &ast.ParenType{
		Lparen: b.pos(ctx),
		Rparen: b.endPos(ctx),
	}

	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			paren.X = typ.(ast.Type)
		}
	}

	return paren
}

// VisitPointerType transforms a pointer type.
func (b *ASTBuilder) VisitPointerType(ctx *PointerTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	ptr := &ast.PointerType{
		Star: b.pos(ctx),
	}

	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			ptr.Base = typ.(ast.Type)
		}
	}

	return ptr
}

// VisitSliceType transforms a slice type.
func (b *ASTBuilder) VisitSliceType(ctx *SliceTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	slice := &ast.SliceType{
		Lbrack: b.pos(ctx),
	}

	if elemCtx := ctx.ElementType(); elemCtx != nil {
		if elem := b.VisitElementType(elemCtx); elem != nil {
			slice.Elem = elem.(ast.Type)
		}
	}

	return slice
}

// VisitElementType transforms an element type.
func (b *ASTBuilder) VisitElementType(ctx *ElementTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	if typeCtx := ctx.Type_(); typeCtx != nil {
		return b.VisitType_(typeCtx)
	}

	return nil
}

// VisitArrayType transforms an array type.
func (b *ASTBuilder) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	array := &ast.ArrayType{
		Lbrack: b.pos(ctx),
	}

	if lenCtx := ctx.ArrayLength(); lenCtx != nil {
		if length := b.VisitArrayLength(lenCtx); length != nil {
			array.Len = length.(ast.Expr)
		}
	}

	if elemCtx := ctx.ElementType(); elemCtx != nil {
		if elem := b.VisitElementType(elemCtx); elem != nil {
			array.Elem = elem.(ast.Type)
		}
	}

	return array
}

// VisitArrayLength transforms an array length expression.
func (b *ASTBuilder) VisitArrayLength(ctx *ArrayLengthContext) interface{} {
	if ctx == nil {
		return nil
	}

	if exprCtx := ctx.Expression(); exprCtx != nil {
		return b.VisitExpression(exprCtx)
	}

	return nil
}

// VisitStructType transforms a struct type.
func (b *ASTBuilder) VisitStructType(ctx *StructTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	structType := &ast.StructType{
		Struct: b.pos(ctx),
		Lbrace: b.pos(ctx),
		Rbrace: b.endPos(ctx),
		Fields: &ast.FieldList{
			Opening: b.pos(ctx),
			Closing: b.endPos(ctx),
		},
	}

	// Add fields
	for _, fieldCtx := range ctx.AllFieldDecl() {
		if field := b.VisitFieldDecl(fieldCtx); field != nil {
			structType.Fields.List = append(structType.Fields.List, field.(*ast.Field))
		}
	}

	return structType
}

// VisitFieldDecl transforms a struct field declaration.
func (b *ASTBuilder) VisitFieldDecl(ctx *FieldDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	field := &ast.Field{}

	// Field names (if present)
	if idListCtx := ctx.IdentifierList(); idListCtx != nil {
		field.Names = b.visitIdentifierList(idListCtx)
	}

	// Field type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			field.Type = typ.(ast.Type)
		}
	}

	// Field tag (if present)
	if tagCtx := ctx.Tag_(); tagCtx != nil {
		if tag := b.VisitTag_(tagCtx); tag != nil {
			field.Tag = tag.(*ast.BasicLit)
		}
	}

	return field
}

// VisitTag_ transforms a struct field tag.
func (b *ASTBuilder) VisitTag_(ctx *Tag_Context) interface{} {
	if ctx == nil {
		return nil
	}

	if str := ctx.String_(); str != nil {
		return b.VisitString_(str)
	}

	return nil
}

// VisitInterfaceType transforms an interface type.
func (b *ASTBuilder) VisitInterfaceType(ctx *InterfaceTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	iface := &ast.InterfaceType{
		Interface: b.pos(ctx),
		Lbrace:    b.pos(ctx),
		Rbrace:    b.endPos(ctx),
		Methods:   &ast.FieldList{
			Opening: b.pos(ctx),
			Closing: b.endPos(ctx),
		},
	}

	// Add interface elements (methods and embedded types)
	for _, elemCtx := range ctx.AllInterfaceElem() {
		if elem := b.VisitInterfaceElem(elemCtx); elem != nil {
			iface.Methods.List = append(iface.Methods.List, elem.(*ast.Field))
		}
	}

	return iface
}

// VisitInterfaceElem transforms an interface element.
func (b *ASTBuilder) VisitInterfaceElem(ctx *InterfaceElemContext) interface{} {
	if ctx == nil {
		return nil
	}

	if methCtx := ctx.MethodElem(); methCtx != nil {
		return b.VisitMethodElem(methCtx)
	}

	if typeCtx := ctx.TypeElem(); typeCtx != nil {
		return b.VisitTypeElem(typeCtx)
	}

	return nil
}

// VisitMethodElem transforms an interface method element.
func (b *ASTBuilder) VisitMethodElem(ctx *MethodElemContext) interface{} {
	if ctx == nil {
		return nil
	}

	field := &ast.Field{}

	// Method name
	if ident := ctx.IDENTIFIER(); ident != nil {
		field.Names = []*ast.Ident{b.visitIdentifier(ident)}
	}

	// Method signature
	if sigCtx := ctx.Signature(); sigCtx != nil {
		if sig := b.VisitSignature(sigCtx); sig != nil {
			field.Type = sig.(ast.Type)
		}
	}

	return field
}

// VisitTypeElem transforms an interface type element (embedded type).
func (b *ASTBuilder) VisitTypeElem(ctx *TypeElemContext) interface{} {
	if ctx == nil {
		return nil
	}

	field := &ast.Field{}

	// Embedded type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			field.Type = typ.(ast.Type)
		}
	}

	return field
}

// VisitMapType transforms a map type.
func (b *ASTBuilder) VisitMapType(ctx *MapTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	mapType := &ast.MapType{
		Map:    b.pos(ctx),
		Lbrack: b.pos(ctx),
	}

	// Key type
	if keyCtx := ctx.Type_(0); keyCtx != nil {
		if key := b.VisitType_(keyCtx); key != nil {
			mapType.Key = key.(ast.Type)
		}
	}

	// Value type
	if valCtx := ctx.Type_(1); valCtx != nil {
		if val := b.VisitType_(valCtx); val != nil {
			mapType.Value = val.(ast.Type)
		}
	}

	return mapType
}

// VisitChannelType transforms a channel type.
func (b *ASTBuilder) VisitChannelType(ctx *ChannelTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	chanType := &ast.ChanType{
		Begin: b.pos(ctx),
		Dir:   ast.ChanBoth,
	}

	// Determine channel direction
	// This depends on the specific rule structure in your grammar
	// For now, assume bidirectional

	// Channel element type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			chanType.Value = typ.(ast.Type)
		}
	}

	return chanType
}

// VisitFunctionType transforms a function type.
func (b *ASTBuilder) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	funcType := &ast.FuncType{
		Func: b.pos(ctx),
	}

	if sigCtx := ctx.Signature(); sigCtx != nil {
		if sig := b.VisitSignature(sigCtx); sig != nil {
			// Signature returns a FuncType
			if ft, ok := sig.(*ast.FuncType); ok {
				funcType.Params = ft.Params
				funcType.Results = ft.Results
			}
		}
	}

	return funcType
}

// VisitSignature transforms a function signature.
func (b *ASTBuilder) VisitSignature(ctx *SignatureContext) interface{} {
	if ctx == nil {
		return nil
	}

	funcType := &ast.FuncType{}

	// Parameters
	if paramsCtx := ctx.Parameters(); paramsCtx != nil {
		if params := b.VisitParameters(paramsCtx); params != nil {
			funcType.Params = params.(*ast.FieldList)
		}
	}

	// Results
	if resultCtx := ctx.Result(); resultCtx != nil {
		if result := b.VisitResult(resultCtx); result != nil {
			funcType.Results = result.(*ast.FieldList)
		}
	}

	return funcType
}

// VisitParameters transforms function parameters.
func (b *ASTBuilder) VisitParameters(ctx *ParametersContext) interface{} {
	if ctx == nil {
		return nil
	}

	fieldList := &ast.FieldList{
		Opening: b.pos(ctx),
		Closing: b.endPos(ctx),
	}

	// Add parameter declarations
	for _, paramCtx := range ctx.AllParameterDecl() {
		if param := b.VisitParameterDecl(paramCtx); param != nil {
			fieldList.List = append(fieldList.List, param.(*ast.Field))
		}
	}

	return fieldList
}

// VisitParameterDecl transforms a parameter declaration.
func (b *ASTBuilder) VisitParameterDecl(ctx *ParameterDeclContext) interface{} {
	if ctx == nil {
		return nil
	}

	field := &ast.Field{}

	// Parameter names (may be empty)
	if idListCtx := ctx.IdentifierList(); idListCtx != nil {
		field.Names = b.visitIdentifierList(idListCtx)
	}

	// Parameter type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		if typ := b.VisitType_(typeCtx); typ != nil {
			field.Type = typ.(ast.Type)
		}
	}

	return field
}

// VisitResult transforms function results.
func (b *ASTBuilder) VisitResult(ctx *ResultContext) interface{} {
	if ctx == nil {
		return nil
	}

	// If result has parameters (named or unnamed), visit them
	if paramsCtx := ctx.Parameters(); paramsCtx != nil {
		return b.VisitParameters(paramsCtx)
	}

	// If result is a single type
	if typeCtx := ctx.Type_(); typeCtx != nil {
		fieldList := &ast.FieldList{}
		if typ := b.VisitType_(typeCtx); typ != nil {
			fieldList.List = []*ast.Field{
				{Type: typ.(ast.Type)},
			}
		}
		return fieldList
	}

	return nil
}

// VisitConstType transforms a const type (Moxie feature).
func (b *ASTBuilder) VisitConstType(ctx *ConstTypeContext) interface{} {
	if ctx == nil {
		return nil
	}

	// For now, treat const types as regular types
	// We'll need to mark them as const in semantic analysis
	if typeCtx := ctx.Type_(); typeCtx != nil {
		return b.VisitType_(typeCtx)
	}

	return nil
}

// VisitQualifiedIdent transforms a qualified identifier (package.Name).
func (b *ASTBuilder) VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Package name
	var pkgIdent *ast.Ident
	if pkgToken := ctx.IDENTIFIER(0); pkgToken != nil {
		pkgIdent = b.visitIdentifier(pkgToken)
	}

	// Type name
	var nameIdent *ast.Ident
	if nameToken := ctx.IDENTIFIER(1); nameToken != nil {
		nameIdent = b.visitIdentifier(nameToken)
	}

	// Return as selector expression
	if pkgIdent != nil && nameIdent != nil {
		return &ast.SelectorExpr{
			X:   pkgIdent,
			Sel: nameIdent,
		}
	}

	return nameIdent
}
