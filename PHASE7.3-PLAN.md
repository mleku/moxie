# Phase 7.3: Advanced Tools - LSP & IDE Integration

**Started**: 2025-11-10
**Status**: In Progress
**Dependencies**: Phase 7.1, 7.2

## Overview

Phase 7.3 adds advanced IDE integration through:
1. Language Server Protocol (LSP) server for Moxie
2. VS Code extension with full IDE support
3. Symbol indexing and navigation
4. Real-time diagnostics
5. Code intelligence features

## Goals

### Primary Goals
- ✅ Full LSP server implementation
- ✅ VS Code extension with syntax highlighting
- ✅ Symbol indexing (document and workspace)
- ✅ Real-time diagnostics (integration with `moxie vet`)
- ✅ Go to definition
- ✅ Find references
- ✅ Hover information
- ✅ Code completion

### Secondary Goals
- Rename refactoring
- Code actions (quick fixes)
- Formatting integration (`moxie fmt`)
- Debugging support
- Semantic tokens (advanced syntax highlighting)

## Architecture

### LSP Server (`moxie lsp`)

```
cmd/moxie/lsp/
├── server.go       # LSP server main implementation
├── protocol.go     # LSP protocol types and messages
├── handlers.go     # LSP request handlers
├── symbols.go      # Symbol indexing and search
├── diagnostics.go  # Diagnostic integration with vet
├── completion.go   # Code completion
└── analysis.go     # Code analysis utilities
```

**Features**:
1. **Symbol Provider**
   - Document symbols (functions, types, variables)
   - Workspace symbols (project-wide search)
   - Symbol hierarchy

2. **Navigation**
   - Go to definition
   - Find references
   - Find implementations
   - Go to type definition

3. **Diagnostics**
   - Real-time error checking
   - Integration with `moxie vet`
   - Syntax error reporting
   - Semantic error reporting

4. **Code Intelligence**
   - Hover information (type info, documentation)
   - Code completion (keywords, identifiers, snippets)
   - Signature help (function parameters)

5. **Code Actions**
   - Quick fixes for common errors
   - Refactoring suggestions
   - Import management

6. **Formatting**
   - Integration with `moxie fmt`
   - Format on save
   - Range formatting

### VS Code Extension

```
editors/vscode/
├── package.json           # Extension manifest
├── tsconfig.json          # TypeScript configuration
├── src/
│   ├── extension.ts       # Extension entry point
│   ├── client.ts          # LSP client
│   └── commands.ts        # VS Code commands
├── syntaxes/
│   └── moxie.tmLanguage.json  # TextMate grammar
├── language-configuration.json # Language config
├── snippets/
│   └── moxie.json         # Code snippets
└── README.md              # Extension documentation
```

**Features**:
1. **Syntax Highlighting**
   - TextMate grammar for Moxie
   - Semantic token provider (from LSP)
   - Moxie-specific syntax (channel literals, endianness tuples)

2. **Language Configuration**
   - Comment toggling
   - Bracket matching
   - Auto-indentation
   - Folding ranges

3. **LSP Integration**
   - LSP client connection
   - Server lifecycle management
   - Progress reporting

4. **Commands**
   - Build, Run, Test (integrated with Moxie CLI)
   - Format document/selection
   - Run vet (linter)
   - Clean cache

5. **Snippets**
   - Common code patterns
   - Function templates
   - Channel patterns
   - Error handling patterns

## Implementation Plan

### Week 1: LSP Server Core (Days 1-3)

**Day 1: LSP Infrastructure**
- [x] Create LSP server package structure
- [ ] Implement JSON-RPC 2.0 protocol handler
- [ ] Implement LSP lifecycle (initialize, initialized, shutdown, exit)
- [ ] Add basic logging and error handling
- [ ] Test basic LSP connection

**Day 2: Document Management**
- [ ] Implement document synchronization (didOpen, didChange, didClose)
- [ ] Document cache management
- [ ] Parse Moxie files into AST
- [ ] Track document versions
- [ ] Test document sync

**Day 3: Symbol Provider**
- [ ] Implement document symbol provider
- [ ] Extract symbols from AST (functions, types, variables, constants)
- [ ] Workspace symbol provider
- [ ] Symbol hierarchy
- [ ] Test symbol indexing

### Week 1: LSP Server Features (Days 4-7)

**Day 4: Navigation Features**
- [ ] Go to definition (uses AST and type tracking)
- [ ] Find references (workspace-wide search)
- [ ] Hover provider (type information, documentation)
- [ ] Test navigation

**Day 5: Diagnostics**
- [ ] Integration with `moxie vet`
- [ ] Syntax error reporting (from parser)
- [ ] Semantic error reporting (from type checker)
- [ ] Publish diagnostics on document change
- [ ] Test diagnostics

**Day 6: Code Completion**
- [ ] Keyword completion
- [ ] Identifier completion (from scope)
- [ ] Type completion
- [ ] Snippet completion
- [ ] Test completion

**Day 7: Formatting & Testing**
- [ ] Format document (integration with `moxie fmt`)
- [ ] Format range
- [ ] Integration tests for all LSP features
- [ ] Performance testing

### Week 2: VS Code Extension (Days 1-3)

**Day 1: Extension Setup**
- [ ] Create extension structure
- [ ] package.json configuration
- [ ] TypeScript setup
- [ ] Build system (esbuild/webpack)
- [ ] Test basic extension loading

**Day 2: Syntax Highlighting**
- [ ] Create TextMate grammar for Moxie
- [ ] Syntax rules for keywords, types, functions
- [ ] Moxie-specific syntax (channels, endianness)
- [ ] String literals, comments
- [ ] Test syntax highlighting

**Day 3: Language Configuration**
- [ ] Comment toggling (line and block)
- [ ] Bracket pairs and auto-closing
- [ ] Indentation rules
- [ ] Folding markers
- [ ] Test language features

### Week 2: VS Code Extension Integration (Days 4-7)

**Day 4: LSP Client**
- [ ] LSP client setup
- [ ] Server activation and deactivation
- [ ] Error handling and reconnection
- [ ] Test LSP connection

**Day 5: Commands**
- [ ] Build command
- [ ] Run command
- [ ] Test command
- [ ] Format command
- [ ] Vet command
- [ ] Clean command
- [ ] Test commands

**Day 6: Snippets & Polish**
- [ ] Code snippets for common patterns
- [ ] Extension icon and branding
- [ ] README and documentation
- [ ] Configuration options
- [ ] Test snippets

**Day 7: Package & Release**
- [ ] Package extension (.vsix)
- [ ] Test installation
- [ ] Publish to VS Code marketplace (optional)
- [ ] Documentation updates

## Technical Details

### LSP Server Implementation

**Protocol**: JSON-RPC 2.0 over stdio
**Transport**: stdin/stdout (standard LSP transport)
**Language**: Go (integrated with existing Moxie toolchain)

**Dependencies**:
- `golang.org/x/tools/internal/lsp/protocol` - LSP protocol types
- `golang.org/x/tools/internal/jsonrpc2` - JSON-RPC 2.0
- Existing Moxie parser and AST
- Existing type tracker
- Existing vet infrastructure

**Key Types**:
```go
type Server struct {
    documents map[DocumentURI]*Document
    workspace *Workspace
    conn      *jsonrpc2.Conn
}

type Document struct {
    URI     DocumentURI
    Version int32
    Content string
    AST     *ast.File
    Tokens  *token.FileSet
}

type Workspace struct {
    rootURI DocumentURI
    index   *SymbolIndex
    cache   *BuildCache
}
```

### VS Code Extension Implementation

**Language**: TypeScript
**Framework**: VS Code Extension API
**LSP Client**: `vscode-languageclient`

**Build**: esbuild for fast compilation
**Package**: vsce (VS Code Extension CLI)

**Key Features**:
```typescript
export function activate(context: vscode.ExtensionContext) {
    // Start LSP server
    const serverOptions: ServerOptions = {
        command: 'moxie',
        args: ['lsp']
    };

    // LSP client options
    const clientOptions: LanguageClientOptions = {
        documentSelector: [
            { scheme: 'file', language: 'moxie' }
        ]
    };

    // Create and start client
    const client = new LanguageClient(
        'moxie-lsp',
        'Moxie Language Server',
        serverOptions,
        clientOptions
    );

    client.start();
}
```

### TextMate Grammar

**Scopes**:
- `keyword.control.moxie` - if, for, switch, etc.
- `keyword.declaration.moxie` - func, type, const, var
- `storage.type.moxie` - int, string, bool, etc.
- `entity.name.function.moxie` - function names
- `entity.name.type.moxie` - type names
- `variable.other.moxie` - variables
- `constant.language.moxie` - true, false, nil
- `string.quoted.double.moxie` - strings
- `comment.line.double-slash.moxie` - comments

**Special Moxie Syntax**:
- `meta.channel-literal.moxie` - `&chan T{}`
- `meta.endianness-cast.moxie` - `(*[]T, Endian)(s)`
- `support.function.builtin.moxie` - clone(), free(), grow()

## Testing Strategy

### LSP Server Tests
1. **Unit Tests**
   - Symbol extraction from AST
   - Completion item generation
   - Diagnostic generation

2. **Integration Tests**
   - Full LSP request/response cycle
   - Document synchronization
   - Workspace indexing

3. **Performance Tests**
   - Large file handling
   - Workspace symbol search
   - Real-time diagnostics

### VS Code Extension Tests
1. **Manual Testing**
   - Install extension in VS Code
   - Open Moxie project
   - Test all features interactively

2. **Automated Tests**
   - Extension activation
   - Command registration
   - LSP client connection

## Success Criteria

### LSP Server
- ✅ Implements core LSP capabilities
- ✅ Symbol indexing works for all Moxie constructs
- ✅ Diagnostics integrate with `moxie vet`
- ✅ Performance: < 100ms for symbol search
- ✅ Performance: < 500ms for diagnostics

### VS Code Extension
- ✅ Syntax highlighting for all Moxie syntax
- ✅ IntelliSense (completion, hover, definition)
- ✅ Commands work correctly
- ✅ Extension size < 1MB
- ✅ Startup time < 1s

## Deliverables

### Code
1. `cmd/moxie/lsp/` - LSP server implementation (~1,500 lines)
2. `editors/vscode/` - VS Code extension (~800 lines)
3. Tests for LSP and extension

### Documentation
1. LSP server documentation
2. VS Code extension README
3. User guide for IDE features
4. Developer guide for LSP protocol

### Artifacts
1. VS Code extension package (.vsix)
2. Extension icon and assets
3. Example workspace configurations

## Timeline

- **Week 1**: LSP server implementation (7 days)
- **Week 2**: VS Code extension (7 days)
- **Total**: 14 days for complete Phase 7.3

## References

- [Language Server Protocol](https://microsoft.github.io/language-server-protocol/)
- [VS Code Extension API](https://code.visualstudio.com/api)
- [TextMate Grammar](https://macromates.com/manual/en/language_grammars)
- [golang.org/x/tools/gopls](https://pkg.go.dev/golang.org/x/tools/gopls) - Reference LSP implementation
