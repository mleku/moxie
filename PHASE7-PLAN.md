# Phase 7: Tooling - Implementation Plan

**Status**: 🚧 In Progress
**Started**: 2025-11-10
**Dependencies**: Phases 1-6 (All core language features complete)

## Overview

Phase 7 focuses on developer tooling to enhance the Moxie development experience. Since Moxie is a transpiler to Go, we can leverage existing Go tooling while adding Moxie-specific enhancements.

## Objectives

1. **Improve developer productivity** with better tooling
2. **Provide IDE/editor support** via LSP
3. **Enhance error reporting** with source mapping
4. **Add code quality tools** (formatter, linter)
5. **Streamline build process** with watch mode and caching

## Priority Order

Based on immediate value and implementation complexity:

### Priority 1: Essential Tools (High Value, Low-Medium Complexity)
- ✅ **Source mapping** for error messages (.go errors → .mx files)
- ✅ **Formatter** (`moxie fmt`) - Format Moxie source code
- ✅ **Watch mode** (`moxie watch`) - Auto-rebuild on file changes

### Priority 2: Quality Tools (High Value, Medium Complexity)
- ⏳ **Linter** (`moxie vet`) - Static analysis for Moxie-specific patterns
- ⏳ **Enhanced error messages** - Better transpiler diagnostics
- ⏳ **Build caching** - Incremental transpilation

### Priority 3: Advanced Tools (Medium Value, High Complexity)
- ⏳ **LSP server** - Language Server Protocol for IDE integration
- ⏳ **Documentation generator** (`moxie doc`) - Generate docs from code
- ⏳ **Dependency analysis** (`moxie mod`) - Module management helpers

### Priority 4: Optional Enhancements (Nice to Have)
- ⏳ **REPL** - Interactive Moxie environment
- ⏳ **Playground** - Web-based Moxie playground
- ⏳ **Profiling tools** - Performance analysis helpers

## Detailed Implementation Plans

---

### 7.1: Source Mapping for Error Messages

**Goal**: Map Go compiler errors back to original .mx source files

**Current Problem**:
```
./test_example.go:15:2: undefined: foo
```
User sees `.go` file and line number, not original `.mx` file.

**Solution**:
1. Track line number mapping during transpilation
2. Intercept Go compiler output
3. Translate error locations back to `.mx` files

**Implementation**:
- Create `sourcemap.go` to track transformations
- Enhance transpiler to record line mappings
- Add error message post-processor in build commands
- Format: `test_example.mx:12:2: undefined: foo`

**Deliverables**:
- `cmd/moxie/sourcemap.go` (150-200 lines)
- Updated build/run/test commands with error mapping
- Tests for various error scenarios

**Estimated Effort**: 1-2 days

---

### 7.2: Formatter (`moxie fmt`)

**Goal**: Format Moxie source code consistently

**Approach**:
1. Leverage Go's `go/format` package
2. Add Moxie-specific formatting rules
3. Respect `.mx` file syntax (channel literals, endianness tuples)

**Commands**:
```bash
moxie fmt <file>           # Format single file
moxie fmt <directory>      # Format all .mx files in directory
moxie fmt -w <path>        # Write changes in-place
moxie fmt -d <path>        # Display diff instead of writing
moxie fmt -l <path>        # List files that need formatting
```

**Implementation**:
- Parse `.mx` file to AST
- Apply Go formatting via `go/format`
- Reverse-preprocess to restore Moxie syntax
- Write formatted output

**Special Handling**:
- Channel literals: `&chan T{}` formatting
- Endianness tuples: `(*[]T, BigEndian)(x)` formatting
- String literals: `*[]byte` display preferences

**Deliverables**:
- `cmd/moxie/format.go` (200-300 lines)
- Format command in main.go
- Tests for various code patterns
- Documentation for formatting rules

**Estimated Effort**: 2-3 days

---

### 7.3: Watch Mode (`moxie watch`)

**Goal**: Auto-rebuild and rerun on file changes

**Commands**:
```bash
moxie watch <directory>              # Watch and rebuild
moxie watch --run <file>             # Watch and auto-run
moxie watch --test <package>         # Watch and auto-test
moxie watch --exec <command>         # Watch and run custom command
```

**Implementation**:
- Use `fsnotify` package for file system watching
- Debounce file change events (avoid rapid rebuilds)
- Clear terminal and show build status
- Colorize output (errors in red, success in green)

**Features**:
- Watch `.mx` files recursively
- Ignore build directories and `.git`
- Show build time and status
- Handle keyboard interrupts gracefully

**Deliverables**:
- `cmd/moxie/watch.go` (300-400 lines)
- Watch command in main.go
- Configuration file support (`.moxie-watch.toml`)
- Tests for file watching logic

**Dependencies**: `github.com/fsnotify/fsnotify`

**Estimated Effort**: 2-3 days

---

### 7.4: Linter (`moxie vet`)

**Goal**: Static analysis for common Moxie errors and anti-patterns

**Checks**:
1. **Memory Management**:
   - Unused `clone()` allocations
   - Missing `free()` calls for allocated resources
   - Double `free()` detection

2. **Channel Safety**:
   - Unbuffered channels that might deadlock
   - Channel leaks (not closed)
   - Send on closed channel

3. **Type Safety**:
   - Unsafe type coercions
   - Endianness mismatches in network code
   - Integer overflow risks

4. **Const Correctness**:
   - Const values that should be variables
   - Variables that could be const

5. **Error Handling**:
   - Unchecked errors
   - Error shadowing

**Commands**:
```bash
moxie vet <package>           # Vet single package
moxie vet ./...               # Vet all packages
moxie vet --checks=memory     # Run specific checks
```

**Implementation**:
- AST-based analysis using `go/ast` and `go/types`
- Plugin system for check modules
- JSON output for IDE integration
- Configurable severity levels

**Deliverables**:
- `cmd/moxie/vet/` package (500-800 lines)
  - `vet.go` - Main vet command
  - `memory.go` - Memory management checks
  - `channels.go` - Channel safety checks
  - `types.go` - Type safety checks
  - `const.go` - Const correctness checks
- Vet command in main.go
- Configuration file (`.moxie-vet.toml`)
- Tests for each check category

**Estimated Effort**: 5-7 days

---

### 7.5: LSP Server

**Goal**: Language Server Protocol implementation for IDE support

**Features**:
- **Diagnostics**: Real-time error checking
- **Completion**: Code completion for Moxie syntax
- **Hover**: Type information and documentation
- **Goto Definition**: Jump to definitions
- **Find References**: Find all references to symbol
- **Rename**: Safe symbol renaming
- **Formatting**: Integrate with `moxie fmt`
- **Code Actions**: Quick fixes for common issues

**Implementation Strategy**:
1. Use `golang.org/x/tools/gopls` as reference
2. Build on top of Go LSP features
3. Add Moxie-specific extensions

**Architecture**:
```
moxie-lsp (binary)
├── server/          LSP server implementation
│   ├── protocol.go  LSP protocol handling
│   ├── handlers.go  Request handlers
│   └── cache.go     File/AST caching
├── analysis/        Code analysis
│   ├── completion.go
│   ├── hover.go
│   └── diagnostics.go
└── config/          Configuration
    └── settings.go
```

**Deliverables**:
- `cmd/moxie-lsp/` package (1500-2000 lines)
- Editor plugins:
  - VS Code extension
  - Vim/Neovim plugin
  - Emacs mode
- Documentation for setup
- Test suite for LSP features

**Dependencies**:
- `golang.org/x/tools/lsp`
- `github.com/sourcegraph/go-lsp`

**Estimated Effort**: 2-3 weeks

---

### 7.6: Documentation Generator (`moxie doc`)

**Goal**: Generate documentation from Moxie source code

**Features**:
- Extract doc comments from code
- Generate HTML/Markdown documentation
- Cross-reference symbols
- Include examples from tests

**Commands**:
```bash
moxie doc <package>              # Show package documentation
moxie doc <package>.<symbol>     # Show symbol documentation
moxie doc -html <package>        # Generate HTML docs
moxie doc -json <package>        # Output JSON for custom tools
```

**Implementation**:
- Parse doc comments using `go/doc`
- Generate documentation in various formats
- Include Moxie-specific syntax highlighting
- Cross-link to Go stdlib documentation

**Deliverables**:
- `cmd/moxie/doc.go` (400-600 lines)
- HTML templates for documentation
- Markdown generator
- Tests for doc extraction

**Estimated Effort**: 3-4 days

---

### 7.7: Build Caching

**Goal**: Speed up builds with incremental transpilation

**Features**:
- Cache transpiled `.go` files
- Detect changes in `.mx` files
- Invalidate cache on syntax changes
- Share cache across projects (optional)

**Implementation**:
- Hash-based caching (SHA256 of source + dependencies)
- Store in `~/.moxie/cache/` or `.moxie-cache/`
- Respect `MOXIE_CACHE` environment variable
- `moxie clean` command to clear cache

**Cache Structure**:
```
.moxie-cache/
├── transpiled/     Cached .go files
├── metadata/       Dependency metadata
└── checksums/      File checksums
```

**Deliverables**:
- `cmd/moxie/cache.go` (300-400 lines)
- Updated build commands to use cache
- `moxie clean` command
- Tests for cache invalidation

**Estimated Effort**: 3-4 days

---

## Implementation Roadmap

### Phase 7.1: Essential Tools (Week 1-2)
- ✅ Source mapping for error messages
- ✅ Basic formatter (`moxie fmt`)
- ✅ Watch mode (`moxie watch`)

**Milestone**: Developers can format code and get fast feedback

### Phase 7.2: Quality Tools (Week 3-4)
- ⏳ Linter (`moxie vet`)
- ⏳ Enhanced error messages
- ⏳ Build caching

**Milestone**: Code quality and build performance improved

### Phase 7.3: Advanced Tools (Week 5-8)
- ⏳ LSP server MVP
- ⏳ VS Code extension
- ⏳ Documentation generator

**Milestone**: IDE integration available

### Phase 7.4: Polish (Week 9+)
- ⏳ Additional editor plugins
- ⏳ Performance optimizations
- ⏳ Documentation and tutorials

**Milestone**: Production-ready tooling

---

## Testing Strategy

### Unit Tests
- Test each tool independently
- Mock file system operations
- Test error handling paths

### Integration Tests
- Test tool interactions
- End-to-end workflows
- Real-world code examples

### Performance Tests
- Measure build times with/without caching
- LSP response times
- Memory usage

### User Acceptance Tests
- Developer workflow testing
- Editor integration testing
- Documentation clarity

---

## Success Metrics

1. **Build Performance**:
   - 50%+ faster builds with caching
   - Sub-second watch mode rebuilds

2. **Developer Experience**:
   - Consistent code formatting across projects
   - Real-time error feedback in editors
   - Accurate error messages pointing to .mx files

3. **Code Quality**:
   - Fewer runtime errors caught by linter
   - Improved const correctness
   - Better memory management patterns

4. **Adoption**:
   - IDE plugins installed and used
   - Formatter integrated in CI/CD
   - Positive developer feedback

---

## Dependencies

### Go Packages
- `golang.org/x/tools` - Go tooling libraries
- `github.com/fsnotify/fsnotify` - File system watching
- `github.com/sourcegraph/go-lsp` - LSP implementation
- `github.com/fatih/color` - Terminal colors

### External Tools
- Go toolchain (1.21+)
- Git (for version control integration)

---

## Notes

- **Leverage Go ecosystem**: Reuse existing Go tools where possible
- **Moxie-specific enhancements**: Focus on features unique to Moxie
- **Incremental approach**: Ship useful tools early, iterate based on feedback
- **Documentation first**: Every tool needs clear documentation
- **Backward compatibility**: Maintain compatibility with existing .mx files

---

## Future Enhancements (Post-Phase 7)

- **REPL**: Interactive Moxie environment
- **Playground**: Web-based code editor
- **Package registry**: Centralized Moxie package repository
- **Code coverage**: Test coverage visualization
- **Benchmarking**: Performance comparison tools
- **Migration tools**: Automated Go → Moxie conversion

---

## References

- Go tooling: `golang.org/x/tools`
- LSP specification: `microsoft.github.io/language-server-protocol`
- gopls source: `golang.org/x/tools/gopls`
- File watching: `github.com/fsnotify/fsnotify`
