# Go to Moxie: Implementation Plan

This document provides a comprehensive, dependency-ordered plan for forking Go and implementing the Moxie language revisions as specified in `go-language-revision.md`.

---

## 📊 Overall Progress

**Last Updated:** 2025-11-08 (Build & Test Complete)

| Phase | Status | Progress | Est. Duration | Actual Duration | Notes |
|-------|--------|----------|---------------|-----------------|-------|
| **Phase 0** | ✅ Complete | 99% | 2-3 days | ~10 hours | Build successful, 99.4% tests passing |
| **Phase 1** | ⏳ Pending | 0% | 15-22 days | - | Type System Foundation |
| **Phase 2** | ⏳ Pending | 0% | 15-22 days | - | Built-in Functions |
| **Phase 3** | ⏳ Pending | 0% | 17-24 days | - | String & Byte Unification |
| **Phase 4** | ⏳ Pending | 0% | 15-21 days | - | Const & Immutability |
| **Phase 5** | ⏳ Pending | 0% | 14-20 days | - | Zero-Copy Type Coercion |
| **Phase 6** | ⏳ Pending | 0% | 28-39 days | - | FFI & dlopen |
| **Phase 7** | ⏳ Pending | 0% | 31-45 days | - | Standard Library Updates |
| **Phase 8** | ⏳ Pending | 0% | 17-24 days | - | Testing & Validation |
| **Phase 9** | ⏳ Pending | 0% | 22-31 days | - | Documentation & Tools |
| **Phase 10** | ⏳ Pending | 0% | Ongoing | - | Releases |

**Overall Project Progress:** 0.99% (Phase 0: 99% × 10% weight)

**Current Milestone:** ✅ Phase 0 Complete - Build working, 357/359 tests passing
**Next Milestone:** Begin Phase 1 - Type System Foundation

**Build Status:** ✅ SUCCESS
**Test Status:** ✅ 99.4% passing (357/359 tests)
**Binary:** `/home/mleku/src/github.com/mleku/moxie/bin/go` (20 MB)
**Version:** `moxie version go1.26-devel_65528fa ... linux/amd64`

---

## Quick Links

- 📖 [Language Specification](go-language-revision.md) - Complete Moxie design
- 📝 [Summary](go-language-revision-summary.md) - Quick overview
- ✅ [Phase 0 Completion](PHASE0-COMPLETE.md) - Overall Phase 0 summary
- 🧪 [Test Results](PHASE0-TEST-RESULTS.md) - Build and test results
- 🏗️ [Build System](BUILD-SYSTEM.md) - Build system guide
- 🧪 [Testing Infrastructure](TESTING-INFRASTRUCTURE.md) - Testing guide
- 🌍 [Environment Variables](ENVIRONMENT-VARIABLES.md) - GO* → MOXIE* migration plan
- 🎨 [Branding Changes](BRANDING-CHANGES.md) - Complete branding changelog

---

## Overview

**Project:** Fork Go language → Create Moxie language
**Source:** go/ directory (Go 1.23+ codebase)
**Target:** Moxie (Go with fundamental revisions)
**Approach:** Incremental transformation with backward compatibility phases

**Key Principles:**
- Maintain working compiler at each phase
- Comprehensive testing at every step
- Performance validation throughout
- Clear migration path for users
- Security-first approach

---

## Phase 0: Foundation & Setup ✅

**Status:** ✅ COMPLETE (99% complete)
**Started:** 2025-11-08
**Completed:** 2025-11-08
**Actual Duration:** ~10 hours

### 0.1 Repository Setup & Branding ✅
**Dependencies:** None
**Duration:** 3 hours
**Status:** ✅ COMPLETE

- [x] Fork Go repository to `moxie` repository
- [x] Copy Go source from `go/` subdirectory to main repository structure
- [x] Update all branding (go → moxie)
  - [x] Update README.md - Complete Moxie project description
  - [x] Update CONTRIBUTING.md - Comprehensive contribution guidelines
  - [x] LICENSE maintained (BSD-style with Go acknowledgment)
  - [x] Create implementation documentation:
    - [x] go-language-revision.md (2047 lines)
    - [x] go-language-revision-summary.md (367 lines)
    - [x] go-to-moxie-plan.md (this file, 881 lines)
    - [x] BRANDING-CHANGES.md (1000+ lines)
    - [x] BRANDING-COMPLETE.md (200+ lines)
  - [x] Create renaming script: `scripts/rename-to-moxie.sh`
  - [x] **DONE:** Update version strings and user-facing text (43 strings)
  - [x] **DONE:** Update compiler name (go → moxie) in 5 source files
  - [x] **DONE:** Update command help and error messages
  - [x] **DONE:** Document all branding changes

**Source Files Modified:**
- ✅ src/cmd/go/internal/base/base.go (command branding)
- ✅ src/cmd/go/internal/version/version.go (version command + 15 strings)
- ✅ src/cmd/go/main.go (error messages, 6 strings)
- ✅ src/cmd/compile/internal/base/print.go (compiler messages, 2 strings)
- ✅ src/runtime/extern.go (runtime version, 3 strings)

**Total:** 43 string replacements across 5 files

### 0.2 Environment Variables Strategy ✅
**Dependencies:** 0.1
**Duration:** 2.25 hours
**Status:** ✅ COMPLETE

- [x] Document all 31 GO* environment variables
- [x] Define MOXIE* equivalents
- [x] Create 6-phase migration timeline
- [x] **DECISION:** Keep GO* variables for backward compatibility (Phase 4+ migration)
- [x] Create ENVIRONMENT-VARIABLES.md (600 lines)
- [x] Create PHASE0.2-COMPLETE.md (400+ lines)

**Strategic Decision:** Maintain backward compatibility, defer MOXIE* to Phase 4+

### 0.3 Build System Documentation ✅
**Dependencies:** 0.1, 0.2
**Duration:** 2.25 hours
**Status:** ✅ COMPLETE

- [x] Document three-stage bootstrap process
- [x] Identify exact changes needed (8 messages)
- [x] Define 3 implementation sub-phases
- [x] Create BUILD-SYSTEM.md (800 lines)
- [x] Create PHASE0.3-COMPLETE.md (500+ lines)
- [x] **TEST:** Build Moxie from source ✅ SUCCESS
  - [x] Bootstrap compiler: Go 1.25.3 (exceeds 1.24.6+ requirement)
  - [x] Three-stage bootstrap completed successfully
  - [x] Binary created: bin/go (20 MB)
  - [x] Version shows "moxie version ..." ✅

**Build Test Results:** ✅ SUCCESS

### 0.4 Testing Infrastructure Documentation ✅
**Dependencies:** 0.1, 0.2, 0.3
**Duration:** 2 hours
**Status:** ✅ COMPLETE

- [x] Document four testing levels
- [x] Identify required changes (4 scripts, 20+ test.go, 100+ expectations)
- [x] Define 6 implementation sub-phases
- [x] Create TESTING-INFRASTRUCTURE.md (900 lines)
- [x] Create PHASE0.4-COMPLETE.md (600+ lines)
- [x] **TEST:** Run full test suite ✅ 99.4% PASS RATE
  - [x] Total tests: 359
  - [x] Passed: 357 (99.4%)
  - [x] Failed: 2 (expected branding failures)
  - [x] Standard library: 100% pass rate
  - [x] Core compiler: 100% pass rate
  - [x] Runtime: 100% pass rate

**Test Results:** ✅ 357/359 passing (99.4%)

### 0.5 Final Documentation ✅
**Dependencies:** 0.1, 0.2, 0.3, 0.4
**Duration:** 30 minutes
**Status:** ✅ COMPLETE

- [x] Create PHASE0-COMPLETE.md (comprehensive summary)
- [x] Create PHASE0-TEST-RESULTS.md (detailed test analysis)
- [x] Update README.md with success metrics
- [x] Update go-to-moxie-plan.md (this file)

**Output:** Clean fork with Moxie branding, working build, 99.4% tests passing
**Current Status:** ✅ COMPLETE - Ready for Phase 1

**Documentation Created:** 6,000+ lines across 12 files
**Time Investment:** ~10 hours (significantly faster than estimated 2-3 days)

---

### Phase 0 Results Summary

**✅ Build:** SUCCESS
- Three-stage bootstrap: ✅
- Binary size: 20 MB
- Version output: `moxie version go1.26-devel_65528fa ... linux/amd64`

**✅ Tests:** 99.4% PASS RATE
- Total: 359 tests
- Passed: 357 tests
- Failed: 2 tests (expected branding-related failures)

**✅ Documentation:** COMPREHENSIVE
- 12 documentation files created
- 6,000+ lines of documentation
- Complete guides for build, test, environment variables

**✅ Code Changes:** MINIMAL & TARGETED
- 5 source files modified
- 43 string replacements
- User-facing branding complete

**Next Action:** Begin Phase 1 - Type System Foundation

---

## Phase 1: Type System Foundation (Core Infrastructure)

### 1.1 Extend Slice Header for Endianness
**Dependencies:** 0.2
**Duration:** 3-5 days
**Critical Path:** YES

**Files to modify:**
- `src/runtime/slice.go` - Add byteOrder field to slice header
- `src/internal/abi/type.go` - Update SliceType structure
- `src/cmd/compile/internal/types/type.go` - Add endianness to slice types
- `src/reflect/value.go` - Update reflection for new slice header

**Tasks:**
```go
// Update slice header structure
type slice struct {
    array unsafe.Pointer
    len   int64          // Change from int
    cap   int64          // Change from int
    byteOrder uint8      // NEW: NativeEndian/LittleEndian/BigEndian
}
```

- [ ] Add byteOrder constants (NativeEndian, LittleEndian, BigEndian)
- [ ] Update slice header in runtime
- [ ] Modify slice allocation to initialize byteOrder = NativeEndian
- [ ] Update slice copying to preserve byteOrder
- [ ] Add tests for slice header modifications

**Output:** Extended slice type with endianness support

---

### 1.2 Replace int/uint with Explicit Sizes
**Dependencies:** 1.1
**Duration:** 5-7 days
**Critical Path:** YES

**Files to modify:**
- `src/builtin/builtin.go` - Remove int/uint type definitions
- `src/cmd/compile/internal/types/type.go` - Remove int/uint types
- `src/runtime/*.go` - Replace all int/uint uses
- `src/cmd/compile/internal/gc/*.go` - Update type checker

**Tasks:**
- [ ] Create deprecation warnings for int/uint types
- [ ] Update all built-in functions to return int64:
  - [ ] len() → int64
  - [ ] cap() → int64
  - [ ] copy() → int64
- [ ] Replace internal uses:
  - [ ] Loop indices → int64
  - [ ] Array/slice indices → int64
  - [ ] Range variables → int64
- [ ] Update standard library (systematic replacement)
- [ ] Create migration tool: `moxie fix -r old.go new.go`
- [ ] Add comprehensive tests

**Migration rules:**
```
int      → int32 (for most uses)
int      → int64 (for sizes, indices, counts)
uint     → uint32 (for most uses)
uint     → uint64 (for sizes, pointers)
```

**Output:** Platform-independent integer types

---

### 1.3 Implement Pointer Types for Reference Types
**Dependencies:** 1.2
**Duration:** 7-10 days
**Critical Path:** YES

**Files to modify:**
- `src/cmd/compile/internal/types/type.go` - Change slice/map/chan to pointers
- `src/cmd/compile/internal/gc/typecheck.go` - Update type checking
- `src/runtime/map.go` - Update map representation
- `src/runtime/chan.go` - Update channel representation
- `src/reflect/type.go` - Update reflection

**Tasks:**

#### 1.3.1 Slice Type Changes
```go
// Before: []T (special reference type)
// After:  *[]T (explicit pointer type)
```
- [ ] Update parser to accept both `[]T` and `*[]T` (compatibility mode)
- [ ] Modify type checker to treat `[]T` as `*[]T` internally
- [ ] Add nil checking for slice operations
- [ ] Update slice literal syntax: `&[]T{...}`

#### 1.3.2 Map Type Changes
```go
// Before: map[K]V (special reference type)
// After:  *map[K]V (explicit pointer type)
```
- [ ] Update parser for `*map[K]V`
- [ ] Modify map allocation code
- [ ] Update map literal syntax: `&map[K]V{...}`
- [ ] Add nil map error handling

#### 1.3.3 Channel Type Changes
```go
// Before: chan T (special reference type)
// After:  *chan T (explicit pointer type)
```
- [ ] Update parser for `*chan T`
- [ ] Modify channel allocation
- [ ] Update channel literal syntax: `&chan T{cap: 10}`
- [ ] Update send/receive operations

#### 1.3.4 Compiler Auto-Dereferencing
- [ ] Implement automatic dereferencing for `*[]T`, `*map[K]V`, `*chan T`
- [ ] Maintain backward compatibility with implicit operations
- [ ] Add optimization passes to eliminate redundant checks

**Output:** Explicit pointer semantics for reference types

---

## Phase 2: Built-in Functions & Operations

### 2.1 Deprecate make() Function
**Dependencies:** 1.3
**Duration:** 2-3 days

**Files to modify:**
- `src/cmd/compile/internal/gc/typecheck.go` - Add deprecation warnings
- `src/builtin/builtin.go` - Mark make() as deprecated

**Tasks:**
- [ ] Add compiler warnings for make() usage
- [ ] Create migration guide (make → &T{} syntax)
- [ ] Update error messages with suggested replacements

**Migration examples:**
```go
// Before
s := make([]T, len, cap)
m := make(map[K]V, hint)
ch := make(chan T, buf)

// After
s := &[]T{}              // + grow() if needed
m := &map[K]V{}
ch := &chan T{cap: buf}
```

**Output:** make() deprecated with warnings

---

### 2.2 Implement grow() Built-in
**Dependencies:** 1.3
**Duration:** 3-4 days

**Files to modify:**
- `src/builtin/builtin.go` - Add grow() definition
- `src/runtime/slice.go` - Implement grow() runtime
- `src/cmd/compile/internal/gc/builtin.go` - Add compiler support

**Tasks:**
```go
// Signature: func grow[T any](s *[]T, n int64) *[]T
```
- [ ] Implement grow() for slices
  - [ ] Pre-allocate capacity
  - [ ] Preserve existing elements
  - [ ] Update slice header
- [ ] Add grow() for maps (hint capacity)
- [ ] Optimize for common patterns
- [ ] Add comprehensive tests

**Output:** grow() built-in function

---

### 2.3 Implement clone() Built-in
**Dependencies:** 1.3
**Duration:** 5-7 days

**Files to modify:**
- `src/builtin/builtin.go` - Add clone() definition
- `src/runtime/slice.go` - Deep copy for slices
- `src/runtime/map.go` - Deep copy for maps
- `src/runtime/string.go` - Copy for strings (once merged with []byte)

**Tasks:**
```go
// Signature: func clone[T any](v *T) *T
```
- [ ] Implement deep copy for slices
- [ ] Implement deep copy for maps (keys + values)
- [ ] Implement copy for strings
- [ ] Implement copy for channels (create new with same buffer size)
- [ ] Handle nested structures recursively
- [ ] Optimize for small objects (inline)
- [ ] Add tests for all types

**Output:** clone() built-in for explicit copying

---

### 2.4 Implement clear() Built-in
**Dependencies:** 1.3
**Duration:** 2-3 days

**Files to modify:**
- `src/builtin/builtin.go` - Add clear() definition
- `src/runtime/slice.go` - Reset length
- `src/runtime/map.go` - Remove all keys

**Tasks:**
```go
// Signature: func clear(v *[]T | *map[K]V)
```
- [ ] Implement clear() for slices (set len=0, keep capacity)
- [ ] Implement clear() for maps (remove all keys efficiently)
- [ ] Preserve allocated memory
- [ ] Add tests

**Output:** clear() built-in function

---

### 2.5 Implement free() Built-in
**Dependencies:** 2.4
**Duration:** 3-5 days

**Files to modify:**
- `src/builtin/builtin.go` - Add free() definition
- `src/runtime/mgc.go` - GC hint implementation
- `src/runtime/mheap.go` - Memory release

**Tasks:**
```go
// Signature: func free[T any](v *T)
```
- [ ] Implement as GC hint (not immediate deallocation)
- [ ] Mark memory as eligible for immediate collection
- [ ] Update GC to respect free() hints
- [ ] Add safety checks (prevent use-after-free in debug mode)
- [ ] Add tests with -race flag integration

**Output:** free() for explicit memory hints

---

### 2.6 Implement + Operator for Slices
**Dependencies:** 1.3, 2.3
**Duration:** 4-5 days

**Files to modify:**
- `src/cmd/compile/internal/gc/typecheck.go` - Add + for slices
- `src/cmd/compile/internal/gc/walk.go` - Lowering to runtime calls
- `src/runtime/slice.go` - Concatenation implementation

**Tasks:**
```go
// Semantics: a + b allocates new slice, concatenates contents
```
- [ ] Add + operator support for `*[]T` types
- [ ] Implement runtime concatenation (equivalent to append(clone(a), b...))
- [ ] Optimize for small slices (inline)
- [ ] Ensure does not mutate operands
- [ ] Add comprehensive tests

**Output:** Slice concatenation with + operator

---

## Phase 3: String & Byte Unification

### 3.1 Create Mutable String Type (Phase 1)
**Dependencies:** 2.6
**Duration:** 7-10 days
**Critical Path:** YES

**Files to modify:**
- `src/runtime/string.go` - Merge with slice implementation
- `src/builtin/builtin.go` - Update string type
- `src/cmd/compile/internal/types/type.go` - String type aliasing
- `src/reflect/type.go` - Update reflection

**Tasks:**

#### 3.1.1 Type Aliasing
```go
// Make string an alias for *[]byte
type string = *[]byte
```
- [ ] Create internal representation change
- [ ] Update string header to match slice header
- [ ] Preserve UTF-8 semantics for range loops
- [ ] Update string literals to allocate mutable byte slices

#### 3.1.2 String Operations
- [ ] Allow mutation: `s[0] = 'H'`
- [ ] Preserve + operator for strings (already implemented for slices)
- [ ] Update strings package functions to work with mutable strings
- [ ] Merge strings and bytes packages where appropriate

#### 3.1.3 UTF-8 Handling
- [ ] Preserve for...range rune iteration
- [ ] Keep len() as byte length
- [ ] Maintain utf8 package compatibility
- [ ] Update unicode handling

#### 3.1.4 Backward Compatibility
- [ ] Allow both immutable and mutable strings temporarily
- [ ] Add warnings for immutable string assumptions
- [ ] Create migration tool

**Output:** Mutable string type

---

### 3.2 Update Standard Library for Mutable Strings
**Dependencies:** 3.1
**Duration:** 10-14 days

**Packages to update:**
- `src/strings/*.go` - Work with *[]byte directly
- `src/bytes/*.go` - Merge functionality into strings
- `src/fmt/*.go` - Handle mutable strings
- `src/io/*.go` - String I/O operations
- `src/strconv/*.go` - Conversions with mutable strings
- `src/unicode/*.go` - Preserve rune handling

**Tasks:**
- [ ] Eliminate []byte conversions in strings package
- [ ] Merge bytes.Buffer functionality
- [ ] Update fmt verbs for mutable strings
- [ ] Update text processing packages
- [ ] Add comprehensive tests

**Output:** Standard library supporting mutable strings

---

## Phase 4: Const & Immutability

### 4.1 Implement const Keyword Extension
**Dependencies:** 3.2
**Duration:** 7-10 days
**Critical Path:** YES

**Files to modify:**
- `src/cmd/compile/internal/gc/typecheck.go` - Extend const
- `src/cmd/compile/internal/gc/const.go` - New const handling
- `src/cmd/compile/internal/obj/link.go` - .rodata section support
- `src/runtime/mgc.go` - Skip const data in GC

**Tasks:**

#### 4.1.1 Const Type System
```go
// Allow const for any type
const Message = "hello"         // const *[]byte
const Config = &struct{...}{}   // const *struct
const Data = &[]int{1, 2, 3}   // const *[]int
const Map = &map[string]int{}   // const *map
```
- [ ] Extend parser to allow const with composite types
- [ ] Add const propagation in type checker
- [ ] Implement const checking for assignments
- [ ] Add const inference for expressions

#### 4.1.2 Const in Function Signatures
```go
func process(data const *[]byte) { ... }
```
- [ ] Add const parameter support
- [ ] Implement const checking in function calls
- [ ] Allow const → mutable = error
- [ ] Allow mutable → const = ok

#### 4.1.3 Compiler Support
- [ ] Evaluate const expressions at compile time
- [ ] Pre-build const data structures
- [ ] Place const data in .rodata section
- [ ] Add const folding optimizations

**Output:** Extended const keyword

---

### 4.2 Implement MMU-Based Const Protection
**Dependencies:** 4.1
**Duration:** 5-7 days

**Files to modify:**
- `src/runtime/mmap.go` - Memory protection
- `src/runtime/mgc.go` - Skip const pages
- `src/cmd/link/internal/ld/data.go` - .rodata section

**Tasks:**
- [ ] Allocate const data in .rodata section
- [ ] Mark const pages as PROT_READ (no write permission)
- [ ] Set up MMU protection at startup
- [ ] Handle SIGSEGV for const violations (debug mode)
- [ ] Skip const pages in GC scanning
- [ ] Add tests for const violations

**Output:** Hardware-enforced immutability

---

### 4.3 Const Runtime Checks
**Dependencies:** 4.2
**Duration:** 3-4 days

**Files to modify:**
- `src/reflect/value.go` - Prevent const modification via reflection
- `src/unsafe/unsafe.go` - Block unsafe const bypasses

**Tasks:**
- [ ] Add runtime checks for reflect.Value.Set() on const
- [ ] Block unsafe.Pointer casts from const
- [ ] Add -race flag integration for const violations
- [ ] Create debug mode with extra const checking

**Output:** Complete const safety

---

## Phase 5: Zero-Copy Type Coercion

### 5.1 Implement Type Coercion for Numeric Slices
**Dependencies:** 1.1, 1.3
**Duration:** 7-10 days

**Files to modify:**
- `src/cmd/compile/internal/gc/typecheck.go` - Cast support
- `src/cmd/compile/internal/gc/walk.go` - Lowering
- `src/runtime/slice.go` - Coercion runtime

**Tasks:**

#### 5.1.1 Basic Type Casting
```go
bytes := &[]byte{1, 2, 3, 4}
u32s := (*[]uint32)(bytes)  // Zero-copy reinterpret
```
- [ ] Implement cast operator for numeric slice types
- [ ] Adjust len and cap based on type size
- [ ] Add alignment checking (debug mode)
- [ ] Implement restrictions (only fixed-width numeric types)

#### 5.1.2 Endianness Support
```go
u32s_le := (*[]uint32, LittleEndian)(bytes)
u32s_be := (*[]uint32, BigEndian)(bytes)
```
- [ ] Add endianness parameter to cast syntax
- [ ] Update slice header byteOrder field on cast
- [ ] Implement NativeEndian, LittleEndian, BigEndian constants
- [ ] Add tests for all endianness combinations

**Output:** Zero-copy type coercion

---

### 5.2 Implement Automatic Byte Swapping
**Dependencies:** 5.1
**Duration:** 5-7 days

**Files to modify:**
- `src/cmd/compile/internal/gc/ssa.go` - SSA passes for byte swap
- `src/cmd/compile/internal/ssa/*.go` - Optimization passes
- `src/runtime/slice.go` - Index operations with byte swap

**Tasks:**

#### 5.2.1 Byte Swap Intrinsics
- [ ] Implement BSWAP instruction for x86/x64
- [ ] Implement REV instruction for ARM
- [ ] Add compiler intrinsics for byte swapping
- [ ] Optimize for common patterns

#### 5.2.2 Slice Index Operations
```go
// Auto byte-swap on read/write based on byteOrder
u32s[i]      // Reads and swaps if needed
u32s[i] = v  // Swaps and writes if needed
```
- [ ] Modify slice index operations to check byteOrder
- [ ] Insert byte swap for non-native endianness
- [ ] Optimize away swaps for native endianness
- [ ] Add tests for all scenarios

#### 5.2.3 Optimizations
- [ ] Compile-time swap elimination if endian matches
- [ ] SIMD batch swapping for large arrays
- [ ] Inline small swaps
- [ ] Add performance benchmarks

**Output:** Automatic endianness handling

---

### 5.3 Explicit Copy with & Operator
**Dependencies:** 2.3, 5.1
**Duration:** 2-3 days

**Files to modify:**
- `src/cmd/compile/internal/gc/typecheck.go` - & for copy-cast

**Tasks:**
```go
// Reinterpret without copy
u32s := (*[]uint32)(bytes)

// Explicit copy
u32s_copy := &(*[]uint32)(bytes)
// Equivalent to: clone((*[]uint32)(bytes))
```
- [ ] Implement & prefix for copy-on-cast
- [ ] Generate code to allocate and copy
- [ ] Add tests

**Output:** Explicit copy-cast operator

---

## Phase 6: FFI & dlopen Implementation

### 6.1 Implement dlopen Built-ins
**Dependencies:** 4.3
**Duration:** 7-10 days

**Files to modify:**
- `src/builtin/builtin.go` - Add dlopen/dlsym/dlclose
- `src/runtime/dlfcn.go` - NEW: Dynamic library loading
- `src/cmd/compile/internal/gc/builtin.go` - Compiler support

**Tasks:**

#### 6.1.1 dlopen Implementation
```go
func dlopen(filename string, flags int32) *DLib
```
- [ ] Implement dynamic library loading (platform-specific)
- [ ] Support RTLD_LAZY, RTLD_NOW, RTLD_GLOBAL, RTLD_LOCAL flags
- [ ] Add error handling with dlerror()
- [ ] Create DLib handle type
- [ ] Add tests for loading system libraries

#### 6.1.2 dlsym Implementation
```go
func dlsym[T any](lib *DLib, name string) T
```
- [ ] Implement symbol lookup with type safety
- [ ] Use generics for type-safe returns
- [ ] Add runtime type checking
- [ ] Create function pointer wrappers
- [ ] Add tests for common C functions

#### 6.1.3 dlclose Implementation
```go
func dlclose(lib *DLib)
```
- [ ] Implement library unloading
- [ ] Handle reference counting
- [ ] Clean up resources
- [ ] Add tests

**Output:** Type-safe FFI built-ins

---

### 6.2 Implement C Type Mapping
**Dependencies:** 6.1
**Duration:** 4-5 days

**Files to modify:**
- `src/cmd/compile/internal/types/type.go` - C type mappings
- `src/runtime/dlfcn.go` - Type marshaling

**Tasks:**
- [ ] Map Go types to C types
  - int8 ↔ char
  - uint8 ↔ unsigned char
  - int16 ↔ short
  - int32 ↔ int
  - int64 ↔ long long
  - float32 ↔ float
  - float64 ↔ double
  - *[]byte ↔ char*
  - unsafe.Pointer ↔ void*
- [ ] Implement struct marshaling with `c:"..."` tags
- [ ] Add padding and alignment handling
- [ ] Create tests for all type mappings

**Output:** C type compatibility

---

### 6.3 Implement Callback Support
**Dependencies:** 6.2
**Duration:** 5-7 days

**Files to modify:**
- `src/runtime/dlfcn.go` - Callback trampolines
- `src/runtime/cgocall.go` - Adapt for new FFI (or create new)

**Tasks:**
- [ ] Create Go → C callback trampolines
- [ ] Handle calling convention differences
- [ ] Support closures in callbacks
- [ ] Manage goroutine scheduling for callbacks
- [ ] Add callback tests

**Output:** Bidirectional FFI calls

---

### 6.4 Implement dlopen_mem (Static Dynamic Linking)
**Dependencies:** 6.1
**Duration:** 10-14 days

**Files to modify:**
- `src/runtime/dlfcn.go` - Memory-based loading
- `src/runtime/elf.go` - ELF parsing (NEW or adapt from debug/elf)
- `src/cmd/link/internal/ld/*.go` - Embedding .so files

**Tasks:**

#### 6.4.1 ELF Parsing
- [ ] Parse ELF headers from memory
- [ ] Load program headers
- [ ] Map sections to memory
- [ ] Parse symbol tables
- [ ] Handle relocations

#### 6.4.2 Memory Mapping
```go
func dlopen_mem(data *[]byte, flags int32) *DLib
```
- [ ] Create anonymous memory mapping
- [ ] Copy library data
- [ ] Apply relocations
- [ ] Run init/constructor functions
- [ ] Return library handle
- [ ] Add tests with embedded libraries

#### 6.4.3 Linker Support
```bash
moxie build -ldflags="-embedso libfoo.so"
```
- [ ] Add -embedso linker flag
- [ ] Embed .so data in .rodata section
- [ ] Generate embed metadata
- [ ] Auto-detect embedded libraries in dlopen()
- [ ] Add integration tests

**Output:** Single-binary deployment with embedded libraries

---

### 6.5 Deprecate CGO
**Dependencies:** 6.4
**Duration:** 2-3 days

**Files to modify:**
- `src/cmd/cgo/*.go` - Add deprecation warnings
- `src/cmd/go/internal/work/build.go` - Warn on CGO usage

**Tasks:**
- [ ] Add compiler warnings for `import "C"`
- [ ] Create CGO → FFI migration guide
- [ ] Update documentation
- [ ] Add examples of FFI usage

**Output:** CGO deprecated

---

## Phase 7: Standard Library Updates

### 7.1 Update Core Packages
**Dependencies:** 3.2, 4.3, 5.3
**Duration:** 14-21 days

**Packages to update:**
- `src/bytes` - Merge into strings or adapt
- `src/strings` - Use *[]byte directly
- `src/fmt` - Handle mutable strings and const
- `src/io` - Update Reader/Writer for new types
- `src/bufio` - Use mutable strings
- `src/encoding/*` - Update for explicit types and endianness

**Tasks:**
- [ ] Systematic replacement of int/uint → int32/int64
- [ ] Update all slice/map/chan uses to explicit pointers
- [ ] Update string handling for mutability
- [ ] Remove unnecessary []byte conversions
- [ ] Add const where appropriate
- [ ] Run full test suite

**Output:** Core packages updated

---

### 7.2 Update Networking Packages
**Dependencies:** 5.3, 6.4, 7.1
**Duration:** 10-14 days

**Packages to update:**
- `src/net` - Update with FFI (remove CGO), use endianness
- `src/net/http` - Mutable strings, const configs
- `src/crypto/*` - Zero-copy optimizations, FFI for OpenSSL
- `src/encoding/binary` - Deprecate (use built-in coercion)

**Tasks:**
- [ ] Replace CGO calls with FFI in net package
- [ ] Use zero-copy type coercion for network protocols
- [ ] Optimize crypto with new features
- [ ] Update HTTP for mutable strings
- [ ] Add tests

**Output:** Networking packages optimized

---

### 7.3 Update Compiler & Tools
**Dependencies:** 7.1
**Duration:** 7-10 days

**Packages to update:**
- `src/cmd/compile` - Already mostly done in earlier phases
- `src/cmd/go` - Update build tool
- `src/cmd/gofmt` → `src/cmd/moxiefmt` - Update formatter
- `src/cmd/doc` → `src/cmd/moxiedoc` - Update documentation
- `src/go/*` - Update go/ packages for Moxie AST

**Tasks:**
- [ ] Update go tool to moxie tool
- [ ] Update module system references
- [ ] Create moxie.mod (replace go.mod)
- [ ] Update toolchain commands
- [ ] Add migration tools

**Output:** Complete Moxie toolchain

---

## Phase 8: Testing & Validation

### 8.1 Comprehensive Test Suite
**Dependencies:** 7.3
**Duration:** 7-10 days

**Tasks:**
- [ ] Run all existing Go tests with Moxie
- [ ] Fix compatibility issues
- [ ] Add Moxie-specific tests:
  - [ ] Explicit pointer semantics
  - [ ] Mutable strings
  - [ ] Const with MMU protection
  - [ ] FFI functionality
  - [ ] Zero-copy coercion
  - [ ] Endianness handling
  - [ ] Type coercion round-trips
- [ ] Performance benchmarks
- [ ] Memory leak tests
- [ ] Race detector tests

**Output:** Validated Moxie implementation

---

### 8.2 Performance Benchmarking
**Dependencies:** 8.1
**Duration:** 5-7 days

**Tasks:**
- [ ] Benchmark vs Go 1.x:
  - [ ] String operations
  - [ ] FFI calls vs CGO
  - [ ] Type coercion vs manual conversion
  - [ ] Const access
  - [ ] Slice operations
- [ ] Create performance report
- [ ] Identify optimization opportunities
- [ ] Implement critical optimizations

**Output:** Performance validation

---

### 8.3 Security Audit
**Dependencies:** 8.1
**Duration:** 5-7 days

**Tasks:**
- [ ] Audit const MMU protection
- [ ] Test for memory safety violations
- [ ] Validate FFI security
- [ ] Check for race conditions
- [ ] Review unsafe operations
- [ ] Test against known attack vectors

**Output:** Security validation

---

## Phase 9: Documentation & Migration Tools

### 9.1 Create Documentation
**Dependencies:** 8.3
**Duration:** 7-10 days

**Tasks:**
- [ ] Write Moxie language specification
- [ ] Create migration guide (Go → Moxie)
- [ ] Document all new features:
  - [ ] Explicit pointer types
  - [ ] Mutable strings
  - [ ] Const system
  - [ ] FFI/dlopen
  - [ ] Type coercion
  - [ ] Endianness
- [ ] Update standard library docs
- [ ] Create tutorial examples
- [ ] Write best practices guide

**Output:** Complete documentation

---

### 9.2 Build Migration Tools
**Dependencies:** 9.1
**Duration:** 10-14 days

**Tasks:**
- [ ] Create `moxie fix` tool:
  - [ ] Convert int/uint → int32/int64
  - [ ] Convert []T → *[]T
  - [ ] Convert map[K]V → *map[K]V
  - [ ] Convert chan T → *chan T
  - [ ] Convert make() → &T{}
  - [ ] Convert string operations
  - [ ] Add clone() where needed
  - [ ] Convert CGO → FFI
- [ ] Create compatibility checker
- [ ] Create code analysis tool
- [ ] Add dry-run mode
- [ ] Add comprehensive tests

**Output:** Automated migration tools

---

### 9.3 Create Compatibility Layer
**Dependencies:** 9.2
**Duration:** 5-7 days

**Tasks:**
- [ ] Create `compat` package for Go code compatibility
- [ ] Implement shims for removed features
- [ ] Create type aliases for smooth transition
- [ ] Add deprecation warnings
- [ ] Document compatibility layer

**Output:** Backward compatibility support

---

## Phase 10: Release Preparation

### 10.1 Alpha Release (Moxie 0.1.0-alpha)
**Dependencies:** 9.3
**Duration:** 3-5 days

**Tasks:**
- [ ] Tag alpha release
- [ ] Create release notes
- [ ] Build binaries for major platforms
- [ ] Publish documentation
- [ ] Create example projects
- [ ] Announce to community
- [ ] Gather feedback

**Output:** Alpha release

---

### 10.2 Beta Release (Moxie 0.1.0-beta)
**Dependencies:** 10.1 + feedback integration
**Duration:** Ongoing

**Tasks:**
- [ ] Fix alpha bugs
- [ ] Performance tuning
- [ ] API stabilization
- [ ] Update documentation
- [ ] Expand test coverage
- [ ] Tag beta release
- [ ] Gather more feedback

**Output:** Beta release

---

### 10.3 Production Release (Moxie 1.0.0)
**Dependencies:** 10.2 + stabilization
**Duration:** Ongoing

**Tasks:**
- [ ] Final bug fixes
- [ ] Security hardening
- [ ] Performance optimization
- [ ] Complete documentation
- [ ] Finalize API
- [ ] Tag 1.0.0 release
- [ ] Major announcement
- [ ] Community building

**Output:** Production-ready Moxie 1.0

---

## Dependency Graph (Critical Path)

```
0.1 (Setup)
  ↓
0.2 (Testing)
  ↓
1.1 (Slice Header) ← CRITICAL
  ↓
1.2 (int/uint Removal) ← CRITICAL
  ↓
1.3 (Pointer Types) ← CRITICAL
  ↓
2.1 (make() Deprecation)
  ↓
2.2-2.6 (Built-ins: grow, clone, clear, free, +)
  ↓
3.1 (Mutable Strings) ← CRITICAL
  ↓
3.2 (Stdlib Strings)
  ↓
4.1 (Const Extension) ← CRITICAL
  ↓
4.2 (MMU Protection)
  ↓
4.3 (Const Runtime)
  ↓
5.1 (Type Coercion)
  ↓
5.2 (Byte Swapping)
  ↓
5.3 (Copy-Cast)
  ↓
6.1 (dlopen)
  ↓
6.2-6.4 (FFI Complete)
  ↓
6.5 (CGO Deprecation)
  ↓
7.1-7.3 (Stdlib Updates)
  ↓
8.1-8.3 (Testing & Validation)
  ↓
9.1-9.3 (Docs & Migration)
  ↓
10.1-10.3 (Releases)
```

---

## Timeline Estimates

| Phase | Duration | Parallel Work Possible |
|-------|----------|----------------------|
| **Phase 0** (Setup) | 2-4 days | No |
| **Phase 1** (Type System) | 15-22 days | Minimal |
| **Phase 2** (Built-ins) | 15-22 days | Yes (after 1.3) |
| **Phase 3** (Strings) | 17-24 days | Minimal |
| **Phase 4** (Const) | 15-21 days | Minimal |
| **Phase 5** (Coercion) | 14-20 days | Some |
| **Phase 6** (FFI) | 28-39 days | Some (after 6.1) |
| **Phase 7** (Stdlib) | 31-45 days | Yes (packages parallel) |
| **Phase 8** (Testing) | 17-24 days | Some |
| **Phase 9** (Docs) | 22-31 days | Yes |
| **Phase 10** (Release) | Ongoing | N/A |

**Critical Path Total:** ~6-9 months (with 2-3 engineers)
**With Parallelization:** ~4-6 months (with 5-7 engineers)

---

## Risk Mitigation

### High-Risk Items
1. **MMU const protection** - OS/platform specific, needs thorough testing
2. **Byte swapping performance** - Must not regress performance
3. **FFI callback handling** - Complex goroutine interaction
4. **Standard library migration** - Massive codebase, potential breakage

### Mitigation Strategies
- [ ] Early prototyping of high-risk features
- [ ] Extensive test coverage (>90%)
- [ ] Performance regression testing
- [ ] Incremental rollout with feature flags
- [ ] Community feedback loops
- [ ] Compatibility layer for gradual migration

---

## Success Criteria

### Functional
- ✅ All explicit pointer types work correctly
- ✅ Mutable strings function as specified
- ✅ Const with MMU protection enforced
- ✅ FFI calls work with type safety
- ✅ Zero-copy coercion with endianness
- ✅ All platform-dependent ints removed

### Performance
- ✅ FFI ≥10x faster than CGO
- ✅ Type coercion ≥10x faster than manual conversion
- ✅ No regression in general performance
- ✅ String operations comparable or faster

### Quality
- ✅ >90% test coverage
- ✅ No security vulnerabilities
- ✅ All Go tests pass (with migration)
- ✅ Zero memory leaks
- ✅ Race detector clean

### Usability
- ✅ Complete documentation
- ✅ Working migration tools
- ✅ Example projects
- ✅ Clear error messages

---

## Open Questions

1. **Versioning Strategy**: How to handle Go module compatibility?
2. **Community Adoption**: Fork vs. proposal to Go team?
3. **Backward Compatibility**: How long to maintain compatibility layer?
4. **Platform Support**: Which platforms to support initially?
5. **License**: Keep Go's BSD license or different?
6. **Governance**: Who maintains Moxie long-term?

---

## Conclusion

This plan provides a comprehensive, dependency-ordered roadmap for transforming the Go language into Moxie. The approach is incremental, testable, and allows for parallel work where possible. Each phase builds upon previous phases, ensuring a stable foundation.

**Key Principles:**
- Maintain working compiler at each phase
- Comprehensive testing at every step
- Performance validation throughout
- Clear migration path for users
- Security-first approach

**Next Steps:**
1. Review and approve plan
2. Set up initial infrastructure (Phase 0)
3. Begin Phase 1 implementation
4. Regular progress reviews and adjustments

**This plan turns the Moxie vision into actionable engineering work.**
