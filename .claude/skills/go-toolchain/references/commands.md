# Go Command Reference

Comprehensive documentation of all `go` commands and their options.

## Core Commands

### go build

Compile packages and dependencies.

**Syntax:**
```bash
go build [-o output] [build flags] [packages]
```

**Common Flags:**
- `-o <file>` - Force build to write resulting executable or object to named output file
- `-a` - Force rebuilding of packages that are already up-to-date
- `-n` - Print commands but do not run them
- `-v` - Print names of packages as they are compiled
- `-x` - Print commands as they are executed
- `-race` - Enable data race detection (sets `-installsuffix race`)
- `-msan` - Enable memory sanitizer (requires C compiler)
- `-asan` - Enable address sanitizer
- `-cover` - Enable code coverage instrumentation
- `-work` - Print temporary work directory name and keep it
- `-buildmode <mode>` - Build mode to use (default, archive, c-archive, c-shared, shared, exe, pie, plugin)
- `-buildvcs <auto|true|false>` - Whether to stamp binaries with version control info
- `-compiler <name>` - Name of compiler to use (gccgo or gc)
- `-gccgoflags '[pattern=]arg list'` - Arguments to pass on each gccgo compiler/linker invocation
- `-gcflags '[pattern=]arg list'` - Arguments to pass on each go tool compile invocation
- `-ldflags '[pattern=]arg list'` - Arguments to pass on each go tool link invocation
- `-linkshared` - Build code that will be linked against shared libraries
- `-mod <mode>` - Module download mode (readonly, vendor, mod)
- `-modcacherw` - Leave newly-created directories in module cache read-write
- `-modfile <file>` - Read and write alternate go.mod file
- `-overlay <file>` - Read build overlay from file
- `-pgo <file>` - Profile-guided optimization
- `-pkgdir <dir>` - Install and load all packages from dir
- `-tags <tag,list>` - Comma-separated list of build tags
- `-trimpath` - Remove all file system paths from resulting executable
- `-toolexec 'cmd args'` - Use program to invoke toolchain programs

**Examples:**
```bash
# Basic build
go build

# Build with output name
go build -o myapp

# Build with race detector
go build -race

# Build with verbose output
go build -v

# Build with custom linker flags (e.g., strip debug info)
go build -ldflags="-s -w"

# Build for specific OS/arch
GOOS=linux GOARCH=amd64 go build

# Build with build tags
go build -tags production

# Build with trimmed paths (reproducible builds)
go build -trimpath
```

### go run

Compile and run Go program.

**Syntax:**
```bash
go run [build flags] [-exec xprog] package [arguments...]
```

**Common Flags:**
- All build flags apply
- `-exec <xprog>` - Run the binary using xprog

**Examples:**
```bash
# Run main package
go run main.go

# Run with arguments
go run main.go arg1 arg2

# Run with build flags
go run -race main.go
```

### go test

Test packages.

**Syntax:**
```bash
go test [build/test flags] [packages] [build/test flags & test binary flags]
```

**Common Flags:**
- `-v` - Verbose output
- `-run <regexp>` - Run only tests matching regexp
- `-bench <regexp>` - Run benchmarks matching regexp
- `-benchmem` - Print memory allocation statistics for benchmarks
- `-benchtime <t>` - Run benchmarks for duration t (e.g., 10s, 100x)
- `-count <n>` - Run tests and benchmarks n times
- `-cover` - Enable coverage analysis
- `-covermode <set|count|atomic>` - Set coverage mode
- `-coverpkg <pattern>` - Apply coverage to packages matching pattern
- `-coverprofile <file>` - Write coverage profile to file
- `-cpu <list>` - Specify GOMAXPROCS values (e.g., 1,2,4)
- `-failfast` - Do not start new tests after first failure
- `-json` - Output in JSON format
- `-list <regexp>` - List tests, benchmarks, or examples matching regexp
- `-parallel <n>` - Allow parallel execution of test functions (default: GOMAXPROCS)
- `-short` - Tell long-running tests to shorten their run time
- `-timeout <d>` - Panic if test runs longer than duration d (default: 10m)
- `-shuffle <off|on|N>` - Randomize execution order
- `-trace <file>` - Write execution trace to file
- `-vet <list>` - Configure the invocation of "go vet" during "go test"

**Examples:**
```bash
# Run all tests
go test

# Run tests with coverage
go test -cover

# Run tests with coverage profile
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Run specific test
go test -run TestMyFunction

# Run benchmarks
go test -bench=.

# Run benchmarks with memory stats
go test -bench=. -benchmem

# Run tests in all subdirectories
go test ./...

# Run tests with race detector
go test -race ./...

# Run tests with timeout
go test -timeout 30s

# Run tests in parallel
go test -parallel 4
```

### go install

Compile and install packages and dependencies.

**Syntax:**
```bash
go install [build flags] [packages]
```

**Common Flags:**
- All build flags apply
- Installs to `$GOPATH/bin` or `$GOBIN`

**Examples:**
```bash
# Install current package
go install

# Install specific package
go install example.com/pkg@latest

# Install specific version
go install example.com/pkg@v1.2.3

# Install with custom GOBIN
GOBIN=/custom/path go install
```

### go get

Add dependencies to current module and install them.

**Syntax:**
```bash
go get [-d] [-t] [-u] [build flags] [packages]
```

**Common Flags:**
- `-d` - Download only (don't install)
- `-t` - Consider test dependencies
- `-u` - Update packages and dependencies
- `-u=patch` - Update dependencies to patch releases

**Examples:**
```bash
# Add dependency
go get example.com/pkg

# Add specific version
go get example.com/pkg@v1.2.3

# Add latest version
go get example.com/pkg@latest

# Update all dependencies
go get -u ./...

# Update to latest patch versions
go get -u=patch ./...
```

### go mod

Module maintenance.

**Subcommands:**

#### go mod init
```bash
go mod init [module-path]
```
Initialize new module in current directory.

#### go mod tidy
```bash
go mod tidy [-e] [-v] [-x] [-go=version] [-compat=version]
```
Add missing and remove unused modules.

Flags:
- `-e` - Attempt to proceed despite errors
- `-v` - Verbose output
- `-x` - Print commands
- `-go=version` - Set go version directive
- `-compat=version` - Set compatibility version

#### go mod download
```bash
go mod download [-x] [-json] [modules]
```
Download modules to local cache.

Flags:
- `-x` - Print commands
- `-json` - Print JSON output

#### go mod verify
```bash
go mod verify
```
Verify dependencies have expected content.

#### go mod graph
```bash
go mod graph [-go=version]
```
Print module requirement graph.

#### go mod why
```bash
go mod why [-m] [-vendor] packages...
```
Explain why packages or modules are needed.

Flags:
- `-m` - Treat arguments as modules
- `-vendor` - Exclude tests of dependencies

#### go mod edit
```bash
go mod edit [editing flags] [go.mod]
```
Edit go.mod programmatically.

Flags:
- `-fmt` - Reformat go.mod
- `-module path` - Change module path
- `-go=version` - Set Go version
- `-require=path@version` - Add requirement
- `-droprequire=path` - Drop requirement
- `-replace=old[@v]=new[@v]` - Add replacement
- `-dropreplace=old[@v]` - Drop replacement
- `-retract=version` - Add retraction
- `-dropretract=version` - Drop retraction

#### go mod vendor
```bash
go mod vendor [-e] [-v] [-o outdir]
```
Make vendored copy of dependencies.

Flags:
- `-e` - Attempt to proceed despite errors
- `-v` - Verbose output
- `-o outdir` - Output to different directory

**Examples:**
```bash
# Initialize module
go mod init github.com/user/repo

# Add missing dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Show dependency graph
go mod graph

# Explain why package is needed
go mod why example.com/pkg

# Update go version
go mod edit -go=1.21

# Add replace directive
go mod edit -replace=old@v1.0.0=new@v2.0.0

# Create vendor directory
go mod vendor
```

## Additional Commands

### go clean

Remove object files and cached files.

**Syntax:**
```bash
go clean [clean flags] [packages]
```

**Flags:**
- `-i` - Remove corresponding installed archive or binary
- `-n` - Print commands but do not run them
- `-r` - Apply recursively to dependencies
- `-x` - Print commands as executed
- `-cache` - Remove entire go build cache
- `-testcache` - Expire all test results in go build cache
- `-modcache` - Remove entire module download cache
- `-fuzzcache` - Remove files in fuzz cache

**Examples:**
```bash
# Clean current package
go clean

# Clean build cache
go clean -cache

# Clean test cache
go clean -testcache

# Clean module cache
go clean -modcache

# Clean all
go clean -i -r -x
```

### go fmt

Format Go source files.

**Syntax:**
```bash
go fmt [-n] [-x] [packages]
```

**Flags:**
- `-n` - Print commands but do not run them
- `-x` - Print commands as executed

**Examples:**
```bash
# Format current package
go fmt

# Format all packages
go fmt ./...

# Use gofmt directly for more control
gofmt -w .
```

### go vet

Examine Go source code and report suspicious constructs.

**Syntax:**
```bash
go vet [-n] [-x] [build flags] [vet flags] [packages]
```

**Common Checks:**
- Useless assignments
- Unreachable code
- Printf format strings
- Struct tags
- And many more

**Examples:**
```bash
# Vet current package
go vet

# Vet all packages
go vet ./...

# Vet with specific checks
go vet -printf ./...
```

### go generate

Generate Go files by processing source.

**Syntax:**
```bash
go generate [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
```

**Flags:**
- `-run <regexp>` - Only run directives matching regexp
- `-n` - Print commands but do not run them
- `-v` - Verbose output
- `-x` - Print commands as executed

**Examples:**
```bash
# Generate for current package
go generate

# Generate for all packages
go generate ./...

# Generate with pattern
go generate -run stringer ./...
```

### go env

Print Go environment information.

**Syntax:**
```bash
go env [-json] [-u] [-w] [var ...]
```

**Flags:**
- `-json` - Print in JSON format
- `-u <var>` - Unset environment variable
- `-w <var>=<value>` - Set environment variable

**Examples:**
```bash
# Print all environment variables
go env

# Print specific variable
go env GOPATH

# Print in JSON
go env -json

# Set environment variable
go env -w GOBIN=/custom/bin

# Unset environment variable
go env -u GOBIN
```

### go list

List packages or modules.

**Syntax:**
```bash
go list [-f format] [-json] [-m] [list flags] [packages]
```

**Flags:**
- `-f <format>` - Use template format
- `-json` - Print in JSON format
- `-m` - List modules instead of packages
- `-u` - Add available upgrade information
- `-versions` - List available versions

**Examples:**
```bash
# List all packages
go list ./...

# List with JSON output
go list -json ./...

# List modules
go list -m all

# List module with versions
go list -m -versions example.com/pkg

# List with custom format
go list -f '{{.ImportPath}}'
```

### go doc

Show documentation for package or symbol.

**Syntax:**
```bash
go doc [doc flags] [package|[package.]symbol[.methodOrField]]
```

**Flags:**
- `-all` - Show all documentation
- `-c` - Case-sensitive matching
- `-cmd` - Show documentation for commands
- `-short` - One-line representation for each symbol
- `-src` - Show source code
- `-u` - Show unexported symbols

**Examples:**
```bash
# Show package documentation
go doc fmt

# Show function documentation
go doc fmt.Println

# Show all documentation
go doc -all fmt

# Show with source
go doc -src fmt.Println
```

### go version

Print Go version.

**Syntax:**
```bash
go version [-m] [-v] [file ...]
```

**Flags:**
- `-m` - Print module information
- `-v` - Verbose output

**Examples:**
```bash
# Print Go version
go version

# Print version of binary
go version /path/to/binary

# Print with module info
go version -m /path/to/binary
```

### go work

Workspace maintenance (Go 1.18+).

**Subcommands:**
- `go work init` - Initialize workspace file
- `go work use` - Add modules to workspace
- `go work edit` - Edit go.work file
- `go work sync` - Sync workspace dependencies

**Examples:**
```bash
# Initialize workspace
go work init ./module1 ./module2

# Add module to workspace
go work use ./module3

# Sync workspace
go work sync
```

## Special Commands

### go fix

Update packages to use new APIs.

```bash
go fix [packages]
```

### go tool

Run specified go tool.

```bash
go tool <command> [args...]
```

Common tools:
- `go tool compile` - Run compiler
- `go tool link` - Run linker
- `go tool cover` - Coverage tool
- `go tool pprof` - CPU profiler
- `go tool trace` - Execution tracer
- `go tool nm` - Symbol table
- `go tool objdump` - Disassembler

### go bug

Start a bug report.

```bash
go bug
```
