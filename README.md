# Moxie Programming Language

**A consistent, explicit, and pragmatic evolution of Go**

Moxie is a programming language that addresses Go's design inconsistencies while preserving its philosophy of simplicity. The Moxie transpiler converts `.mx` files to standard Go code, allowing full interoperability with the Go ecosystem.

## Status

🟡 **Phase 2 - 75% Complete** (Active Development)

- ✅ Phase 0: Foundation (100%)
- ✅ Phase 1: Name Transformation Infrastructure (100%)
- 🟡 Phase 2: Syntax Transformations (75%)
- ⏳ Phase 3+: Additional Language Features (Planned)

See [IMPLEMENTATION-STATUS.md](IMPLEMENTATION-STATUS.md) for detailed progress.

## Quick Start

```bash
# Build the transpiler
go build -o moxie ./cmd/moxie

# Run a Moxie program
./moxie run examples/phase2/main.mx

# Build a Moxie program
./moxie build examples/phase2

# See all commands
./moxie help
```

## Key Features

### ✅ Explicit Pointer Types (Working)

In Moxie, slices, maps, and channels are explicitly pointer types:

```go
// Moxie - Explicit pointers
s := &[]int{1, 2, 3}
m := &map[string]int{"one": 1}
ch := &chan int{}

// Automatically transpiled to valid Go
s := &[]int{1, 2, 3}
m := &map[string]int{"one": 1}
ch := make(chan int)  // Special handling
```

### ✅ Improved Built-in Functions (Working)

```go
// Moxie - Works with pointer types
s = append(s, 4, 5)    // Transpiles to: *s = append(*s, 4, 5)
clear(m)               // Transpiles to: clear(*m)
```

### ✅ No make() Function (Working)

```go
// Moxie - Use composite literals
s := &[]int{}           // Instead of make([]int, 0)
m := &map[string]int{}  // Instead of make(map[string]int)

// Using make() produces a helpful error message
```

### 🟡 New Built-in Functions (AST transformation ready, module setup pending)

```go
// Moxie - Additional built-ins
s = grow(s, 100)       // Pre-allocate capacity
s2 := clone(s)         // Deep copy
free(s2)               // GC hint
```

### ⏳ Future Features (Planned)

- **Mutable strings**: `string` = `*[]byte`
- **True const**: MMU-protected immutability
- **Native FFI**: Built-in `dlopen`/`dlsym`
- **Zero-copy coercion**: Reinterpret slices with endianness control

## Language Principles

1. **Explicitness** - No hidden behavior
2. **Consistency** - Uniform type system
3. **Pragmatism** - Solve real problems
4. **Simplicity** - Easy to learn and use
5. **Compatibility** - Full Go interoperability

## Naming Conventions

Moxie uses **PascalCase** and **camelCase**, matching Go conventions:

```go
// Exported (public)
type User struct { }
func ProcessData() { }
const MaxConnections = 100

// Unexported (private)
type internalState struct { }
func parseRequest() { }
var errorCount = 0
```

**No snake_case** - This maintains compatibility with the Go ecosystem.

## Documentation

### Language
- [MOXIE-LANGUAGE-SPEC.md](MOXIE-LANGUAGE-SPEC.md) - Complete language specification
- [go-language-revision.md](go-language-revision.md) - Design rationale and proposals

### Implementation
- [IMPLEMENTATION-STATUS.md](IMPLEMENTATION-STATUS.md) - Detailed progress tracking
- [STATUS-SUMMARY.md](STATUS-SUMMARY.md) - Quick status overview
- [PHASE2-PROGRESS.md](PHASE2-PROGRESS.md) - Current phase details

### Completed Phases
- [PHASE1.2-COMPLETE.md](PHASE1.2-COMPLETE.md) - Type name transformations
- [PHASE1.3-COMPLETE.md](PHASE1.3-COMPLETE.md) - Function name transformations
- [PHASE1.4-COMPLETE.md](PHASE1.4-COMPLETE.md) - Variable name transformations

## Examples

### Hello World

```go
// examples/hello/main.mx
package main

import "fmt"

func main() {
    fmt.Println("Hello from Moxie!")
}
```

### Working with Slices and Maps

```go
// examples/phase2/main.mx
package main

import "fmt"

func main() {
    // Explicit pointer types
    s := &[]int{1, 2, 3}
    m := &map[string]int{"one": 1, "two": 2}

    // Automatic transformation
    s = append(s, 4, 5, 6)
    clear(m)

    fmt.Println("Slice:", s)
    fmt.Println("Map:", m)
}
```

## Project Structure

```
moxie/
├── cmd/moxie/           # Transpiler source
│   ├── main.go          # Main transpiler logic
│   ├── syntax.go        # Syntax transformations (Phase 2)
│   ├── naming.go        # Name conversion utilities
│   ├── pkgmap.go        # Package mapping
│   ├── typemap.go       # Type transformations
│   ├── funcmap.go       # Function transformations
│   ├── varmap.go        # Variable transformations
│   └── *_test.go        # Test files (330+ tests)
├── runtime/             # Moxie runtime support
│   └── builtins.go      # Runtime functions (grow, clone, free)
├── examples/            # Example programs
│   ├── hello/           # Hello world
│   ├── webserver/       # HTTP server
│   ├── json-api/        # JSON API
│   └── phase2/          # Phase 2 demos
├── docs/                # Additional documentation
└── *.md                 # Documentation files
```

## Requirements

- Go 1.21+ (for `clear()` built-in)
- Go 1.18+ (for generics in runtime)

## Building

```bash
# Build the transpiler
go build -o moxie ./cmd/moxie

# Run tests
go test ./cmd/moxie/...

# Build all examples
./moxie build examples/hello
./moxie build examples/webserver
./moxie build examples/json-api
./moxie build examples/phase2
```

## Usage

```bash
# Run a Moxie program directly
./moxie run main.mx

# Build a Moxie program
./moxie build [directory]
./moxie build -o output main.mx

# Install a Moxie program
./moxie install [directory]

# Run tests
./moxie test [directory]

# Show version
./moxie version
```

## Differences from Go

### What's Different

1. **Explicit pointers**: `*[]T`, `*map[K]V`, `*chan T` instead of implicit references
2. **No make()**: Use `&[]T{}`, `&map[K]V{}`, `&chan T{}` instead
3. **New built-ins**: `grow()`, `clone()`, `free()`
4. **Future**: Mutable strings, true const, native FFI, zero-copy coercion

### What's the Same

- All Go syntax and features work
- Same type system (with explicit pointers)
- Same standard library
- Same toolchain (transpiles to Go)
- Same performance characteristics
- Full interoperability

## Development Roadmap

### ✅ Completed
- Phase 0: Foundation and transpiler infrastructure
- Phase 1: Name transformation infrastructure (disabled by default)
- Phase 2: Core syntax transformations (75%)

### 🟡 In Progress
- Phase 2: Complete runtime module resolution
- Phase 2: Channel literal transformation
- Phase 2: Comprehensive test suite

### ⏳ Planned
- Phase 3: String mutability (`string = *[]byte`)
- Phase 4: True const with MMU protection
- Phase 5: Native FFI (dlopen, dlsym)
- Phase 6: Zero-copy type coercion
- Phase 7+: Tooling (LSP, formatter, linter)
- Future: Self-hosting (bootstrap)

## Performance

- **Transpilation**: ~2-5ms per file
- **Runtime overhead**: Minimal (mostly direct Go code)
- **Memory overhead**: Small (AST transformation only)
- **Execution**: Same as Go (transpiles to standard Go)

## License

BSD-style license. See LICENSE file for details.

## Credits

The Moxie Authors

## Links

- Language Specification: [MOXIE-LANGUAGE-SPEC.md](MOXIE-LANGUAGE-SPEC.md)
- Design Rationale: [go-language-revision.md](go-language-revision.md)
- Implementation Status: [IMPLEMENTATION-STATUS.md](IMPLEMENTATION-STATUS.md)
- Current Phase: [PHASE2-PROGRESS.md](PHASE2-PROGRESS.md)

---

**Moxie** - Making Go more consistent, one transformation at a time.
