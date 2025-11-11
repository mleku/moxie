package antlr

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

// TestPrintAST parses example.x and prints the parse tree
func TestPrintAST(t *testing.T) {
	// Read the example file
	content, err := os.ReadFile("../../moxie-intellij-plugin/example.x")
	if err != nil {
		t.Fatalf("Failed to read example.x: %v", err)
	}

	// Create input stream
	is := antlr.NewInputStream(string(content))

	// Create lexer
	lexer := NewMoxieLexer(is)

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser
	parser := NewMoxieParser(stream)

	// Add error listener
	errorListener := &CustomErrorListener{}
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	// Parse the source file
	tree := parser.SourceFile()

	// Check for errors
	if len(errorListener.errors) > 0 {
		t.Errorf("Parse errors detected:")
		for _, err := range errorListener.errors {
			t.Errorf("  Line %d:%d - %s", err.line, err.column, err.msg)
		}
	}

	// Print the parse tree (built-in method)
	fmt.Println("\n=== PARSE TREE (compact) ===")
	fmt.Println(tree.ToStringTree([]string{}, parser))

	// Print a more readable tree with limited depth
	fmt.Println("\n=== READABLE PARSE TREE (depth limited to 5) ===")
	printTreeLimited(tree, parser, 0, 5)
}

// printTreeLimited recursively prints the parse tree in a readable format with depth limit
func printTreeLimited(tree antlr.Tree, parser *MoxieParser, indent int, maxDepth int) {
	if tree == nil || indent > maxDepth {
		if indent > maxDepth {
			fmt.Printf("%s...\n", strings.Repeat("  ", indent))
		}
		return
	}

	indentStr := strings.Repeat("  ", indent)

	// Get node text
	var nodeText string
	switch t := tree.(type) {
	case *antlr.TerminalNodeImpl:
		// Terminal node (token)
		token := t.GetSymbol()
		tokenType := token.GetTokenType()
		var tokenName string

		// Safely get token name
		if tokenType >= 0 {
			literalNames := parser.GetLiteralNames()
			symbolicNames := parser.GetSymbolicNames()

			if tokenType < len(literalNames) && literalNames[tokenType] != "" {
				tokenName = literalNames[tokenType]
			} else if tokenType < len(symbolicNames) && symbolicNames[tokenType] != "" {
				tokenName = symbolicNames[tokenType]
			} else {
				tokenName = fmt.Sprintf("TOKEN_%d", tokenType)
			}
		} else {
			tokenName = "EOF"
		}

		nodeText = fmt.Sprintf("TOKEN[%s]: %q", tokenName, token.GetText())
	case antlr.RuleNode:
		// Rule node
		ruleIndex := t.GetRuleContext().GetRuleIndex()
		ruleName := parser.GetRuleNames()[ruleIndex]
		nodeText = fmt.Sprintf("RULE[%s]", ruleName)
	default:
		nodeText = "UNKNOWN"
	}

	fmt.Printf("%s%s\n", indentStr, nodeText)

	// Recurse for children
	if ctx, ok := tree.(antlr.RuleContext); ok {
		for i := 0; i < ctx.GetChildCount(); i++ {
			printTreeLimited(ctx.GetChild(i), parser, indent+1, maxDepth)
		}
	}
}

// TestPrintASTWithListener demonstrates using a listener pattern
func TestPrintASTWithListener(t *testing.T) {
	// Read the example file
	content, err := os.ReadFile("../../moxie-intellij-plugin/example.x")
	if err != nil {
		t.Fatalf("Failed to read example.x: %v", err)
	}

	// Parse
	is := antlr.NewInputStream(string(content))
	lexer := NewMoxieLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewMoxieParser(stream)

	tree := parser.SourceFile()

	// Print summary using listener
	fmt.Println("\n=== AST SUMMARY ===")
	listener := &ASTSummaryListener{}
	antlr.NewParseTreeWalker().Walk(listener, tree)
	listener.PrintSummary()
}

// ASTSummaryListener collects statistics about the parsed AST
type ASTSummaryListener struct {
	*BaseMoxieListener

	packageName string
	imports     []string
	functions   []string
	types       []string
	constants   []string
	variables   []string
}

func (l *ASTSummaryListener) EnterPackageClause(ctx *PackageClauseContext) {
	l.packageName = ctx.IDENTIFIER().GetText()
}

func (l *ASTSummaryListener) EnterImportSpec(ctx *ImportSpecContext) {
	if str := ctx.String_(); str != nil {
		l.imports = append(l.imports, str.GetText())
	}
}

func (l *ASTSummaryListener) EnterFunctionDecl(ctx *FunctionDeclContext) {
	funcName := ctx.IDENTIFIER().GetText()
	l.functions = append(l.functions, funcName)
}

func (l *ASTSummaryListener) EnterTypeAlias(ctx *TypeAliasContext) {
	typeName := ctx.IDENTIFIER().GetText()
	l.types = append(l.types, typeName)
}

func (l *ASTSummaryListener) EnterTypeDef(ctx *TypeDefContext) {
	typeName := ctx.IDENTIFIER().GetText()
	l.types = append(l.types, typeName)
}

func (l *ASTSummaryListener) EnterConstSpec(ctx *ConstSpecContext) {
	if idList := ctx.IdentifierList(); idList != nil {
		for _, id := range idList.AllIDENTIFIER() {
			l.constants = append(l.constants, id.GetText())
		}
	}
}

func (l *ASTSummaryListener) EnterVarSpec(ctx *VarSpecContext) {
	if idList := ctx.IdentifierList(); idList != nil {
		for _, id := range idList.AllIDENTIFIER() {
			l.variables = append(l.variables, id.GetText())
		}
	}
}

func (l *ASTSummaryListener) PrintSummary() {
	fmt.Printf("Package: %s\n", l.packageName)
	fmt.Printf("\nImports (%d):\n", len(l.imports))
	for _, imp := range l.imports {
		fmt.Printf("  - %s\n", imp)
	}
	fmt.Printf("\nTypes (%d):\n", len(l.types))
	for _, t := range l.types {
		fmt.Printf("  - %s\n", t)
	}
	fmt.Printf("\nConstants (%d):\n", len(l.constants))
	for _, c := range l.constants {
		fmt.Printf("  - %s\n", c)
	}
	fmt.Printf("\nFunctions (%d):\n", len(l.functions))
	for _, f := range l.functions {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Printf("\nVariables (%d):\n", len(l.variables))
	for _, vr := range l.variables {
		fmt.Printf("  - %s\n", vr)
	}
}
