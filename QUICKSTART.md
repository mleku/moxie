# Moxie Quick Start Guide

Get started with Moxie in 5 minutes!

## 1. Build the Transpiler

```bash
cd /home/mleku/src/github.com/mleku/moxie
go build -o moxie ./cmd/moxie
```

## 2. Create Your First Moxie Program

Create a file called `hello.mx`:

```go
package main

import "github.com/mleku/moxie/src/fmt"

func main() {
    fmt.Println("Hello from Moxie!")
}
```

## 3. Run It!

```bash
./moxie run hello.mx
```

Output: `Hello from Moxie!`

## 4. Build a Binary

```bash
./moxie build hello.mx
./hello
```

## 5. Try the Examples

### Simple Hello World

```bash
./moxie run examples/hello/main.mx
```

### Web Server

```bash
./moxie build examples/webserver
./webserver
# In another terminal: curl http://localhost:8080/test
```

### JSON API

```bash
./moxie build examples/json-api -o api
./api
# In another terminal: curl http://localhost:8080/api/status
```

## Key Concepts

### File Extension

Moxie files use `.mx` extension:
- Write: `main.mx`
- Transpiles to: `main.go` (automatic)

### Import Paths

Use Moxie import paths in your code:

```go
import "github.com/mleku/moxie/src/fmt"           // → fmt
import "github.com/mleku/moxie/internal/os"       // → os
import "github.com/mleku/moxie/src/net/http"      // → net/http
```

The transpiler automatically converts these to standard Go imports.

### Commands

All Go commands work:

```bash
moxie build      # Build packages
moxie run        # Run programs
moxie install    # Install to $GOBIN
moxie test       # Run tests
```

All Go flags work too:

```bash
moxie build -o myapp -ldflags="-s -w"
moxie test -v -race ./...
```

## What Happens Under the Hood

1. Your `.mx` files are parsed
2. Import paths are rewritten to standard Go
3. Standard `.go` files are generated in a temp directory
4. The regular `go` command builds them
5. The binary is copied back to you
6. Temp files are cleaned up

You get all the power of Go, with Moxie conventions!

## Next Steps

- Read [README.md](README.md) for overview
- See [TRANSPILER.md](TRANSPILER.md) for details
- Check [IMPLEMENTATION.md](IMPLEMENTATION.md) for architecture
- Explore `examples/` for more code

## Need Help?

Check the documentation or examine the examples directory!
