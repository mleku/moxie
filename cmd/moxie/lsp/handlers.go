// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"os/exec"
	"strings"
)

// handleDocumentSymbol handles textDocument/documentSymbol request
func (s *Server) handleDocumentSymbol(ctx context.Context, req *Request) (interface{}, error) {
	var params DocumentSymbolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.mu.RLock()
	doc, ok := s.documents[params.TextDocument.URI]
	s.mu.RUnlock()

	if !ok {
		return []DocumentSymbol{}, nil
	}

	if doc.AST == nil {
		return []DocumentSymbol{}, nil
	}

	symbols := extractDocumentSymbols(doc)
	return symbols, nil
}

// handleWorkspaceSymbol handles workspace/symbol request
func (s *Server) handleWorkspaceSymbol(ctx context.Context, req *Request) (interface{}, error) {
	var params WorkspaceSymbolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	query := strings.ToLower(params.Query)
	var results []SymbolInformation

	s.workspace.Index.mu.RLock()
	defer s.workspace.Index.mu.RUnlock()

	// Search all symbols
	for _, symbols := range s.workspace.Index.symbols {
		for _, sym := range symbols {
			if strings.Contains(strings.ToLower(sym.Name), query) {
				results = append(results, SymbolInformation{
					Name:     sym.Name,
					Kind:     sym.Kind,
					Location: sym.Location,
				})
			}
		}
	}

	return results, nil
}

// handleHover handles textDocument/hover request
func (s *Server) handleHover(ctx context.Context, req *Request) (interface{}, error) {
	var params HoverParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.mu.RLock()
	doc, ok := s.documents[params.TextDocument.URI]
	s.mu.RUnlock()

	if !ok || doc.AST == nil {
		return nil, nil
	}

	// Find symbol at position
	var hoverText string
	ast.Inspect(doc.AST, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		// Check if node contains the position
		if doc.Fset.Position(n.Pos()).Line != params.Position.Line+1 {
			return true
		}

		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Name != nil {
				hoverText = fmt.Sprintf("```moxie\nfunc %s\n```", node.Name.Name)
				return false
			}
		case *ast.TypeSpec:
			if node.Name != nil {
				hoverText = fmt.Sprintf("```moxie\ntype %s\n```", node.Name.Name)
				return false
			}
		case *ast.ValueSpec:
			if len(node.Names) > 0 {
				hoverText = fmt.Sprintf("```moxie\nvar %s\n```", node.Names[0].Name)
				return false
			}
		case *ast.Ident:
			hoverText = fmt.Sprintf("```moxie\n%s\n```", node.Name)
			return false
		}

		return true
	})

	if hoverText == "" {
		return nil, nil
	}

	return &Hover{
		Contents: MarkupContent{
			Kind:  MarkupKindMarkdown,
			Value: hoverText,
		},
	}, nil
}

// handleDefinition handles textDocument/definition request
func (s *Server) handleDefinition(ctx context.Context, req *Request) (interface{}, error) {
	var params DefinitionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.mu.RLock()
	doc, ok := s.documents[params.TextDocument.URI]
	s.mu.RUnlock()

	if !ok || doc.AST == nil {
		return nil, nil
	}

	// Find identifier at position
	var identName string
	ast.Inspect(doc.AST, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		if ident, ok := n.(*ast.Ident); ok {
			pos := doc.Fset.Position(ident.Pos())
			if pos.Line == params.Position.Line+1 {
				identName = ident.Name
				return false
			}
		}
		return true
	})

	if identName == "" {
		return nil, nil
	}

	// Find definition
	s.workspace.Index.mu.RLock()
	defer s.workspace.Index.mu.RUnlock()

	for _, symbols := range s.workspace.Index.symbols {
		for _, sym := range symbols {
			if sym.Name == identName {
				return sym.Location, nil
			}
		}
	}

	return nil, nil
}

// handleReferences handles textDocument/references request
func (s *Server) handleReferences(ctx context.Context, req *Request) (interface{}, error) {
	var params ReferenceParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.mu.RLock()
	doc, ok := s.documents[params.TextDocument.URI]
	s.mu.RUnlock()

	if !ok || doc.AST == nil {
		return []Location{}, nil
	}

	// Find identifier at position
	var identName string
	ast.Inspect(doc.AST, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		if ident, ok := n.(*ast.Ident); ok {
			pos := doc.Fset.Position(ident.Pos())
			if pos.Line == params.Position.Line+1 {
				identName = ident.Name
				return false
			}
		}
		return true
	})

	if identName == "" {
		return []Location{}, nil
	}

	// Find all references
	var locations []Location

	s.mu.RLock()
	for uri, d := range s.documents {
		if d.AST == nil {
			continue
		}

		ast.Inspect(d.AST, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok {
				if ident.Name == identName {
					pos := d.Fset.Position(ident.Pos())
					locations = append(locations, Location{
						URI: uri,
						Range: Range{
							Start: Position{Line: pos.Line - 1, Character: pos.Column - 1},
							End:   Position{Line: pos.Line - 1, Character: pos.Column - 1 + len(identName)},
						},
					})
				}
			}
			return true
		})
	}
	s.mu.RUnlock()

	return locations, nil
}

// handleCompletion handles textDocument/completion request
func (s *Server) handleCompletion(ctx context.Context, req *Request) (interface{}, error) {
	var params CompletionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	// Return keywords and common identifiers
	items := []CompletionItem{
		// Keywords
		{Label: "func", Kind: CompletionItemKindKeyword},
		{Label: "type", Kind: CompletionItemKindKeyword},
		{Label: "const", Kind: CompletionItemKindKeyword},
		{Label: "var", Kind: CompletionItemKindKeyword},
		{Label: "if", Kind: CompletionItemKindKeyword},
		{Label: "else", Kind: CompletionItemKindKeyword},
		{Label: "for", Kind: CompletionItemKindKeyword},
		{Label: "switch", Kind: CompletionItemKindKeyword},
		{Label: "case", Kind: CompletionItemKindKeyword},
		{Label: "default", Kind: CompletionItemKindKeyword},
		{Label: "return", Kind: CompletionItemKindKeyword},
		{Label: "break", Kind: CompletionItemKindKeyword},
		{Label: "continue", Kind: CompletionItemKindKeyword},
		{Label: "goto", Kind: CompletionItemKindKeyword},
		{Label: "package", Kind: CompletionItemKindKeyword},
		{Label: "import", Kind: CompletionItemKindKeyword},
		{Label: "struct", Kind: CompletionItemKindKeyword},
		{Label: "interface", Kind: CompletionItemKindKeyword},
		{Label: "map", Kind: CompletionItemKindKeyword},
		{Label: "chan", Kind: CompletionItemKindKeyword},
		{Label: "select", Kind: CompletionItemKindKeyword},
		{Label: "defer", Kind: CompletionItemKindKeyword},
		{Label: "go", Kind: CompletionItemKindKeyword},
		{Label: "range", Kind: CompletionItemKindKeyword},
		{Label: "true", Kind: CompletionItemKindKeyword},
		{Label: "false", Kind: CompletionItemKindKeyword},
		{Label: "nil", Kind: CompletionItemKindKeyword},

		// Builtin functions
		{Label: "clone", Kind: CompletionItemKindFunction, Detail: "clone(x) - Deep copy value"},
		{Label: "free", Kind: CompletionItemKindFunction, Detail: "free(x) - Free allocated memory"},
		{Label: "grow", Kind: CompletionItemKindFunction, Detail: "grow(slice, n) - Grow slice capacity"},
		{Label: "append", Kind: CompletionItemKindFunction},
		{Label: "len", Kind: CompletionItemKindFunction},
		{Label: "cap", Kind: CompletionItemKindFunction},
		{Label: "make", Kind: CompletionItemKindFunction},
		{Label: "new", Kind: CompletionItemKindFunction},
		{Label: "delete", Kind: CompletionItemKindFunction},
		{Label: "copy", Kind: CompletionItemKindFunction},
		{Label: "clear", Kind: CompletionItemKindFunction},
		{Label: "close", Kind: CompletionItemKindFunction},
		{Label: "panic", Kind: CompletionItemKindFunction},
		{Label: "recover", Kind: CompletionItemKindFunction},
		{Label: "print", Kind: CompletionItemKindFunction},
		{Label: "println", Kind: CompletionItemKindFunction},
	}

	// Add symbols from workspace
	s.workspace.Index.mu.RLock()
	for _, symbols := range s.workspace.Index.symbols {
		for _, sym := range symbols {
			kind := CompletionItemKindText
			switch sym.Kind {
			case SymbolKindFunction:
				kind = CompletionItemKindFunction
			case SymbolKindVariable:
				kind = CompletionItemKindVariable
			case SymbolKindConstant:
				kind = CompletionItemKindConstant
			case SymbolKindStruct:
				kind = CompletionItemKindStruct
			case SymbolKindInterface:
				kind = CompletionItemKindInterface
			}

			items = append(items, CompletionItem{
				Label:  sym.Name,
				Kind:   kind,
				Detail: sym.Detail,
			})
		}
	}
	s.workspace.Index.mu.RUnlock()

	return CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

// handleFormatting handles textDocument/formatting request
func (s *Server) handleFormatting(ctx context.Context, req *Request) (interface{}, error) {
	var params DocumentFormattingParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.mu.RLock()
	doc, ok := s.documents[params.TextDocument.URI]
	s.mu.RUnlock()

	if !ok {
		return nil, nil
	}

	// Use moxie fmt to format
	path := uriToPath(params.TextDocument.URI)
	cmd := exec.Command("moxie", "fmt", path)
	output, err := cmd.Output()
	if err != nil {
		s.logger.Printf("Format error: %v", err)
		return nil, err
	}

	// Calculate range for entire document
	lines := strings.Split(doc.Content, "\n")
	endLine := len(lines) - 1
	endChar := 0
	if endLine >= 0 {
		endChar = len(lines[endLine])
	}

	return []TextEdit{
		{
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: endLine, Character: endChar},
			},
			NewText: string(output),
		},
	}, nil
}
