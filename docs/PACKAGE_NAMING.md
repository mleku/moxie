# Moxie Package Naming Conventions

## Overview

Moxie uses a transparent package naming system that maintains compatibility with Go's stdlib while allowing for future divergence.

## Current Strategy (Phase 1.1)

### Standard Library Packages

Moxie stdlib packages use the **same names** as Go stdlib packages:
- `fmt` → `fmt`
- `os` → `os`
- `net/http` → `http` (when imported via Moxie paths)
- `encoding/json` → `json` (when imported via Moxie paths)

### Custom/User Packages

User-defined packages follow Go conventions:
- Package names should be lowercase
- Single-word package names preferred
- No underscores or mixedCaps

## Import Path Translation

Moxie uses special import paths that map to standard Go packages:

```go
// Moxie code
import "github.com/mleku/moxie/src/fmt"      // → fmt
import "github.com/mleku/moxie/internal/os"  // → os
import "github.com/mleku/moxie/src/net/http" // → net/http
```

The transpiler automatically converts these paths to standard Go imports.

## Package Name Mapping

The transpiler maintains a bidirectional mapping between Moxie and Go package names:

```go
type PackageMapping struct {
    moxieToGo map[string]string  // Moxie → Go
    goToMoxie map[string]string  // Go → Moxie
}
```

### Mapped Packages

The following stdlib packages have explicit mappings (currently 1:1):

**Core Packages:**
- `fmt`, `os`, `io`, `errors`, `context`

**Networking:**
- `net`, `http`, `url`

**Time & Concurrency:**
- `time`, `sync`

**String/Byte Manipulation:**
- `strings`, `bytes`, `unicode`, `utf8`

**Encoding:**
- `json`, `xml`, `base64`, `hex`

**Data Structures:**
- `sort`, `container/*`

**Math:**
- `math`, `rand`, `big`

**I/O:**
- `bufio`, `io/fs`, `io/ioutil`

**Testing:**
- `testing`, `testing/quick`

**System:**
- `runtime`, `debug`, `pprof`, `trace`, `unsafe`, `reflect`

**Archives & Compression:**
- `tar`, `zip`, `gzip`, `zlib`

**Cryptography:**
- `crypto`, `md5`, `sha1`, `sha256`, `sha512`, `aes`, `des`, `rsa`, `ecdsa`, `tls`, `x509`

**File Paths:**
- `path`, `filepath`

**Logging:**
- `log`, `slog`

**Text Processing:**
- `regexp`, `scanner`, `template`

**Database:**
- `sql`, `sql/driver`

## Conflict Resolution

### Naming Conflicts

If a conflict arises between package names:

1. **Stdlib takes precedence** - Standard library packages have priority
2. **Use full path** - User packages should use full import paths
3. **Alias imports** - Use import aliases when necessary:
   ```go
   import myhttp "example.com/mypackage/http"
   ```

### Reserved Names

The following package names are reserved for stdlib:
- All current Go stdlib package names
- Future Moxie stdlib additions

### Custom Package Guidelines

For user-defined packages:
- Avoid names that conflict with stdlib
- Use descriptive, single-word names
- Use domain-based import paths for uniqueness
- Follow Go package naming conventions

## Future Considerations

### Potential Divergence

Future Moxie versions may introduce:
- Different package names for clarity
- Additional stdlib packages
- Renamed packages for consistency

The mapping system is designed to handle these changes transparently.

### Versioning

Package mappings may be version-specific:
- Go 1.24 stdlib → Moxie 0.1 stdlib
- Different Go versions may have different mappings
- Mapping tables can be versioned

## Implementation Details

### Package Declaration Transformation

When transpiling `.mx` files:

```go
// Input: main.mx
package main

// Output: main.go (unchanged)
package main
```

Currently, package declarations are passed through unchanged since Moxie uses the same package names as Go.

### Import Statement Transformation

Import paths are transformed but package names remain the same:

```go
// Input: main.mx
import "github.com/mleku/moxie/src/fmt"

// Output: main.go
import "fmt"
```

### Package Scoping

Package names are scoped and don't affect:
- Type names
- Function names
- Variable names
- Constants

These will have their own transformation rules in later phases.

## Examples

### Example 1: Simple Program

```go
// hello.mx
package main

import "github.com/mleku/moxie/src/fmt"

func main() {
    fmt.Println("Hello, Moxie!")
}
```

Transpiles to:

```go
// hello.go (generated)
package main

import "fmt"

func main() {
    fmt.Println("Hello, Moxie!")
}
```

### Example 2: Multiple Imports

```go
// server.mx
package main

import (
    "github.com/mleku/moxie/src/fmt"
    "github.com/mleku/moxie/internal/net/http"
    "github.com/mleku/moxie/src/encoding/json"
)

func main() {
    // Package names used as 'fmt', 'http', 'json'
    fmt.Println("Starting server...")
    http.HandleFunc("/", handler)
}
```

Transpiles to:

```go
// server.go (generated)
package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

func main() {
    // Package names remain the same
    fmt.Println("Starting server...")
    http.HandleFunc("/", handler)
}
```

## Testing

Package name transformation is tested via:
- Unit tests in `pkgmap_test.go`
- Integration tests with example programs
- Round-trip testing (Moxie → Go → compile → run)

## Summary

Phase 1.1 establishes:
✅ Bidirectional package name mapping
✅ 1:1 mapping for stdlib (Moxie = Go names)
✅ Infrastructure for future divergence
✅ Conflict resolution strategy
✅ Clear documentation and examples

Future phases will build on this foundation to add:
- Type name transformations (Phase 1.2)
- Function name transformations (Phase 1.3)
- Variable name transformations (Phase 1.4)
