# The Moxie Programming Language

**Systems programming with spirit.**

![Phase 0: 85%](https://img.shields.io/badge/Phase%200-85%25-yellow)
![Overall: 0.85%](https://img.shields.io/badge/Overall-0.85%25-red)
![Status: In Progress](https://img.shields.io/badge/Status-In%20Progress-blue)

Moxie is a fork of the Go programming language with fundamental revisions based on 8 years of real-world experience. It addresses core design inconsistencies while preserving Go's philosophy of simplicity and explicitness.

## What is Moxie?

Moxie eliminates the complexity and race-prone design of Go's reference types, adds native FFI support to replace CGO, introduces mutable strings, implements hardware-enforced immutability, and provides zero-copy type coercion with endianness handling for high-performance systems programming.

### Key Improvements Over Go

1. **Explicit Reference Types** - `[]T` → `*[]T`, `map[K]V` → `*map[K]V`, `chan T` → `*chan T`
2. **Mutable Strings** - Merge `string` and `[]byte` into unified `*[]byte` type
3. **True Immutability** - `const` with MMU memory protection (hardware-enforced)
4. **Native FFI** - Replace CGO with type-safe `dlopen()` built-ins (20x faster)
5. **Zero-Copy Type Coercion** - Cast between numeric slice types with automatic endianness handling
6. **Explicit Integer Sizes** - Eliminate platform-dependent `int`/`uint` types

For complete details, see [go-language-revision.md](go-language-revision.md).

## Repository

This repository is a fork of the Go programming language.

Original Go repository: https://go.googlesource.com/go
Original Go mirror: https://github.com/golang/go

Unless otherwise noted, the Moxie source files are distributed under the BSD-style license found in the LICENSE file.

## Installation

### Install From Source

Moxie is currently in development (Phase 0 - Foundation & Setup).

To build from source:

```bash
cd src
./all.bash
```

For detailed installation instructions, visit the documentation (coming soon).

### Current Status

**Phase 0: Foundation & Setup** (In Progress)
- ✅ Repository structure established
- 🔄 Branding updates in progress
- ⏳ Build system modifications pending
- ⏳ Testing infrastructure pending

See [go-to-moxie-plan.md](go-to-moxie-plan.md) for the complete implementation roadmap.

## Contributing

Moxie is being actively developed. We welcome contributions!

To contribute:
1. Read the [implementation plan](go-to-moxie-plan.md)
2. Check the current phase and open tasks
3. Submit issues and pull requests

### Design Goals

- **Simplicity** - 30% complexity reduction over Go
- **Safety** - Explicit semantics, hardware-enforced immutability
- **Performance** - Zero-copy operations, faster FFI, optimized string handling
- **Portability** - Platform-independent types, explicit endianness

## Documentation

- [Complete Language Specification](go-language-revision.md) - Full Moxie design
- [Summary](go-language-revision-summary.md) - Quick overview of changes
- [Implementation Plan](go-to-moxie-plan.md) - Development roadmap

## Acknowledgments

Moxie is built on the foundation of the Go programming language, which is the work of thousands of contributors. We are grateful for their efforts.

Original Go authors and the Go team at Google.

## License

BSD-style license - see [LICENSE](LICENSE) file.

---

*Moxie: The Go we should have had from the beginning.*
