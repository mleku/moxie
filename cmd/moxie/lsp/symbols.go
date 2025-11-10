// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"go/ast"
)

// extractSymbols extracts symbols from a document
func extractSymbols(doc *Document) []Symbol {
	if doc.AST == nil {
		return nil
	}

	var symbols []Symbol

	// Extract package-level declarations
	for _, decl := range doc.AST.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Function declaration
			if d.Name != nil {
				pos := doc.Fset.Position(d.Name.Pos())
				end := doc.Fset.Position(d.Name.End())
				symbols = append(symbols, Symbol{
					Name: d.Name.Name,
					Kind: SymbolKindFunction,
					Location: Location{
						URI: doc.URI,
						Range: Range{
							Start: Position{Line: pos.Line - 1, Character: pos.Column - 1},
							End:   Position{Line: end.Line - 1, Character: end.Column - 1},
						},
					},
					Detail: "function",
				})
			}

		case *ast.GenDecl:
			// Type, const, var declarations
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					// Type declaration
					if s.Name != nil {
						pos := doc.Fset.Position(s.Name.Pos())
						end := doc.Fset.Position(s.Name.End())

						kind := SymbolKindStruct
						detail := "type"

						// Determine specific type kind
						switch s.Type.(type) {
						case *ast.StructType:
							kind = SymbolKindStruct
							detail = "struct"
						case *ast.InterfaceType:
							kind = SymbolKindInterface
							detail = "interface"
						}

						symbols = append(symbols, Symbol{
							Name: s.Name.Name,
							Kind: kind,
							Location: Location{
								URI: doc.URI,
								Range: Range{
									Start: Position{Line: pos.Line - 1, Character: pos.Column - 1},
									End:   Position{Line: end.Line - 1, Character: end.Column - 1},
								},
							},
							Detail: detail,
						})
					}

				case *ast.ValueSpec:
					// Const or var declaration
					for _, name := range s.Names {
						pos := doc.Fset.Position(name.Pos())
						end := doc.Fset.Position(name.End())

						kind := SymbolKindVariable
						detail := "var"

						// Check if it's a const
						if d.Tok.String() == "const" {
							kind = SymbolKindConstant
							detail = "const"
						}

						symbols = append(symbols, Symbol{
							Name: name.Name,
							Kind: kind,
							Location: Location{
								URI: doc.URI,
								Range: Range{
									Start: Position{Line: pos.Line - 1, Character: pos.Column - 1},
									End:   Position{Line: end.Line - 1, Character: end.Column - 1},
								},
							},
							Detail: detail,
						})
					}
				}
			}
		}
	}

	return symbols
}

// extractDocumentSymbols extracts hierarchical symbols for textDocument/documentSymbol
func extractDocumentSymbols(doc *Document) []DocumentSymbol {
	if doc.AST == nil {
		return nil
	}

	var symbols []DocumentSymbol

	// Extract package-level declarations
	for _, decl := range doc.AST.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Function declaration
			if d.Name != nil {
				namePos := doc.Fset.Position(d.Name.Pos())
				nameEnd := doc.Fset.Position(d.Name.End())
				declPos := doc.Fset.Position(d.Pos())
				declEnd := doc.Fset.Position(d.End())

				symbols = append(symbols, DocumentSymbol{
					Name: d.Name.Name,
					Kind: SymbolKindFunction,
					Range: Range{
						Start: Position{Line: declPos.Line - 1, Character: declPos.Column - 1},
						End:   Position{Line: declEnd.Line - 1, Character: declEnd.Column - 1},
					},
					SelectionRange: Range{
						Start: Position{Line: namePos.Line - 1, Character: namePos.Column - 1},
						End:   Position{Line: nameEnd.Line - 1, Character: nameEnd.Column - 1},
					},
				})
			}

		case *ast.GenDecl:
			// Type, const, var declarations
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					// Type declaration
					if s.Name != nil {
						namePos := doc.Fset.Position(s.Name.Pos())
						nameEnd := doc.Fset.Position(s.Name.End())
						specPos := doc.Fset.Position(s.Pos())
						specEnd := doc.Fset.Position(s.End())

						kind := SymbolKindStruct
						var children []DocumentSymbol

						// Determine specific type kind and extract fields
						switch typeNode := s.Type.(type) {
						case *ast.StructType:
							kind = SymbolKindStruct
							if typeNode.Fields != nil {
								for _, field := range typeNode.Fields.List {
									for _, fieldName := range field.Names {
										fieldPos := doc.Fset.Position(fieldName.Pos())
										fieldEnd := doc.Fset.Position(fieldName.End())
										children = append(children, DocumentSymbol{
											Name: fieldName.Name,
											Kind: SymbolKindField,
											Range: Range{
												Start: Position{Line: fieldPos.Line - 1, Character: fieldPos.Column - 1},
												End:   Position{Line: fieldEnd.Line - 1, Character: fieldEnd.Column - 1},
											},
											SelectionRange: Range{
												Start: Position{Line: fieldPos.Line - 1, Character: fieldPos.Column - 1},
												End:   Position{Line: fieldEnd.Line - 1, Character: fieldEnd.Column - 1},
											},
										})
									}
								}
							}
						case *ast.InterfaceType:
							kind = SymbolKindInterface
							if typeNode.Methods != nil {
								for _, method := range typeNode.Methods.List {
									for _, methodName := range method.Names {
										methodPos := doc.Fset.Position(methodName.Pos())
										methodEnd := doc.Fset.Position(methodName.End())
										children = append(children, DocumentSymbol{
											Name: methodName.Name,
											Kind: SymbolKindMethod,
											Range: Range{
												Start: Position{Line: methodPos.Line - 1, Character: methodPos.Column - 1},
												End:   Position{Line: methodEnd.Line - 1, Character: methodEnd.Column - 1},
											},
											SelectionRange: Range{
												Start: Position{Line: methodPos.Line - 1, Character: methodPos.Column - 1},
												End:   Position{Line: methodEnd.Line - 1, Character: methodEnd.Column - 1},
											},
										})
									}
								}
							}
						}

						symbols = append(symbols, DocumentSymbol{
							Name: s.Name.Name,
							Kind: kind,
							Range: Range{
								Start: Position{Line: specPos.Line - 1, Character: specPos.Column - 1},
								End:   Position{Line: specEnd.Line - 1, Character: specEnd.Column - 1},
							},
							SelectionRange: Range{
								Start: Position{Line: namePos.Line - 1, Character: namePos.Column - 1},
								End:   Position{Line: nameEnd.Line - 1, Character: nameEnd.Column - 1},
							},
							Children: children,
						})
					}

				case *ast.ValueSpec:
					// Const or var declaration
					for _, name := range s.Names {
						namePos := doc.Fset.Position(name.Pos())
						nameEnd := doc.Fset.Position(name.End())
						specPos := doc.Fset.Position(s.Pos())
						specEnd := doc.Fset.Position(s.End())

						kind := SymbolKindVariable
						if d.Tok.String() == "const" {
							kind = SymbolKindConstant
						}

						symbols = append(symbols, DocumentSymbol{
							Name: name.Name,
							Kind: kind,
							Range: Range{
								Start: Position{Line: specPos.Line - 1, Character: specPos.Column - 1},
								End:   Position{Line: specEnd.Line - 1, Character: specEnd.Column - 1},
							},
							SelectionRange: Range{
								Start: Position{Line: namePos.Line - 1, Character: namePos.Column - 1},
								End:   Position{Line: nameEnd.Line - 1, Character: nameEnd.Column - 1},
							},
						})
					}
				}
			}
		}
	}

	return symbols
}
