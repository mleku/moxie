# Phase 7.2: Quality Tools - Implementation Plan

**Status**: 🚧 In Progress
**Started**: 2025-11-10
**Dependencies**: Phase 7.1 (Essential Tools)
**Parent**: PHASE7-PLAN.md

## Overview

Phase 7.2 focuses on code quality tools that help developers write better Moxie code through static analysis, improved error reporting, and build optimization.

## Goals

1. **Catch errors early** - Static analysis before runtime
2. **Improve code quality** - Detect anti-patterns and potential bugs
3. **Speed up builds** - Incremental compilation with caching
4. **Better diagnostics** - Clear, actionable error messages

## Priority Implementation Order

### 1. Linter (`moxie vet`) - Priority 1 ✅
**Value**: High - Catches bugs before they reach production
**Complexity**: Medium - AST analysis with multiple check modules
**Estimated Effort**: 5-7 days

### 2. Enhanced Error Messages - Priority 2 ⏳
**Value**: High - Improves developer experience
**Complexity**: Low-Medium - Error message formatting
**Estimated Effort**: 2-3 days

### 3. Build Caching - Priority 3 ⏳
**Value**: Medium-High - Speeds up development
**Complexity**: Medium - Cache invalidation logic
**Estimated Effort**: 3-4 days

---

## 7.2.1: Linter (`moxie vet`)

### Architecture

```
cmd/moxie/vet/
├── vet.go           Main vet command and orchestration
├── analyzer.go      AST analysis framework
├── checks.go        Check registration and management
├── memory.go        Memory management checks
├── channels.go      Channel safety checks
├── types.go         Type safety checks
├── const.go         Const correctness checks
├── errors.go        Error handling checks
├── report.go        Issue reporting and formatting
└── config.go        Configuration handling
```

### Check Categories

#### 1. Memory Management Checks (Phase 7.2.1a)
**Priority**: Highest - Moxie-specific

- **UnusedClone**: Detect `clone()` calls with unused results
  ```go
  // Bad
  clone(data)  // Result not used

  // Good
  newData := clone(data)
  ```

- **MissingFree**: Detect allocated resources without corresponding `free()`
  ```go
  // Bad
  data := clone(original)
  // No free() call - memory leak

  // Good
  data := clone(original)
  defer free(data)
  ```

- **DoubleFree**: Detect multiple `free()` calls on same resource
  ```go
  // Bad
  free(data)
  free(data)  // Double free!

  // Good
  free(data)
  data = nil
  ```

- **FreeAfterReturn**: Detect `free()` on value after return
  ```go
  // Bad
  func getData() *[]byte {
      data := clone(source)
      free(data)  // Freed before return!
      return data
  }
  ```

#### 2. Channel Safety Checks (Phase 7.2.1b)
**Priority**: High - Common source of bugs

- **UnbufferedDeadlock**: Detect potential deadlocks with unbuffered channels
  ```go
  // Bad
  ch := &chan int{}
  ch <- 1  // Blocks forever if no receiver

  // Good
  ch := &chan int{1}  // Buffered
  go func() { <-ch }()
  ch <- 1
  ```

- **ChannelLeak**: Detect channels that are never closed
  ```go
  // Bad
  ch := &chan int{10}
  // Never closed - goroutine leak

  // Good
  ch := &chan int{10}
  defer close(ch)
  ```

- **SendOnClosed**: Detect potential sends on closed channels
  ```go
  // Bad
  close(ch)
  ch <- 1  // Panic!
  ```

#### 3. Type Safety Checks (Phase 7.2.1c)
**Priority**: Medium - Prevent runtime errors

- **UnsafeCoercion**: Detect risky type coercions
  ```go
  // Warning: Size mismatch
  bytes := &[]byte{1, 2, 3}
  ints := (*[]int64)(bytes)  // 3 bytes != 8 bytes
  ```

- **EndiannessMismatch**: Detect missing endianness in network code
  ```go
  // Bad - network protocols need explicit endianness
  func parsePacket(data *[]byte) uint32 {
      return (*(*[]uint32)(data))[0]
  }

  // Good
  func parsePacket(data *[]byte) uint32 {
      return (*(*[]uint32, BigEndian)(data))[0]
  }
  ```

- **IntegerOverflow**: Detect potential integer overflows
  ```go
  // Warning
  var x uint8 = 255
  x = x + 1  // Overflow to 0
  ```

#### 4. Const Correctness Checks (Phase 7.2.1d)
**Priority**: Low-Medium - Code quality

- **MutableConst**: Detect const values that are modified (compile-time check already exists)
- **ShouldBeConst**: Suggest variables that could be const
  ```go
  // Could be const
  var pi = 3.14159  // Never modified
  ```

#### 5. Error Handling Checks (Phase 7.2.1e)
**Priority**: Medium - Reliability

- **UncheckedError**: Detect ignored error return values
  ```go
  // Bad
  file.Write(data)  // Error ignored

  // Good
  if err := file.Write(data); err != nil {
      return err
  }
  ```

- **ErrorShadow**: Detect error variable shadowing
  ```go
  // Bad
  err := doSomething()
  if err != nil {
      err := doSomethingElse()  // Shadows outer err
  }
  ```

### Implementation Phases

#### Phase 7.2.1a: Framework + Memory Checks (Days 1-2)
- Basic vet command structure
- AST analysis framework
- Memory management checks (highest priority)
- Initial test suite

#### Phase 7.2.1b: Channel Safety (Days 3-4)
- Channel safety checks
- Control flow analysis for deadlocks
- Expand test coverage

#### Phase 7.2.1c: Type Safety (Day 5)
- Type coercion checks
- Endianness validation
- Integer overflow detection

#### Phase 7.2.1d: Polish + Config (Days 6-7)
- Const correctness checks
- Error handling checks
- Configuration file support
- JSON output for IDE integration
- Documentation

### Commands

```bash
# Basic usage
moxie vet file.x              # Vet single file
moxie vet ./...               # Vet all files recursively
moxie vet package/path        # Vet package

# Specific checks
moxie vet --checks=memory ./...           # Only memory checks
moxie vet --checks=memory,channels ./...  # Multiple categories
moxie vet --checks=all ./...              # All checks (default)

# Output formats
moxie vet --format=text ./...    # Human-readable (default)
moxie vet --format=json ./...    # JSON for IDE integration
moxie vet --format=github ./...  # GitHub Actions format

# Severity filtering
moxie vet --min-severity=error ./...   # Only errors
moxie vet --min-severity=warning ./... # Warnings and errors

# Configuration
moxie vet --config=.moxie-vet.toml ./...  # Custom config
```

### Configuration File (.moxie-vet.toml)

```toml
[vet]
# Enable/disable check categories
checks = ["memory", "channels", "types", "const", "errors"]

# Minimum severity to report
min_severity = "warning"  # info, warning, error

[checks.memory]
enabled = true
unused_clone = "error"
missing_free = "warning"
double_free = "error"
free_after_return = "error"

[checks.channels]
enabled = true
unbuffered_deadlock = "warning"
channel_leak = "warning"
send_on_closed = "error"

[checks.types]
enabled = true
unsafe_coercion = "warning"
endianness_mismatch = "warning"
integer_overflow = "warning"

[checks.const]
enabled = true
mutable_const = "error"
should_be_const = "info"

[checks.errors]
enabled = true
unchecked_error = "warning"
error_shadow = "warning"
```

### Output Format

```
examples/test.x:15:2: [memory/missing_free] potential memory leak
  data := clone(original)
  ^
  help: add 'defer free(data)' after allocation

examples/test.x:23:5: [channels/unbuffered_deadlock] potential deadlock
  ch <- value
  ^
  help: use buffered channel or start receiver goroutine first

Summary: 2 errors, 3 warnings, 1 info
```

### JSON Output (for IDE integration)

```json
{
  "issues": [
    {
      "file": "examples/test.x",
      "line": 15,
      "column": 2,
      "severity": "warning",
      "category": "memory",
      "check": "missing_free",
      "message": "potential memory leak",
      "help": "add 'defer free(data)' after allocation"
    }
  ],
  "summary": {
    "errors": 2,
    "warnings": 3,
    "info": 1
  }
}
```

---

## 7.2.2: Enhanced Error Messages

### Goals
1. Map Go errors back to Moxie source
2. Add context and suggestions
3. Colorize output
4. Show code snippets

### Implementation
- Source line mapping during transpilation
- Error message post-processing
- Terminal color support
- Code snippet extraction

### Example

**Before**:
```
./test_example.go:15:2: undefined: foo
```

**After**:
```
test_example.x:12:2: error: undefined identifier 'foo'
   10 | func main() {
   11 |     x := 42
   12 |     foo(x)
      |     ^^^ not found in this scope
   13 | }
help: did you mean 'bar'?
```

---

## 7.2.3: Build Caching

### Goals
1. Cache transpiled .go files
2. Detect source changes
3. Invalidate on dependency changes
4. Share cache across projects (optional)

### Implementation
- Content-based hashing (SHA256)
- Dependency tracking
- Cache directory: `~/.moxie/cache/` or `.moxie-cache/`
- `moxie clean` command

### Cache Structure
```
.moxie-cache/
├── transpiled/          # Cached .go files
│   └── <hash>.go
├── metadata/           # Dependency metadata
│   └── <hash>.json
└── checksums/          # File checksums
    └── <file>.sha256
```

### Commands
```bash
moxie build --no-cache ./...   # Skip cache
moxie clean                    # Clear cache
moxie clean --cache            # Clear cache only
```

---

## Success Metrics

### Linter
- Catches 80%+ of memory management bugs
- Zero false positives on test suite
- Configurable with sensible defaults
- Fast (< 1s for typical project)

### Error Messages
- All Go errors mapped to Moxie source
- Helpful suggestions for common errors
- Clear, colorized output

### Build Caching
- 50%+ faster on incremental builds
- Correct invalidation (no stale builds)
- Minimal disk usage (< 100MB typical)

---

## Testing Strategy

### Linter Tests
- Unit tests for each check
- False positive prevention
- Performance benchmarks
- Integration with CI/CD

### Error Message Tests
- Test all error types
- Verify source mapping accuracy
- Terminal output validation

### Cache Tests
- Invalidation correctness
- Concurrent access
- Cleanup and maintenance

---

## Timeline

**Week 1**: Linter framework + Memory checks
**Week 2**: Channel + Type safety checks
**Week 3**: Error messages + Polish
**Week 4**: Build caching

**Total**: 4 weeks for complete Phase 7.2

---

## Dependencies

### Go Packages
- `go/ast` - AST analysis
- `go/types` - Type checking
- `go/parser` - Parsing
- `golang.org/x/tools/go/analysis` - Analysis framework
- `github.com/fatih/color` - Terminal colors

### Internal
- `cmd/moxie/preprocess.go` - Syntax preprocessing
- `cmd/moxie/typetrack.go` - Type tracking

---

## Notes

- Start with most valuable checks (memory management)
- Keep false positive rate low
- Provide actionable suggestions
- Make it fast enough for watch mode integration
- Consider adding quick fixes for common issues
- Plan for LSP integration (Phase 7.3)
