# Go Environment Variables

Comprehensive guide to Go environment variables that control build behavior.

## Core Environment Variables

### GOROOT

Location of Go installation.

**Default:** Installation directory (e.g., `/usr/local/go`)

**Usage:**
```bash
export GOROOT=/usr/local/go
```

**When to set:** Rarely needed; only when using custom Go installation.

### GOPATH

Workspace location for Go code (pre-modules).

**Default:** `$HOME/go`

**Usage:**
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

**Structure:**
```
$GOPATH/
├── bin/       # Compiled binaries
├── pkg/       # Package objects
└── src/       # Source code (pre-modules)
```

**With modules:** GOPATH still used for:
- Module cache: `$GOPATH/pkg/mod`
- Installed binaries: `$GOPATH/bin`

### GOBIN

Directory where `go install` places binaries.

**Default:** `$GOPATH/bin` (or `$HOME/go/bin` if GOPATH not set)

**Usage:**
```bash
export GOBIN=/usr/local/bin
go install github.com/user/tool@latest
```

### GOOS and GOARCH

Target operating system and architecture.

**Common GOOS values:**
- `linux` - Linux
- `darwin` - macOS
- `windows` - Windows
- `freebsd` - FreeBSD
- `openbsd` - OpenBSD
- `netbsd` - NetBSD
- `android` - Android
- `ios` - iOS
- `js` - JavaScript (WASM)
- `plan9` - Plan 9
- `aix` - AIX
- `dragonfly` - DragonFly BSD
- `illumos` - illumos
- `solaris` - Solaris

**Common GOARCH values:**
- `amd64` - x86-64
- `386` - x86-32
- `arm` - ARM 32-bit
- `arm64` - ARM 64-bit
- `ppc64` - PowerPC 64-bit big-endian
- `ppc64le` - PowerPC 64-bit little-endian
- `mips` - MIPS 32-bit big-endian
- `mipsle` - MIPS 32-bit little-endian
- `mips64` - MIPS 64-bit big-endian
- `mips64le` - MIPS 64-bit little-endian
- `riscv64` - RISC-V 64-bit
- `s390x` - IBM z/Architecture
- `wasm` - WebAssembly

**Usage:**
```bash
# Cross-compile for Linux on AMD64
GOOS=linux GOARCH=amd64 go build

# Cross-compile for Windows on ARM64
GOOS=windows GOARCH=arm64 go build

# Cross-compile for macOS on Apple Silicon
GOOS=darwin GOARCH=arm64 go build

# List all supported platforms
go tool dist list
```

**ARM variants (GOARM):**
```bash
# ARMv5
GOOS=linux GOARCH=arm GOARM=5 go build

# ARMv6
GOOS=linux GOARCH=arm GOARM=6 go build

# ARMv7
GOOS=linux GOARCH=arm GOARM=7 go build
```

**MIPS variants (GOMIPS, GOMIPS64):**
```bash
# MIPS32 with hardware float
GOOS=linux GOARCH=mips GOMIPS=hardfloat go build

# MIPS32 with soft float
GOOS=linux GOARCH=mips GOMIPS=softfloat go build
```

**AMD64 variants (GOAMD64):**
```bash
# v1: baseline (default)
GOARCH=amd64 GOAMD64=v1 go build

# v2: +POPCNT, +SSE3, +SSE4.1, +SSE4.2, +SSSE3
GOARCH=amd64 GOAMD64=v2 go build

# v3: +AVX, +AVX2, +BMI1, +BMI2, +FMA
GOARCH=amd64 GOAMD64=v3 go build

# v4: +AVX512
GOARCH=amd64 GOAMD64=v4 go build
```

## Build Control Variables

### CGO_ENABLED

Enable or disable cgo.

**Values:** `0` (disabled) or `1` (enabled)

**Default:** `1` (enabled)

**Usage:**
```bash
# Disable CGO for static builds
CGO_ENABLED=0 go build

# Enable CGO
CGO_ENABLED=1 go build
```

**Effects:**
- `CGO_ENABLED=0`: Pure Go, static binaries, no C dependencies
- `CGO_ENABLED=1`: Can use C code, may have dynamic dependencies

**Cross-compilation:** CGO disabled by default when GOOS/GOARCH differ from host.

### CGO_CFLAGS, CGO_CPPFLAGS, CGO_CXXFLAGS, CGO_LDFLAGS

C compiler and linker flags for cgo.

**Usage:**
```bash
# Set C compiler flags
export CGO_CFLAGS="-I/opt/include -O2"

# Set C++ compiler flags
export CGO_CXXFLAGS="-std=c++17"

# Set linker flags
export CGO_LDFLAGS="-L/opt/lib -lmylib"

go build
```

### CC, CXX, AR

C, C++, and archive tool commands for cgo.

**Usage:**
```bash
# Use specific C compiler
export CC=clang
export CXX=clang++
go build

# Cross-compile with specific toolchain
export CC=aarch64-linux-gnu-gcc
export CXX=aarch64-linux-gnu-g++
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build
```

### GOTOOLCHAIN

Control which Go toolchain to use (Go 1.21+).

**Values:**
- `auto` - Automatically download required version
- `local` - Use only local toolchain
- `<version>` - Specific version (e.g., `go1.21.0`)
- `<version>+auto` - Prefer version, download if needed
- `<version>+path` - Prefer version, search PATH

**Usage:**
```bash
# Use only local toolchain
export GOTOOLCHAIN=local

# Allow automatic downloads
export GOTOOLCHAIN=auto

# Require specific version
export GOTOOLCHAIN=go1.21.0
```

**go.mod directive:**
```go
toolchain go1.21.5
```

### GOEXPERIMENT

Enable experimental features.

**Usage:**
```bash
# Enable arenas experiment
GOEXPERIMENT=arenas go build

# Multiple experiments
GOEXPERIMENT=arenas,rangefunc go build
```

**List experiments:**
```bash
go tool compile -V=full
```

## Module Variables

### GO111MODULE

Module mode control (legacy, mostly deprecated).

**Values:**
- `on` - Force module mode
- `off` - Force GOPATH mode
- `auto` - Auto-detect (default)

**Usage:**
```bash
export GO111MODULE=on
```

**Note:** As of Go 1.16+, module mode is default. This variable is rarely needed.

### GOMODCACHE

Location of module cache.

**Default:** `$GOPATH/pkg/mod`

**Usage:**
```bash
export GOMODCACHE=/custom/cache/path
```

### GONOPROXY, GONOSUMDB, GOPRIVATE

Private module configuration.

**GOPRIVATE:** Comma-separated list of glob patterns for private modules.
```bash
export GOPRIVATE=github.com/mycompany/*,gitlab.internal.com/*
```

**GONOPROXY:** Modules that should not use proxy.
```bash
export GONOPROXY=github.com/mycompany/*
```

**GONOSUMDB:** Modules that should not use checksum database.
```bash
export GONOSUMDB=github.com/mycompany/*
```

**Shorthand:** Setting GOPRIVATE sets both GONOPROXY and GONOSUMDB.

### GOPROXY

Module proxy configuration.

**Default:** `https://proxy.golang.org,direct`

**Usage:**
```bash
# Use custom proxy
export GOPROXY=https://goproxy.io,direct

# Multiple proxies with fallback
export GOPROXY=https://proxy1.com,https://proxy2.com,direct

# Direct only (no proxy)
export GOPROXY=direct

# Disable module downloads
export GOPROXY=off
```

**Syntax:**
- `https://proxy.com` - Use proxy
- `direct` - Direct VCS access
- `off` - Disable downloads
- `|` - Return error if left side fails (vs `,` which falls through)

### GOSUMDB

Checksum database for module authentication.

**Default:** `sum.golang.org`

**Usage:**
```bash
# Use custom checksum database
export GOSUMDB=sum.golang.google.cn

# Disable checksum verification
export GOSUMDB=off
```

### GOINSECURE

Allow insecure module downloads (use with caution).

**Usage:**
```bash
# Allow insecure downloads for specific module
export GOINSECURE=example.com/insecure/module

# Multiple modules
export GOINSECURE=example.com/*,other.com/pkg
```

## Build Cache Variables

### GOCACHE

Location of build cache.

**Default:** Platform-specific (e.g., `$HOME/.cache/go-build`)

**Usage:**
```bash
# Use custom cache location
export GOCACHE=/tmp/go-build

# Disable caching
export GOCACHE=off
```

### GODEBUG

Runtime and build debugging options.

**Common values:**

**Memory/GC:**
- `gctrace=1` - Print GC trace
- `gcpacertrace=1` - Print pacer trace
- `madvdontneed=0` - Use MADV_FREE instead of MADV_DONTNEED
- `scavtrace=1` - Print scavenger trace

**HTTP:**
- `http2client=0` - Disable HTTP/2 client
- `http2server=0` - Disable HTTP/2 server
- `http2debug=1` or `2` - Enable HTTP/2 debug logging

**Net:**
- `netdns=go` - Force pure Go DNS resolver
- `netdns=cgo` - Force cgo DNS resolver

**TLS:**
- `tls13=1` - Enable TLS 1.3
- `tlsrsakex=1` - Enable RSA key exchange

**Other:**
- `asyncpreemptoff=1` - Disable async preemption
- `cgocheck=0|1|2` - Control cgo pointer checks
- `cpu.*=off` - Disable CPU features
- `inittrace=1` - Trace init execution
- `panicnil=1` - Panic on nil pointer

**Usage:**
```bash
# Enable GC trace
GODEBUG=gctrace=1 go run main.go

# Multiple options
GODEBUG=gctrace=1,gcpacertrace=1 go run main.go

# At build time
GODEBUG=http2client=0 go build
```

## Compiler/Linker Variables

### GCCGO

gccgo compiler command.

**Usage:**
```bash
export GCCGO=/usr/bin/gccgo-12
go build -compiler=gccgo
```

### GO_EXTLINK_ENABLED

Force external linking.

**Values:** `0` or `1`

**Usage:**
```bash
export GO_EXTLINK_ENABLED=1
go build
```

## Test Variables

### GORACE

Race detector options.

**Options:**
- `log_path=<path>` - Write logs to file
- `exitcode=<code>` - Exit code when race detected
- `strip_path_prefix=<prefix>` - Strip prefix from paths
- `history_size=<size>` - Per-goroutine history size
- `halt_on_error=<0|1>` - Stop after first race

**Usage:**
```bash
GORACE="log_path=/tmp/race.log halt_on_error=1" go test -race
```

### GOMAXPROCS

Maximum number of CPUs for execution.

**Default:** Number of CPUs

**Usage:**
```bash
# Use 4 CPUs
GOMAXPROCS=4 go test

# Use 1 CPU (useful for debugging)
GOMAXPROCS=1 go run main.go
```

## Debugging Variables

### GOTRACEBACK

Control stack trace on panic.

**Values:**
- `none` - No stack traces
- `single` - Single goroutine (default)
- `all` - All goroutines
- `system` - All goroutines + system frames
- `crash` - All goroutines + system frames + OS dump

**Usage:**
```bash
GOTRACEBACK=all go run main.go
```

### GOROOT_FINAL

Expected GOROOT after installation (rarely needed).

**Usage:**
```bash
export GOROOT_FINAL=/usr/local/go
```

## Performance Variables

### GOGC

Garbage collection target percentage.

**Default:** `100` (GC when heap doubles)

**Usage:**
```bash
# More aggressive GC
GOGC=50 go run main.go

# Less aggressive GC
GOGC=200 go run main.go

# Disable GC (not recommended)
GOGC=off go run main.go
```

### GOMEMLIMIT

Memory limit for Go runtime (Go 1.19+).

**Usage:**
```bash
# Set 1 GB limit
GOMEMLIMIT=1GiB go run main.go

# Set 512 MB limit
GOMEMLIMIT=512MiB go run main.go
```

## Common Configurations

### Development Environment
```bash
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
export CGO_ENABLED=1
export GOTOOLCHAIN=local
```

### CI/CD Environment
```bash
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOCACHE=/tmp/go-cache
export GOMODCACHE=/tmp/go-mod
export GOTOOLCHAIN=local
```

### Private Module Setup
```bash
export GOPRIVATE=github.com/mycompany/*
export GOPROXY=https://proxy.golang.org,direct
export GONOSUMDB=github.com/mycompany/*
```

### Cross-Compilation Setup
```bash
# Linux ARM64
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=0

# Windows AMD64
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0
```

### Performance Tuning
```bash
# Production runtime
export GOMAXPROCS=8
export GOGC=100
export GOMEMLIMIT=4GiB
export GODEBUG=madvdontneed=0

# Development
export GOMAXPROCS=4
export GOGC=50
```

## Viewing Current Environment

```bash
# View all Go environment variables
go env

# View specific variable
go env GOPATH

# View in JSON format
go env -json

# View effective build settings
go env GOOS GOARCH CGO_ENABLED
```

## Setting Environment Permanently

### Per-User (go env)
```bash
# Set GOBIN permanently
go env -w GOBIN=/usr/local/bin

# Set private modules
go env -w GOPRIVATE=github.com/mycompany/*

# Unset variable
go env -u GOBIN
```

### Shell Profile
Add to `~/.bashrc`, `~/.zshrc`, etc.:
```bash
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
export GOPRIVATE=github.com/mycompany/*
```

### System-Wide
Edit `/etc/environment` or `/etc/profile.d/golang.sh`:
```bash
GOPATH=/opt/go
GOBIN=/opt/go/bin
```
