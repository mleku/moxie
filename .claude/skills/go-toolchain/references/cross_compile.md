# Cross-Compilation Guide

Comprehensive guide to cross-compiling Go programs for different platforms.

## Quick Reference

### Common Target Platforms

```bash
# Linux AMD64 (most common server)
GOOS=linux GOARCH=amd64 go build

# Linux ARM64 (Raspberry Pi 3+, AWS Graviton)
GOOS=linux GOARCH=arm64 go build

# Linux ARMv7 (Raspberry Pi 2)
GOOS=linux GOARCH=arm GOARM=7 go build

# macOS AMD64 (Intel Mac)
GOOS=darwin GOARCH=amd64 go build

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build

# Windows AMD64
GOOS=windows GOARCH=amd64 go build

# Windows ARM64
GOOS=windows GOARCH=arm64 go build

# FreeBSD AMD64
GOOS=freebsd GOARCH=amd64 go build

# WebAssembly
GOOS=js GOARCH=wasm go build -o main.wasm
```

### List All Supported Platforms

```bash
go tool dist list
```

## Cross-Compilation Basics

### Pure Go Programs

Pure Go programs (no cgo) cross-compile trivially:

```bash
# Build for Linux from any platform
GOOS=linux GOARCH=amd64 go build -o myapp-linux

# Build for Windows from any platform
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# Build for macOS from any platform
GOOS=darwin GOARCH=amd64 go build -o myapp-macos
```

### Programs Using CGO

CGO requires cross-compilation toolchain:

```bash
# Example: Cross-compile for Linux ARM64 with CGO
export CC=aarch64-linux-gnu-gcc
export CXX=aarch64-linux-gnu-g++
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build
```

**Recommendation:** Avoid CGO when possible for easier cross-compilation.

## Multi-Platform Build Scripts

### Makefile Example

```makefile
# Build for all platforms
.PHONY: all
all: linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

# Variables
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"
BUILD_DIR := dist

# Linux AMD64
.PHONY: linux-amd64
linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/myapp-linux-amd64

# Linux ARM64
.PHONY: linux-arm64
linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/myapp-linux-arm64

# macOS AMD64
.PHONY: darwin-amd64
darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/myapp-darwin-amd64

# macOS ARM64 (Apple Silicon)
.PHONY: darwin-arm64
darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/myapp-darwin-arm64

# Windows AMD64
.PHONY: windows-amd64
windows-amd64:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/myapp-windows-amd64.exe

# Clean
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
```

### Bash Script Example

```bash
#!/bin/bash
# build-all.sh

set -e

VERSION=$(git describe --tags --always --dirty)
BINARY_NAME="myapp"
BUILD_DIR="dist"

# Platforms to build
platforms=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm/7"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Create build directory
mkdir -p "$BUILD_DIR"

# Build for each platform
for platform in "${platforms[@]}"; do
    # Split platform into OS, ARCH, and optional ARM version
    IFS='/' read -r -a parts <<< "$platform"
    GOOS="${parts[0]}"
    GOARCH="${parts[1]}"
    GOARM="${parts[2]:-}"

    # Set output name
    output_name="$BINARY_NAME-$GOOS-$GOARCH"
    if [ -n "$GOARM" ]; then
        output_name="$output_name-v$GOARM"
    fi
    if [ "$GOOS" = "windows" ]; then
        output_name="$output_name.exe"
    fi

    # Build
    echo "Building $output_name..."
    env GOOS="$GOOS" GOARCH="$GOARCH" GOARM="$GOARM" \
        go build -trimpath \
        -ldflags="-s -w -X main.Version=$VERSION" \
        -o "$BUILD_DIR/$output_name" \
        .
done

echo "Build complete! Binaries in $BUILD_DIR/"
```

### Go Script Example (Using go:generate)

```go
//go:build ignore
// +build ignore

// build.go - Run with: go run build.go

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Platform struct {
	OS   string
	Arch string
	Arm  string
}

func main() {
	version := getVersion()
	buildDir := "dist"
	binaryName := "myapp"

	platforms := []Platform{
		{"linux", "amd64", ""},
		{"linux", "arm64", ""},
		{"linux", "arm", "7"},
		{"darwin", "amd64", ""},
		{"darwin", "arm64", ""},
		{"windows", "amd64", ""},
	}

	os.MkdirAll(buildDir, 0755)

	for _, p := range platforms {
		outputName := fmt.Sprintf("%s-%s-%s", binaryName, p.OS, p.Arch)
		if p.Arm != "" {
			outputName = fmt.Sprintf("%s-v%s", outputName, p.Arm)
		}
		if p.OS == "windows" {
			outputName += ".exe"
		}

		outputPath := filepath.Join(buildDir, outputName)
		fmt.Printf("Building %s...\n", outputName)

		cmd := exec.Command("go", "build",
			"-trimpath",
			"-ldflags=-s -w -X main.Version="+version,
			"-o", outputPath,
			".",
		)
		cmd.Env = append(os.Environ(),
			"GOOS="+p.OS,
			"GOARCH="+p.Arch,
			"GOARM="+p.Arm,
			"CGO_ENABLED=0",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error building %s: %v\n", outputName, err)
			os.Exit(1)
		}
	}

	fmt.Println("Build complete!")
}

func getVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	out, err := cmd.Output()
	if err != nil {
		return "dev"
	}
	return string(out)
}
```

## Platform-Specific Considerations

### macOS Code Signing

macOS binaries may need signing:

```bash
# Build
GOOS=darwin GOARCH=arm64 go build -o myapp

# Sign (on macOS)
codesign -s "Developer ID Application: Your Name" myapp

# Verify
codesign -v myapp
```

### Windows Manifests

Add Windows manifest for admin privileges:

```xml
<!-- app.manifest -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="requireAdministrator"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

Embed using windres or go-winres.

### Linux Static Binaries

For maximum Linux compatibility:

```bash
# Pure static binary (no CGO)
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o myapp

# With CGO (requires musl-gcc or static glibc)
CC=musl-gcc CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o myapp
```

## Docker-Based Cross-Compilation

### Dockerfile for Multi-Platform Builds

```dockerfile
# Multi-stage build for multiple platforms
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

# Build for multiple platforms
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/myapp-linux-amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/myapp-linux-arm64
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/myapp-windows-amd64.exe

# Extract artifacts
FROM scratch AS artifacts
COPY --from=builder /app/dist/* /
```

### Docker Buildx

```bash
# Build for multiple platforms
docker buildx build \
    --platform linux/amd64,linux/arm64,linux/arm/v7 \
    -t myapp:latest \
    .

# Build and push
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    -t myregistry/myapp:latest \
    --push \
    .
```

## CGO Cross-Compilation

### Installing Cross-Compilers

**Ubuntu/Debian:**
```bash
# ARM64
sudo apt-get install gcc-aarch64-linux-gnu

# ARM
sudo apt-get install gcc-arm-linux-gnueabihf

# Windows (MinGW)
sudo apt-get install gcc-mingw-w64
```

**macOS (via Homebrew):**
```bash
# ARM64
brew install arm-linux-gnueabihf-binutils

# Windows
brew install mingw-w64
```

### Cross-Compiling with CGO

**Linux to Windows:**
```bash
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
go build -o myapp.exe
```

**Linux to ARM64:**
```bash
CC=aarch64-linux-gnu-gcc \
CXX=aarch64-linux-gnu-g++ \
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm64 \
go build -o myapp-arm64
```

**Using Docker for CGO Cross-Compilation:**
```dockerfile
FROM golang:1.21

RUN apt-get update && apt-get install -y \
    gcc-aarch64-linux-gnu \
    g++-aarch64-linux-gnu \
    gcc-mingw-w64 \
    g++-mingw-w64

WORKDIR /app
COPY . .

# Build for ARM64 Linux
RUN CC=aarch64-linux-gnu-gcc \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64 \
    go build -o dist/myapp-linux-arm64

# Build for Windows AMD64
RUN CC=x86_64-w64-mingw32-gcc \
    CGO_ENABLED=1 \
    GOOS=windows \
    GOARCH=amd64 \
    go build -o dist/myapp-windows-amd64.exe
```

## Platform-Specific Build Tags

Use build tags for platform-specific code:

**File naming:**
```
file.go           # All platforms
file_linux.go     # Linux only
file_darwin.go    # macOS only
file_windows.go   # Windows only
file_unix.go      # Unix-like systems
file_amd64.go     # AMD64 only
file_arm64.go     # ARM64 only
```

**Build tag comments:**
```go
//go:build linux
// +build linux

package main

// Linux-specific code
```

```go
//go:build darwin && arm64
// +build darwin,arm64

package main

// macOS ARM64-specific code
```

## Testing Cross-Compiled Binaries

### QEMU User Mode

```bash
# Install QEMU
sudo apt-get install qemu-user-static

# Build for ARM64
GOOS=linux GOARCH=arm64 go build -o myapp-arm64

# Run with QEMU
qemu-aarch64-static myapp-arm64
```

### Docker Emulation

```bash
# Enable multiarch
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

# Run ARM64 container on AMD64
docker run --platform linux/arm64 alpine uname -m
```

## Common Cross-Compilation Patterns

### Build Matrix

```bash
# All Linux architectures
for arch in amd64 386 arm arm64 ppc64le s390x; do
    GOOS=linux GOARCH=$arch go build -o dist/myapp-linux-$arch
done

# All platforms for AMD64
for os in linux darwin windows freebsd openbsd netbsd; do
    ext=""
    [[ $os == "windows" ]] && ext=".exe"
    GOOS=$os GOARCH=amd64 go build -o dist/myapp-$os-amd64$ext
done
```

### Reproducible Builds

```bash
# Reproducible cross-platform builds
build() {
    local goos=$1
    local goarch=$2

    CGO_ENABLED=0 \
    GOOS=$goos \
    GOARCH=$goarch \
    go build \
        -trimpath \
        -ldflags="-s -w -X main.version=$VERSION -buildid=" \
        -o dist/myapp-$goos-$goarch
}

# Build
VERSION=$(git describe --tags)
build linux amd64
build darwin arm64
build windows amd64
```

### Parallel Builds

```bash
#!/bin/bash
# Parallel cross-compilation

platforms=(
    "linux:amd64"
    "linux:arm64"
    "darwin:amd64"
    "darwin:arm64"
    "windows:amd64"
)

build_platform() {
    local platform=$1
    local os=${platform%:*}
    local arch=${platform#*:}

    echo "Building $os/$arch..."
    GOOS=$os GOARCH=$arch go build -o dist/myapp-$os-$arch
}

export -f build_platform

# Build in parallel (4 at a time)
printf '%s\n' "${platforms[@]}" | xargs -P 4 -I {} bash -c 'build_platform "$@"' _ {}
```

## Troubleshooting

### Issue: Missing libc

**Problem:** Binary requires glibc on target system.

**Solution:**
```bash
# Build fully static
CGO_ENABLED=0 go build -ldflags '-extldflags "-static"'
```

### Issue: CGO cross-compilation fails

**Problem:** C compiler not configured correctly.

**Solution:**
```bash
# Verify cross-compiler
which aarch64-linux-gnu-gcc

# Test compilation
echo "int main() { return 0; }" | aarch64-linux-gnu-gcc -x c - -o test

# Use correct environment
export CC=aarch64-linux-gnu-gcc
export CXX=aarch64-linux-gnu-g++
```

### Issue: Binary not executable on target

**Problem:** Wrong GOOS/GOARCH or architecture mismatch.

**Solution:**
```bash
# Verify target platform
uname -s  # OS
uname -m  # Architecture

# Check binary format
file myapp

# Ensure correct variables
echo $GOOS $GOARCH
```

### Issue: Import cycle with platform-specific files

**Problem:** Build tags create circular dependencies.

**Solution:** Use `//go:build` with proper logic:
```go
//go:build !windows
// +build !windows

package main

// Unix implementation
```

```go
//go:build windows
// +build windows

package main

// Windows implementation
```
