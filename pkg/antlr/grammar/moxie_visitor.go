// Code generated from grammar/Moxie.g4 by ANTLR 4.13.2. DO NOT EDIT.

package antlr // Moxie
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by MoxieParser.
type MoxieVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by MoxieParser#sourceFile.
	VisitSourceFile(ctx *SourceFileContext) interface{}

	// Visit a parse tree produced by MoxieParser#packageClause.
	VisitPackageClause(ctx *PackageClauseContext) interface{}

	// Visit a parse tree produced by MoxieParser#importDecl.
	VisitImportDecl(ctx *ImportDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// Visit a parse tree produced by MoxieParser#topLevelDecl.
	VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#declaration.
	VisitDeclaration(ctx *DeclarationContext) interface{}

	// Visit a parse tree produced by MoxieParser#constDecl.
	VisitConstDecl(ctx *ConstDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#constSpec.
	VisitConstSpec(ctx *ConstSpecContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeDecl.
	VisitTypeDecl(ctx *TypeDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#TypeAlias.
	VisitTypeAlias(ctx *TypeAliasContext) interface{}

	// Visit a parse tree produced by MoxieParser#TypeDef.
	VisitTypeDef(ctx *TypeDefContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeParameters.
	VisitTypeParameters(ctx *TypeParametersContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeParameterDecl.
	VisitTypeParameterDecl(ctx *TypeParameterDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeConstraint.
	VisitTypeConstraint(ctx *TypeConstraintContext) interface{}

	// Visit a parse tree produced by MoxieParser#varDecl.
	VisitVarDecl(ctx *VarDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#varSpec.
	VisitVarSpec(ctx *VarSpecContext) interface{}

	// Visit a parse tree produced by MoxieParser#functionDecl.
	VisitFunctionDecl(ctx *FunctionDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#methodDecl.
	VisitMethodDecl(ctx *MethodDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#receiver.
	VisitReceiver(ctx *ReceiverContext) interface{}

	// Visit a parse tree produced by MoxieParser#NamedType.
	VisitNamedType(ctx *NamedTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#TypeLiteral.
	VisitTypeLiteral(ctx *TypeLiteralContext) interface{}

	// Visit a parse tree produced by MoxieParser#ParenType.
	VisitParenType(ctx *ParenTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#ConstType.
	VisitConstType(ctx *ConstTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeName.
	VisitTypeName(ctx *TypeNameContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeArgs.
	VisitTypeArgs(ctx *TypeArgsContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeLit.
	VisitTypeLit(ctx *TypeLitContext) interface{}

	// Visit a parse tree produced by MoxieParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#arrayLength.
	VisitArrayLength(ctx *ArrayLengthContext) interface{}

	// Visit a parse tree produced by MoxieParser#elementType.
	VisitElementType(ctx *ElementTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#sliceType.
	VisitSliceType(ctx *SliceTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#structType.
	VisitStructType(ctx *StructTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#fieldDecl.
	VisitFieldDecl(ctx *FieldDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#embeddedField.
	VisitEmbeddedField(ctx *EmbeddedFieldContext) interface{}

	// Visit a parse tree produced by MoxieParser#tag_.
	VisitTag_(ctx *Tag_Context) interface{}

	// Visit a parse tree produced by MoxieParser#pointerType.
	VisitPointerType(ctx *PointerTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#functionType.
	VisitFunctionType(ctx *FunctionTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#signature.
	VisitSignature(ctx *SignatureContext) interface{}

	// Visit a parse tree produced by MoxieParser#result.
	VisitResult(ctx *ResultContext) interface{}

	// Visit a parse tree produced by MoxieParser#parameters.
	VisitParameters(ctx *ParametersContext) interface{}

	// Visit a parse tree produced by MoxieParser#parameterDecl.
	VisitParameterDecl(ctx *ParameterDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#interfaceType.
	VisitInterfaceType(ctx *InterfaceTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#interfaceElem.
	VisitInterfaceElem(ctx *InterfaceElemContext) interface{}

	// Visit a parse tree produced by MoxieParser#methodElem.
	VisitMethodElem(ctx *MethodElemContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeElem.
	VisitTypeElem(ctx *TypeElemContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeTerm.
	VisitTypeTerm(ctx *TypeTermContext) interface{}

	// Visit a parse tree produced by MoxieParser#mapType.
	VisitMapType(ctx *MapTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#SendRecvChan.
	VisitSendRecvChan(ctx *SendRecvChanContext) interface{}

	// Visit a parse tree produced by MoxieParser#RecvOnlyChan.
	VisitRecvOnlyChan(ctx *RecvOnlyChanContext) interface{}

	// Visit a parse tree produced by MoxieParser#SendRecvChanCompat.
	VisitSendRecvChanCompat(ctx *SendRecvChanCompatContext) interface{}

	// Visit a parse tree produced by MoxieParser#RecvOnlyChanCompat.
	VisitRecvOnlyChanCompat(ctx *RecvOnlyChanCompatContext) interface{}

	// Visit a parse tree produced by MoxieParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by MoxieParser#statementList.
	VisitStatementList(ctx *StatementListContext) interface{}

	// Visit a parse tree produced by MoxieParser#DeclStmt.
	VisitDeclStmt(ctx *DeclStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#SimpleStatement.
	VisitSimpleStatement(ctx *SimpleStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#LabeledStatement.
	VisitLabeledStatement(ctx *LabeledStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#GoStatement.
	VisitGoStatement(ctx *GoStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#ReturnStatement.
	VisitReturnStatement(ctx *ReturnStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#BreakStatement.
	VisitBreakStatement(ctx *BreakStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#ContinueStatement.
	VisitContinueStatement(ctx *ContinueStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#GotoStatement.
	VisitGotoStatement(ctx *GotoStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#FallthroughStatement.
	VisitFallthroughStatement(ctx *FallthroughStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#BlockStatement.
	VisitBlockStatement(ctx *BlockStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#IfStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#SwitchStatement.
	VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#SelectStatement.
	VisitSelectStatement(ctx *SelectStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#ForStatement.
	VisitForStatement(ctx *ForStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#DeferStatement.
	VisitDeferStatement(ctx *DeferStatementContext) interface{}

	// Visit a parse tree produced by MoxieParser#simpleStmt.
	VisitSimpleStmt(ctx *SimpleStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#expressionStmt.
	VisitExpressionStmt(ctx *ExpressionStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#sendStmt.
	VisitSendStmt(ctx *SendStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#incDecStmt.
	VisitIncDecStmt(ctx *IncDecStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#assignment.
	VisitAssignment(ctx *AssignmentContext) interface{}

	// Visit a parse tree produced by MoxieParser#assign_op.
	VisitAssign_op(ctx *Assign_opContext) interface{}

	// Visit a parse tree produced by MoxieParser#shortVarDecl.
	VisitShortVarDecl(ctx *ShortVarDeclContext) interface{}

	// Visit a parse tree produced by MoxieParser#labeledStmt.
	VisitLabeledStmt(ctx *LabeledStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#returnStmt.
	VisitReturnStmt(ctx *ReturnStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#breakStmt.
	VisitBreakStmt(ctx *BreakStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#continueStmt.
	VisitContinueStmt(ctx *ContinueStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#gotoStmt.
	VisitGotoStmt(ctx *GotoStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#fallthroughStmt.
	VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#deferStmt.
	VisitDeferStmt(ctx *DeferStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#ifStmt.
	VisitIfStmt(ctx *IfStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#switchStmt.
	VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#exprSwitchStmt.
	VisitExprSwitchStmt(ctx *ExprSwitchStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#exprCaseClause.
	VisitExprCaseClause(ctx *ExprCaseClauseContext) interface{}

	// Visit a parse tree produced by MoxieParser#exprSwitchCase.
	VisitExprSwitchCase(ctx *ExprSwitchCaseContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeSwitchStmt.
	VisitTypeSwitchStmt(ctx *TypeSwitchStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeSwitchGuard.
	VisitTypeSwitchGuard(ctx *TypeSwitchGuardContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeCaseClause.
	VisitTypeCaseClause(ctx *TypeCaseClauseContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeSwitchCase.
	VisitTypeSwitchCase(ctx *TypeSwitchCaseContext) interface{}

	// Visit a parse tree produced by MoxieParser#typeList.
	VisitTypeList(ctx *TypeListContext) interface{}

	// Visit a parse tree produced by MoxieParser#selectStmt.
	VisitSelectStmt(ctx *SelectStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#commClause.
	VisitCommClause(ctx *CommClauseContext) interface{}

	// Visit a parse tree produced by MoxieParser#commCase.
	VisitCommCase(ctx *CommCaseContext) interface{}

	// Visit a parse tree produced by MoxieParser#recvStmt.
	VisitRecvStmt(ctx *RecvStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#forStmt.
	VisitForStmt(ctx *ForStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#forClause.
	VisitForClause(ctx *ForClauseContext) interface{}

	// Visit a parse tree produced by MoxieParser#rangeClause.
	VisitRangeClause(ctx *RangeClauseContext) interface{}

	// Visit a parse tree produced by MoxieParser#goStmt.
	VisitGoStmt(ctx *GoStmtContext) interface{}

	// Visit a parse tree produced by MoxieParser#MultiplicativeExpr.
	VisitMultiplicativeExpr(ctx *MultiplicativeExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#ConcatenationExpr.
	VisitConcatenationExpr(ctx *ConcatenationExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#LogicalOrExpr.
	VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#AdditiveExpr.
	VisitAdditiveExpr(ctx *AdditiveExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#UnaryExpression.
	VisitUnaryExpression(ctx *UnaryExpressionContext) interface{}

	// Visit a parse tree produced by MoxieParser#RelationalExpr.
	VisitRelationalExpr(ctx *RelationalExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#LogicalAndExpr.
	VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#SelectorExpr.
	VisitSelectorExpr(ctx *SelectorExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#TypeAssertionExpr.
	VisitTypeAssertionExpr(ctx *TypeAssertionExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#ConversionExpr.
	VisitConversionExpr(ctx *ConversionExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#PrimaryOperand.
	VisitPrimaryOperand(ctx *PrimaryOperandContext) interface{}

	// Visit a parse tree produced by MoxieParser#SliceExpr.
	VisitSliceExpr(ctx *SliceExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#CallExpr.
	VisitCallExpr(ctx *CallExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#MethodExpression.
	VisitMethodExpression(ctx *MethodExpressionContext) interface{}

	// Visit a parse tree produced by MoxieParser#IndexExpr.
	VisitIndexExpr(ctx *IndexExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#unaryExpr.
	VisitUnaryExpr(ctx *UnaryExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#SimpleConversion.
	VisitSimpleConversion(ctx *SimpleConversionContext) interface{}

	// Visit a parse tree produced by MoxieParser#SliceCastExpr.
	VisitSliceCastExpr(ctx *SliceCastExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#SliceCastEndianExpr.
	VisitSliceCastEndianExpr(ctx *SliceCastEndianExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#SliceCastCopyExpr.
	VisitSliceCastCopyExpr(ctx *SliceCastCopyExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#SliceCastCopyEndianExpr.
	VisitSliceCastCopyEndianExpr(ctx *SliceCastCopyEndianExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#endianness.
	VisitEndianness(ctx *EndiannessContext) interface{}

	// Visit a parse tree produced by MoxieParser#LiteralOperand.
	VisitLiteralOperand(ctx *LiteralOperandContext) interface{}

	// Visit a parse tree produced by MoxieParser#NameOperand.
	VisitNameOperand(ctx *NameOperandContext) interface{}

	// Visit a parse tree produced by MoxieParser#ParenOperand.
	VisitParenOperand(ctx *ParenOperandContext) interface{}

	// Visit a parse tree produced by MoxieParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by MoxieParser#basicLit.
	VisitBasicLit(ctx *BasicLitContext) interface{}

	// Visit a parse tree produced by MoxieParser#string_.
	VisitString_(ctx *String_Context) interface{}

	// Visit a parse tree produced by MoxieParser#operandName.
	VisitOperandName(ctx *OperandNameContext) interface{}

	// Visit a parse tree produced by MoxieParser#qualifiedIdent.
	VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{}

	// Visit a parse tree produced by MoxieParser#compositeLit.
	VisitCompositeLit(ctx *CompositeLitContext) interface{}

	// Visit a parse tree produced by MoxieParser#literalType.
	VisitLiteralType(ctx *LiteralTypeContext) interface{}

	// Visit a parse tree produced by MoxieParser#literalValue.
	VisitLiteralValue(ctx *LiteralValueContext) interface{}

	// Visit a parse tree produced by MoxieParser#elementList.
	VisitElementList(ctx *ElementListContext) interface{}

	// Visit a parse tree produced by MoxieParser#keyedElement.
	VisitKeyedElement(ctx *KeyedElementContext) interface{}

	// Visit a parse tree produced by MoxieParser#key.
	VisitKey(ctx *KeyContext) interface{}

	// Visit a parse tree produced by MoxieParser#element.
	VisitElement(ctx *ElementContext) interface{}

	// Visit a parse tree produced by MoxieParser#functionLit.
	VisitFunctionLit(ctx *FunctionLitContext) interface{}

	// Visit a parse tree produced by MoxieParser#selector.
	VisitSelector(ctx *SelectorContext) interface{}

	// Visit a parse tree produced by MoxieParser#index.
	VisitIndex(ctx *IndexContext) interface{}

	// Visit a parse tree produced by MoxieParser#slice_.
	VisitSlice_(ctx *Slice_Context) interface{}

	// Visit a parse tree produced by MoxieParser#typeAssertion.
	VisitTypeAssertion(ctx *TypeAssertionContext) interface{}

	// Visit a parse tree produced by MoxieParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}

	// Visit a parse tree produced by MoxieParser#methodExpr.
	VisitMethodExpr(ctx *MethodExprContext) interface{}

	// Visit a parse tree produced by MoxieParser#mul_op.
	VisitMul_op(ctx *Mul_opContext) interface{}

	// Visit a parse tree produced by MoxieParser#add_op.
	VisitAdd_op(ctx *Add_opContext) interface{}

	// Visit a parse tree produced by MoxieParser#rel_op.
	VisitRel_op(ctx *Rel_opContext) interface{}

	// Visit a parse tree produced by MoxieParser#unary_op.
	VisitUnary_op(ctx *Unary_opContext) interface{}

	// Visit a parse tree produced by MoxieParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by MoxieParser#identifierList.
	VisitIdentifierList(ctx *IdentifierListContext) interface{}

	// Visit a parse tree produced by MoxieParser#eos.
	VisitEos(ctx *EosContext) interface{}
}
