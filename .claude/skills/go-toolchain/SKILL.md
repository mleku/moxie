---
name: go-toolchain
description: This skill should be used when working with the Go toolchain, including building, testing, cross-compiling, managing modules, or configuring Go build settings. Use when executing go commands, setting compiler/linker flags, configuring environment variables, or troubleshooting build issues.
---

# Go Toolchain

## Overview

Comprehensive knowledge of the Go toolchain, including all `go` commands, compiler and linker flags, environment variables, and cross-compilation techniques. This skill provides detailed reference documentation to help execute Go build tasks correctly and efficiently.

## When to Use This Skill

Use this skill when:
- Building Go programs with specific compiler or linker flags
- Cross-compiling Go programs for different platforms (GOOS/GOARCH)
- Managing Go modules (go.mod, dependencies, versioning)
- Running Go tests with specific options (coverage, race detection, benchmarks)
- Configuring Go environment variables
- Troubleshooting build issues or understanding go command options
- Setting up CI/CD pipelines for Go projects
- Optimizing binary size or build performance
- Injecting build-time variables (version, commit hash, build date)
- Creating static binaries or using CGO

## Core Go Commands

The Go toolchain provides comprehensive commands for all development tasks. For detailed information on any command, consult `references/commands.md`.

### Essential Commands Quick Reference

**Building and Running:**
- `go build` - Compile packages and dependencies
- `go run` - Compile and run Go program
- `go install` - Compile and install packages
- `go clean` - Remove object files and cached files

**Testing:**
- `go test` - Test packages with various options
  - Coverage: `-cover`, `-coverprofile`
  - Benchmarks: `-bench`, `-benchmem`
  - Race detection: `-race`
  - Parallel execution: `-parallel`

**Module Management:**
- `go mod init` - Initialize new module
- `go mod tidy` - Add missing and remove unused modules
- `go mod download` - Download modules to local cache
- `go mod vendor` - Make vendored copy of dependencies
- `go get` - Add dependencies to current module

**Code Quality:**
- `go fmt` - Format Go source files
- `go vet` - Examine code and report suspicious constructs
- `go generate` - Generate Go files by processing source

**Information:**
- `go env` - Print Go environment information
- `go list` - List packages or modules
- `go doc` - Show documentation for package or symbol
- `go version` - Print Go version

### Finding the Right Command

To find detailed information about any command:

1. **Read `references/commands.md`** for comprehensive documentation of all commands, flags, and examples
2. Use `go help <command>` for built-in help
3. Use `go <command> -h` for command-specific flags

## Build Flags and Optimization

Go provides extensive compiler and linker flags for controlling build behavior. Consult `references/build_flags.md` for complete documentation.

### Common Build Patterns

**Development Build (fast compilation, full debug):**
```bash
go build -gcflags='all=-N -l' -race
```

**Production Build (optimized, stripped, versioned):**
```bash
go build -trimpath -ldflags='-s -w -X main.Version=1.0.0'
```

**Static Binary (no external dependencies):**
```bash
CGO_ENABLED=0 go build -ldflags='-s -w -extldflags "-static"'
```

**Size-Optimized Build:**
```bash
go build -ldflags='-s -w' -trimpath
```

### Compiler Flags (-gcflags)

Apply compiler flags using `-gcflags`:

**Disable optimizations for debugging:**
```bash
go build -gcflags='-N -l'
```

**Show optimization decisions:**
```bash
go build -gcflags='-m'
```

**Apply to specific package:**
```bash
go build -gcflags='example.com/pkg=-N -l'
```

### Linker Flags (-ldflags)

Apply linker flags using `-ldflags`:

**Strip debug information (reduce binary size):**
```bash
go build -ldflags='-s -w'
```

**Inject build-time variables:**
```bash
go build -ldflags="-X main.Version=1.2.3 -X main.BuildTime=$(date -u +%Y-%m-%d)"
```

**Static linking:**
```bash
go build -ldflags='-linkmode external -extldflags "-static"'
```

### Flag Reference Location

For comprehensive flag documentation:
- **Compiler flags:** See `references/build_flags.md` → "Compiler Flags (-gcflags)"
- **Linker flags:** See `references/build_flags.md` → "Linker Flags (-ldflags)"
- **Assembly flags:** See `references/build_flags.md` → "Assembly Flags (-asmflags)"

## Environment Variables

Go behavior is controlled by numerous environment variables. Consult `references/environment.md` for complete documentation.

### Critical Environment Variables

**Build Configuration:**
- `GOOS` - Target operating system (linux, darwin, windows, etc.)
- `GOARCH` - Target architecture (amd64, arm64, 386, etc.)
- `CGO_ENABLED` - Enable/disable cgo (0 or 1)

**Module Configuration:**
- `GOPRIVATE` - Private module patterns
- `GOPROXY` - Module proxy URL
- `GOMODCACHE` - Module cache location

**Build Cache:**
- `GOCACHE` - Build cache location
- `GOMODCACHE` - Module download cache

**Performance:**
- `GOMAXPROCS` - Maximum number of CPUs
- `GOGC` - Garbage collection target percentage
- `GOMEMLIMIT` - Memory limit for Go runtime

### Setting Environment Variables

**Temporarily (single command):**
```bash
GOOS=linux GOARCH=amd64 go build
```

**Permanently (go env):**
```bash
go env -w GOPRIVATE=github.com/mycompany/*
```

**View current settings:**
```bash
go env
go env GOOS GOARCH CGO_ENABLED
```

### Environment Variable Reference Location

For comprehensive environment variable documentation:
- See `references/environment.md` for detailed explanations of all variables
- Includes common configurations for development, CI/CD, and production

## Cross-Compilation

Go excels at cross-platform compilation. Consult `references/cross_compile.md` for comprehensive cross-compilation guide.

### Quick Cross-Compilation

**Pure Go programs (no cgo):**
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build

# Windows AMD64
GOOS=windows GOARCH=amd64 go build

# Linux ARM64 (Raspberry Pi, AWS Graviton)
GOOS=linux GOARCH=arm64 go build
```

**List all supported platforms:**
```bash
go tool dist list
```

### CGO Cross-Compilation

CGO requires cross-compilation toolchain:

```bash
# Linux to ARM64 with CGO
export CC=aarch64-linux-gnu-gcc
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build

# Linux to Windows with CGO
export CC=x86_64-w64-mingw32-gcc
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build
```

### Multi-Platform Build Scripts

See `references/cross_compile.md` for:
- Makefile examples for multi-platform builds
- Bash scripts for parallel compilation
- Go scripts for automated builds
- Docker-based cross-compilation
- Platform-specific considerations

### Cross-Compilation Reference Location

For detailed cross-compilation documentation:
- See `references/cross_compile.md` for complete guide
- Includes build scripts, Docker examples, CGO setup, and troubleshooting

## Workflow Decision Tree

Choose the appropriate workflow based on the task:

### Build Tasks
1. **Simple build** → `go build`
2. **Production build** → Read `references/build_flags.md` → Apply `-ldflags='-s -w'` and `-trimpath`
3. **Debug build** → Read `references/build_flags.md` → Apply `-gcflags='all=-N -l'`
4. **Cross-platform build** → Read `references/cross_compile.md` → Set GOOS/GOARCH
5. **Static binary** → Read `references/build_flags.md` → Set CGO_ENABLED=0

### Test Tasks
1. **Run tests** → `go test`
2. **With coverage** → Read `references/commands.md` → Use `-cover -coverprofile`
3. **Race detection** → `go test -race`
4. **Benchmarks** → `go test -bench=. -benchmem`
5. **Specific test** → `go test -run TestName`

### Module Tasks
1. **New module** → `go mod init <module-path>`
2. **Add dependency** → `go get package@version`
3. **Update dependencies** → `go get -u ./...`
4. **Clean dependencies** → `go mod tidy`
5. **Vendor dependencies** → `go mod vendor`

### Environment Tasks
1. **Check settings** → `go env`
2. **Set private modules** → Read `references/environment.md` → Set GOPRIVATE
3. **Configure proxy** → Read `references/environment.md` → Set GOPROXY
4. **Tune performance** → Read `references/environment.md` → Set GOMAXPROCS, GOGC

## Using the Reference Documentation

This skill includes four comprehensive reference documents:

### 1. references/commands.md
**Use when:** Need detailed information about any `go` command

**Contains:**
- All `go` commands with full flag documentation
- Examples for every command
- Common usage patterns
- Special commands and tools

**When to read:**
- Before using unfamiliar `go` commands
- To find specific flags or options
- To understand command interactions
- To see practical examples

### 2. references/build_flags.md
**Use when:** Building Go programs with specific compiler or linker behavior

**Contains:**
- Complete compiler flags (-gcflags)
- Complete linker flags (-ldflags)
- Assembly flags (-asmflags)
- Flag combination patterns
- Variable injection techniques
- Optimization strategies

**When to read:**
- Setting up production builds
- Debugging compilation issues
- Optimizing binary size
- Injecting version information
- Understanding optimization levels

### 3. references/environment.md
**Use when:** Configuring Go environment or troubleshooting environment issues

**Contains:**
- All Go environment variables
- Default values and effects
- Platform-specific variables
- Module configuration variables
- Performance tuning variables
- Common configurations

**When to read:**
- Setting up new development environment
- Configuring CI/CD pipelines
- Managing private modules
- Cross-compilation setup
- Performance tuning

### 4. references/cross_compile.md
**Use when:** Building Go programs for different platforms

**Contains:**
- Platform quick reference
- Pure Go cross-compilation
- CGO cross-compilation setup
- Multi-platform build scripts
- Docker-based builds
- Platform-specific considerations
- Troubleshooting guide

**When to read:**
- Building for multiple platforms
- Setting up release pipelines
- Creating static binaries
- Handling platform-specific code
- CGO cross-compilation

## Common Workflows

### 1. Production Binary Release

**Goal:** Build optimized, stripped binaries for multiple platforms with version information.

**Steps:**
1. Read `references/cross_compile.md` for multi-platform build scripts
2. Read `references/build_flags.md` for optimization flags
3. Use build script pattern:

```bash
VERSION=$(git describe --tags --always)

# Linux AMD64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.Version=$VERSION" \
    -o dist/myapp-linux-amd64

# macOS ARM64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
    go build -trimpath -ldflags="-s -w -X main.Version=$VERSION" \
    -o dist/myapp-darwin-arm64

# Windows AMD64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.Version=$VERSION" \
    -o dist/myapp-windows-amd64.exe
```

### 2. Debug Session Setup

**Goal:** Build with full debug symbols and race detection.

**Steps:**
1. Read `references/build_flags.md` for debug flags
2. Build with debugging enabled:

```bash
go build -gcflags='all=-N -l' -race -o myapp-debug
```

3. Run with debug settings:
```bash
GODEBUG=gctrace=1 GOMAXPROCS=1 ./myapp-debug
```

### 3. Module Dependency Management

**Goal:** Update dependencies safely and verify module integrity.

**Steps:**
1. Read `references/commands.md` → "go mod" section
2. Update and verify:

```bash
# Update to latest compatible versions
go get -u ./...

# Or update to latest patch versions only
go get -u=patch ./...

# Clean up unused dependencies
go mod tidy

# Verify dependencies
go mod verify

# Check for available updates
go list -u -m all
```

### 4. CI/CD Pipeline Setup

**Goal:** Configure reproducible builds in CI/CD.

**Steps:**
1. Read `references/environment.md` → "CI/CD Environment"
2. Read `references/build_flags.md` → "Production Build"
3. Use environment configuration:

```bash
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOCACHE=/tmp/go-cache
export GOMODCACHE=/tmp/go-mod
export GOTOOLCHAIN=local

# Build
go build -trimpath -ldflags="-s -w -buildid="
```

## Troubleshooting

### Build Issues

**Problem:** Binary too large
**Solution:** Read `references/build_flags.md` → "Size-Optimized Build"
```bash
go build -ldflags='-s -w' -trimpath
```

**Problem:** Cross-compilation with CGO fails
**Solution:** Read `references/cross_compile.md` → "CGO Cross-Compilation"
- Install cross-compiler toolchain
- Set CC/CXX environment variables

**Problem:** Module download fails
**Solution:** Read `references/environment.md` → "Private Module Setup"
- Configure GOPRIVATE for private modules
- Set GOPROXY if behind firewall

### Environment Issues

**Problem:** Wrong Go version being used
**Solution:** Read `references/environment.md` → "GOTOOLCHAIN"
```bash
go env -w GOTOOLCHAIN=local
```

**Problem:** Builds not reproducible
**Solution:** Use `-trimpath` and empty buildid:
```bash
go build -trimpath -ldflags='-buildid='
```

### Testing Issues

**Problem:** Tests timing out
**Solution:** Read `references/commands.md` → "go test" → `-timeout`
```bash
go test -timeout 30s
```

**Problem:** Race detector reporting issues
**Solution:** Read `references/environment.md` → "GORACE"
```bash
GORACE="log_path=/tmp/race.log halt_on_error=1" go test -race
```

## Best Practices

### Always

- Use `-trimpath` for reproducible builds
- Use `go mod tidy` regularly to clean dependencies
- Use `go test -race` to catch race conditions
- Use `go vet` to catch common mistakes
- Set `-ldflags='-s -w'` for production binaries

### Development

- Enable race detector: `go test -race`
- Disable optimizations for debugging: `-gcflags='all=-N -l'`
- Use `go build -n` to see what commands will be executed
- Use `go build -x` to see commands as they execute

### Production

- Strip debug info: `-ldflags='-s -w'`
- Use `-trimpath` for reproducible builds
- Inject version info: `-ldflags='-X main.Version=...'`
- For static binaries: `CGO_ENABLED=0`
- Consider UPX compression for smallest size

### Cross-Compilation

- Disable CGO unless required: `CGO_ENABLED=0`
- Test on target platform when possible
- Use Docker for CGO cross-compilation
- Verify platform with `go tool dist list`

## Quick Command Reference

### Most Used Commands

```bash
# Build current package
go build

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Update dependencies
go get -u ./...

# Clean dependencies
go mod tidy

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build

# Production build
go build -trimpath -ldflags='-s -w'

# Check environment
go env

# Format code
go fmt ./...

# Vet code
go vet ./...
```

For detailed information on any command, consult the appropriate reference document in the `references/` directory.
