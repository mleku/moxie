// Code generated from grammar/Moxie.g4 by ANTLR 4.13.2. DO NOT EDIT.

package antlr // Moxie
import "github.com/antlr4-go/antlr/v4"

type BaseMoxieVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseMoxieVisitor) VisitSourceFile(ctx *SourceFileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitPackageClause(ctx *PackageClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitDeclaration(ctx *DeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitConstDecl(ctx *ConstDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitConstSpec(ctx *ConstSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeDecl(ctx *TypeDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeAlias(ctx *TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeDef(ctx *TypeDefContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeParameters(ctx *TypeParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeParameterDecl(ctx *TypeParameterDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeConstraint(ctx *TypeConstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitVarDecl(ctx *VarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitVarSpec(ctx *VarSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitFunctionDecl(ctx *FunctionDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMethodDecl(ctx *MethodDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitReceiver(ctx *ReceiverContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitNamedType(ctx *NamedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeLiteral(ctx *TypeLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitParenType(ctx *ParenTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitConstType(ctx *ConstTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeName(ctx *TypeNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeArgs(ctx *TypeArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeLit(ctx *TypeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitArrayLength(ctx *ArrayLengthContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitElementType(ctx *ElementTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSliceType(ctx *SliceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitStructType(ctx *StructTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitFieldDecl(ctx *FieldDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitEmbeddedField(ctx *EmbeddedFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTag_(ctx *Tag_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitPointerType(ctx *PointerTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSignature(ctx *SignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitResult(ctx *ResultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitParameters(ctx *ParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitParameterDecl(ctx *ParameterDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitInterfaceType(ctx *InterfaceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitInterfaceElem(ctx *InterfaceElemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMethodElem(ctx *MethodElemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeElem(ctx *TypeElemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeTerm(ctx *TypeTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMapType(ctx *MapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSendRecvChan(ctx *SendRecvChanContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitRecvOnlyChan(ctx *RecvOnlyChanContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSendRecvChanCompat(ctx *SendRecvChanCompatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitRecvOnlyChanCompat(ctx *RecvOnlyChanCompatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitStatementList(ctx *StatementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitDeclStmt(ctx *DeclStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSimpleStatement(ctx *SimpleStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLabeledStatement(ctx *LabeledStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitGoStatement(ctx *GoStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitGotoStatement(ctx *GotoStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitFallthroughStatement(ctx *FallthroughStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitBlockStatement(ctx *BlockStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSwitchStatement(ctx *SwitchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSelectStatement(ctx *SelectStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitForStatement(ctx *ForStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitDeferStatement(ctx *DeferStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSimpleStmt(ctx *SimpleStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitExpressionStmt(ctx *ExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSendStmt(ctx *SendStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitIncDecStmt(ctx *IncDecStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitAssignment(ctx *AssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitAssign_op(ctx *Assign_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitShortVarDecl(ctx *ShortVarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLabeledStmt(ctx *LabeledStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitBreakStmt(ctx *BreakStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitContinueStmt(ctx *ContinueStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitGotoStmt(ctx *GotoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitExprSwitchStmt(ctx *ExprSwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitExprCaseClause(ctx *ExprCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitExprSwitchCase(ctx *ExprSwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeSwitchStmt(ctx *TypeSwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeSwitchGuard(ctx *TypeSwitchGuardContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeCaseClause(ctx *TypeCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeSwitchCase(ctx *TypeSwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeList(ctx *TypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSelectStmt(ctx *SelectStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitCommClause(ctx *CommClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitCommCase(ctx *CommCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitRecvStmt(ctx *RecvStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitForStmt(ctx *ForStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitForClause(ctx *ForClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitRangeClause(ctx *RangeClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitGoStmt(ctx *GoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMultiplicativeExpr(ctx *MultiplicativeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitConcatenationExpr(ctx *ConcatenationExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitAdditiveExpr(ctx *AdditiveExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitUnaryExpression(ctx *UnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitRelationalExpr(ctx *RelationalExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSelectorExpr(ctx *SelectorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeAssertionExpr(ctx *TypeAssertionExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitConversionExpr(ctx *ConversionExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitPrimaryOperand(ctx *PrimaryOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSliceExpr(ctx *SliceExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitCallExpr(ctx *CallExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMethodExpression(ctx *MethodExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitIndexExpr(ctx *IndexExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitUnaryExpr(ctx *UnaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSimpleConversion(ctx *SimpleConversionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSliceCastExpr(ctx *SliceCastExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSliceCastEndianExpr(ctx *SliceCastEndianExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSliceCastCopyExpr(ctx *SliceCastCopyExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSliceCastCopyEndianExpr(ctx *SliceCastCopyEndianExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitEndianness(ctx *EndiannessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLiteralOperand(ctx *LiteralOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitNameOperand(ctx *NameOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitParenOperand(ctx *ParenOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitBasicLit(ctx *BasicLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitString_(ctx *String_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitOperandName(ctx *OperandNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitCompositeLit(ctx *CompositeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLiteralType(ctx *LiteralTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitElementList(ctx *ElementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitKeyedElement(ctx *KeyedElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitKey(ctx *KeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitElement(ctx *ElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitFunctionLit(ctx *FunctionLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSelector(ctx *SelectorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitIndex(ctx *IndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitSlice_(ctx *Slice_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitTypeAssertion(ctx *TypeAssertionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMethodExpr(ctx *MethodExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitMul_op(ctx *Mul_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitAdd_op(ctx *Add_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitRel_op(ctx *Rel_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitUnary_op(ctx *Unary_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitIdentifierList(ctx *IdentifierListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMoxieVisitor) VisitEos(ctx *EosContext) interface{} {
	return v.VisitChildren(ctx)
}
