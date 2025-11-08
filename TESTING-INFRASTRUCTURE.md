# Moxie Testing Infrastructure - Strategy & Documentation

**Date:** 2025-11-08
**Phase:** Phase 0.4 - Testing Infrastructure
**Status:** DOCUMENTED - Implementation Deferred to Post-Build

---

## Executive Summary

**Decision:** Document testing infrastructure and required changes, but **defer implementation** until we have a working build and have completed Phase 0.3.1 (baseline build test).

**Rationale:**
1. **Dependencies** - Can't test without a working build
2. **Risk Management** - Test infrastructure depends on build system working
3. **Incremental Approach** - Build first, test second
4. **Clear Documentation** - Know exactly what to change when ready
5. **Backward Compatibility** - Maintain test compatibility during branding transition

---

## Current Testing Infrastructure Overview

### Testing Levels

**1. Unit Tests**
- Location: Throughout codebase in `*_test.go` files
- Run with: `go test ./...` or `go test std cmd`
- Coverage: Individual packages and functions
- Execution time: Fast (seconds to minutes)

**2. Integration Tests**
- Location: `test/` directory (~2000+ test files)
- Run with: `../bin/go test cmd/internal/testdir`
- Coverage: Compiler, runtime, language features
- Execution time: Medium (minutes)

**3. Toolchain Tests**
- Location: `src/cmd/*/testdata/` directories
- Run with: `go test cmd/compile`, `go test cmd/link`, etc.
- Coverage: Compiler, linker, assembler, other tools
- Execution time: Medium (minutes)

**4. Full Test Suite**
- Location: All tests combined
- Run with: `cd src && ./all.bash`
- Coverage: Everything
- Execution time: Slow (10-30 minutes)

### Key Scripts

| Script | Purpose | Changes Needed |
|--------|---------|----------------|
| `src/all.bash` | Build + full test suite | Update messages |
| `src/run.bash` | Run tests only (no rebuild) | Update messages, binary paths |
| `src/race.bash` | Run race detector tests | Update messages |
| `src/cmd/dist/test.go` | Test orchestration (1870 lines) | Update messages, binary names |

---

## Test Infrastructure Components

### 1. Test Runner: cmd/dist test

**File:** `src/cmd/dist/test.go` (1870 lines)

**Purpose:**
- Orchestrates all tests
- Handles test sharding for parallel execution
- Manages test timeouts and retries
- Reports test results (text or JSON)

**Key Functions:**
- `registerTests()` - Registers all standard library and cmd tests
- `registerStdTest()` - Registers individual package tests
- `registerCgoTests()` - CGO-specific tests
- `registerRaceTests()` - Race detector tests

**Test Categories Registered:**
1. Standard library tests (`go test std`)
2. Command tests (`go test cmd`)
3. Special variant tests (purego, osusergo, FIPS, etc.)
4. Platform-specific tests (iOS, Android, etc.)
5. Race detector tests
6. CGO tests
7. Test directory tests (compiler/runtime regression tests)

### 2. Test Directory: test/

**Location:** `/home/mleku/src/github.com/mleku/moxie/test/`

**Contents:**
- ~2000+ individual test files
- Subdirectories for specific test categories:
  - `test/abi/` - ABI (Application Binary Interface) tests
  - `test/chan/` - Channel tests
  - `test/arenas/` - Memory arena tests
  - Various individual `.go` files for language feature tests

**Test Types:**
- **Run tests:** Tests that should compile and run successfully
- **Error tests:** Tests that should fail to compile with specific errors
- **Regression tests:** Tests for historical bugs
- **Feature tests:** Tests for specific language features

**Test Format:**
- Single file tests: `test/foo.go`
- Directory tests: `test/foo.dir/` (multiple files)
- Expected output: `test/foo.out` (for tests with output)

### 3. Test Execution: cmd/internal/testdir

**Package:** `cmd/internal/testdir`

**Purpose:**
- Reads test files from `test/` directory
- Parses special directives in comments:
  - `// run` - Compile and run
  - `// errorcheck` - Should fail with specific errors
  - `// compile` - Should compile successfully
  - `// runoutput` - Run and compare output
  - `// build` - Should build successfully
- Executes tests based on directives
- Compares actual vs. expected results

**Test Sharding:**
- Tests are sharded for parallel execution
- Default: 10 shards on builders, 1 shard locally
- Controlled by `GO_TEST_SHARDS` environment variable

---

## Changes Required

### Phase 0.4.1: Document Test Strategy (Complete)

**Objective:** Understand testing infrastructure (this document)

**Status:** ✅ COMPLETE

### Phase 0.4.2: Update Test Messages (After Build Works)

**Files to Update:**

#### src/all.bash

**Line 8:** Error message
```bash
# Before
echo 'all.bash must be run from $GOROOT/src' 1>&2

# After
echo 'all.bash must be run from $MOXIEROOT/src' 1>&2
# Note: Keep GOROOT for backward compatibility, mention both
```

**Line 13:** Banner command
```bash
# Before
../bin/go tool dist banner # print build info

# After
../bin/moxie tool dist banner # print build info
# Note: This happens AFTER binary rename in Phase 0.3.3
```

#### src/run.bash

**Line 27-28:** Error message
```bash
# Before
if [ ! -f ../bin/go ]; then
	echo 'run.bash must be run from $GOROOT/src after installing cmd/go' 1>&2

# After
if [ ! -f ../bin/moxie ]; then
	echo 'run.bash must be run from $MOXIEROOT/src after installing cmd/moxie' 1>&2
# Note: This happens AFTER binary rename
```

**Line 33:** Binary reference
```bash
# Before
eval $(../bin/go tool dist env)

# After
eval $(../bin/moxie tool dist env)
# Note: This happens AFTER binary rename
```

**Line 53:** Binary reference
```bash
# Before
exec ../bin/go tool dist test -rebuild "$@"

# After
exec ../bin/moxie tool dist test -rebuild "$@"
# Note: This happens AFTER binary rename
```

#### src/cmd/dist/test.go

**Multiple References:**

**Line 126:** Binary variable reference
```go
// Before
cmd := exec.Command(gorootBinGo, "env", "CGO_ENABLED")

// After
cmd := exec.Command(gorootBinMoxie, "env", "CGO_ENABLED")
// Note: Variable name change in dist tool
```

**Line 657:** List command
```go
// Before
cmd := exec.Command(gorootBinGo, "list")

// After
cmd := exec.Command(gorootBinMoxie, "list")
```

**Line 312:** Comment
```go
// Before
// goTest represents all options to a "go test" command. The final command will

// After
// goTest represents all options to a "moxie test" command. The final command will
```

**Line 620:** Comment
```go
// Before
// These tests *must* be able to run normally as part of "go test std cmd",

// After
// These tests *must* be able to run normally as part of "moxie test std cmd",
```

**And ~20+ other references to "go test", "go list", "go build" in comments and strings**

---

### Phase 0.4.3: Update Binary References (After Phase 0.3.3)

**Deferred Because:**
- Depends on binary rename (Phase 0.3.3)
- Need working build first
- Many interdependent changes
- Can be done systematically after core branding

**When Implemented:**

**Global Changes in test.go:**
```go
// Variable definitions (top of file)
var (
	gorootBinGo    = filepath.Join(runtime.GOROOT(), "bin/go")        // Before
	gorootBinMoxie = filepath.Join(runtime.GOROOT(), "bin/moxie")     // After
)
```

**Command executions:**
```go
// Before
exec.Command(gorootBinGo, "test", ...)
exec.Command(gorootBinGo, "build", ...)
exec.Command(gorootBinGo, "list", ...)

// After
exec.Command(gorootBinMoxie, "test", ...)
exec.Command(gorootBinMoxie, "build", ...)
exec.Command(gorootBinMoxie, "list", ...)
```

---

### Phase 0.4.4: Update Test Expectations (Future - Phase 1+)

**Challenge:** Tests expect specific output containing "go" references

**Examples:**

1. **Version tests** expect:
   ```
   go version go1.24.6 linux/amd64
   ```
   Need to update to:
   ```
   moxie version moxie0.1.0 linux/amd64
   ```

2. **Error message tests** expect:
   ```
   go: cannot find package...
   ```
   Need to update to:
   ```
   moxie: cannot find package...
   ```

3. **Help text tests** expect:
   ```
   Go is a tool for managing Go source code.
   ```
   Need to update to:
   ```
   Moxie is a tool for managing Moxie source code.
   ```

**Deferred Because:**
- Requires grep through ~2000+ test files
- Many test `.out` files need updates
- Can be done systematically after core changes work
- Lower risk than build/runtime changes

**Implementation Strategy:**
1. Search for all test files with "go" in expected output
2. Update expected outputs to match new branding
3. Run tests and fix failures iteratively
4. Document any tests that need special handling

---

## Environment Variables Used by Tests

### Test-Specific Variables

| Variable | Purpose | Change Status |
|----------|---------|---------------|
| `GO_TEST_SHARDS` | Number of test shards for parallel execution | Keep for now |
| `GO_BUILDER_NAME` | Builder name for CI/CD systems | Keep for now |
| `GO_TEST_SHORT` | Run tests in short mode | Keep for now |
| `GO_TEST_TIMEOUT_SCALE` | Scale test timeouts | Keep for now |
| `GOENV` | Go environment file location | Keep (Phase 4+) |
| `GOTRACEBACK` | Traceback verbosity | Keep for now |

**Note:** Per Phase 0.2 decision, all GO* variables remain unchanged for backward compatibility.

---

## Test Execution Flow

### Standard Test Run (all.bash)

```bash
cd src
./all.bash
```

**What Happens:**
1. Sources `make.bash` to build toolchain
2. Runs `bash run.bash --no-rebuild` to execute tests
3. Prints banner with build info

**Expected Output (Current):**
```
Building Go cmd/dist using .../go1.24.6
Building Go toolchain1 using .../go1.24.6
Building Go bootstrap cmd/go (go_bootstrap) using Go toolchain1
Building Go toolchain2 using go_bootstrap and Go toolchain1
Building Go toolchain3 using go_bootstrap and Go toolchain2
---
Installed Go for linux/amd64 in /home/mleku/src/github.com/mleku/moxie

##### Testing packages.
ok      archive/tar     0.123s
ok      archive/zip     0.456s
...
##### PASS
```

**Expected Output (After Changes):**
```
Building Moxie cmd/dist using .../go1.24.6
Building Moxie toolchain1 using .../go1.24.6
Building Moxie bootstrap cmd/go (go_bootstrap) using Go toolchain1
Building Moxie toolchain2 using go_bootstrap and Go toolchain1
Building Moxie toolchain3 using go_bootstrap and Go toolchain2
---
Installed Moxie for linux/amd64 in /home/mleku/src/github.com/mleku/moxie

##### Testing packages.
ok      archive/tar     0.123s
ok      archive/zip     0.456s
...
##### PASS
```

### Test-Only Run (run.bash)

```bash
cd src
./run.bash
```

**What Happens:**
1. Verifies `../bin/go` exists (later `../bin/moxie`)
2. Sets up environment
3. Executes `../bin/go tool dist test -rebuild`

### Specific Package Test

```bash
../bin/go test archive/tar
../bin/go test cmd/compile
../bin/go test cmd/internal/testdir
```

### Test Directory Tests

```bash
../bin/go test cmd/internal/testdir
# Or with sharding:
../bin/go test cmd/internal/testdir -shard=1 -shards=10
```

---

## Known Issues & Risks

### Risk 1: Test Count and Scope
**Issue:** ~2000+ test files in test/ directory, thousands more in packages
**Impact:** Large surface area for branding changes
**Mitigation:** Systematic approach, automated search/replace where safe
**Status:** Manageable but time-consuming

### Risk 2: Test Output Dependencies
**Issue:** Many tests check exact output strings
**Impact:** Tests will fail if output doesn't match expectations
**Mitigation:** Update .out files alongside source changes
**Status:** Known issue, will be addressed systematically

### Risk 3: Third-Party Test Integration
**Issue:** External tools may depend on "go test" command structure
**Impact:** Tools may break if command semantics change
**Mitigation:** Maintain command-line compatibility, only change branding
**Status:** Low risk - only changing output messages, not behavior

### Risk 4: CI/CD Integration
**Issue:** Build systems expect GO* environment variables
**Impact:** CI systems may need updates
**Mitigation:** Support both GO* and MOXIE* variables during transition
**Status:** Deferred to Phase 4+

---

## Testing Strategy

### Phase 0.4.1: Baseline Understanding (Complete)
**Status:** ✅ COMPLETE (This document)
**Duration:** Documentation only
**Risk:** None

### Phase 0.4.2: Test Current Build (After Phase 0.3.1)
**Status:** ⏳ PENDING (waiting for successful build)
**Duration:** 10-30 minutes
**Risk:** Low (no changes, just testing)

**Steps:**
```bash
cd /home/mleku/src/github.com/mleku/moxie/src

# Run full test suite
./all.bash

# Or just quick smoke tests
./make.bash
./run.bash --no-rebuild -run=^archive/
```

**Success Criteria:**
- [ ] Build completes successfully
- [ ] Tests run
- [ ] Document pass/fail rate
- [ ] Note any failures related to branding changes from Phase 0.1

**Expected Issues:**
- Version tests may fail (we changed version output in Phase 0.1)
- Error message tests may fail (we changed error messages in Phase 0.1)
- Most core tests should pass

### Phase 0.4.3: Fix Branding-Related Test Failures (After Phase 0.4.2)
**Status:** ⏳ PENDING
**Duration:** 1-2 hours
**Risk:** Low (isolated to test expectations)

**Steps:**
1. Run tests and capture failures
2. Identify failures caused by Phase 0.1 branding changes
3. Update test expectations to match new output
4. Re-run tests to verify fixes
5. Document any remaining failures

**Example Fixes:**
```go
// test/version/version.go or similar
// Before expected output:
// go version go1.24.6 linux/amd64

// After expected output:
// moxie version moxie0.1.0 linux/amd64
```

### Phase 0.4.4: Update Test Script Messages (Future)
**Status:** ⏳ DEFERRED (until after Phase 0.3.2)
**Duration:** 30 minutes
**Risk:** Low (cosmetic changes only)

**Changes:**
- Update all.bash error messages
- Update run.bash error messages
- Update race.bash messages
- Test that scripts still work

### Phase 0.4.5: Update Binary References (Future - After Phase 0.3.3)
**Status:** ⏳ DEFERRED (until after binary rename)
**Duration:** 1-2 hours
**Risk:** Medium (affects test execution)

**Changes:**
- Update gorootBinGo → gorootBinMoxie variable
- Update all `../bin/go` → `../bin/moxie` references
- Update test.go command executions
- Full test suite run to verify

---

## Test Categories and Estimated Update Effort

### 1. Standard Library Package Tests
- **Count:** 150+ packages
- **Location:** Throughout src/
- **Update Effort:** Low - mostly automatic via `go test`
- **Changes:** None needed - tests are package-focused

### 2. Command Tests
- **Count:** 30+ commands
- **Location:** src/cmd/*/
- **Update Effort:** Medium - some tests check command output
- **Changes:** Update expected outputs for version, help text, errors

### 3. Test Directory Tests
- **Count:** 2000+ files
- **Location:** test/
- **Update Effort:** Medium-High - many test outputs to verify
- **Changes:** Update .out files for tests that check exact output

### 4. Build/Toolchain Tests
- **Count:** 100+ tests
- **Location:** src/cmd/dist/, src/cmd/compile/, etc.
- **Update Effort:** Medium - some output checking
- **Changes:** Update build message expectations

### 5. Integration Tests
- **Count:** 50+ tests
- **Location:** Various
- **Update Effort:** Medium - check end-to-end workflows
- **Changes:** Update expected outputs and error messages

**Total Estimated Effort:** 4-8 hours of systematic updates

---

## Implementation Checklist

### Phase 0.4.1: Documentation (Complete)
- [x] Document testing infrastructure
- [x] Identify test categories
- [x] List required changes
- [x] Define testing strategy
- [x] Create this document

### Phase 0.4.2: Baseline Test (After Phase 0.3.1)
- [ ] Verify build completes (from Phase 0.3.1)
- [ ] Run full test suite: `cd src && ./all.bash`
- [ ] Document test results
- [ ] Identify branding-related failures
- [ ] Document other failures (if any)

### Phase 0.4.3: Fix Branding Test Failures (After Phase 0.4.2)
- [ ] Update version test expectations
- [ ] Update error message test expectations
- [ ] Update help text test expectations
- [ ] Re-run tests to verify fixes
- [ ] Achieve clean test run (or document known issues)

### Phase 0.4.4: Update Test Scripts (After Phase 0.3.2)
- [ ] Update all.bash line 8 (error message)
- [ ] Update run.bash lines 27-28 (error message)
- [ ] Update run.bash line 33 (binary reference)
- [ ] Update run.bash line 53 (binary reference)
- [ ] Test scripts work correctly

### Phase 0.4.5: Update test.go (After Phase 0.3.3)
- [ ] Update variable gorootBinGo → gorootBinMoxie
- [ ] Update all exec.Command references
- [ ] Update comments referencing "go test"
- [ ] Update comments referencing "go build", "go list"
- [ ] Rebuild and verify tests run

### Phase 0.4.6: Comprehensive Test Update (Future - Phase 1+)
- [ ] Search all test/ files for "go version" output checks
- [ ] Search all test/ files for "go:" error message checks
- [ ] Update all .out files with expected outputs
- [ ] Run full test suite
- [ ] Fix any remaining failures
- [ ] Document test suite status

---

## Cross-Platform Considerations

### Unix/Linux/macOS ✅
- Primary focus for Phase 0
- `all.bash`, `run.bash`, `race.bash` scripts
- Well-tested, most developers use

### Windows ⏳
- Equivalent `.bat` scripts
- Will need similar updates
- Deferred to post-Unix success

### Plan 9 ⏳
- Equivalent `.rc` scripts
- Will need similar updates
- Minimal user base, lowest priority

**Strategy:** Unix first, Windows second, Plan 9 last

---

## Success Criteria

### Phase 0.4.2 Success
- [ ] Full test suite runs
- [ ] Pass/fail results documented
- [ ] Branding-related failures identified
- [ ] Ready for fixes

### Phase 0.4.3 Success
- [ ] All branding-related test failures fixed
- [ ] Clean test run (all tests pass)
- [ ] Or: Known failures documented with reasons

### Phase 0.4.4 Success
- [ ] Test script messages updated
- [ ] Scripts execute correctly
- [ ] Tests still run and pass

### Phase 0.4.5 Success
- [ ] Binary references updated in test.go
- [ ] Test orchestration works
- [ ] Full test suite passes

### Overall Phase 0.4 Success
- [ ] Testing infrastructure understood
- [ ] Strategy documented
- [ ] Changes identified and categorized
- [ ] Baseline tests run successfully
- [ ] Path forward clear

---

## Rollback Plan

### If Tests Fail After Changes

1. **Identify Scope**
   - Determine if failure is test-related or code-related
   - Check if failure existed before changes
   - Isolate to specific test category

2. **Revert Changes**
   ```bash
   # Revert specific file
   git checkout -- src/all.bash

   # Or revert specific commit
   git revert <commit-hash>
   ```

3. **Verify Baseline**
   ```bash
   # Run tests again
   cd src
   ./all.bash
   ```

4. **Document and Fix**
   - Document what broke
   - Understand root cause
   - Fix properly, not just make tests pass

### Version Control Strategy

```bash
# Before test updates
git add -A
git commit -m "Phase 0.4: Pre-test-updates checkpoint"

# After script updates
git add src/*.bash
git commit -m "Phase 0.4: Update test script messages"

# After test.go updates
git add src/cmd/dist/test.go
git commit -m "Phase 0.4: Update test.go binary references"

# After test expectation updates
git add test/
git commit -m "Phase 0.4: Update test expectations for branding"
```

---

## Test Execution Times

### Full Test Suite
- **Command:** `cd src && ./all.bash`
- **Time:** 10-30 minutes (depends on hardware)
- **Coverage:** Build + all tests

### Test Only (No Rebuild)
- **Command:** `cd src && ./run.bash`
- **Time:** 5-20 minutes
- **Coverage:** All tests, skip build

### Quick Smoke Test
- **Command:** `cd src && ./run.bash -run=^archive/`
- **Time:** 10-30 seconds
- **Coverage:** Just archive packages

### Standard Library Only
- **Command:** `../bin/go test std`
- **Time:** 3-10 minutes
- **Coverage:** Standard library packages

### Commands Only
- **Command:** `../bin/go test cmd`
- **Time:** 2-5 minutes
- **Coverage:** Command-line tools

### Single Package
- **Command:** `../bin/go test archive/tar`
- **Time:** 1-5 seconds
- **Coverage:** One package

---

## Current Status

**Phase 0.4.1:** ✅ COMPLETE (Documentation)
**Phase 0.4.2:** ⏳ PENDING (Baseline test - after Phase 0.3.1)
**Phase 0.4.3:** ⏳ PENDING (Fix failures - after baseline)
**Phase 0.4.4:** ⏳ PENDING (Script updates - after Phase 0.3.2)
**Phase 0.4.5:** ⏳ DEFERRED (Binary references - after Phase 0.3.3)
**Phase 0.4.6:** ⏳ DEFERRED (Comprehensive - Phase 1+)

---

## Dependencies

### Prerequisites
- ✅ Phase 0.1: Branding (COMPLETE)
- ✅ Phase 0.2: Environment variables strategy (COMPLETE)
- ✅ Phase 0.3: Build system documentation (COMPLETE)
- ⏳ Phase 0.3.1: Successful build (PENDING - manual step)

### Blockers
- **Cannot run tests** until Phase 0.3.1 completes (working build)
- **Cannot update binary references** until Phase 0.3.3 (binary rename)

### Enables
- Phase 0 completion
- Confidence in build quality
- Foundation for Phase 1 development

---

## Next Steps

### Immediate (After Phase 0.3.1 Completes)
1. Run baseline test suite
2. Document results
3. Identify any branding-related failures from Phase 0.1

### After Baseline Tests
1. Fix branding-related test failures
2. Get clean test run
3. Document test suite status

### After Phase 0.3.2 (Build Messages)
1. Update test script messages
2. Verify scripts work

### After Phase 0.3.3 (Binary Rename)
1. Update binary references in test.go
2. Full test suite verification

### Future (Phase 1+)
1. Comprehensive test expectation updates
2. Full test suite validation
3. Document any language-change-related test updates

---

## Conclusion

The testing infrastructure is well-documented and ready for systematic updates. The pragmatic approach is:

1. **First:** Run baseline tests to verify current state
2. **Second:** Fix branding-related test failures
3. **Third:** Update test script messages (after build messages update)
4. **Fourth:** Update binary references (after binary rename)
5. **Fifth:** Comprehensive test updates (as language changes require)

This minimizes risk and allows incremental validation.

**Phase 0.4 Status:** DOCUMENTED - Ready for Baseline Test (after build)

**Next Action:** Wait for Phase 0.3.1 (build test), then run test suite

---

## Files Requiring Changes (Summary)

### Immediate Updates (After Build Works)
1. None - tests should run as-is

### After Branding Test Failures Identified
2. Various test expectation files in test/

### After Phase 0.3.2 (Message Updates)
3. **src/all.bash** - 1 change (error message)
4. **src/run.bash** - 3 changes (error message, binary refs)

### After Phase 0.3.3 (Binary Rename)
5. **src/cmd/dist/test.go** - Variable + ~20+ command references

### Future (Phase 1+)
6. **test/*.out** - Expected output files (~100+ files)
7. **Various *_test.go** - Test files that check exact output

**Total:** 6+ files for Phase 0, 100+ files for comprehensive updates

---

**Status:** Phase 0.4 ✅ DOCUMENTATION COMPLETE

**Outcome:** Clear strategy, comprehensive guide, ready for testing when build completes

**Impact:** Zero risk (no code changes), maximum preparation, clear dependencies
