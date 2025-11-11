// Code generated from grammar/Moxie.g4 by ANTLR 4.13.2. DO NOT EDIT.

package antlr // Moxie
import "github.com/antlr4-go/antlr/v4"

// BaseMoxieListener is a complete listener for a parse tree produced by MoxieParser.
type BaseMoxieListener struct{}

var _ MoxieListener = &BaseMoxieListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseMoxieListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseMoxieListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseMoxieListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseMoxieListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSourceFile is called when production sourceFile is entered.
func (s *BaseMoxieListener) EnterSourceFile(ctx *SourceFileContext) {}

// ExitSourceFile is called when production sourceFile is exited.
func (s *BaseMoxieListener) ExitSourceFile(ctx *SourceFileContext) {}

// EnterPackageClause is called when production packageClause is entered.
func (s *BaseMoxieListener) EnterPackageClause(ctx *PackageClauseContext) {}

// ExitPackageClause is called when production packageClause is exited.
func (s *BaseMoxieListener) ExitPackageClause(ctx *PackageClauseContext) {}

// EnterImportDecl is called when production importDecl is entered.
func (s *BaseMoxieListener) EnterImportDecl(ctx *ImportDeclContext) {}

// ExitImportDecl is called when production importDecl is exited.
func (s *BaseMoxieListener) ExitImportDecl(ctx *ImportDeclContext) {}

// EnterImportSpec is called when production importSpec is entered.
func (s *BaseMoxieListener) EnterImportSpec(ctx *ImportSpecContext) {}

// ExitImportSpec is called when production importSpec is exited.
func (s *BaseMoxieListener) ExitImportSpec(ctx *ImportSpecContext) {}

// EnterTopLevelDecl is called when production topLevelDecl is entered.
func (s *BaseMoxieListener) EnterTopLevelDecl(ctx *TopLevelDeclContext) {}

// ExitTopLevelDecl is called when production topLevelDecl is exited.
func (s *BaseMoxieListener) ExitTopLevelDecl(ctx *TopLevelDeclContext) {}

// EnterDeclaration is called when production declaration is entered.
func (s *BaseMoxieListener) EnterDeclaration(ctx *DeclarationContext) {}

// ExitDeclaration is called when production declaration is exited.
func (s *BaseMoxieListener) ExitDeclaration(ctx *DeclarationContext) {}

// EnterConstDecl is called when production constDecl is entered.
func (s *BaseMoxieListener) EnterConstDecl(ctx *ConstDeclContext) {}

// ExitConstDecl is called when production constDecl is exited.
func (s *BaseMoxieListener) ExitConstDecl(ctx *ConstDeclContext) {}

// EnterConstSpec is called when production constSpec is entered.
func (s *BaseMoxieListener) EnterConstSpec(ctx *ConstSpecContext) {}

// ExitConstSpec is called when production constSpec is exited.
func (s *BaseMoxieListener) ExitConstSpec(ctx *ConstSpecContext) {}

// EnterTypeDecl is called when production typeDecl is entered.
func (s *BaseMoxieListener) EnterTypeDecl(ctx *TypeDeclContext) {}

// ExitTypeDecl is called when production typeDecl is exited.
func (s *BaseMoxieListener) ExitTypeDecl(ctx *TypeDeclContext) {}

// EnterTypeAlias is called when production TypeAlias is entered.
func (s *BaseMoxieListener) EnterTypeAlias(ctx *TypeAliasContext) {}

// ExitTypeAlias is called when production TypeAlias is exited.
func (s *BaseMoxieListener) ExitTypeAlias(ctx *TypeAliasContext) {}

// EnterTypeDef is called when production TypeDef is entered.
func (s *BaseMoxieListener) EnterTypeDef(ctx *TypeDefContext) {}

// ExitTypeDef is called when production TypeDef is exited.
func (s *BaseMoxieListener) ExitTypeDef(ctx *TypeDefContext) {}

// EnterTypeParameters is called when production typeParameters is entered.
func (s *BaseMoxieListener) EnterTypeParameters(ctx *TypeParametersContext) {}

// ExitTypeParameters is called when production typeParameters is exited.
func (s *BaseMoxieListener) ExitTypeParameters(ctx *TypeParametersContext) {}

// EnterTypeParameterDecl is called when production typeParameterDecl is entered.
func (s *BaseMoxieListener) EnterTypeParameterDecl(ctx *TypeParameterDeclContext) {}

// ExitTypeParameterDecl is called when production typeParameterDecl is exited.
func (s *BaseMoxieListener) ExitTypeParameterDecl(ctx *TypeParameterDeclContext) {}

// EnterTypeConstraint is called when production typeConstraint is entered.
func (s *BaseMoxieListener) EnterTypeConstraint(ctx *TypeConstraintContext) {}

// ExitTypeConstraint is called when production typeConstraint is exited.
func (s *BaseMoxieListener) ExitTypeConstraint(ctx *TypeConstraintContext) {}

// EnterVarDecl is called when production varDecl is entered.
func (s *BaseMoxieListener) EnterVarDecl(ctx *VarDeclContext) {}

// ExitVarDecl is called when production varDecl is exited.
func (s *BaseMoxieListener) ExitVarDecl(ctx *VarDeclContext) {}

// EnterVarSpec is called when production varSpec is entered.
func (s *BaseMoxieListener) EnterVarSpec(ctx *VarSpecContext) {}

// ExitVarSpec is called when production varSpec is exited.
func (s *BaseMoxieListener) ExitVarSpec(ctx *VarSpecContext) {}

// EnterFunctionDecl is called when production functionDecl is entered.
func (s *BaseMoxieListener) EnterFunctionDecl(ctx *FunctionDeclContext) {}

// ExitFunctionDecl is called when production functionDecl is exited.
func (s *BaseMoxieListener) ExitFunctionDecl(ctx *FunctionDeclContext) {}

// EnterMethodDecl is called when production methodDecl is entered.
func (s *BaseMoxieListener) EnterMethodDecl(ctx *MethodDeclContext) {}

// ExitMethodDecl is called when production methodDecl is exited.
func (s *BaseMoxieListener) ExitMethodDecl(ctx *MethodDeclContext) {}

// EnterReceiver is called when production receiver is entered.
func (s *BaseMoxieListener) EnterReceiver(ctx *ReceiverContext) {}

// ExitReceiver is called when production receiver is exited.
func (s *BaseMoxieListener) ExitReceiver(ctx *ReceiverContext) {}

// EnterNamedType is called when production NamedType is entered.
func (s *BaseMoxieListener) EnterNamedType(ctx *NamedTypeContext) {}

// ExitNamedType is called when production NamedType is exited.
func (s *BaseMoxieListener) ExitNamedType(ctx *NamedTypeContext) {}

// EnterTypeLiteral is called when production TypeLiteral is entered.
func (s *BaseMoxieListener) EnterTypeLiteral(ctx *TypeLiteralContext) {}

// ExitTypeLiteral is called when production TypeLiteral is exited.
func (s *BaseMoxieListener) ExitTypeLiteral(ctx *TypeLiteralContext) {}

// EnterParenType is called when production ParenType is entered.
func (s *BaseMoxieListener) EnterParenType(ctx *ParenTypeContext) {}

// ExitParenType is called when production ParenType is exited.
func (s *BaseMoxieListener) ExitParenType(ctx *ParenTypeContext) {}

// EnterConstType is called when production ConstType is entered.
func (s *BaseMoxieListener) EnterConstType(ctx *ConstTypeContext) {}

// ExitConstType is called when production ConstType is exited.
func (s *BaseMoxieListener) ExitConstType(ctx *ConstTypeContext) {}

// EnterTypeName is called when production typeName is entered.
func (s *BaseMoxieListener) EnterTypeName(ctx *TypeNameContext) {}

// ExitTypeName is called when production typeName is exited.
func (s *BaseMoxieListener) ExitTypeName(ctx *TypeNameContext) {}

// EnterTypeArgs is called when production typeArgs is entered.
func (s *BaseMoxieListener) EnterTypeArgs(ctx *TypeArgsContext) {}

// ExitTypeArgs is called when production typeArgs is exited.
func (s *BaseMoxieListener) ExitTypeArgs(ctx *TypeArgsContext) {}

// EnterTypeLit is called when production typeLit is entered.
func (s *BaseMoxieListener) EnterTypeLit(ctx *TypeLitContext) {}

// ExitTypeLit is called when production typeLit is exited.
func (s *BaseMoxieListener) ExitTypeLit(ctx *TypeLitContext) {}

// EnterArrayType is called when production arrayType is entered.
func (s *BaseMoxieListener) EnterArrayType(ctx *ArrayTypeContext) {}

// ExitArrayType is called when production arrayType is exited.
func (s *BaseMoxieListener) ExitArrayType(ctx *ArrayTypeContext) {}

// EnterArrayLength is called when production arrayLength is entered.
func (s *BaseMoxieListener) EnterArrayLength(ctx *ArrayLengthContext) {}

// ExitArrayLength is called when production arrayLength is exited.
func (s *BaseMoxieListener) ExitArrayLength(ctx *ArrayLengthContext) {}

// EnterElementType is called when production elementType is entered.
func (s *BaseMoxieListener) EnterElementType(ctx *ElementTypeContext) {}

// ExitElementType is called when production elementType is exited.
func (s *BaseMoxieListener) ExitElementType(ctx *ElementTypeContext) {}

// EnterSliceType is called when production sliceType is entered.
func (s *BaseMoxieListener) EnterSliceType(ctx *SliceTypeContext) {}

// ExitSliceType is called when production sliceType is exited.
func (s *BaseMoxieListener) ExitSliceType(ctx *SliceTypeContext) {}

// EnterStructType is called when production structType is entered.
func (s *BaseMoxieListener) EnterStructType(ctx *StructTypeContext) {}

// ExitStructType is called when production structType is exited.
func (s *BaseMoxieListener) ExitStructType(ctx *StructTypeContext) {}

// EnterFieldDecl is called when production fieldDecl is entered.
func (s *BaseMoxieListener) EnterFieldDecl(ctx *FieldDeclContext) {}

// ExitFieldDecl is called when production fieldDecl is exited.
func (s *BaseMoxieListener) ExitFieldDecl(ctx *FieldDeclContext) {}

// EnterEmbeddedField is called when production embeddedField is entered.
func (s *BaseMoxieListener) EnterEmbeddedField(ctx *EmbeddedFieldContext) {}

// ExitEmbeddedField is called when production embeddedField is exited.
func (s *BaseMoxieListener) ExitEmbeddedField(ctx *EmbeddedFieldContext) {}

// EnterTag_ is called when production tag_ is entered.
func (s *BaseMoxieListener) EnterTag_(ctx *Tag_Context) {}

// ExitTag_ is called when production tag_ is exited.
func (s *BaseMoxieListener) ExitTag_(ctx *Tag_Context) {}

// EnterPointerType is called when production pointerType is entered.
func (s *BaseMoxieListener) EnterPointerType(ctx *PointerTypeContext) {}

// ExitPointerType is called when production pointerType is exited.
func (s *BaseMoxieListener) ExitPointerType(ctx *PointerTypeContext) {}

// EnterFunctionType is called when production functionType is entered.
func (s *BaseMoxieListener) EnterFunctionType(ctx *FunctionTypeContext) {}

// ExitFunctionType is called when production functionType is exited.
func (s *BaseMoxieListener) ExitFunctionType(ctx *FunctionTypeContext) {}

// EnterSignature is called when production signature is entered.
func (s *BaseMoxieListener) EnterSignature(ctx *SignatureContext) {}

// ExitSignature is called when production signature is exited.
func (s *BaseMoxieListener) ExitSignature(ctx *SignatureContext) {}

// EnterResult is called when production result is entered.
func (s *BaseMoxieListener) EnterResult(ctx *ResultContext) {}

// ExitResult is called when production result is exited.
func (s *BaseMoxieListener) ExitResult(ctx *ResultContext) {}

// EnterParameters is called when production parameters is entered.
func (s *BaseMoxieListener) EnterParameters(ctx *ParametersContext) {}

// ExitParameters is called when production parameters is exited.
func (s *BaseMoxieListener) ExitParameters(ctx *ParametersContext) {}

// EnterParameterDecl is called when production parameterDecl is entered.
func (s *BaseMoxieListener) EnterParameterDecl(ctx *ParameterDeclContext) {}

// ExitParameterDecl is called when production parameterDecl is exited.
func (s *BaseMoxieListener) ExitParameterDecl(ctx *ParameterDeclContext) {}

// EnterInterfaceType is called when production interfaceType is entered.
func (s *BaseMoxieListener) EnterInterfaceType(ctx *InterfaceTypeContext) {}

// ExitInterfaceType is called when production interfaceType is exited.
func (s *BaseMoxieListener) ExitInterfaceType(ctx *InterfaceTypeContext) {}

// EnterInterfaceElem is called when production interfaceElem is entered.
func (s *BaseMoxieListener) EnterInterfaceElem(ctx *InterfaceElemContext) {}

// ExitInterfaceElem is called when production interfaceElem is exited.
func (s *BaseMoxieListener) ExitInterfaceElem(ctx *InterfaceElemContext) {}

// EnterMethodElem is called when production methodElem is entered.
func (s *BaseMoxieListener) EnterMethodElem(ctx *MethodElemContext) {}

// ExitMethodElem is called when production methodElem is exited.
func (s *BaseMoxieListener) ExitMethodElem(ctx *MethodElemContext) {}

// EnterTypeElem is called when production typeElem is entered.
func (s *BaseMoxieListener) EnterTypeElem(ctx *TypeElemContext) {}

// ExitTypeElem is called when production typeElem is exited.
func (s *BaseMoxieListener) ExitTypeElem(ctx *TypeElemContext) {}

// EnterTypeTerm is called when production typeTerm is entered.
func (s *BaseMoxieListener) EnterTypeTerm(ctx *TypeTermContext) {}

// ExitTypeTerm is called when production typeTerm is exited.
func (s *BaseMoxieListener) ExitTypeTerm(ctx *TypeTermContext) {}

// EnterMapType is called when production mapType is entered.
func (s *BaseMoxieListener) EnterMapType(ctx *MapTypeContext) {}

// ExitMapType is called when production mapType is exited.
func (s *BaseMoxieListener) ExitMapType(ctx *MapTypeContext) {}

// EnterSendRecvChan is called when production SendRecvChan is entered.
func (s *BaseMoxieListener) EnterSendRecvChan(ctx *SendRecvChanContext) {}

// ExitSendRecvChan is called when production SendRecvChan is exited.
func (s *BaseMoxieListener) ExitSendRecvChan(ctx *SendRecvChanContext) {}

// EnterRecvOnlyChan is called when production RecvOnlyChan is entered.
func (s *BaseMoxieListener) EnterRecvOnlyChan(ctx *RecvOnlyChanContext) {}

// ExitRecvOnlyChan is called when production RecvOnlyChan is exited.
func (s *BaseMoxieListener) ExitRecvOnlyChan(ctx *RecvOnlyChanContext) {}

// EnterSendRecvChanCompat is called when production SendRecvChanCompat is entered.
func (s *BaseMoxieListener) EnterSendRecvChanCompat(ctx *SendRecvChanCompatContext) {}

// ExitSendRecvChanCompat is called when production SendRecvChanCompat is exited.
func (s *BaseMoxieListener) ExitSendRecvChanCompat(ctx *SendRecvChanCompatContext) {}

// EnterRecvOnlyChanCompat is called when production RecvOnlyChanCompat is entered.
func (s *BaseMoxieListener) EnterRecvOnlyChanCompat(ctx *RecvOnlyChanCompatContext) {}

// ExitRecvOnlyChanCompat is called when production RecvOnlyChanCompat is exited.
func (s *BaseMoxieListener) ExitRecvOnlyChanCompat(ctx *RecvOnlyChanCompatContext) {}

// EnterBlock is called when production block is entered.
func (s *BaseMoxieListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BaseMoxieListener) ExitBlock(ctx *BlockContext) {}

// EnterStatementList is called when production statementList is entered.
func (s *BaseMoxieListener) EnterStatementList(ctx *StatementListContext) {}

// ExitStatementList is called when production statementList is exited.
func (s *BaseMoxieListener) ExitStatementList(ctx *StatementListContext) {}

// EnterDeclStmt is called when production DeclStmt is entered.
func (s *BaseMoxieListener) EnterDeclStmt(ctx *DeclStmtContext) {}

// ExitDeclStmt is called when production DeclStmt is exited.
func (s *BaseMoxieListener) ExitDeclStmt(ctx *DeclStmtContext) {}

// EnterSimpleStatement is called when production SimpleStatement is entered.
func (s *BaseMoxieListener) EnterSimpleStatement(ctx *SimpleStatementContext) {}

// ExitSimpleStatement is called when production SimpleStatement is exited.
func (s *BaseMoxieListener) ExitSimpleStatement(ctx *SimpleStatementContext) {}

// EnterLabeledStatement is called when production LabeledStatement is entered.
func (s *BaseMoxieListener) EnterLabeledStatement(ctx *LabeledStatementContext) {}

// ExitLabeledStatement is called when production LabeledStatement is exited.
func (s *BaseMoxieListener) ExitLabeledStatement(ctx *LabeledStatementContext) {}

// EnterGoStatement is called when production GoStatement is entered.
func (s *BaseMoxieListener) EnterGoStatement(ctx *GoStatementContext) {}

// ExitGoStatement is called when production GoStatement is exited.
func (s *BaseMoxieListener) ExitGoStatement(ctx *GoStatementContext) {}

// EnterReturnStatement is called when production ReturnStatement is entered.
func (s *BaseMoxieListener) EnterReturnStatement(ctx *ReturnStatementContext) {}

// ExitReturnStatement is called when production ReturnStatement is exited.
func (s *BaseMoxieListener) ExitReturnStatement(ctx *ReturnStatementContext) {}

// EnterBreakStatement is called when production BreakStatement is entered.
func (s *BaseMoxieListener) EnterBreakStatement(ctx *BreakStatementContext) {}

// ExitBreakStatement is called when production BreakStatement is exited.
func (s *BaseMoxieListener) ExitBreakStatement(ctx *BreakStatementContext) {}

// EnterContinueStatement is called when production ContinueStatement is entered.
func (s *BaseMoxieListener) EnterContinueStatement(ctx *ContinueStatementContext) {}

// ExitContinueStatement is called when production ContinueStatement is exited.
func (s *BaseMoxieListener) ExitContinueStatement(ctx *ContinueStatementContext) {}

// EnterGotoStatement is called when production GotoStatement is entered.
func (s *BaseMoxieListener) EnterGotoStatement(ctx *GotoStatementContext) {}

// ExitGotoStatement is called when production GotoStatement is exited.
func (s *BaseMoxieListener) ExitGotoStatement(ctx *GotoStatementContext) {}

// EnterFallthroughStatement is called when production FallthroughStatement is entered.
func (s *BaseMoxieListener) EnterFallthroughStatement(ctx *FallthroughStatementContext) {}

// ExitFallthroughStatement is called when production FallthroughStatement is exited.
func (s *BaseMoxieListener) ExitFallthroughStatement(ctx *FallthroughStatementContext) {}

// EnterBlockStatement is called when production BlockStatement is entered.
func (s *BaseMoxieListener) EnterBlockStatement(ctx *BlockStatementContext) {}

// ExitBlockStatement is called when production BlockStatement is exited.
func (s *BaseMoxieListener) ExitBlockStatement(ctx *BlockStatementContext) {}

// EnterIfStatement is called when production IfStatement is entered.
func (s *BaseMoxieListener) EnterIfStatement(ctx *IfStatementContext) {}

// ExitIfStatement is called when production IfStatement is exited.
func (s *BaseMoxieListener) ExitIfStatement(ctx *IfStatementContext) {}

// EnterSwitchStatement is called when production SwitchStatement is entered.
func (s *BaseMoxieListener) EnterSwitchStatement(ctx *SwitchStatementContext) {}

// ExitSwitchStatement is called when production SwitchStatement is exited.
func (s *BaseMoxieListener) ExitSwitchStatement(ctx *SwitchStatementContext) {}

// EnterSelectStatement is called when production SelectStatement is entered.
func (s *BaseMoxieListener) EnterSelectStatement(ctx *SelectStatementContext) {}

// ExitSelectStatement is called when production SelectStatement is exited.
func (s *BaseMoxieListener) ExitSelectStatement(ctx *SelectStatementContext) {}

// EnterForStatement is called when production ForStatement is entered.
func (s *BaseMoxieListener) EnterForStatement(ctx *ForStatementContext) {}

// ExitForStatement is called when production ForStatement is exited.
func (s *BaseMoxieListener) ExitForStatement(ctx *ForStatementContext) {}

// EnterDeferStatement is called when production DeferStatement is entered.
func (s *BaseMoxieListener) EnterDeferStatement(ctx *DeferStatementContext) {}

// ExitDeferStatement is called when production DeferStatement is exited.
func (s *BaseMoxieListener) ExitDeferStatement(ctx *DeferStatementContext) {}

// EnterSimpleStmt is called when production simpleStmt is entered.
func (s *BaseMoxieListener) EnterSimpleStmt(ctx *SimpleStmtContext) {}

// ExitSimpleStmt is called when production simpleStmt is exited.
func (s *BaseMoxieListener) ExitSimpleStmt(ctx *SimpleStmtContext) {}

// EnterExpressionStmt is called when production expressionStmt is entered.
func (s *BaseMoxieListener) EnterExpressionStmt(ctx *ExpressionStmtContext) {}

// ExitExpressionStmt is called when production expressionStmt is exited.
func (s *BaseMoxieListener) ExitExpressionStmt(ctx *ExpressionStmtContext) {}

// EnterSendStmt is called when production sendStmt is entered.
func (s *BaseMoxieListener) EnterSendStmt(ctx *SendStmtContext) {}

// ExitSendStmt is called when production sendStmt is exited.
func (s *BaseMoxieListener) ExitSendStmt(ctx *SendStmtContext) {}

// EnterIncDecStmt is called when production incDecStmt is entered.
func (s *BaseMoxieListener) EnterIncDecStmt(ctx *IncDecStmtContext) {}

// ExitIncDecStmt is called when production incDecStmt is exited.
func (s *BaseMoxieListener) ExitIncDecStmt(ctx *IncDecStmtContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseMoxieListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseMoxieListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterAssign_op is called when production assign_op is entered.
func (s *BaseMoxieListener) EnterAssign_op(ctx *Assign_opContext) {}

// ExitAssign_op is called when production assign_op is exited.
func (s *BaseMoxieListener) ExitAssign_op(ctx *Assign_opContext) {}

// EnterShortVarDecl is called when production shortVarDecl is entered.
func (s *BaseMoxieListener) EnterShortVarDecl(ctx *ShortVarDeclContext) {}

// ExitShortVarDecl is called when production shortVarDecl is exited.
func (s *BaseMoxieListener) ExitShortVarDecl(ctx *ShortVarDeclContext) {}

// EnterLabeledStmt is called when production labeledStmt is entered.
func (s *BaseMoxieListener) EnterLabeledStmt(ctx *LabeledStmtContext) {}

// ExitLabeledStmt is called when production labeledStmt is exited.
func (s *BaseMoxieListener) ExitLabeledStmt(ctx *LabeledStmtContext) {}

// EnterReturnStmt is called when production returnStmt is entered.
func (s *BaseMoxieListener) EnterReturnStmt(ctx *ReturnStmtContext) {}

// ExitReturnStmt is called when production returnStmt is exited.
func (s *BaseMoxieListener) ExitReturnStmt(ctx *ReturnStmtContext) {}

// EnterBreakStmt is called when production breakStmt is entered.
func (s *BaseMoxieListener) EnterBreakStmt(ctx *BreakStmtContext) {}

// ExitBreakStmt is called when production breakStmt is exited.
func (s *BaseMoxieListener) ExitBreakStmt(ctx *BreakStmtContext) {}

// EnterContinueStmt is called when production continueStmt is entered.
func (s *BaseMoxieListener) EnterContinueStmt(ctx *ContinueStmtContext) {}

// ExitContinueStmt is called when production continueStmt is exited.
func (s *BaseMoxieListener) ExitContinueStmt(ctx *ContinueStmtContext) {}

// EnterGotoStmt is called when production gotoStmt is entered.
func (s *BaseMoxieListener) EnterGotoStmt(ctx *GotoStmtContext) {}

// ExitGotoStmt is called when production gotoStmt is exited.
func (s *BaseMoxieListener) ExitGotoStmt(ctx *GotoStmtContext) {}

// EnterFallthroughStmt is called when production fallthroughStmt is entered.
func (s *BaseMoxieListener) EnterFallthroughStmt(ctx *FallthroughStmtContext) {}

// ExitFallthroughStmt is called when production fallthroughStmt is exited.
func (s *BaseMoxieListener) ExitFallthroughStmt(ctx *FallthroughStmtContext) {}

// EnterDeferStmt is called when production deferStmt is entered.
func (s *BaseMoxieListener) EnterDeferStmt(ctx *DeferStmtContext) {}

// ExitDeferStmt is called when production deferStmt is exited.
func (s *BaseMoxieListener) ExitDeferStmt(ctx *DeferStmtContext) {}

// EnterIfStmt is called when production ifStmt is entered.
func (s *BaseMoxieListener) EnterIfStmt(ctx *IfStmtContext) {}

// ExitIfStmt is called when production ifStmt is exited.
func (s *BaseMoxieListener) ExitIfStmt(ctx *IfStmtContext) {}

// EnterSwitchStmt is called when production switchStmt is entered.
func (s *BaseMoxieListener) EnterSwitchStmt(ctx *SwitchStmtContext) {}

// ExitSwitchStmt is called when production switchStmt is exited.
func (s *BaseMoxieListener) ExitSwitchStmt(ctx *SwitchStmtContext) {}

// EnterExprSwitchStmt is called when production exprSwitchStmt is entered.
func (s *BaseMoxieListener) EnterExprSwitchStmt(ctx *ExprSwitchStmtContext) {}

// ExitExprSwitchStmt is called when production exprSwitchStmt is exited.
func (s *BaseMoxieListener) ExitExprSwitchStmt(ctx *ExprSwitchStmtContext) {}

// EnterExprCaseClause is called when production exprCaseClause is entered.
func (s *BaseMoxieListener) EnterExprCaseClause(ctx *ExprCaseClauseContext) {}

// ExitExprCaseClause is called when production exprCaseClause is exited.
func (s *BaseMoxieListener) ExitExprCaseClause(ctx *ExprCaseClauseContext) {}

// EnterExprSwitchCase is called when production exprSwitchCase is entered.
func (s *BaseMoxieListener) EnterExprSwitchCase(ctx *ExprSwitchCaseContext) {}

// ExitExprSwitchCase is called when production exprSwitchCase is exited.
func (s *BaseMoxieListener) ExitExprSwitchCase(ctx *ExprSwitchCaseContext) {}

// EnterTypeSwitchStmt is called when production typeSwitchStmt is entered.
func (s *BaseMoxieListener) EnterTypeSwitchStmt(ctx *TypeSwitchStmtContext) {}

// ExitTypeSwitchStmt is called when production typeSwitchStmt is exited.
func (s *BaseMoxieListener) ExitTypeSwitchStmt(ctx *TypeSwitchStmtContext) {}

// EnterTypeSwitchGuard is called when production typeSwitchGuard is entered.
func (s *BaseMoxieListener) EnterTypeSwitchGuard(ctx *TypeSwitchGuardContext) {}

// ExitTypeSwitchGuard is called when production typeSwitchGuard is exited.
func (s *BaseMoxieListener) ExitTypeSwitchGuard(ctx *TypeSwitchGuardContext) {}

// EnterTypeCaseClause is called when production typeCaseClause is entered.
func (s *BaseMoxieListener) EnterTypeCaseClause(ctx *TypeCaseClauseContext) {}

// ExitTypeCaseClause is called when production typeCaseClause is exited.
func (s *BaseMoxieListener) ExitTypeCaseClause(ctx *TypeCaseClauseContext) {}

// EnterTypeSwitchCase is called when production typeSwitchCase is entered.
func (s *BaseMoxieListener) EnterTypeSwitchCase(ctx *TypeSwitchCaseContext) {}

// ExitTypeSwitchCase is called when production typeSwitchCase is exited.
func (s *BaseMoxieListener) ExitTypeSwitchCase(ctx *TypeSwitchCaseContext) {}

// EnterTypeList is called when production typeList is entered.
func (s *BaseMoxieListener) EnterTypeList(ctx *TypeListContext) {}

// ExitTypeList is called when production typeList is exited.
func (s *BaseMoxieListener) ExitTypeList(ctx *TypeListContext) {}

// EnterSelectStmt is called when production selectStmt is entered.
func (s *BaseMoxieListener) EnterSelectStmt(ctx *SelectStmtContext) {}

// ExitSelectStmt is called when production selectStmt is exited.
func (s *BaseMoxieListener) ExitSelectStmt(ctx *SelectStmtContext) {}

// EnterCommClause is called when production commClause is entered.
func (s *BaseMoxieListener) EnterCommClause(ctx *CommClauseContext) {}

// ExitCommClause is called when production commClause is exited.
func (s *BaseMoxieListener) ExitCommClause(ctx *CommClauseContext) {}

// EnterCommCase is called when production commCase is entered.
func (s *BaseMoxieListener) EnterCommCase(ctx *CommCaseContext) {}

// ExitCommCase is called when production commCase is exited.
func (s *BaseMoxieListener) ExitCommCase(ctx *CommCaseContext) {}

// EnterRecvStmt is called when production recvStmt is entered.
func (s *BaseMoxieListener) EnterRecvStmt(ctx *RecvStmtContext) {}

// ExitRecvStmt is called when production recvStmt is exited.
func (s *BaseMoxieListener) ExitRecvStmt(ctx *RecvStmtContext) {}

// EnterForStmt is called when production forStmt is entered.
func (s *BaseMoxieListener) EnterForStmt(ctx *ForStmtContext) {}

// ExitForStmt is called when production forStmt is exited.
func (s *BaseMoxieListener) ExitForStmt(ctx *ForStmtContext) {}

// EnterForClause is called when production forClause is entered.
func (s *BaseMoxieListener) EnterForClause(ctx *ForClauseContext) {}

// ExitForClause is called when production forClause is exited.
func (s *BaseMoxieListener) ExitForClause(ctx *ForClauseContext) {}

// EnterRangeClause is called when production rangeClause is entered.
func (s *BaseMoxieListener) EnterRangeClause(ctx *RangeClauseContext) {}

// ExitRangeClause is called when production rangeClause is exited.
func (s *BaseMoxieListener) ExitRangeClause(ctx *RangeClauseContext) {}

// EnterGoStmt is called when production goStmt is entered.
func (s *BaseMoxieListener) EnterGoStmt(ctx *GoStmtContext) {}

// ExitGoStmt is called when production goStmt is exited.
func (s *BaseMoxieListener) ExitGoStmt(ctx *GoStmtContext) {}

// EnterMultiplicativeExpr is called when production MultiplicativeExpr is entered.
func (s *BaseMoxieListener) EnterMultiplicativeExpr(ctx *MultiplicativeExprContext) {}

// ExitMultiplicativeExpr is called when production MultiplicativeExpr is exited.
func (s *BaseMoxieListener) ExitMultiplicativeExpr(ctx *MultiplicativeExprContext) {}

// EnterConcatenationExpr is called when production ConcatenationExpr is entered.
func (s *BaseMoxieListener) EnterConcatenationExpr(ctx *ConcatenationExprContext) {}

// ExitConcatenationExpr is called when production ConcatenationExpr is exited.
func (s *BaseMoxieListener) ExitConcatenationExpr(ctx *ConcatenationExprContext) {}

// EnterLogicalOrExpr is called when production LogicalOrExpr is entered.
func (s *BaseMoxieListener) EnterLogicalOrExpr(ctx *LogicalOrExprContext) {}

// ExitLogicalOrExpr is called when production LogicalOrExpr is exited.
func (s *BaseMoxieListener) ExitLogicalOrExpr(ctx *LogicalOrExprContext) {}

// EnterAdditiveExpr is called when production AdditiveExpr is entered.
func (s *BaseMoxieListener) EnterAdditiveExpr(ctx *AdditiveExprContext) {}

// ExitAdditiveExpr is called when production AdditiveExpr is exited.
func (s *BaseMoxieListener) ExitAdditiveExpr(ctx *AdditiveExprContext) {}

// EnterUnaryExpression is called when production UnaryExpression is entered.
func (s *BaseMoxieListener) EnterUnaryExpression(ctx *UnaryExpressionContext) {}

// ExitUnaryExpression is called when production UnaryExpression is exited.
func (s *BaseMoxieListener) ExitUnaryExpression(ctx *UnaryExpressionContext) {}

// EnterRelationalExpr is called when production RelationalExpr is entered.
func (s *BaseMoxieListener) EnterRelationalExpr(ctx *RelationalExprContext) {}

// ExitRelationalExpr is called when production RelationalExpr is exited.
func (s *BaseMoxieListener) ExitRelationalExpr(ctx *RelationalExprContext) {}

// EnterLogicalAndExpr is called when production LogicalAndExpr is entered.
func (s *BaseMoxieListener) EnterLogicalAndExpr(ctx *LogicalAndExprContext) {}

// ExitLogicalAndExpr is called when production LogicalAndExpr is exited.
func (s *BaseMoxieListener) ExitLogicalAndExpr(ctx *LogicalAndExprContext) {}

// EnterSelectorExpr is called when production SelectorExpr is entered.
func (s *BaseMoxieListener) EnterSelectorExpr(ctx *SelectorExprContext) {}

// ExitSelectorExpr is called when production SelectorExpr is exited.
func (s *BaseMoxieListener) ExitSelectorExpr(ctx *SelectorExprContext) {}

// EnterTypeAssertionExpr is called when production TypeAssertionExpr is entered.
func (s *BaseMoxieListener) EnterTypeAssertionExpr(ctx *TypeAssertionExprContext) {}

// ExitTypeAssertionExpr is called when production TypeAssertionExpr is exited.
func (s *BaseMoxieListener) ExitTypeAssertionExpr(ctx *TypeAssertionExprContext) {}

// EnterConversionExpr is called when production ConversionExpr is entered.
func (s *BaseMoxieListener) EnterConversionExpr(ctx *ConversionExprContext) {}

// ExitConversionExpr is called when production ConversionExpr is exited.
func (s *BaseMoxieListener) ExitConversionExpr(ctx *ConversionExprContext) {}

// EnterPrimaryOperand is called when production PrimaryOperand is entered.
func (s *BaseMoxieListener) EnterPrimaryOperand(ctx *PrimaryOperandContext) {}

// ExitPrimaryOperand is called when production PrimaryOperand is exited.
func (s *BaseMoxieListener) ExitPrimaryOperand(ctx *PrimaryOperandContext) {}

// EnterSliceExpr is called when production SliceExpr is entered.
func (s *BaseMoxieListener) EnterSliceExpr(ctx *SliceExprContext) {}

// ExitSliceExpr is called when production SliceExpr is exited.
func (s *BaseMoxieListener) ExitSliceExpr(ctx *SliceExprContext) {}

// EnterCallExpr is called when production CallExpr is entered.
func (s *BaseMoxieListener) EnterCallExpr(ctx *CallExprContext) {}

// ExitCallExpr is called when production CallExpr is exited.
func (s *BaseMoxieListener) ExitCallExpr(ctx *CallExprContext) {}

// EnterMethodExpression is called when production MethodExpression is entered.
func (s *BaseMoxieListener) EnterMethodExpression(ctx *MethodExpressionContext) {}

// ExitMethodExpression is called when production MethodExpression is exited.
func (s *BaseMoxieListener) ExitMethodExpression(ctx *MethodExpressionContext) {}

// EnterIndexExpr is called when production IndexExpr is entered.
func (s *BaseMoxieListener) EnterIndexExpr(ctx *IndexExprContext) {}

// ExitIndexExpr is called when production IndexExpr is exited.
func (s *BaseMoxieListener) ExitIndexExpr(ctx *IndexExprContext) {}

// EnterUnaryExpr is called when production unaryExpr is entered.
func (s *BaseMoxieListener) EnterUnaryExpr(ctx *UnaryExprContext) {}

// ExitUnaryExpr is called when production unaryExpr is exited.
func (s *BaseMoxieListener) ExitUnaryExpr(ctx *UnaryExprContext) {}

// EnterSimpleConversion is called when production SimpleConversion is entered.
func (s *BaseMoxieListener) EnterSimpleConversion(ctx *SimpleConversionContext) {}

// ExitSimpleConversion is called when production SimpleConversion is exited.
func (s *BaseMoxieListener) ExitSimpleConversion(ctx *SimpleConversionContext) {}

// EnterSliceCastExpr is called when production SliceCastExpr is entered.
func (s *BaseMoxieListener) EnterSliceCastExpr(ctx *SliceCastExprContext) {}

// ExitSliceCastExpr is called when production SliceCastExpr is exited.
func (s *BaseMoxieListener) ExitSliceCastExpr(ctx *SliceCastExprContext) {}

// EnterSliceCastEndianExpr is called when production SliceCastEndianExpr is entered.
func (s *BaseMoxieListener) EnterSliceCastEndianExpr(ctx *SliceCastEndianExprContext) {}

// ExitSliceCastEndianExpr is called when production SliceCastEndianExpr is exited.
func (s *BaseMoxieListener) ExitSliceCastEndianExpr(ctx *SliceCastEndianExprContext) {}

// EnterSliceCastCopyExpr is called when production SliceCastCopyExpr is entered.
func (s *BaseMoxieListener) EnterSliceCastCopyExpr(ctx *SliceCastCopyExprContext) {}

// ExitSliceCastCopyExpr is called when production SliceCastCopyExpr is exited.
func (s *BaseMoxieListener) ExitSliceCastCopyExpr(ctx *SliceCastCopyExprContext) {}

// EnterSliceCastCopyEndianExpr is called when production SliceCastCopyEndianExpr is entered.
func (s *BaseMoxieListener) EnterSliceCastCopyEndianExpr(ctx *SliceCastCopyEndianExprContext) {}

// ExitSliceCastCopyEndianExpr is called when production SliceCastCopyEndianExpr is exited.
func (s *BaseMoxieListener) ExitSliceCastCopyEndianExpr(ctx *SliceCastCopyEndianExprContext) {}

// EnterEndianness is called when production endianness is entered.
func (s *BaseMoxieListener) EnterEndianness(ctx *EndiannessContext) {}

// ExitEndianness is called when production endianness is exited.
func (s *BaseMoxieListener) ExitEndianness(ctx *EndiannessContext) {}

// EnterLiteralOperand is called when production LiteralOperand is entered.
func (s *BaseMoxieListener) EnterLiteralOperand(ctx *LiteralOperandContext) {}

// ExitLiteralOperand is called when production LiteralOperand is exited.
func (s *BaseMoxieListener) ExitLiteralOperand(ctx *LiteralOperandContext) {}

// EnterNameOperand is called when production NameOperand is entered.
func (s *BaseMoxieListener) EnterNameOperand(ctx *NameOperandContext) {}

// ExitNameOperand is called when production NameOperand is exited.
func (s *BaseMoxieListener) ExitNameOperand(ctx *NameOperandContext) {}

// EnterParenOperand is called when production ParenOperand is entered.
func (s *BaseMoxieListener) EnterParenOperand(ctx *ParenOperandContext) {}

// ExitParenOperand is called when production ParenOperand is exited.
func (s *BaseMoxieListener) ExitParenOperand(ctx *ParenOperandContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BaseMoxieListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BaseMoxieListener) ExitLiteral(ctx *LiteralContext) {}

// EnterBasicLit is called when production basicLit is entered.
func (s *BaseMoxieListener) EnterBasicLit(ctx *BasicLitContext) {}

// ExitBasicLit is called when production basicLit is exited.
func (s *BaseMoxieListener) ExitBasicLit(ctx *BasicLitContext) {}

// EnterString_ is called when production string_ is entered.
func (s *BaseMoxieListener) EnterString_(ctx *String_Context) {}

// ExitString_ is called when production string_ is exited.
func (s *BaseMoxieListener) ExitString_(ctx *String_Context) {}

// EnterOperandName is called when production operandName is entered.
func (s *BaseMoxieListener) EnterOperandName(ctx *OperandNameContext) {}

// ExitOperandName is called when production operandName is exited.
func (s *BaseMoxieListener) ExitOperandName(ctx *OperandNameContext) {}

// EnterQualifiedIdent is called when production qualifiedIdent is entered.
func (s *BaseMoxieListener) EnterQualifiedIdent(ctx *QualifiedIdentContext) {}

// ExitQualifiedIdent is called when production qualifiedIdent is exited.
func (s *BaseMoxieListener) ExitQualifiedIdent(ctx *QualifiedIdentContext) {}

// EnterCompositeLit is called when production compositeLit is entered.
func (s *BaseMoxieListener) EnterCompositeLit(ctx *CompositeLitContext) {}

// ExitCompositeLit is called when production compositeLit is exited.
func (s *BaseMoxieListener) ExitCompositeLit(ctx *CompositeLitContext) {}

// EnterLiteralType is called when production literalType is entered.
func (s *BaseMoxieListener) EnterLiteralType(ctx *LiteralTypeContext) {}

// ExitLiteralType is called when production literalType is exited.
func (s *BaseMoxieListener) ExitLiteralType(ctx *LiteralTypeContext) {}

// EnterLiteralValue is called when production literalValue is entered.
func (s *BaseMoxieListener) EnterLiteralValue(ctx *LiteralValueContext) {}

// ExitLiteralValue is called when production literalValue is exited.
func (s *BaseMoxieListener) ExitLiteralValue(ctx *LiteralValueContext) {}

// EnterElementList is called when production elementList is entered.
func (s *BaseMoxieListener) EnterElementList(ctx *ElementListContext) {}

// ExitElementList is called when production elementList is exited.
func (s *BaseMoxieListener) ExitElementList(ctx *ElementListContext) {}

// EnterKeyedElement is called when production keyedElement is entered.
func (s *BaseMoxieListener) EnterKeyedElement(ctx *KeyedElementContext) {}

// ExitKeyedElement is called when production keyedElement is exited.
func (s *BaseMoxieListener) ExitKeyedElement(ctx *KeyedElementContext) {}

// EnterKey is called when production key is entered.
func (s *BaseMoxieListener) EnterKey(ctx *KeyContext) {}

// ExitKey is called when production key is exited.
func (s *BaseMoxieListener) ExitKey(ctx *KeyContext) {}

// EnterElement is called when production element is entered.
func (s *BaseMoxieListener) EnterElement(ctx *ElementContext) {}

// ExitElement is called when production element is exited.
func (s *BaseMoxieListener) ExitElement(ctx *ElementContext) {}

// EnterFunctionLit is called when production functionLit is entered.
func (s *BaseMoxieListener) EnterFunctionLit(ctx *FunctionLitContext) {}

// ExitFunctionLit is called when production functionLit is exited.
func (s *BaseMoxieListener) ExitFunctionLit(ctx *FunctionLitContext) {}

// EnterSelector is called when production selector is entered.
func (s *BaseMoxieListener) EnterSelector(ctx *SelectorContext) {}

// ExitSelector is called when production selector is exited.
func (s *BaseMoxieListener) ExitSelector(ctx *SelectorContext) {}

// EnterIndex is called when production index is entered.
func (s *BaseMoxieListener) EnterIndex(ctx *IndexContext) {}

// ExitIndex is called when production index is exited.
func (s *BaseMoxieListener) ExitIndex(ctx *IndexContext) {}

// EnterSlice_ is called when production slice_ is entered.
func (s *BaseMoxieListener) EnterSlice_(ctx *Slice_Context) {}

// ExitSlice_ is called when production slice_ is exited.
func (s *BaseMoxieListener) ExitSlice_(ctx *Slice_Context) {}

// EnterTypeAssertion is called when production typeAssertion is entered.
func (s *BaseMoxieListener) EnterTypeAssertion(ctx *TypeAssertionContext) {}

// ExitTypeAssertion is called when production typeAssertion is exited.
func (s *BaseMoxieListener) ExitTypeAssertion(ctx *TypeAssertionContext) {}

// EnterArguments is called when production arguments is entered.
func (s *BaseMoxieListener) EnterArguments(ctx *ArgumentsContext) {}

// ExitArguments is called when production arguments is exited.
func (s *BaseMoxieListener) ExitArguments(ctx *ArgumentsContext) {}

// EnterMethodExpr is called when production methodExpr is entered.
func (s *BaseMoxieListener) EnterMethodExpr(ctx *MethodExprContext) {}

// ExitMethodExpr is called when production methodExpr is exited.
func (s *BaseMoxieListener) ExitMethodExpr(ctx *MethodExprContext) {}

// EnterMul_op is called when production mul_op is entered.
func (s *BaseMoxieListener) EnterMul_op(ctx *Mul_opContext) {}

// ExitMul_op is called when production mul_op is exited.
func (s *BaseMoxieListener) ExitMul_op(ctx *Mul_opContext) {}

// EnterAdd_op is called when production add_op is entered.
func (s *BaseMoxieListener) EnterAdd_op(ctx *Add_opContext) {}

// ExitAdd_op is called when production add_op is exited.
func (s *BaseMoxieListener) ExitAdd_op(ctx *Add_opContext) {}

// EnterRel_op is called when production rel_op is entered.
func (s *BaseMoxieListener) EnterRel_op(ctx *Rel_opContext) {}

// ExitRel_op is called when production rel_op is exited.
func (s *BaseMoxieListener) ExitRel_op(ctx *Rel_opContext) {}

// EnterUnary_op is called when production unary_op is entered.
func (s *BaseMoxieListener) EnterUnary_op(ctx *Unary_opContext) {}

// ExitUnary_op is called when production unary_op is exited.
func (s *BaseMoxieListener) ExitUnary_op(ctx *Unary_opContext) {}

// EnterExpressionList is called when production expressionList is entered.
func (s *BaseMoxieListener) EnterExpressionList(ctx *ExpressionListContext) {}

// ExitExpressionList is called when production expressionList is exited.
func (s *BaseMoxieListener) ExitExpressionList(ctx *ExpressionListContext) {}

// EnterIdentifierList is called when production identifierList is entered.
func (s *BaseMoxieListener) EnterIdentifierList(ctx *IdentifierListContext) {}

// ExitIdentifierList is called when production identifierList is exited.
func (s *BaseMoxieListener) ExitIdentifierList(ctx *IdentifierListContext) {}

// EnterEos is called when production eos is entered.
func (s *BaseMoxieListener) EnterEos(ctx *EosContext) {}

// ExitEos is called when production eos is exited.
func (s *BaseMoxieListener) ExitEos(ctx *EosContext) {}
