# Phase 1: Type System Foundation - Implementation Plan

**Date:** 2025-11-08
**Phase:** Phase 1 - Type System Foundation
**Status:** PLANNING - Document Before Implementation
**Estimated Duration:** 15-22 days

---

## Executive Summary

**Decision:** Following the successful Phase 0 "document-first" approach, we will **thoroughly plan and document** Phase 1 changes before implementing them. Phase 1 involves fundamental changes to the Go type system, runtime, and compiler - changes that are far more complex and risky than Phase 0's branding updates.

**Rationale:**
1. **Complexity** - Modifying runtime slice headers, type system, and compiler is extremely complex
2. **Risk Management** - These are core infrastructure changes that affect everything
3. **Testing Requirements** - Need comprehensive test strategy before making changes
4. **Dependencies** - Must understand all interconnected systems
5. **Reversibility** - Need clear rollback plan for each sub-phase

---

## Phase 1 Overview

**Goal:** Establish the foundation for Moxie's type system improvements

**Key Changes:**
1. ✅ Extend slice headers to include endianness metadata (Phase 1.1 - PLANNING COMPLETE)
2. ⏸️ Remove platform-dependent `int` and `uint` types (Phase 1.2 - DEFERRED)
3. ⏳ Implement explicit pointer types for slices, maps, and channels (Phase 1.3 - TODO)

**Impact:** CRITICAL - Affects every part of the compiler and runtime

**Status:** Phase 1.1 planning complete, ready for implementation

---

## Phase 1.1: Extend Slice Header for Endianness

### Current State

**File:** `src/runtime/slice.go` (line 15-19)
```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

**Size:** 24 bytes on 64-bit systems (8 + 8 + 8)

### Target State

```go
type slice struct {
	array     unsafe.Pointer  // 8 bytes
	len       int64           // 8 bytes
	cap       int64           // 8 bytes
	byteOrder uint8           // 1 byte
	_         [7]byte         // 7 bytes padding (for alignment)
}
```

**New Size:** 32 bytes on 64-bit systems (aligned)

**ByteOrder Values:**
- `0` - NativeEndian (default, zero-cost)
- `1` - LittleEndian
- `2` - BigEndian

### Files Requiring Changes

#### Runtime Files
1. **src/runtime/slice.go**
   - Update `slice` struct definition
   - Add `byteOrder` field
   - Update all slice construction functions
   - Add endianness validation

2. **src/runtime/mbarrier.go**
   - Update write barriers for new slice size
   - Handle 32-byte slices in GC

3. **src/runtime/mgc.go**
   - Update GC scanning for new slice layout
   - Ensure byteOrder field is ignored by GC

#### Type System Files
4. **src/internal/abi/type.go**
   - Update `SliceType` struct (line 479-482)
   - Add endianness metadata

5. **src/cmd/compile/internal/types/type.go**
   - Update compiler's slice type representation
   - Add byteOrder tracking

6. **src/cmd/compile/internal/types/size.go**
   - Update slice size calculations (24 → 32 bytes)
   - Update alignment requirements

#### Compiler Files
7. **src/cmd/compile/internal/ssagen/ssa.go**
   - Update SSA generation for slice operations
   - Generate code to set/check byteOrder

8. **src/cmd/compile/internal/walk/builtin.go**
   - Update built-in function handling (len, cap, append, etc.)
   - Preserve byteOrder across operations

9. **src/cmd/compile/internal/walk/assign.go**
   - Update slice assignment handling
   - Handle byteOrder propagation

#### Reflection Files
10. **src/reflect/value.go**
    - Update reflection for new slice layout
    - Expose byteOrder through reflection API

### Implementation Strategy

**Phase 1.1.1: Add Field (Compatibility Mode)**
- Add `byteOrder` field to slice struct
- Initialize to 0 (NativeEndian) everywhere
- Ensure all existing code continues to work
- **Test:** Full test suite should pass

**Phase 1.1.2: Wire Up Infrastructure**
- Update all slice creation points to set byteOrder
- Add byteOrder parameter to make/append
- Update reflection to handle new field
- **Test:** Verify byteOrder is preserved

**Phase 1.1.3: Implement Endianness Operations**
- Add runtime functions for endianness conversion
- Implement (*[]T, LittleEndian)(s) syntax
- Add compile-time checks
- **Test:** Endianness conversion tests

### Known Challenges

**Challenge 1: ABI Compatibility**
- **Issue:** Changing slice size breaks ABI compatibility
- **Impact:** Can't call Go code from Moxie (or vice versa)
- **Mitigation:** This is acceptable for Moxie (it's a fork, not an extension)
- **Status:** Document as breaking change

**Challenge 2: Memory Overhead**
- **Issue:** 33% increase in slice header size (24 → 32 bytes)
- **Impact:** More memory usage for slice headers
- **Mitigation:** Modern systems have plenty of RAM, safety is worth it
- **Status:** Acceptable tradeoff

**Challenge 3: Performance Impact**
- **Issue:** Larger slices, more memory to copy
- **Impact:** Potential slowdown in slice operations
- **Mitigation:** Benchmark and optimize hot paths
- **Status:** Needs measurement

**Challenge 4: Existing Code**
- **Issue:** All existing Go code uses 24-byte slices
- **Impact:** Need to recompile everything
- **Mitigation:** Clean rebuild, no gradual migration possible
- **Status:** Expected for language fork

### Testing Strategy

**Unit Tests:**
- Slice creation with explicit endianness
- Endianness preservation across operations
- Endianness conversion correctness

**Integration Tests:**
- Interoperability between different endianness slices
- Performance benchmarks (before/after)
- Memory usage validation

**Stress Tests:**
- Large slices with endianness
- Concurrent slice operations
- GC behavior with new slice layout

### Success Criteria

- [ ] Slice struct updated with byteOrder field
- [ ] All existing tests pass with NativeEndian
- [ ] Can create slices with LittleEndian/BigEndian
- [ ] Endianness preserved across append/copy
- [ ] Reflection correctly reports byteOrder
- [ ] Performance impact < 5% for common operations
- [ ] Memory overhead documented and acceptable

---

## Phase 1.2: Remove Platform-Dependent int/uint

### Current State

**Built-in Types:**
- `int` - platform-dependent (int32 or int64)
- `uint` - platform-dependent (uint32 or uint64)

**Used In:**
- Array/slice indices
- Loop variables
- len() and cap() return types
- Most integer math in user code

### Target State

**Removed Types:**
- ❌ `int` - Must use `int32` or `int64` explicitly
- ❌ `uint` - Must use `uint32` or `uint64` explicitly

**Standard Sizes:**
- Indices, lengths, capacities: `int64`
- General integers: User's choice (int32/int64)

### Migration Strategy

**Phase 1.2.1: Add Deprecation Warnings**
- Compiler warns on `int` and `uint` usage
- Suggest `int32`/`int64` replacement
- Grace period for migration

**Phase 1.2.2: Update Standard Library**
- Systematically replace all `int`/`uint` in stdlib
- Use `int64` for sizes/indices
- Use `int32` for most other uses

**Phase 1.2.3: Update Built-in Functions**
- `len()` returns `int64`
- `cap()` returns `int64`
- `copy()` returns `int64`
- Array indices become `int64`

**Phase 1.2.4: Remove Types**
- Remove `int` and `uint` from type system
- Update compiler to reject these types
- **Breaking change** - no backward compatibility

### Files Requiring Changes

#### Built-in Definitions
1. **src/builtin/builtin.go**
   - Remove `int` and `uint` type definitions
   - Update len(), cap(), copy() signatures

2. **src/cmd/compile/internal/types/universe.go**
   - Remove `int` and `uint` from predeclared types
   - Update built-in function signatures

#### Type Checker
3. **src/cmd/compile/internal/types/type.go**
   - Remove TINT and TUINT type kinds
   - Update type checking rules

4. **src/cmd/compile/internal/typecheck/typecheck.go**
   - Reject `int` and `uint` in type expressions
   - Provide helpful error messages

#### Runtime
5. **src/runtime/*.go** (hundreds of files)
   - Replace all `int` with `int64` or `int32`
   - Systematic replacement required
   - High risk of bugs

6. **src/runtime/slice.go**
   - Update len/cap to use int64 (already planned in 1.1)

#### Standard Library
7. **src/**/* (all packages)
   - Replace `int` → `int32` or `int64`
   - Careful analysis needed for each use
   - Massive undertaking (~15,000 files)

### Estimated Effort

**Analysis:** 2 days - Identify all int/uint uses
**Standard Library:** 5-7 days - Systematic replacement
**Runtime:** 3-4 days - Critical path, high risk
**Testing:** 2-3 days - Comprehensive validation
**Total:** 12-16 days

### Risks

**Risk 1: Scope** 🔴
- **Issue:** Touching thousands of files
- **Impact:** High chance of introducing bugs
- **Mitigation:** Automated tools, systematic approach
- **Status:** VERY HIGH RISK

**Risk 2: Performance** 🟡
- **Issue:** int64 everywhere might be slower on 32-bit
- **Impact:** Performance hit on 32-bit systems
- **Mitigation:** Acceptable (Moxie targets modern 64-bit)
- **Status:** Low concern

**Risk 3: Testing** 🔴
- **Issue:** Can't test until everything is converted
- **Impact:** Big-bang integration, high risk
- **Mitigation:** Careful analysis, staged approach
- **Status:** HIGH RISK

**Recommendation:** This is the riskiest part of Phase 1. Consider:
1. Creating automated conversion tool
2. Converting in stages (stdlib → runtime → compiler)
3. Extensive testing at each stage

---

## Phase 1.3: Implement Explicit Pointer Types

### Current State

**Reference Types (Implicit Pointers):**
- `[]T` - Slice (internally a pointer to slice header)
- `map[K]V` - Map (internally a pointer to hmap)
- `chan T` - Channel (internally a pointer to hchan)

**Problem:** Hidden pointer semantics, confusing behavior

### Target State

**Explicit Pointer Types:**
- `*[]T` - Explicit slice pointer
- `*map[K]V` - Explicit map pointer
- `*chan T` - Explicit channel pointer

**Syntax Changes:**
```go
// Before (Go)
s := make([]int, 0, 10)
m := make(map[string]int)
ch := make(chan int, 5)

// After (Moxie)
s := &[]int{}
m := &map[string]int{}
ch := &chan int{cap: 5}
```

### Implementation Strategy

**Phase 1.3.1: Parser Updates**
- Accept `*[]T`, `*map[K]V`, `*chan T` syntax
- Maintain compatibility with `[]T` syntax temporarily
- Update type parser

**Phase 1.3.2: Type System Updates**
- Treat `[]T` internally as `*[]T`
- Update type checking to handle explicit pointers
- Add nil checks

**Phase 1.3.3: Compiler Auto-Dereferencing**
- Automatically dereference `*[]T` for indexing
- Maintain ergonomic access (no manual dereferencing)
- Optimize away redundant operations

**Phase 1.3.4: Runtime Updates**
- Maps and channels already pointers internally
- Slices need adjustment (header on heap)
- Update allocation functions

### Files Requiring Changes

#### Parser
1. **src/cmd/compile/internal/syntax/parser.go**
   - Accept `*[]T`, `*map[K]V`, `*chan T` syntax
   - Update type expressions

2. **src/cmd/compile/internal/types2/type.go**
   - Add explicit pointer variants
   - Update type equivalence rules

#### Type Checker
3. **src/cmd/compile/internal/typecheck/**
   - Update slice/map/chan operations
   - Add nil checking
   - Handle auto-dereferencing

#### Compiler Backend
4. **src/cmd/compile/internal/ssagen/**
   - Generate dereference code automatically
   - Optimize away redundant checks
   - Update call conventions

#### Runtime
5. **src/runtime/map.go**
   - Already pointer-based, minimal changes
   - Update allocation functions

6. **src/runtime/chan.go**
   - Already pointer-based, minimal changes
   - Update allocation functions

7. **src/runtime/slice.go**
   - Significant changes needed
   - Heap-allocate slice headers
   - Update all slice operations

### Auto-Dereferencing Rules

**Slice Operations:**
```go
s := &[]int{1, 2, 3}
x := s[0]        // Auto-dereference: (*s)[0]
s[1] = 5         // Auto-dereference: (*s)[1] = 5
len(s)           // Auto-dereference: len(*s)
```

**Map Operations:**
```go
m := &map[string]int{"a": 1}
x := m["a"]      // Auto-dereference: (*m)["a"]
m["b"] = 2       // Auto-dereference: (*m)["b"] = 2
```

**Channel Operations:**
```go
ch := &chan int{cap: 5}
ch <- 1          // Auto-dereference: (*ch) <- 1
x := <-ch        // Auto-dereference: x := <-(*ch)
```

### Risks

**Risk 1: Backward Compatibility** 🔴
- **Issue:** Breaking change from Go
- **Impact:** All Go code needs rewriting
- **Mitigation:** This is a Moxie feature, expected
- **Status:** Acceptable for fork

**Risk 2: Complexity** 🟡
- **Issue:** Auto-dereferencing is complex
- **Impact:** Compiler complexity increases
- **Mitigation:** Careful implementation, testing
- **Status:** Manageable

**Risk 3: Performance** 🟡
- **Issue:** Extra indirection for slices
- **Impact:** Potential slowdown
- **Mitigation:** Benchmark and optimize
- **Status:** Needs measurement

---

## Phase 1 Implementation Order

### Recommended Sequence

**Week 1-2: Phase 1.1 (Slice Headers)**
- Lower risk, foundational
- Can be tested independently
- Enables Phase 1.3

**Week 3: Phase 1.3 (Pointer Types)**
- Medium risk
- Depends on 1.1
- Tests explicit pointer semantics

**Week 4-6: Phase 1.2 (Remove int/uint)**
- Highest risk
- Massive scope
- Do last when other changes stable

### Alternative: Defer Phase 1.2

**Option:** Skip Phase 1.2 for now, do in Phase 4+
- **Rationale:** int/uint removal is extremely risky and large-scope
- **Benefit:** Can ship Phases 1.1 and 1.3 faster
- **Drawback:** Platform-dependent types remain
- **Recommendation:** Consider deferring to reduce Phase 1 risk

---

## Testing Strategy

### Phase 1.1 Testing
- Slice creation with different endianness
- Endianness preservation
- Conversion between endianness
- Performance benchmarks

### Phase 1.2 Testing
- All existing tests with int64
- Range checking
- Overflow detection
- Performance comparison

### Phase 1.3 Testing
- Nil slice/map/chan handling
- Auto-dereferencing correctness
- Memory safety
- Concurrent access

### Integration Testing
- All phases working together
- Standard library compatibility
- Real-world code patterns

---

## Success Criteria

### Phase 1.1
- [ ] Slice header includes byteOrder
- [ ] Endianness preserved across operations
- [ ] Can convert between endianness
- [ ] All tests pass
- [ ] Performance acceptable (<5% overhead)

### Phase 1.2
- [ ] int and uint types removed
- [ ] All stdlib uses int32/int64
- [ ] len()/cap() return int64
- [ ] All tests pass
- [ ] Migration guide complete

### Phase 1.3
- [ ] *[]T, *map[K]V, *chan T syntax works
- [ ] Auto-dereferencing implemented
- [ ] Nil checking works
- [ ] All tests pass
- [ ] Ergonomics validated

---

## Risk Assessment

### Overall Phase 1 Risk: 🔴 VERY HIGH

**Major Risks:**
1. **Scope** - Touching core runtime and compiler
2. **Complexity** - Fundamental type system changes
3. **Testing** - Hard to test until complete
4. **Performance** - Potential regressions
5. **Bugs** - High chance of introducing subtle bugs

### Mitigation Strategy

1. **Document First** ✅ (This document)
2. **Plan Carefully** - Detailed analysis before coding
3. **Test Incrementally** - Test each sub-phase
4. **Automate** - Use tools for systematic changes
5. **Benchmark** - Measure performance continuously
6. **Rollback Plan** - Clear reversion strategy

---

## Recommendation

**PHASE 1 APPROACH:**

Given the extreme complexity and risk of Phase 1, I recommend:

1. **Do Phase 1.1** (Slice headers with endianness)
   - Moderate complexity
   - Clear benefit
   - Testable independently

2. **Do Phase 1.3** (Explicit pointer types)
   - Moderate complexity
   - Big ergonomic win
   - Testable with 1.1

3. **DEFER Phase 1.2** (Remove int/uint) to Phase 4+
   - Extremely high risk
   - Massive scope
   - Can be done later when Moxie is more stable

This reduces Phase 1 from 15-22 days to 8-12 days and significantly lowers risk.

---

## Next Steps

1. **Review this plan** with stakeholders
2. **Decide on scope** (all of Phase 1 or defer 1.2)
3. **Create detailed sub-phase plans** for chosen work
4. **Set up testing infrastructure** for Phase 1
5. **Begin implementation** with Phase 1.1

---

**Status:** PLANNING COMPLETE - Ready for Review

**Recommendation:** Review plan, make scope decision, then proceed

**Estimated Time to Implement:** 8-12 days (if deferring 1.2) or 15-22 days (full Phase 1)
