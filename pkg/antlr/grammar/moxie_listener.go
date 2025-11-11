// Code generated from grammar/Moxie.g4 by ANTLR 4.13.2. DO NOT EDIT.

package antlr // Moxie
import "github.com/antlr4-go/antlr/v4"

// MoxieListener is a complete listener for a parse tree produced by MoxieParser.
type MoxieListener interface {
	antlr.ParseTreeListener

	// EnterSourceFile is called when entering the sourceFile production.
	EnterSourceFile(c *SourceFileContext)

	// EnterPackageClause is called when entering the packageClause production.
	EnterPackageClause(c *PackageClauseContext)

	// EnterImportDecl is called when entering the importDecl production.
	EnterImportDecl(c *ImportDeclContext)

	// EnterImportSpec is called when entering the importSpec production.
	EnterImportSpec(c *ImportSpecContext)

	// EnterTopLevelDecl is called when entering the topLevelDecl production.
	EnterTopLevelDecl(c *TopLevelDeclContext)

	// EnterDeclaration is called when entering the declaration production.
	EnterDeclaration(c *DeclarationContext)

	// EnterConstDecl is called when entering the constDecl production.
	EnterConstDecl(c *ConstDeclContext)

	// EnterConstSpec is called when entering the constSpec production.
	EnterConstSpec(c *ConstSpecContext)

	// EnterTypeDecl is called when entering the typeDecl production.
	EnterTypeDecl(c *TypeDeclContext)

	// EnterTypeAlias is called when entering the TypeAlias production.
	EnterTypeAlias(c *TypeAliasContext)

	// EnterTypeDef is called when entering the TypeDef production.
	EnterTypeDef(c *TypeDefContext)

	// EnterTypeParameters is called when entering the typeParameters production.
	EnterTypeParameters(c *TypeParametersContext)

	// EnterTypeParameterDecl is called when entering the typeParameterDecl production.
	EnterTypeParameterDecl(c *TypeParameterDeclContext)

	// EnterTypeConstraint is called when entering the typeConstraint production.
	EnterTypeConstraint(c *TypeConstraintContext)

	// EnterVarDecl is called when entering the varDecl production.
	EnterVarDecl(c *VarDeclContext)

	// EnterVarSpec is called when entering the varSpec production.
	EnterVarSpec(c *VarSpecContext)

	// EnterFunctionDecl is called when entering the functionDecl production.
	EnterFunctionDecl(c *FunctionDeclContext)

	// EnterMethodDecl is called when entering the methodDecl production.
	EnterMethodDecl(c *MethodDeclContext)

	// EnterReceiver is called when entering the receiver production.
	EnterReceiver(c *ReceiverContext)

	// EnterNamedType is called when entering the NamedType production.
	EnterNamedType(c *NamedTypeContext)

	// EnterTypeLiteral is called when entering the TypeLiteral production.
	EnterTypeLiteral(c *TypeLiteralContext)

	// EnterParenType is called when entering the ParenType production.
	EnterParenType(c *ParenTypeContext)

	// EnterConstType is called when entering the ConstType production.
	EnterConstType(c *ConstTypeContext)

	// EnterTypeName is called when entering the typeName production.
	EnterTypeName(c *TypeNameContext)

	// EnterTypeArgs is called when entering the typeArgs production.
	EnterTypeArgs(c *TypeArgsContext)

	// EnterTypeLit is called when entering the typeLit production.
	EnterTypeLit(c *TypeLitContext)

	// EnterArrayType is called when entering the arrayType production.
	EnterArrayType(c *ArrayTypeContext)

	// EnterArrayLength is called when entering the arrayLength production.
	EnterArrayLength(c *ArrayLengthContext)

	// EnterElementType is called when entering the elementType production.
	EnterElementType(c *ElementTypeContext)

	// EnterSliceType is called when entering the sliceType production.
	EnterSliceType(c *SliceTypeContext)

	// EnterStructType is called when entering the structType production.
	EnterStructType(c *StructTypeContext)

	// EnterFieldDecl is called when entering the fieldDecl production.
	EnterFieldDecl(c *FieldDeclContext)

	// EnterEmbeddedField is called when entering the embeddedField production.
	EnterEmbeddedField(c *EmbeddedFieldContext)

	// EnterTag_ is called when entering the tag_ production.
	EnterTag_(c *Tag_Context)

	// EnterPointerType is called when entering the pointerType production.
	EnterPointerType(c *PointerTypeContext)

	// EnterFunctionType is called when entering the functionType production.
	EnterFunctionType(c *FunctionTypeContext)

	// EnterSignature is called when entering the signature production.
	EnterSignature(c *SignatureContext)

	// EnterResult is called when entering the result production.
	EnterResult(c *ResultContext)

	// EnterParameters is called when entering the parameters production.
	EnterParameters(c *ParametersContext)

	// EnterParameterDecl is called when entering the parameterDecl production.
	EnterParameterDecl(c *ParameterDeclContext)

	// EnterInterfaceType is called when entering the interfaceType production.
	EnterInterfaceType(c *InterfaceTypeContext)

	// EnterInterfaceElem is called when entering the interfaceElem production.
	EnterInterfaceElem(c *InterfaceElemContext)

	// EnterMethodElem is called when entering the methodElem production.
	EnterMethodElem(c *MethodElemContext)

	// EnterTypeElem is called when entering the typeElem production.
	EnterTypeElem(c *TypeElemContext)

	// EnterTypeTerm is called when entering the typeTerm production.
	EnterTypeTerm(c *TypeTermContext)

	// EnterMapType is called when entering the mapType production.
	EnterMapType(c *MapTypeContext)

	// EnterSendRecvChan is called when entering the SendRecvChan production.
	EnterSendRecvChan(c *SendRecvChanContext)

	// EnterRecvOnlyChan is called when entering the RecvOnlyChan production.
	EnterRecvOnlyChan(c *RecvOnlyChanContext)

	// EnterSendRecvChanCompat is called when entering the SendRecvChanCompat production.
	EnterSendRecvChanCompat(c *SendRecvChanCompatContext)

	// EnterRecvOnlyChanCompat is called when entering the RecvOnlyChanCompat production.
	EnterRecvOnlyChanCompat(c *RecvOnlyChanCompatContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterStatementList is called when entering the statementList production.
	EnterStatementList(c *StatementListContext)

	// EnterDeclStmt is called when entering the DeclStmt production.
	EnterDeclStmt(c *DeclStmtContext)

	// EnterSimpleStatement is called when entering the SimpleStatement production.
	EnterSimpleStatement(c *SimpleStatementContext)

	// EnterLabeledStatement is called when entering the LabeledStatement production.
	EnterLabeledStatement(c *LabeledStatementContext)

	// EnterGoStatement is called when entering the GoStatement production.
	EnterGoStatement(c *GoStatementContext)

	// EnterReturnStatement is called when entering the ReturnStatement production.
	EnterReturnStatement(c *ReturnStatementContext)

	// EnterBreakStatement is called when entering the BreakStatement production.
	EnterBreakStatement(c *BreakStatementContext)

	// EnterContinueStatement is called when entering the ContinueStatement production.
	EnterContinueStatement(c *ContinueStatementContext)

	// EnterGotoStatement is called when entering the GotoStatement production.
	EnterGotoStatement(c *GotoStatementContext)

	// EnterFallthroughStatement is called when entering the FallthroughStatement production.
	EnterFallthroughStatement(c *FallthroughStatementContext)

	// EnterBlockStatement is called when entering the BlockStatement production.
	EnterBlockStatement(c *BlockStatementContext)

	// EnterIfStatement is called when entering the IfStatement production.
	EnterIfStatement(c *IfStatementContext)

	// EnterSwitchStatement is called when entering the SwitchStatement production.
	EnterSwitchStatement(c *SwitchStatementContext)

	// EnterSelectStatement is called when entering the SelectStatement production.
	EnterSelectStatement(c *SelectStatementContext)

	// EnterForStatement is called when entering the ForStatement production.
	EnterForStatement(c *ForStatementContext)

	// EnterDeferStatement is called when entering the DeferStatement production.
	EnterDeferStatement(c *DeferStatementContext)

	// EnterSimpleStmt is called when entering the simpleStmt production.
	EnterSimpleStmt(c *SimpleStmtContext)

	// EnterExpressionStmt is called when entering the expressionStmt production.
	EnterExpressionStmt(c *ExpressionStmtContext)

	// EnterSendStmt is called when entering the sendStmt production.
	EnterSendStmt(c *SendStmtContext)

	// EnterIncDecStmt is called when entering the incDecStmt production.
	EnterIncDecStmt(c *IncDecStmtContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterAssign_op is called when entering the assign_op production.
	EnterAssign_op(c *Assign_opContext)

	// EnterShortVarDecl is called when entering the shortVarDecl production.
	EnterShortVarDecl(c *ShortVarDeclContext)

	// EnterLabeledStmt is called when entering the labeledStmt production.
	EnterLabeledStmt(c *LabeledStmtContext)

	// EnterReturnStmt is called when entering the returnStmt production.
	EnterReturnStmt(c *ReturnStmtContext)

	// EnterBreakStmt is called when entering the breakStmt production.
	EnterBreakStmt(c *BreakStmtContext)

	// EnterContinueStmt is called when entering the continueStmt production.
	EnterContinueStmt(c *ContinueStmtContext)

	// EnterGotoStmt is called when entering the gotoStmt production.
	EnterGotoStmt(c *GotoStmtContext)

	// EnterFallthroughStmt is called when entering the fallthroughStmt production.
	EnterFallthroughStmt(c *FallthroughStmtContext)

	// EnterDeferStmt is called when entering the deferStmt production.
	EnterDeferStmt(c *DeferStmtContext)

	// EnterIfStmt is called when entering the ifStmt production.
	EnterIfStmt(c *IfStmtContext)

	// EnterSwitchStmt is called when entering the switchStmt production.
	EnterSwitchStmt(c *SwitchStmtContext)

	// EnterExprSwitchStmt is called when entering the exprSwitchStmt production.
	EnterExprSwitchStmt(c *ExprSwitchStmtContext)

	// EnterExprCaseClause is called when entering the exprCaseClause production.
	EnterExprCaseClause(c *ExprCaseClauseContext)

	// EnterExprSwitchCase is called when entering the exprSwitchCase production.
	EnterExprSwitchCase(c *ExprSwitchCaseContext)

	// EnterTypeSwitchStmt is called when entering the typeSwitchStmt production.
	EnterTypeSwitchStmt(c *TypeSwitchStmtContext)

	// EnterTypeSwitchGuard is called when entering the typeSwitchGuard production.
	EnterTypeSwitchGuard(c *TypeSwitchGuardContext)

	// EnterTypeCaseClause is called when entering the typeCaseClause production.
	EnterTypeCaseClause(c *TypeCaseClauseContext)

	// EnterTypeSwitchCase is called when entering the typeSwitchCase production.
	EnterTypeSwitchCase(c *TypeSwitchCaseContext)

	// EnterTypeList is called when entering the typeList production.
	EnterTypeList(c *TypeListContext)

	// EnterSelectStmt is called when entering the selectStmt production.
	EnterSelectStmt(c *SelectStmtContext)

	// EnterCommClause is called when entering the commClause production.
	EnterCommClause(c *CommClauseContext)

	// EnterCommCase is called when entering the commCase production.
	EnterCommCase(c *CommCaseContext)

	// EnterRecvStmt is called when entering the recvStmt production.
	EnterRecvStmt(c *RecvStmtContext)

	// EnterForStmt is called when entering the forStmt production.
	EnterForStmt(c *ForStmtContext)

	// EnterForClause is called when entering the forClause production.
	EnterForClause(c *ForClauseContext)

	// EnterRangeClause is called when entering the rangeClause production.
	EnterRangeClause(c *RangeClauseContext)

	// EnterGoStmt is called when entering the goStmt production.
	EnterGoStmt(c *GoStmtContext)

	// EnterMultiplicativeExpr is called when entering the MultiplicativeExpr production.
	EnterMultiplicativeExpr(c *MultiplicativeExprContext)

	// EnterConcatenationExpr is called when entering the ConcatenationExpr production.
	EnterConcatenationExpr(c *ConcatenationExprContext)

	// EnterLogicalOrExpr is called when entering the LogicalOrExpr production.
	EnterLogicalOrExpr(c *LogicalOrExprContext)

	// EnterAdditiveExpr is called when entering the AdditiveExpr production.
	EnterAdditiveExpr(c *AdditiveExprContext)

	// EnterUnaryExpression is called when entering the UnaryExpression production.
	EnterUnaryExpression(c *UnaryExpressionContext)

	// EnterRelationalExpr is called when entering the RelationalExpr production.
	EnterRelationalExpr(c *RelationalExprContext)

	// EnterLogicalAndExpr is called when entering the LogicalAndExpr production.
	EnterLogicalAndExpr(c *LogicalAndExprContext)

	// EnterSelectorExpr is called when entering the SelectorExpr production.
	EnterSelectorExpr(c *SelectorExprContext)

	// EnterTypeAssertionExpr is called when entering the TypeAssertionExpr production.
	EnterTypeAssertionExpr(c *TypeAssertionExprContext)

	// EnterConversionExpr is called when entering the ConversionExpr production.
	EnterConversionExpr(c *ConversionExprContext)

	// EnterPrimaryOperand is called when entering the PrimaryOperand production.
	EnterPrimaryOperand(c *PrimaryOperandContext)

	// EnterSliceExpr is called when entering the SliceExpr production.
	EnterSliceExpr(c *SliceExprContext)

	// EnterCallExpr is called when entering the CallExpr production.
	EnterCallExpr(c *CallExprContext)

	// EnterMethodExpression is called when entering the MethodExpression production.
	EnterMethodExpression(c *MethodExpressionContext)

	// EnterIndexExpr is called when entering the IndexExpr production.
	EnterIndexExpr(c *IndexExprContext)

	// EnterUnaryExpr is called when entering the unaryExpr production.
	EnterUnaryExpr(c *UnaryExprContext)

	// EnterSimpleConversion is called when entering the SimpleConversion production.
	EnterSimpleConversion(c *SimpleConversionContext)

	// EnterSliceCastExpr is called when entering the SliceCastExpr production.
	EnterSliceCastExpr(c *SliceCastExprContext)

	// EnterSliceCastEndianExpr is called when entering the SliceCastEndianExpr production.
	EnterSliceCastEndianExpr(c *SliceCastEndianExprContext)

	// EnterSliceCastCopyExpr is called when entering the SliceCastCopyExpr production.
	EnterSliceCastCopyExpr(c *SliceCastCopyExprContext)

	// EnterSliceCastCopyEndianExpr is called when entering the SliceCastCopyEndianExpr production.
	EnterSliceCastCopyEndianExpr(c *SliceCastCopyEndianExprContext)

	// EnterEndianness is called when entering the endianness production.
	EnterEndianness(c *EndiannessContext)

	// EnterLiteralOperand is called when entering the LiteralOperand production.
	EnterLiteralOperand(c *LiteralOperandContext)

	// EnterNameOperand is called when entering the NameOperand production.
	EnterNameOperand(c *NameOperandContext)

	// EnterParenOperand is called when entering the ParenOperand production.
	EnterParenOperand(c *ParenOperandContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterBasicLit is called when entering the basicLit production.
	EnterBasicLit(c *BasicLitContext)

	// EnterString_ is called when entering the string_ production.
	EnterString_(c *String_Context)

	// EnterOperandName is called when entering the operandName production.
	EnterOperandName(c *OperandNameContext)

	// EnterQualifiedIdent is called when entering the qualifiedIdent production.
	EnterQualifiedIdent(c *QualifiedIdentContext)

	// EnterCompositeLit is called when entering the compositeLit production.
	EnterCompositeLit(c *CompositeLitContext)

	// EnterLiteralType is called when entering the literalType production.
	EnterLiteralType(c *LiteralTypeContext)

	// EnterLiteralValue is called when entering the literalValue production.
	EnterLiteralValue(c *LiteralValueContext)

	// EnterElementList is called when entering the elementList production.
	EnterElementList(c *ElementListContext)

	// EnterKeyedElement is called when entering the keyedElement production.
	EnterKeyedElement(c *KeyedElementContext)

	// EnterKey is called when entering the key production.
	EnterKey(c *KeyContext)

	// EnterElement is called when entering the element production.
	EnterElement(c *ElementContext)

	// EnterFunctionLit is called when entering the functionLit production.
	EnterFunctionLit(c *FunctionLitContext)

	// EnterSelector is called when entering the selector production.
	EnterSelector(c *SelectorContext)

	// EnterIndex is called when entering the index production.
	EnterIndex(c *IndexContext)

	// EnterSlice_ is called when entering the slice_ production.
	EnterSlice_(c *Slice_Context)

	// EnterTypeAssertion is called when entering the typeAssertion production.
	EnterTypeAssertion(c *TypeAssertionContext)

	// EnterArguments is called when entering the arguments production.
	EnterArguments(c *ArgumentsContext)

	// EnterMethodExpr is called when entering the methodExpr production.
	EnterMethodExpr(c *MethodExprContext)

	// EnterMul_op is called when entering the mul_op production.
	EnterMul_op(c *Mul_opContext)

	// EnterAdd_op is called when entering the add_op production.
	EnterAdd_op(c *Add_opContext)

	// EnterRel_op is called when entering the rel_op production.
	EnterRel_op(c *Rel_opContext)

	// EnterUnary_op is called when entering the unary_op production.
	EnterUnary_op(c *Unary_opContext)

	// EnterExpressionList is called when entering the expressionList production.
	EnterExpressionList(c *ExpressionListContext)

	// EnterIdentifierList is called when entering the identifierList production.
	EnterIdentifierList(c *IdentifierListContext)

	// EnterEos is called when entering the eos production.
	EnterEos(c *EosContext)

	// ExitSourceFile is called when exiting the sourceFile production.
	ExitSourceFile(c *SourceFileContext)

	// ExitPackageClause is called when exiting the packageClause production.
	ExitPackageClause(c *PackageClauseContext)

	// ExitImportDecl is called when exiting the importDecl production.
	ExitImportDecl(c *ImportDeclContext)

	// ExitImportSpec is called when exiting the importSpec production.
	ExitImportSpec(c *ImportSpecContext)

	// ExitTopLevelDecl is called when exiting the topLevelDecl production.
	ExitTopLevelDecl(c *TopLevelDeclContext)

	// ExitDeclaration is called when exiting the declaration production.
	ExitDeclaration(c *DeclarationContext)

	// ExitConstDecl is called when exiting the constDecl production.
	ExitConstDecl(c *ConstDeclContext)

	// ExitConstSpec is called when exiting the constSpec production.
	ExitConstSpec(c *ConstSpecContext)

	// ExitTypeDecl is called when exiting the typeDecl production.
	ExitTypeDecl(c *TypeDeclContext)

	// ExitTypeAlias is called when exiting the TypeAlias production.
	ExitTypeAlias(c *TypeAliasContext)

	// ExitTypeDef is called when exiting the TypeDef production.
	ExitTypeDef(c *TypeDefContext)

	// ExitTypeParameters is called when exiting the typeParameters production.
	ExitTypeParameters(c *TypeParametersContext)

	// ExitTypeParameterDecl is called when exiting the typeParameterDecl production.
	ExitTypeParameterDecl(c *TypeParameterDeclContext)

	// ExitTypeConstraint is called when exiting the typeConstraint production.
	ExitTypeConstraint(c *TypeConstraintContext)

	// ExitVarDecl is called when exiting the varDecl production.
	ExitVarDecl(c *VarDeclContext)

	// ExitVarSpec is called when exiting the varSpec production.
	ExitVarSpec(c *VarSpecContext)

	// ExitFunctionDecl is called when exiting the functionDecl production.
	ExitFunctionDecl(c *FunctionDeclContext)

	// ExitMethodDecl is called when exiting the methodDecl production.
	ExitMethodDecl(c *MethodDeclContext)

	// ExitReceiver is called when exiting the receiver production.
	ExitReceiver(c *ReceiverContext)

	// ExitNamedType is called when exiting the NamedType production.
	ExitNamedType(c *NamedTypeContext)

	// ExitTypeLiteral is called when exiting the TypeLiteral production.
	ExitTypeLiteral(c *TypeLiteralContext)

	// ExitParenType is called when exiting the ParenType production.
	ExitParenType(c *ParenTypeContext)

	// ExitConstType is called when exiting the ConstType production.
	ExitConstType(c *ConstTypeContext)

	// ExitTypeName is called when exiting the typeName production.
	ExitTypeName(c *TypeNameContext)

	// ExitTypeArgs is called when exiting the typeArgs production.
	ExitTypeArgs(c *TypeArgsContext)

	// ExitTypeLit is called when exiting the typeLit production.
	ExitTypeLit(c *TypeLitContext)

	// ExitArrayType is called when exiting the arrayType production.
	ExitArrayType(c *ArrayTypeContext)

	// ExitArrayLength is called when exiting the arrayLength production.
	ExitArrayLength(c *ArrayLengthContext)

	// ExitElementType is called when exiting the elementType production.
	ExitElementType(c *ElementTypeContext)

	// ExitSliceType is called when exiting the sliceType production.
	ExitSliceType(c *SliceTypeContext)

	// ExitStructType is called when exiting the structType production.
	ExitStructType(c *StructTypeContext)

	// ExitFieldDecl is called when exiting the fieldDecl production.
	ExitFieldDecl(c *FieldDeclContext)

	// ExitEmbeddedField is called when exiting the embeddedField production.
	ExitEmbeddedField(c *EmbeddedFieldContext)

	// ExitTag_ is called when exiting the tag_ production.
	ExitTag_(c *Tag_Context)

	// ExitPointerType is called when exiting the pointerType production.
	ExitPointerType(c *PointerTypeContext)

	// ExitFunctionType is called when exiting the functionType production.
	ExitFunctionType(c *FunctionTypeContext)

	// ExitSignature is called when exiting the signature production.
	ExitSignature(c *SignatureContext)

	// ExitResult is called when exiting the result production.
	ExitResult(c *ResultContext)

	// ExitParameters is called when exiting the parameters production.
	ExitParameters(c *ParametersContext)

	// ExitParameterDecl is called when exiting the parameterDecl production.
	ExitParameterDecl(c *ParameterDeclContext)

	// ExitInterfaceType is called when exiting the interfaceType production.
	ExitInterfaceType(c *InterfaceTypeContext)

	// ExitInterfaceElem is called when exiting the interfaceElem production.
	ExitInterfaceElem(c *InterfaceElemContext)

	// ExitMethodElem is called when exiting the methodElem production.
	ExitMethodElem(c *MethodElemContext)

	// ExitTypeElem is called when exiting the typeElem production.
	ExitTypeElem(c *TypeElemContext)

	// ExitTypeTerm is called when exiting the typeTerm production.
	ExitTypeTerm(c *TypeTermContext)

	// ExitMapType is called when exiting the mapType production.
	ExitMapType(c *MapTypeContext)

	// ExitSendRecvChan is called when exiting the SendRecvChan production.
	ExitSendRecvChan(c *SendRecvChanContext)

	// ExitRecvOnlyChan is called when exiting the RecvOnlyChan production.
	ExitRecvOnlyChan(c *RecvOnlyChanContext)

	// ExitSendRecvChanCompat is called when exiting the SendRecvChanCompat production.
	ExitSendRecvChanCompat(c *SendRecvChanCompatContext)

	// ExitRecvOnlyChanCompat is called when exiting the RecvOnlyChanCompat production.
	ExitRecvOnlyChanCompat(c *RecvOnlyChanCompatContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitStatementList is called when exiting the statementList production.
	ExitStatementList(c *StatementListContext)

	// ExitDeclStmt is called when exiting the DeclStmt production.
	ExitDeclStmt(c *DeclStmtContext)

	// ExitSimpleStatement is called when exiting the SimpleStatement production.
	ExitSimpleStatement(c *SimpleStatementContext)

	// ExitLabeledStatement is called when exiting the LabeledStatement production.
	ExitLabeledStatement(c *LabeledStatementContext)

	// ExitGoStatement is called when exiting the GoStatement production.
	ExitGoStatement(c *GoStatementContext)

	// ExitReturnStatement is called when exiting the ReturnStatement production.
	ExitReturnStatement(c *ReturnStatementContext)

	// ExitBreakStatement is called when exiting the BreakStatement production.
	ExitBreakStatement(c *BreakStatementContext)

	// ExitContinueStatement is called when exiting the ContinueStatement production.
	ExitContinueStatement(c *ContinueStatementContext)

	// ExitGotoStatement is called when exiting the GotoStatement production.
	ExitGotoStatement(c *GotoStatementContext)

	// ExitFallthroughStatement is called when exiting the FallthroughStatement production.
	ExitFallthroughStatement(c *FallthroughStatementContext)

	// ExitBlockStatement is called when exiting the BlockStatement production.
	ExitBlockStatement(c *BlockStatementContext)

	// ExitIfStatement is called when exiting the IfStatement production.
	ExitIfStatement(c *IfStatementContext)

	// ExitSwitchStatement is called when exiting the SwitchStatement production.
	ExitSwitchStatement(c *SwitchStatementContext)

	// ExitSelectStatement is called when exiting the SelectStatement production.
	ExitSelectStatement(c *SelectStatementContext)

	// ExitForStatement is called when exiting the ForStatement production.
	ExitForStatement(c *ForStatementContext)

	// ExitDeferStatement is called when exiting the DeferStatement production.
	ExitDeferStatement(c *DeferStatementContext)

	// ExitSimpleStmt is called when exiting the simpleStmt production.
	ExitSimpleStmt(c *SimpleStmtContext)

	// ExitExpressionStmt is called when exiting the expressionStmt production.
	ExitExpressionStmt(c *ExpressionStmtContext)

	// ExitSendStmt is called when exiting the sendStmt production.
	ExitSendStmt(c *SendStmtContext)

	// ExitIncDecStmt is called when exiting the incDecStmt production.
	ExitIncDecStmt(c *IncDecStmtContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitAssign_op is called when exiting the assign_op production.
	ExitAssign_op(c *Assign_opContext)

	// ExitShortVarDecl is called when exiting the shortVarDecl production.
	ExitShortVarDecl(c *ShortVarDeclContext)

	// ExitLabeledStmt is called when exiting the labeledStmt production.
	ExitLabeledStmt(c *LabeledStmtContext)

	// ExitReturnStmt is called when exiting the returnStmt production.
	ExitReturnStmt(c *ReturnStmtContext)

	// ExitBreakStmt is called when exiting the breakStmt production.
	ExitBreakStmt(c *BreakStmtContext)

	// ExitContinueStmt is called when exiting the continueStmt production.
	ExitContinueStmt(c *ContinueStmtContext)

	// ExitGotoStmt is called when exiting the gotoStmt production.
	ExitGotoStmt(c *GotoStmtContext)

	// ExitFallthroughStmt is called when exiting the fallthroughStmt production.
	ExitFallthroughStmt(c *FallthroughStmtContext)

	// ExitDeferStmt is called when exiting the deferStmt production.
	ExitDeferStmt(c *DeferStmtContext)

	// ExitIfStmt is called when exiting the ifStmt production.
	ExitIfStmt(c *IfStmtContext)

	// ExitSwitchStmt is called when exiting the switchStmt production.
	ExitSwitchStmt(c *SwitchStmtContext)

	// ExitExprSwitchStmt is called when exiting the exprSwitchStmt production.
	ExitExprSwitchStmt(c *ExprSwitchStmtContext)

	// ExitExprCaseClause is called when exiting the exprCaseClause production.
	ExitExprCaseClause(c *ExprCaseClauseContext)

	// ExitExprSwitchCase is called when exiting the exprSwitchCase production.
	ExitExprSwitchCase(c *ExprSwitchCaseContext)

	// ExitTypeSwitchStmt is called when exiting the typeSwitchStmt production.
	ExitTypeSwitchStmt(c *TypeSwitchStmtContext)

	// ExitTypeSwitchGuard is called when exiting the typeSwitchGuard production.
	ExitTypeSwitchGuard(c *TypeSwitchGuardContext)

	// ExitTypeCaseClause is called when exiting the typeCaseClause production.
	ExitTypeCaseClause(c *TypeCaseClauseContext)

	// ExitTypeSwitchCase is called when exiting the typeSwitchCase production.
	ExitTypeSwitchCase(c *TypeSwitchCaseContext)

	// ExitTypeList is called when exiting the typeList production.
	ExitTypeList(c *TypeListContext)

	// ExitSelectStmt is called when exiting the selectStmt production.
	ExitSelectStmt(c *SelectStmtContext)

	// ExitCommClause is called when exiting the commClause production.
	ExitCommClause(c *CommClauseContext)

	// ExitCommCase is called when exiting the commCase production.
	ExitCommCase(c *CommCaseContext)

	// ExitRecvStmt is called when exiting the recvStmt production.
	ExitRecvStmt(c *RecvStmtContext)

	// ExitForStmt is called when exiting the forStmt production.
	ExitForStmt(c *ForStmtContext)

	// ExitForClause is called when exiting the forClause production.
	ExitForClause(c *ForClauseContext)

	// ExitRangeClause is called when exiting the rangeClause production.
	ExitRangeClause(c *RangeClauseContext)

	// ExitGoStmt is called when exiting the goStmt production.
	ExitGoStmt(c *GoStmtContext)

	// ExitMultiplicativeExpr is called when exiting the MultiplicativeExpr production.
	ExitMultiplicativeExpr(c *MultiplicativeExprContext)

	// ExitConcatenationExpr is called when exiting the ConcatenationExpr production.
	ExitConcatenationExpr(c *ConcatenationExprContext)

	// ExitLogicalOrExpr is called when exiting the LogicalOrExpr production.
	ExitLogicalOrExpr(c *LogicalOrExprContext)

	// ExitAdditiveExpr is called when exiting the AdditiveExpr production.
	ExitAdditiveExpr(c *AdditiveExprContext)

	// ExitUnaryExpression is called when exiting the UnaryExpression production.
	ExitUnaryExpression(c *UnaryExpressionContext)

	// ExitRelationalExpr is called when exiting the RelationalExpr production.
	ExitRelationalExpr(c *RelationalExprContext)

	// ExitLogicalAndExpr is called when exiting the LogicalAndExpr production.
	ExitLogicalAndExpr(c *LogicalAndExprContext)

	// ExitSelectorExpr is called when exiting the SelectorExpr production.
	ExitSelectorExpr(c *SelectorExprContext)

	// ExitTypeAssertionExpr is called when exiting the TypeAssertionExpr production.
	ExitTypeAssertionExpr(c *TypeAssertionExprContext)

	// ExitConversionExpr is called when exiting the ConversionExpr production.
	ExitConversionExpr(c *ConversionExprContext)

	// ExitPrimaryOperand is called when exiting the PrimaryOperand production.
	ExitPrimaryOperand(c *PrimaryOperandContext)

	// ExitSliceExpr is called when exiting the SliceExpr production.
	ExitSliceExpr(c *SliceExprContext)

	// ExitCallExpr is called when exiting the CallExpr production.
	ExitCallExpr(c *CallExprContext)

	// ExitMethodExpression is called when exiting the MethodExpression production.
	ExitMethodExpression(c *MethodExpressionContext)

	// ExitIndexExpr is called when exiting the IndexExpr production.
	ExitIndexExpr(c *IndexExprContext)

	// ExitUnaryExpr is called when exiting the unaryExpr production.
	ExitUnaryExpr(c *UnaryExprContext)

	// ExitSimpleConversion is called when exiting the SimpleConversion production.
	ExitSimpleConversion(c *SimpleConversionContext)

	// ExitSliceCastExpr is called when exiting the SliceCastExpr production.
	ExitSliceCastExpr(c *SliceCastExprContext)

	// ExitSliceCastEndianExpr is called when exiting the SliceCastEndianExpr production.
	ExitSliceCastEndianExpr(c *SliceCastEndianExprContext)

	// ExitSliceCastCopyExpr is called when exiting the SliceCastCopyExpr production.
	ExitSliceCastCopyExpr(c *SliceCastCopyExprContext)

	// ExitSliceCastCopyEndianExpr is called when exiting the SliceCastCopyEndianExpr production.
	ExitSliceCastCopyEndianExpr(c *SliceCastCopyEndianExprContext)

	// ExitEndianness is called when exiting the endianness production.
	ExitEndianness(c *EndiannessContext)

	// ExitLiteralOperand is called when exiting the LiteralOperand production.
	ExitLiteralOperand(c *LiteralOperandContext)

	// ExitNameOperand is called when exiting the NameOperand production.
	ExitNameOperand(c *NameOperandContext)

	// ExitParenOperand is called when exiting the ParenOperand production.
	ExitParenOperand(c *ParenOperandContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitBasicLit is called when exiting the basicLit production.
	ExitBasicLit(c *BasicLitContext)

	// ExitString_ is called when exiting the string_ production.
	ExitString_(c *String_Context)

	// ExitOperandName is called when exiting the operandName production.
	ExitOperandName(c *OperandNameContext)

	// ExitQualifiedIdent is called when exiting the qualifiedIdent production.
	ExitQualifiedIdent(c *QualifiedIdentContext)

	// ExitCompositeLit is called when exiting the compositeLit production.
	ExitCompositeLit(c *CompositeLitContext)

	// ExitLiteralType is called when exiting the literalType production.
	ExitLiteralType(c *LiteralTypeContext)

	// ExitLiteralValue is called when exiting the literalValue production.
	ExitLiteralValue(c *LiteralValueContext)

	// ExitElementList is called when exiting the elementList production.
	ExitElementList(c *ElementListContext)

	// ExitKeyedElement is called when exiting the keyedElement production.
	ExitKeyedElement(c *KeyedElementContext)

	// ExitKey is called when exiting the key production.
	ExitKey(c *KeyContext)

	// ExitElement is called when exiting the element production.
	ExitElement(c *ElementContext)

	// ExitFunctionLit is called when exiting the functionLit production.
	ExitFunctionLit(c *FunctionLitContext)

	// ExitSelector is called when exiting the selector production.
	ExitSelector(c *SelectorContext)

	// ExitIndex is called when exiting the index production.
	ExitIndex(c *IndexContext)

	// ExitSlice_ is called when exiting the slice_ production.
	ExitSlice_(c *Slice_Context)

	// ExitTypeAssertion is called when exiting the typeAssertion production.
	ExitTypeAssertion(c *TypeAssertionContext)

	// ExitArguments is called when exiting the arguments production.
	ExitArguments(c *ArgumentsContext)

	// ExitMethodExpr is called when exiting the methodExpr production.
	ExitMethodExpr(c *MethodExprContext)

	// ExitMul_op is called when exiting the mul_op production.
	ExitMul_op(c *Mul_opContext)

	// ExitAdd_op is called when exiting the add_op production.
	ExitAdd_op(c *Add_opContext)

	// ExitRel_op is called when exiting the rel_op production.
	ExitRel_op(c *Rel_opContext)

	// ExitUnary_op is called when exiting the unary_op production.
	ExitUnary_op(c *Unary_opContext)

	// ExitExpressionList is called when exiting the expressionList production.
	ExitExpressionList(c *ExpressionListContext)

	// ExitIdentifierList is called when exiting the identifierList production.
	ExitIdentifierList(c *IdentifierListContext)

	// ExitEos is called when exiting the eos production.
	ExitEos(c *EosContext)
}
