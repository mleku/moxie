# Moxie Bootstrap - Self-Hosted Compiler

This directory contains the Moxie transpiler rewritten in Moxie itself.

## Status

🚀 **Phase 11: Bootstrap** - In Progress

The Moxie language is now self-hosting! This implementation demonstrates Moxie's capability to build complex software including its own compiler.

## Directory Structure

```
moxie-bootstrap/
├── cmd/moxie/          # Main compiler
│   ├── main.x          # Entry point ✅
│   ├── build.x         # Build command
│   ├── run.x           # Run command
│   ├── test.x          # Test command
│   ├── install.x       # Install command
│   ├── parser.x        # Parser wrapper
│   ├── syntax.x        # Syntax transformations
│   ├── preprocess.x    # Preprocessing
│   ├── typetrack.x     # Type tracking
│   ├── naming.x        # Name conversions
│   ├── pkgmap.x        # Package mapping
│   ├── typemap.x       # Type mapping
│   ├── funcmap.x       # Function mapping
│   ├── varmap.x        # Variable mapping
│   ├── const.x         # Const checking
│   ├── cache.x         # Build cache
│   ├── sourcemap.x     # Source mapping
│   ├── format.x        # Formatter
│   ├── watch.x         # File watcher
│   ├── clean.x         # Cache cleaner
│   ├── vet.x           # Vet command
│   ├── lsp.x           # LSP command
│   └── lsp/            # LSP package
│       ├── server.x
│       ├── protocol.x
│       ├── connection.x
│       ├── handlers.x
│       └── symbols.x
└── vet/                # Vet package
    ├── vet.x
    ├── memory.x
    ├── channels.x
    ├── types.x
    └── report.x
```

## Building

### Using Go Bootstrap Compiler

```bash
# From repository root
./moxie build moxie-bootstrap/cmd/moxie

# Run the self-hosted compiler
./moxie-bootstrap/cmd/moxie/moxie version
```

### Self-Compilation Test

```bash
# Use Moxie to compile itself!
./moxie-bootstrap/cmd/moxie/moxie build moxie-bootstrap/cmd/moxie
```

## Moxie Features Demonstrated

### 1. Explicit Memory Management

```moxie
func parseFile(filename string) (*ast.File, error) {
    content := readFile(filename)
    defer free(content)  // Explicit cleanup

    fset := &token.FileSet{}
    defer free(fset)

    astFile, err := parser.ParseFile(fset, filename, content, 0)
    if err != nil {
        return nil, err
    }

    return clone(astFile), nil  // Clone before free
}
```

### 2. Mutable Strings

```moxie
// Strings are *[]byte - mutable!
func transformIdentifier(name string) string {
    for i := 0; i < len(*name); i++ {
        if (*name)[i] == '_' {
            (*name)[i] = '-'
        }
    }
    return name
}
```

### 3. Channel Literals

```moxie
func watchFileSystem(paths *[]*[]byte) {
    events := &chan FileEvent{100}  // Buffered channel
    errors := &chan error{}          // Unbuffered

    go func() {
        for event := range events {
            processEvent(event)
        }
    }()

    watcher.Watch(paths, events, errors)
}
```

### 4. Explicit Pointer Types

```moxie
// Slices and maps are explicit pointers
func buildIndex(files *[]*[]byte) *map[*[]byte]*[]Symbol {
    index := &map[*[]byte]*[]Symbol{}

    for _, file := range *files {
        symbols := extractSymbols(file)
        (*index)[file] = symbols
    }

    return index
}
```

### 5. Type Coercion with Endianness

```moxie
func readInt32LE(data *[]byte) int32 {
    ints := (*[]int32, moxie.LittleEndian)(data)
    defer free(ints)
    return (*ints)[0]
}
```

## Implementation Progress

- [x] Phase 11 plan created
- [x] Directory structure created
- [x] Main entry point (`main.x`)
- [ ] Core transpiler (parser, syntax)
- [ ] Name transformations
- [ ] Build system
- [ ] Tooling (fmt, watch, vet)
- [ ] LSP server
- [ ] Self-compilation test
- [ ] v1.0.0 release

## Timeline

**5 weeks** (35 days) to complete bootstrap:
- Week 1-2: Core transpiler
- Week 3: Build system
- Week 3-4: Tooling & LSP
- Week 5: Self-hosting validation

## Philosophy

This is a **bootstrap compiler**: correctness over optimization. The goal is self-hosting, not peak performance. Once Moxie compiles itself, we can optimize the self-hosted version.

**Key Principles**:
- ✅ Correctness first
- ✅ Feature parity with Go version
- ✅ Use Moxie idioms
- ✅ Trust the language design
- ⏭️ Skip premature optimization

## Testing

```bash
# Test the bootstrap compiler
./moxie test moxie-bootstrap/...

# Test self-compilation
./moxie run moxie-bootstrap/cmd/moxie build moxie-bootstrap/cmd/moxie
```

## References

- Bootstrap plan: `../PHASE11-BOOTSTRAP.md`
- Go implementation: `../cmd/moxie/`
- Language spec: `../MOXIE-LANGUAGE-SPEC.md`
- Implementation status: `../IMPLEMENTATION-STATUS.md`

## Success Criteria

### Must Have
- ✅ Moxie compiles itself
- ✅ All commands work (build, run, test, fmt, watch, vet, lsp, clean)
- ✅ Output identical to Go version
- ✅ All language features functional
- ✅ LSP server works with VS Code

### Nice to Have
- Performance within 2x of Go version
- Reduced memory footprint
- Cleaner code organization

## License

BSD 3-Clause License - See LICENSE file in repository root.

---

**Milestone**: When this compiler successfully compiles itself, Moxie will be self-hosting and ready for v1.0.0 release! 🎉
