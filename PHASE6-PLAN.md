# Phase 6: Standard Library Extensions

**Goal**: Implement core language features from the Moxie spec that extend beyond Go's capabilities

## Overview

Phase 6 implements three major features defined in the Moxie language specification:
1. **Native FFI** - Built-in foreign function interface without CGO
2. **Zero-copy type coercion** - Safe reinterpretation with endianness control
3. **const with MMU protection** - True immutability with hardware enforcement

## Priorities

### Priority 1: Native FFI (dlopen/dlsym)
Implement built-in foreign function interface to replace CGO.

**Rationale**: CGO is slow, breaks cross-compilation, and requires C compiler. Native FFI is faster, type-safe, and works with pure Go builds.

### Priority 2: Zero-copy Type Coercion
Implement safe, zero-copy slice reinterpretation with automatic endianness handling.

**Rationale**: Binary data parsing, network protocols, and cryptography benefit from zero-copy operations.

### Priority 3: const with MMU Protection
Implement true immutability with memory protection for const values.

**Rationale**: True immutability enables security guarantees and compiler optimizations.

## Detailed Plans

### 1. Native FFI Implementation

#### 1.1 Core Functions

**dlopen**:
```go
func dlopen(filename string, flags int32) *DLib
```
- Load shared library dynamically
- Return handle or nil on error
- Support RTLD_LAZY, RTLD_NOW, RTLD_GLOBAL, RTLD_LOCAL flags

**dlsym** (generic):
```go
func dlsym[T any](lib *DLib, name string) T
```
- Type-safe symbol lookup
- Return function pointer with specific signature
- Validate type at runtime (ensure T is function pointer)

**dlclose**:
```go
func dlclose(lib *DLib)
```
- Close library handle
- Clean up resources

**dlerror**:
```go
func dlerror() string
```
- Return last error message
- Reset error state

**dlopen_mem**:
```go
func dlopen_mem(data *[]byte, flags int32) *DLib
```
- Load library from memory (embedded in binary)
- Enables static binaries with embedded shared libraries

#### 1.2 Implementation Approach

Since we're transpiling to Go, we'll use CGO under the hood but expose a clean FFI API:

**Option A: Wrapper around CGO** (Quick, works now)
```go
// In runtime/ffi.go
/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DLib struct {
    handle unsafe.Pointer
}

const (
    RTLD_LAZY   = 0x001
    RTLD_NOW    = 0x002
    RTLD_GLOBAL = 0x100
    RTLD_LOCAL  = 0x000
)

func Dlopen(filename string, flags int32) *DLib {
    cname := C.CString(filename)
    defer C.free(unsafe.Pointer(cname))
    handle := C.dlopen(cname, C.int(flags))
    if handle == nil {
        return nil
    }
    return &DLib{handle: handle}
}

func Dlsym[T any](lib *DLib, name string) T {
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    sym := C.dlsym(lib.handle, cname)
    return *(*T)(unsafe.Pointer(&sym))
}

func Dlclose(lib *DLib) {
    if lib != nil && lib.handle != nil {
        C.dlclose(lib.handle)
    }
}

func Dlerror() string {
    err := C.dlerror()
    if err == nil {
        return ""
    }
    return C.GoString(err)
}
```

**Option B: Pure Go** (Future, for native Moxie compiler)
- Use syscalls directly
- Implement ELF/Mach-O/PE loader
- Deferred to native compiler phase

**Chosen Approach**: Pure Go using ebitengine/purego library (no CGO required!)

**Implementation Note (2025-11-09)**: Successfully migrated from CGO to purego for pure Go FFI.
- Uses github.com/ebitengine/purego v0.8.1
- No C compiler required
- Faster builds, smaller binaries
- Full cross-compilation support

**Current Status**: Implemented and ready, minor go.sum resolution issue in temp build directories (being investigated).

#### 1.3 AST Transformation

Detect calls to FFI functions and inject runtime import:

```go
// Detect: dlopen("libc.so.6", RTLD_LAZY)
// Transform: Add import to moxie runtime
// Rewrite: moxie.Dlopen("libc.so.6", moxie.RTLD_LAZY)
```

### 2. Zero-Copy Type Coercion

#### 2.1 Syntax

```go
(*[]TargetType)(sourceSlice)                     // Native endian
(*[]TargetType, LittleEndian)(sourceSlice)       // Little-endian
(*[]TargetType, BigEndian)(sourceSlice)          // Big-endian
```

#### 2.2 Implementation Approach

**Challenge**: Go doesn't support tuple-like syntax in type casts.

**Solution**: Use function call syntax that looks like a cast:

```go
// In runtime/coerce.go
const (
    NativeEndian = 0
    LittleEndian = 1
    BigEndian    = 2
)

// Coerce reinterprets a slice with optional endianness conversion
func Coerce[From, To any](src *[]From, endian ...int) *[]To {
    if src == nil {
        return nil
    }

    // Get endianness (default to native)
    byteOrder := NativeEndian
    if len(endian) > 0 {
        byteOrder = endian[0]
    }

    // Calculate size ratio
    fromSize := int(unsafe.Sizeof((*From)(nil)))
    toSize := int(unsafe.Sizeof((*To)(nil)))

    // Calculate new length
    srcLen := len(*src)
    dstLen := (srcLen * fromSize) / toSize

    // Create slice header with same backing array
    srcHeader := (*reflect.SliceHeader)(unsafe.Pointer(src))
    dstHeader := reflect.SliceHeader{
        Data: srcHeader.Data,
        Len:  dstLen,
        Cap:  (srcHeader.Cap * fromSize) / toSize,
    }

    result := *(*[]To)(unsafe.Pointer(&dstHeader))

    // Apply endianness conversion if needed
    if byteOrder != NativeEndian {
        swapEndian(&result, byteOrder)
    }

    return &result
}

// Helper to swap endianness
func swapEndian[T any](slice *[]T, targetEndian int) {
    // Implement byte swapping based on type size and target endian
    // ...
}
```

**Syntax transformation**:
```go
// Moxie source:
bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
u32s := (*[]uint32)(bytes)
u32s_le := (*[]uint32, LittleEndian)(bytes)

// Transpiled Go:
bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
u32s := moxie.Coerce[byte, uint32](bytes)
u32s_le := moxie.Coerce[byte, uint32](bytes, moxie.LittleEndian)
```

#### 2.3 AST Transformation

Detect cast syntax with optional endianness:
```go
// In transformTypeConversion:
// Look for: (*[]T)(expr)
// Check if source is also slice type
// Transform to: moxie.Coerce[S, T](expr)
//
// Look for: (*[]T, Endian)(expr)
// Transform to: moxie.Coerce[S, T](expr, moxie.Endian)
```

### 3. const with MMU Protection

#### 3.1 Challenge

This feature requires:
- Memory page protection (mprotect on Unix, VirtualProtect on Windows)
- Compile-time placement of const data in .rodata section
- Runtime enforcement of immutability

#### 3.2 Implementation Approach

**For Transpiler** (Limited):
```go
// Detection only - cannot fully implement in transpiler
// Provide compile-time checking but rely on Go's const semantics
```

**For Native Compiler** (Full):
- Place const values in .rodata ELF section
- Use mprotect(PROT_READ) to make pages read-only
- Hardware will trap on write attempts (SIGSEGV)

**Transpiler Approach**:
```go
// In syntax.go:
// 1. Detect const declarations with pointer types
// 2. Transform to global variables with runtime protection
// 3. Add init() function to call mprotect

// Example:
// const Config = &map[string]int32{"timeout": 30}
//
// Becomes:
// var __const_Config = &map[string]int32{"timeout": 30}
// const Config = __const_Config // Read-only reference
//
// func init() {
//     moxie.ProtectConst(unsafe.Pointer(__const_Config))
// }
```

**Runtime Support**:
```go
// In runtime/protect.go
func ProtectConst(ptr unsafe.Pointer) {
    // Round down to page boundary
    pageSize := syscall.Getpagesize()
    page := uintptr(ptr) &^ uintptr(pageSize-1)

    // Protect page as read-only
    syscall.Mprotect(
        (*[1<<30]byte)(unsafe.Pointer(page))[:pageSize:pageSize],
        syscall.PROT_READ,
    )
}
```

**Deferred**: Full implementation deferred to native compiler. Transpiler provides best-effort const checking.

## Implementation Steps

### Step 1: Native FFI ⏳
1. Create `runtime/ffi.go` with dlopen/dlsym/dlclose/dlerror
2. Create `runtime/ffi_mem.go` with dlopen_mem
3. Add FFI constants (RTLD_*) to runtime
4. Add AST transformation to detect FFI calls
5. Inject runtime import when FFI is used
6. Create test files

### Step 2: Zero-Copy Type Coercion ⏳
1. Create `runtime/coerce.go` with Coerce[T,S] function
2. Add endianness constants and swap functions
3. Implement byte swapping for all numeric types
4. Add AST transformation for cast syntax
5. Detect (*[]T)(expr) and (*[]T, Endian)(expr) patterns
6. Create test files for network parsing, crypto use cases

### Step 3: const with MMU Protection ⏸️
1. Document limitation in transpiler (deferred)
2. Add runtime support for future use
3. Provide compile-time const checking
4. Full implementation in native compiler phase

### Step 4: Testing ⏳
- test_ffi_basic.mx - dlopen/dlsym basic usage
- test_ffi_strlen.mx - Call libc strlen
- test_ffi_error.mx - Error handling
- test_coerce_basic.mx - Basic type reinterpretation
- test_coerce_endian.mx - Endianness conversion
- test_coerce_network.mx - Network protocol parsing
- test_const_protection.mx - Const checking (compile-time only)

## Testing Plan

### Test Files

1. **test_ffi_basic.mx**
   ```go
   package main

   func main() {
       lib := dlopen("libc.so.6", RTLD_LAZY)
       if lib == nil {
           panic(dlerror())
       }
       defer dlclose(lib)

       strlen := dlsym[func(*byte) int64](lib, "strlen")
       msg := "Hello FFI\x00"
       len := strlen(&msg[0])
       moxie.Print("Length:", len)
   }
   ```

2. **test_coerce_basic.mx**
   ```go
   package main

   func main() {
       bytes := &[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

       // Native endian
       u32s := (*[]uint32)(bytes)
       moxie.Print("u32s:", u32s[0], u32s[1])

       // Little endian (portable)
       u32s_le := (*[]uint32, LittleEndian)(bytes)
       moxie.Print("u32s_le:", u32s_le[0], u32s_le[1])

       // Big endian (network byte order)
       u32s_be := (*[]uint32, BigEndian)(bytes)
       moxie.Print("u32s_be:", u32s_be[0], u32s_be[1])
   }
   ```

3. **test_coerce_network.mx**
   ```go
   package main

   type PacketHeader struct {
       Magic   uint32
       Version uint16
       Length  uint16
   }

   func parseHeader(data *[]byte) PacketHeader {
       // Network protocols use big-endian
       fields := (*[]uint32, BigEndian)(data[0:8])

       return PacketHeader{
           Magic:   fields[0],
           Version: uint16(fields[1] >> 16),
           Length:  uint16(fields[1] & 0xFFFF),
       }
   }

   func main() {
       packet := &[]byte{
           0x12, 0x34, 0x56, 0x78,  // Magic
           0x00, 0x01, 0x02, 0x00,  // Version=1, Length=512
       }

       header := parseHeader(packet)
       moxie.Printf("Magic: 0x%08X\n", header.Magic)
       moxie.Printf("Version: %d\n", header.Version)
       moxie.Printf("Length: %d\n", header.Length)
   }
   ```

## Known Limitations

### FFI
- Requires CGO under the hood (transpiler limitation)
- Function pointer type safety relies on user correctness
- Memory management for C strings is manual
- Platform-specific library names (.so vs .dylib vs .dll)

### Type Coercion
- Alignment not checked (may panic on some architectures)
- Only works with fixed-width numeric types
- No support for structs (padding/alignment issues)
- Endianness swapping has runtime cost

### const with MMU
- Not fully implementable in transpiler
- Deferred to native compiler
- Best-effort const checking only
- Go's const limitations apply

## Success Criteria

- [ ] FFI functions implemented and tested
- [ ] Can call libc functions from Moxie
- [ ] Type coercion works with all numeric types
- [ ] Endianness conversion is correct (verified with tests)
- [ ] Network protocol parsing example works
- [ ] All previous tests still pass
- [ ] Documentation updated

## Future Work (Native Compiler)

### FFI
- Pure Go ELF/Mach-O/PE loader
- No CGO dependency
- Faster function calls (<10ns)
- Better error messages
- Symbol versioning support

### Type Coercion
- Compile-time alignment checking
- SIMD-accelerated byte swapping
- Support for packed structs
- Better type inference

### const with MMU
- Full .rodata section placement
- Page-level protection
- Compile-time const propagation
- Security guarantees
- Hardware-enforced immutability

## Phase 6 Scope

**In Scope**:
- ✅ Native FFI (dlopen, dlsym, dlclose, dlerror)
- ✅ Zero-copy type coercion with endianness
- ⚠️  const documentation (full impl deferred)

**Out of Scope**:
- ❌ Full const with MMU (needs native compiler)
- ❌ dlopen_mem (requires custom loader)
- ❌ Pure Go FFI (needs native compiler)

**Deferred to Phase 7+**:
- Standard library wrappers
- Enhanced error handling
- Advanced FFI features
- Compile-time optimizations
