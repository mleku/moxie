# Moxie Transpiler

The Moxie transpiler converts Moxie source code to standard Go, then uses the existing Go toolchain to build it. This approach avoids the complexity of modifying the Go compiler itself.

## How It Works

1. **Parse** - Moxie source files (`.mx` files) are parsed into Go ASTs
2. **Transform** - Import paths are rewritten from Moxie paths to standard Go paths
3. **Generate** - Transformed ASTs are written as `.go` files to a temporary directory
4. **Build** - The standard Go toolchain (`go build`, `go install`, etc.) is invoked on the transpiled code

## Installation

Build the transpiler:

```bash
go build -o moxie ./cmd/moxie
```

## Usage

The `moxie` command works like `go`:

### Build

```bash
moxie build [packages]
```

Transpiles and builds the specified packages (or current directory if not specified).

### Install

```bash
moxie install [packages]
```

Transpiles and installs the specified packages to `$GOBIN`.

### Run

```bash
moxie run [file.go]
```

Transpiles and runs a Moxie program.

### Test

```bash
moxie test [packages]
```

Transpiles and runs tests.

## Import Path Translation

The transpiler automatically translates Moxie import paths to standard Go paths:

| Moxie Import Path | Transpiled To |
|-------------------|---------------|
| `github.com/mleku/moxie/src/fmt` | `fmt` |
| `github.com/mleku/moxie/internal/os` | `os` |
| `github.com/mleku/moxie/src/net/http` | `net/http` |

The pattern is:
- `github.com/mleku/moxie/src/X` → `X`
- `github.com/mleku/moxie/internal/X` → `X`

Other import paths are left unchanged.

## Example

Create a Moxie program:

```go
// hello.mx
package main

import "github.com/mleku/moxie/src/fmt"

func main() {
	fmt.Println("Hello from Moxie!")
}
```

Build and run:

```bash
moxie build
./hello
```

Or run directly:

```bash
moxie run hello.mx
```

## Environment Variables

The transpiler respects all standard Go environment variables:

- `GOPATH` - Go workspace location
- `GOBIN` - Where `moxie install` puts binaries
- `GOOS`, `GOARCH` - Target platform
- `CGO_ENABLED` - C integration
- etc.

All Go build flags are also supported:

```bash
moxie build -o myapp -ldflags="-s -w" ./cmd/myapp
```

## Development Workflow

1. Write code using Moxie import paths (`github.com/mleku/moxie/src/...`)
2. Use `moxie build` instead of `go build`
3. Use `moxie test` instead of `go test`
4. Use `moxie install` to install binaries

The transpiler is transparent - you get all the benefits of the standard Go toolchain.

## Future Enhancements

The transpiler can be extended to support:

- Custom type name transformations
- Code instrumentation
- Moxie-specific syntax extensions
- Performance profiling
- etc.

All of this without modifying the Go compiler!
