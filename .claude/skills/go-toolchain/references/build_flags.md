# Go Build Flags Reference

Comprehensive guide to Go compiler and linker flags.

## Compiler Flags (-gcflags)

Compiler flags are passed to the Go compiler (gc) via `-gcflags`.

**Syntax:**
```bash
go build -gcflags='[pattern=]flags'
```

### Common Compiler Flags

#### Optimization Flags

- `-N` - Disable optimizations
- `-l` - Disable inlining
- `-m` - Print optimization decisions (use `-m -m` for more detail)
- `-m=<level>` - Set optimization level for inlining (0-4)

**Examples:**
```bash
# Disable optimizations for debugging
go build -gcflags='-N -l'

# Show inlining decisions
go build -gcflags='-m'

# Show detailed optimization decisions
go build -gcflags='-m -m'
```

#### Debug Flags

- `-S` - Print assembly listing
- `-live` - Print liveness analysis
- `-wb` - Print write barrier analysis
- `-d=<option>` - Enable debug option

**Examples:**
```bash
# Print assembly
go build -gcflags='-S'

# Debug specific package
go build -gcflags='example.com/pkg=-N -l'
```

#### Other Compiler Flags

- `-p <path>` - Set expected package import path
- `-D <path>` - Set relative path for local imports
- `-I <dir>` - Add directory to import search path
- `-pack` - Write package file instead of object file
- `-e` - Remove limit on error reporting
- `-h` - Panic on error
- `-complete` - Assume package has no non-Go components
- `-nolocalimports` - Reject local (relative) imports
- `-importcfg <file>` - Read import configuration from file
- `-asmhdr <file>` - Write assembly header to file
- `-buildid <id>` - Record build id
- `-c <n>` - Concurrency level (default: 1)
- `-dwarf` - Generate DWARF symbols
- `-dwarflocationlists` - Add location lists to DWARF
- `-dynlink` - Support references to Go symbols defined in other shared libraries
- `-goversion <version>` - Required Go version
- `-importmap <old>=<new>` - Add import map entry
- `-lang <version>` - Set language version
- `-race` - Enable race detector
- `-shared` - Generate code that can be linked into a shared library
- `-spectre <list>` - Enable Spectre mitigations
- `-traceprofile <file>` - Write execution trace to file

### Per-Package Compiler Flags

Apply flags to specific packages using pattern matching:

```bash
# Disable optimizations for specific package
go build -gcflags='example.com/pkg=-N -l'

# Disable optimizations for all packages
go build -gcflags='all=-N -l'

# Apply to stdlib and dependencies
go build -gcflags='std=-m'
```

Patterns:
- `all` - All packages
- `std` - Standard library packages
- `cmd` - Go command packages
- `path/to/pkg` - Specific package

## Linker Flags (-ldflags)

Linker flags are passed to the Go linker via `-ldflags`.

**Syntax:**
```bash
go build -ldflags='[pattern=]flags'
```

### Common Linker Flags

#### Stripping and Compression

- `-s` - Omit symbol table and debug information
- `-w` - Omit DWARF symbol table
- `-compressdwarf` - Compress DWARF debugging information (default: true)

**Examples:**
```bash
# Strip debug information (smaller binaries)
go build -ldflags='-s -w'

# Keep symbols but strip DWARF
go build -ldflags='-w'
```

**Size Comparison:**
```bash
# Normal build: ~10 MB
go build

# Stripped build: ~7 MB
go build -ldflags='-s -w'

# With UPX compression: ~3 MB
go build -ldflags='-s -w' && upx binary
```

#### Variable Injection

Set package variables at build time using `-X`:

```bash
go build -ldflags='-X package.variable=value'
```

**Examples:**
```bash
# Set version information
go build -ldflags="-X main.Version=1.2.3 -X main.BuildTime=$(date -u +%Y-%m-%d_%H:%M:%S)"

# Set multiple variables
go build -ldflags="-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)"
```

**Code Example:**
```go
package main

var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)

func main() {
    fmt.Printf("Version: %s\n", Version)
    fmt.Printf("Built: %s\n", BuildTime)
    fmt.Printf("Commit: %s\n", GitCommit)
}
```

#### External Linking

- `-linkmode <mode>` - Set link mode (internal, external, auto)
- `-extld <linker>` - Set external linker (default: gcc or clang)
- `-extldflags <flags>` - Pass flags to external linker

**Examples:**
```bash
# Force external linking
go build -ldflags='-linkmode external'

# Use specific external linker
go build -ldflags='-extld=clang'

# Pass flags to external linker
go build -ldflags='-extldflags "-static"'

# Build fully static binary
CGO_ENABLED=1 go build -ldflags='-linkmode external -extldflags "-static"'
```

#### Build ID and Reproducibility

- `-buildid <id>` - Record build id
- `-B <gobuildid>` - Add gobuildid to symbol table

**Examples:**
```bash
# Set custom build ID
go build -ldflags='-buildid=custom-id-123'

# Reproducible builds (combine with -trimpath)
go build -trimpath -ldflags='-buildid='
```

#### Other Linker Flags

- `-E <entry>` - Set entry point symbol
- `-H <type>` - Set executable format (darwin, linux, windows, etc.)
- `-I <linker>` - Use linker (not implemented)
- `-L <dir>` - Add directory to library search path
- `-T <address>` - Set text segment address
- `-V` - Print linker version
- `-d` - Disable dynamic executable
- `-f` - Ignore version mismatch
- `-g` - Disable Go package data checks
- `-k <symbol>` - Set field tracking symbol
- `-libgcc <file>` - Set compiler support library
- `-linkshared` - Link against shared libraries
- `-r <path>` - Set ELF dynamic linker search path
- `-race` - Enable race detector
- `-tmpdir <dir>` - Use dir for temporary files
- `-v` - Print link trace

### Per-Package Linker Flags

Similar to gcflags, can apply to specific packages:

```bash
# Apply to specific package
go build -ldflags='example.com/pkg=-s -w'

# Apply to all packages
go build -ldflags='all=-s -w'
```

## Assembly Flags (-asmflags)

Assembly flags are passed to the assembler.

**Syntax:**
```bash
go build -asmflags='[pattern=]flags'
```

### Common Assembly Flags

- `-D <name>=<value>` - Define symbol
- `-I <dir>` - Add include directory
- `-S` - Print assembly listing
- `-V` - Print assembler version
- `-dynlink` - Support references to Go symbols in other shared libraries
- `-gensymabis` - Write symbol ABI information to stdout
- `-o <file>` - Write output to file
- `-p <path>` - Set expected package import path
- `-shared` - Generate code for shared object
- `-spectre <list>` - Enable Spectre mitigations
- `-trimpath <prefix>` - Remove prefix from recorded paths

**Examples:**
```bash
# Print assembly listing
go build -asmflags='-S'

# Define assembly symbol
go build -asmflags='-D MAXSIZE=1024'
```

## Gccgo Flags (-gccgoflags)

Flags for the gccgo compiler (alternative Go compiler).

**Syntax:**
```bash
go build -compiler=gccgo -gccgoflags='flags'
```

**Examples:**
```bash
# Build with gccgo
go build -compiler=gccgo

# Build with gccgo optimizations
go build -compiler=gccgo -gccgoflags='-O3'
```

## Combined Flag Examples

### Development Build
```bash
# Fast compilation, full debug info
go build -gcflags='all=-N -l'
```

### Production Build
```bash
# Optimized, stripped, reproducible
go build -trimpath -ldflags='-s -w -X main.version=1.0.0'
```

### Debug Build
```bash
# No optimizations, race detection
go build -race -gcflags='all=-N -l'
```

### Static Binary
```bash
# Fully static, no external dependencies
CGO_ENABLED=0 go build -ldflags='-s -w -extldflags "-static"'
```

### Profiling Build
```bash
# With CPU profiling support
go build -gcflags='-m=2' -ldflags='-s -w'
```

### Size-Optimized Build
```bash
# Smallest possible binary
go build -ldflags='-s -w' -trimpath
# Optional: compress with UPX
upx --best --lzma binary
```

## Flag Combination Patterns

### Pattern Syntax

Flags can be applied selectively using patterns:

```bash
# Apply to all packages
-gcflags='all=flags'

# Apply to standard library
-gcflags='std=flags'

# Apply to specific package
-gcflags='github.com/user/pkg=flags'

# Apply to multiple patterns
-gcflags='std=-m github.com/user/pkg=-N -l'
```

### Common Combinations

**Development:**
```bash
go build -gcflags='all=-N -l' -race
```

**Production:**
```bash
go build -trimpath -ldflags='-s -w -X main.version=$(VERSION)'
```

**Debugging:**
```bash
go build -gcflags='all=-N -l' -ldflags='' -race -cover
```

**Performance Analysis:**
```bash
go build -gcflags='-m -m' -ldflags='-s -w'
```

**Cross-Platform Static:**
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w'
```

## Flag Precedence and Interaction

### Build Mode Effects

Different build modes affect which flags are valid:

- **default** - Normal executable
- **archive** - Non-main package as .a file
- **c-archive** - C archive (.a) for C programs
- **c-shared** - C shared library (.so/.dll)
- **shared** - Go shared library
- **exe** - Executable (default for main packages)
- **pie** - Position-independent executable
- **plugin** - Go plugin

**Examples:**
```bash
# Build C archive
go build -buildmode=c-archive

# Build shared library
go build -buildmode=shared

# Build plugin
go build -buildmode=plugin
```

### Environment Variable Interaction

Some flags interact with environment variables:

```bash
# CGO_ENABLED affects linking
CGO_ENABLED=0 go build -ldflags='-linkmode external'  # Error!

# GOEXPERIMENT affects available features
GOEXPERIMENT=arenas go build

# GODEBUG affects runtime behavior
GODEBUG=gctrace=1 go build
```

## Debugging Build Flags

### Verify Flag Application

```bash
# Print all compiler invocations
go build -x -gcflags='-m'

# Check what flags are actually used
go build -n  # Print commands without executing
```

### Common Issues

**Issue: Flags not applied**
```bash
# Wrong: quotes around individual flags
go build -ldflags='-s' -ldflags='-w'  # Only -w applied

# Correct: all flags in one string
go build -ldflags='-s -w'
```

**Issue: Pattern not matching**
```bash
# Wrong: relative path
go build -gcflags='./pkg=-N -l'

# Correct: import path
go build -gcflags='github.com/user/repo/pkg=-N -l'
```

**Issue: Variable injection not working**
```bash
# Wrong: missing package path
go build -ldflags='-X Version=1.0.0'

# Correct: full path to variable
go build -ldflags='-X main.Version=1.0.0'
```
