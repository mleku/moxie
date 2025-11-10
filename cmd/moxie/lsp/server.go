// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Server is the LSP server for Moxie
type Server struct {
	mu          sync.RWMutex
	documents   map[string]*Document
	workspace   *Workspace
	conn        *Connection
	initialized bool
	rootURI     string
	logger      *log.Logger
}

// Document represents an open document
type Document struct {
	URI     string
	Version int32
	Content string
	AST     *ast.File
	Fset    *token.FileSet
}

// Workspace represents the workspace state
type Workspace struct {
	RootURI string
	Index   *SymbolIndex
}

// SymbolIndex indexes symbols in the workspace
type SymbolIndex struct {
	mu      sync.RWMutex
	symbols map[string][]Symbol
}

// Symbol represents a code symbol
type Symbol struct {
	Name     string
	Kind     SymbolKind
	Location Location
	Detail   string
}

// NewServer creates a new LSP server
func NewServer() *Server {
	return &Server{
		documents: make(map[string]*Document),
		workspace: &Workspace{
			Index: NewSymbolIndex(),
		},
		logger: log.New(os.Stderr, "[moxie-lsp] ", log.LstdFlags),
	}
}

// NewSymbolIndex creates a new symbol index
func NewSymbolIndex() *SymbolIndex {
	return &SymbolIndex{
		symbols: make(map[string][]Symbol),
	}
}

// Start starts the LSP server
func (s *Server) Start(ctx context.Context, in io.Reader, out io.Writer) error {
	s.conn = NewConnection(in, out, s.logger)
	s.registerHandlers()

	s.logger.Println("Moxie LSP server starting...")
	return s.conn.Run(ctx)
}

// registerHandlers registers all LSP request handlers
func (s *Server) registerHandlers() {
	// Lifecycle
	s.conn.Handle("initialize", s.handleInitialize)
	s.conn.Handle("initialized", s.handleInitialized)
	s.conn.Handle("shutdown", s.handleShutdown)
	s.conn.Handle("exit", s.handleExit)

	// Document synchronization
	s.conn.Handle("textDocument/didOpen", s.handleDidOpen)
	s.conn.Handle("textDocument/didChange", s.handleDidChange)
	s.conn.Handle("textDocument/didSave", s.handleDidSave)
	s.conn.Handle("textDocument/didClose", s.handleDidClose)

	// Language features
	s.conn.Handle("textDocument/documentSymbol", s.handleDocumentSymbol)
	s.conn.Handle("textDocument/hover", s.handleHover)
	s.conn.Handle("textDocument/definition", s.handleDefinition)
	s.conn.Handle("textDocument/references", s.handleReferences)
	s.conn.Handle("textDocument/completion", s.handleCompletion)
	s.conn.Handle("textDocument/formatting", s.handleFormatting)
	s.conn.Handle("workspace/symbol", s.handleWorkspaceSymbol)
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(ctx context.Context, req *Request) (interface{}, error) {
	var params InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.rootURI = params.RootURI
	if s.workspace != nil {
		s.workspace.RootURI = params.RootURI
	}
	s.mu.Unlock()

	s.logger.Printf("Initialize request from %s", params.RootURI)

	// Return server capabilities
	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: TextDocumentSyncOptions{
				OpenClose: true,
				Change:    TextDocumentSyncKindFull,
			},
			DocumentSymbolProvider: true,
			HoverProvider:          true,
			DefinitionProvider:     true,
			ReferencesProvider:     true,
			CompletionProvider: &CompletionOptions{
				TriggerCharacters: []string{".", ":"},
			},
			DocumentFormattingProvider: true,
			WorkspaceSymbolProvider:    true,
		},
		ServerInfo: &ServerInfo{
			Name:    "moxie-lsp",
			Version: "0.12.0",
		},
	}, nil
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized(ctx context.Context, req *Request) (interface{}, error) {
	s.mu.Lock()
	s.initialized = true
	s.mu.Unlock()

	s.logger.Println("Client initialized")

	// Index workspace
	if s.rootURI != "" {
		go s.indexWorkspace()
	}

	return nil, nil
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(ctx context.Context, req *Request) (interface{}, error) {
	s.logger.Println("Shutdown request received")
	return nil, nil
}

// handleExit handles the exit notification
func (s *Server) handleExit(ctx context.Context, req *Request) (interface{}, error) {
	s.logger.Println("Exit notification received")
	os.Exit(0)
	return nil, nil
}

// handleDidOpen handles document open notification
func (s *Server) handleDidOpen(ctx context.Context, req *Request) (interface{}, error) {
	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	uri := params.TextDocument.URI
	content := params.TextDocument.Text
	version := params.TextDocument.Version

	s.logger.Printf("Document opened: %s (version %d)", uri, version)

	doc := &Document{
		URI:     uri,
		Version: version,
		Content: content,
		Fset:    token.NewFileSet(),
	}

	// Parse document
	var err error
	doc.AST, err = parser.ParseFile(doc.Fset, uriToPath(uri), content, parser.ParseComments)
	if err != nil {
		s.logger.Printf("Parse error: %v", err)
		// Send diagnostics for parse errors
		s.publishDiagnostics(uri, []Diagnostic{
			{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 1},
				},
				Severity: DiagnosticSeverityError,
				Message:  err.Error(),
				Source:   "moxie-parser",
			},
		})
	}

	s.mu.Lock()
	s.documents[uri] = doc
	s.mu.Unlock()

	// Index symbols
	if doc.AST != nil {
		s.indexDocument(doc)
	}

	return nil, nil
}

// handleDidChange handles document change notification
func (s *Server) handleDidChange(ctx context.Context, req *Request) (interface{}, error) {
	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	uri := params.TextDocument.URI

	s.mu.Lock()
	doc, ok := s.documents[uri]
	if !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	// Update content (full sync)
	if len(params.ContentChanges) > 0 {
		doc.Content = params.ContentChanges[0].Text
		doc.Version = params.TextDocument.Version
	}
	s.mu.Unlock()

	// Reparse
	doc.Fset = token.NewFileSet()
	var err error
	doc.AST, err = parser.ParseFile(doc.Fset, uriToPath(uri), doc.Content, parser.ParseComments)
	if err != nil {
		s.logger.Printf("Parse error: %v", err)
		s.publishDiagnostics(uri, []Diagnostic{
			{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 1},
				},
				Severity: DiagnosticSeverityError,
				Message:  err.Error(),
				Source:   "moxie-parser",
			},
		})
	} else {
		// Clear diagnostics if parse succeeded
		s.publishDiagnostics(uri, []Diagnostic{})
		s.indexDocument(doc)
	}

	return nil, nil
}

// handleDidSave handles document save notification
func (s *Server) handleDidSave(ctx context.Context, req *Request) (interface{}, error) {
	var params DidSaveTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	s.logger.Printf("Document saved: %s", params.TextDocument.URI)
	return nil, nil
}

// handleDidClose handles document close notification
func (s *Server) handleDidClose(ctx context.Context, req *Request) (interface{}, error) {
	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	uri := params.TextDocument.URI

	s.mu.Lock()
	delete(s.documents, uri)
	s.mu.Unlock()

	s.logger.Printf("Document closed: %s", uri)
	return nil, nil
}

// publishDiagnostics sends diagnostics to the client
func (s *Server) publishDiagnostics(uri string, diagnostics []Diagnostic) {
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	s.conn.Notify("textDocument/publishDiagnostics", params)
}

// indexWorkspace indexes all Moxie files in the workspace
func (s *Server) indexWorkspace() {
	if s.rootURI == "" {
		return
	}

	rootPath := uriToPath(s.rootURI)
	s.logger.Printf("Indexing workspace: %s", rootPath)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() {
			// Skip hidden directories and build artifacts
			if info.Name()[0] == '.' || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Index .mx and .x files
		if filepath.Ext(path) == ".mx" || filepath.Ext(path) == ".x" {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, path, content, parser.ParseComments)
			if err != nil {
				return nil
			}

			doc := &Document{
				URI:     pathToURI(path),
				Content: string(content),
				AST:     astFile,
				Fset:    fset,
			}

			s.indexDocument(doc)
		}

		return nil
	})

	if err != nil {
		s.logger.Printf("Workspace indexing error: %v", err)
	} else {
		s.logger.Println("Workspace indexing complete")
	}
}

// indexDocument indexes symbols in a document
func (s *Server) indexDocument(doc *Document) {
	if doc.AST == nil {
		return
	}

	symbols := extractSymbols(doc)

	s.workspace.Index.mu.Lock()
	s.workspace.Index.symbols[doc.URI] = symbols
	s.workspace.Index.mu.Unlock()
}

// uriToPath converts a file URI to a file path
func uriToPath(uri string) string {
	// Simple conversion: file:///path -> /path
	if len(uri) > 7 && uri[:7] == "file://" {
		return uri[7:]
	}
	return uri
}

// pathToURI converts a file path to a file URI
func pathToURI(path string) string {
	return "file://" + path
}
