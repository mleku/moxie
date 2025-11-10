# Phase 11: Bootstrap - Rewriting Moxie in Moxie

**Started**: 2025-11-10
**Status**: In Progress
**Dependencies**: Phases 1-7 (All tooling complete)

## Overview

Phase 11 is the bootstrap phase where we rewrite the Moxie transpiler in Moxie itself. The current Go implementation serves as the bootstrap compiler. Once complete, Moxie will be self-hosting.

**Philosophy**: This is a bootstrap compiler. We prioritize correctness and completeness over optimization. Documentation exists in the repository. We skip Phase 8 (Optimization) and Phase 9 (Documentation) to focus on self-hosting.

## Goals

### Primary Goal
✅ **Self-Hosting**: Moxie compiler written in Moxie, compiled by itself

### Secondary Goals
- Maintain feature parity with Go implementation
- Demonstrate Moxie's capabilities for complex software
- Validate language design decisions
- Create production-ready v1.0.0 release

## Architecture

The Moxie implementation will mirror the Go structure but use Moxie idioms:

```
moxie-bootstrap/           # New Moxie implementation
├── cmd/
│   └── moxie/
│       ├── main.x         # Main entry point
│       ├── build.x        # Build command
│       ├── run.x          # Run command
│       ├── test.x         # Test command
│       ├── fmt.x          # Format command
│       ├── watch.x        # Watch command
│       ├── vet.x          # Vet command
│       ├── clean.x        # Clean command
│       ├── lsp.x          # LSP command
│       │
│       ├── parser.x       # Parser (wraps go/parser)
│       ├── syntax.x       # Syntax transformations
│       ├── preprocess.x   # Preprocessing
│       ├── typetrack.x    # Type tracking
│       │
│       ├── naming.x       # Name transformations
│       ├── pkgmap.x       # Package mapping
│       ├── typemap.x      # Type mapping
│       ├── funcmap.x      # Function mapping
│       ├── varmap.x       # Variable mapping
│       │
│       ├── const.x        # Const checking
│       ├── cache.x        # Build cache
│       ├── sourcemap.x    # Source mapping
│       │
│       ├── format.x       # Formatter
│       ├── watch.x        # File watcher
│       │
│       └── lsp/           # LSP package
│           ├── server.x
│           ├── protocol.x
│           ├── connection.x
│           ├── handlers.x
│           └── symbols.x
│
├── vet/                   # Vet package
│   ├── vet.x
│   ├── memory.x
│   ├── channels.x
│   ├── types.x
│   └── report.x
│
└── runtime/               # Runtime (already exists)
    ├── builtins.go
    ├── coerce.go
    └── ffi.go
```

## Implementation Strategy

### Phase 11.1: Core Transpiler (Week 1-2)

**Day 1-2: Project Setup**
- [x] Create `moxie-bootstrap/` directory structure
- [ ] Set up module system
- [ ] Create main.x entry point
- [ ] Implement command routing

**Day 3-5: Parser Integration**
- [ ] Wrap Go's `go/parser` and `go/ast`
- [ ] Implement AST inspection utilities
- [ ] Create token file set management
- [ ] Test parsing of basic Moxie files

**Day 6-10: Syntax Transformations**
- [ ] Port syntax.go to syntax.x (~1,500 lines)
- [ ] Implement preprocessing (preprocess.x)
- [ ] Type tracking system (typetrack.x)
- [ ] Channel literal transformation
- [ ] Endianness coercion
- [ ] clone() / free() / grow() transformations
- [ ] Test all Phase 2-6 transformations

### Phase 11.2: Name Transformations (Week 2)

**Day 11-12: Package & Type Mapping**
- [ ] Port pkgmap.go to pkgmap.x
- [ ] Port typemap.go to typemap.x
- [ ] Port naming.go to naming.x
- [ ] Test package/type mappings

**Day 13-14: Function & Variable Mapping**
- [ ] Port funcmap.go to funcmap.x
- [ ] Port varmap.go to varmap.x
- [ ] Test function/variable mappings

### Phase 11.3: Build System (Week 3)

**Day 15-17: Build Commands**
- [ ] Implement build command (build.x)
- [ ] Implement run command (run.x)
- [ ] Implement test command (test.x)
- [ ] Implement install command (install.x)
- [ ] Temporary directory management
- [ ] Go compiler invocation
- [ ] Module resolution

**Day 18-19: Caching & Source Mapping**
- [ ] Port cache.go to cache.x
- [ ] Port sourcemap.go to sourcemap.x
- [ ] Port clean.go to clean.x
- [ ] Implement const checker (const.x)

### Phase 11.4: Tooling (Week 3-4)

**Day 20-21: Formatter & Watch**
- [ ] Port format.go to format.x
- [ ] Port watch.go to watch.x
- [ ] fsnotify integration
- [ ] Test formatting and watching

**Day 22-24: Linter**
- [ ] Port vet package to Moxie
- [ ] Memory management checks
- [ ] Report generation (text, JSON, GitHub)
- [ ] Test linter on Moxie code

### Phase 11.5: LSP Server (Week 4)

**Day 25-28: LSP Implementation**
- [ ] Port lsp package to Moxie
- [ ] JSON-RPC 2.0 handler
- [ ] Document synchronization
- [ ] Symbol extraction and indexing
- [ ] Navigation features (definition, references, hover)
- [ ] Completion provider
- [ ] Diagnostics integration
- [ ] Test with VS Code extension

### Phase 11.6: Self-Hosting & Validation (Week 5)

**Day 29: Bootstrap Build**
- [ ] Use Go compiler to build Moxie bootstrap
- [ ] Verify moxie-bootstrap/cmd/moxie/main.x compiles
- [ ] Test basic commands (version, help)

**Day 30: Self-Compilation Test**
- [ ] Use moxie-bootstrap to compile itself
- [ ] Compare output with Go version
- [ ] Verify all commands work

**Day 31-32: Comprehensive Testing**
- [ ] Run all phase tests
- [ ] Test all commands (build, run, test, fmt, watch, vet, lsp)
- [ ] Verify LSP server works
- [ ] Test VS Code extension
- [ ] Performance comparison

**Day 33-35: Final Validation & Release**
- [ ] Fix any remaining bugs
- [ ] Update documentation
- [ ] Create v1.0.0 release
- [ ] Announce self-hosting milestone

## Key Moxie Features to Use

### Memory Management
```moxie
// Use explicit memory management
func parseFile(filename string) *ast.File {
    content := readFile(filename)
    defer free(content)

    fset := &token.FileSet{}
    defer free(fset)

    astFile := parser.ParseFile(fset, filename, content, 0)
    return clone(astFile) // Clone before freeing
}
```

### String Handling
```moxie
// Mutable strings (string = *[]byte)
func transformCode(source string) string {
    // Direct string mutation
    for i := 0; i < len(*source); i++ {
        if (*source)[i] == '.' {
            (*source)[i] = '_'
        }
    }
    return source
}
```

### Channels
```moxie
// Channel literals
func watchFiles(paths []string) {
    events := &chan FileEvent{100}  // Buffered channel
    errors := &chan error{}          // Unbuffered channel

    go func() {
        for event := range events {
            processEvent(event)
        }
    }()
}
```

### FFI (for Go stdlib access)
```moxie
import "moxie"

// Access Go's os/exec
var execCommand func(name string, args ...*[]byte) *exec.Cmd

func init() {
    lib := moxie.Dlopen("libgo.so", moxie.RTLD_LAZY)
    execCommand = moxie.Dlsym[func(string, ...*[]byte) *exec.Cmd](lib, "os/exec.Command")
}
```

### Type Coercion
```moxie
// Endianness handling
func readInt32LE(data *[]byte) int32 {
    ints := (*[]int32, moxie.LittleEndian)(data)
    return (*ints)[0]
}
```

## Challenges & Solutions

### Challenge 1: Circular Dependency
**Problem**: Moxie compiler written in Moxie needs Moxie compiler
**Solution**: Use Go implementation as bootstrap compiler

### Challenge 2: Go Standard Library Access
**Problem**: Need Go's go/parser, go/ast, os, io, etc.
**Solution**:
- Use FFI (Dlopen/Dlsym) for Go functions
- Or transpile without modification (Go code stays as-is in Moxie)
- Import Go packages directly (Moxie is syntactic sugar)

### Challenge 3: Error Handling
**Problem**: Moxie doesn't change Go's error handling
**Solution**: Use standard Go error patterns:
```moxie
result, err := someOperation()
if err != nil {
    return err
}
```

### Challenge 4: Generics
**Problem**: Runtime functions use Go generics
**Solution**: Keep using them! Moxie is Go-compatible:
```moxie
func clone[T any](x T) T {
    // Implementation
}
```

## Testing Strategy

### Unit Tests
- Test each module independently
- Compare output with Go version
- Validate transformations

### Integration Tests
- Full compile cycle tests
- Self-compilation tests
- LSP server tests

### Validation Tests
- Compile all existing Go code
- Verify identical output
- Performance benchmarks

## Success Criteria

### Must Have
- ✅ Moxie compiler compiles itself
- ✅ All commands functional (build, run, test, fmt, watch, vet, lsp, clean)
- ✅ Output identical to Go version
- ✅ All language features working
- ✅ LSP server functional with VS Code

### Nice to Have
- Performance within 2x of Go version
- Reduced memory footprint
- Cleaner code organization

## Timeline

**Total**: 5 weeks (35 days)

- **Week 1-2**: Core transpiler & name transformations (Days 1-14)
- **Week 3**: Build system & caching (Days 15-19)
- **Week 3-4**: Tooling & LSP (Days 20-28)
- **Week 5**: Self-hosting & validation (Days 29-35)

## Deliverables

### Code
1. `moxie-bootstrap/` - Complete Moxie implementation
2. All tests passing
3. Self-compilation working

### Documentation
1. Bootstrap guide
2. Moxie language examples
3. v1.0.0 release notes

### Release
1. `moxie` v1.0.0 binary (self-hosted)
2. VS Code extension v1.0.0
3. Announcement: Moxie is self-hosting!

## Post-Bootstrap

After successful bootstrap:
- Archive Go implementation (`cmd/moxie-go/`)
- Make Moxie version canonical
- Continue development in Moxie
- Build ecosystem (packages, tools)

## Notes

- **Keep it Simple**: Don't optimize prematurely
- **Trust the Design**: Moxie's syntax transformations are proven
- **Test Continuously**: Verify each module as you go
- **Use Moxie Features**: Demonstrate the language's power
- **Stay Compatible**: Maintain Go compatibility for ecosystem access

## References

- Current Go implementation: `cmd/moxie/`
- Language spec: `MOXIE-LANGUAGE-SPEC.md`
- Implementation status: `IMPLEMENTATION-STATUS.md`
- Phase plans: `PHASE*.md`
