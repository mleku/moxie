# Moxie

A Go transpiler that converts Moxie-flavored Go code to standard Go.

## Overview

Moxie is a transpiler-based approach to creating a Go variant. Instead of modifying the Go compiler itself, Moxie transpiles source code with Moxie conventions into standard Go, then uses the existing Go toolchain for building.

This approach provides:
- ✅ **Simplicity** - No compiler modifications needed
- ✅ **Compatibility** - Uses standard Go toolchain
- ✅ **Transparency** - All Go features and flags work
- ✅ **Extensibility** - Easy to add new transformations

## Quick Start

### Build the transpiler

```bash
go build -o moxie ./cmd/moxie
```

### Write Moxie code

Moxie files use the `.mx` extension:

```go
// main.mx
package main

import "github.com/mleku/moxie/src/fmt"

func main() {
    fmt.Println("Hello from Moxie!")
}
```

### Build and run

```bash
./moxie build
./moxie run main.mx
```

## Features

### Import Path Translation

Moxie automatically translates import paths:

```go
// Write this:
import "github.com/mleku/moxie/src/fmt"
import "github.com/mleku/moxie/internal/net/http"

// Transpiles to:
import "fmt"
import "net/http"
```

### All Go Commands Supported

- `moxie build` - Build packages
- `moxie install` - Install packages
- `moxie run` - Run programs
- `moxie test` - Run tests

All standard Go flags are supported:

```bash
moxie build -o myapp -ldflags="-s -w" ./cmd/myapp
moxie test -v -race ./...
```

## Examples

See the `examples/` directory:

- **hello** - Simple "Hello World" example
- **webserver** - HTTP server using `net/http`

## Documentation

- [TRANSPILER.md](TRANSPILER.md) - Detailed transpiler documentation
- [examples/](examples/) - Example programs

## How It Works

1. **Parse** - Source files are parsed into Go ASTs
2. **Transform** - Import paths are rewritten
3. **Generate** - Transformed code is written to a temp directory
4. **Build** - Standard `go` command is invoked

## Environment Variables

Moxie respects all standard Go environment variables:
- `GOPATH`, `GOBIN`, `GOROOT`
- `GOOS`, `GOARCH`
- `CGO_ENABLED`
- All others

## Future Enhancements

The transpiler architecture allows for:
- Custom syntax transformations
- Code instrumentation
- Performance profiling
- Type system extensions
- And more...

All without modifying the Go compiler!

## License

BSD-style license (see LICENSE file)

## Status

This is the transpiler-based implementation of Moxie, designed as a simpler alternative to modifying the Go compiler directly.
