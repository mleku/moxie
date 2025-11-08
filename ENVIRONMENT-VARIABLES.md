# Moxie Environment Variables - Strategy & Mapping

**Date:** 2025-11-08
**Phase:** Phase 0.2 - Environment Variables
**Status:** DOCUMENTED - Implementation Deferred

---

## Executive Summary

**Decision:** Environment variable names will remain as `GOROOT`, `GOPATH`, etc. during Phase 0 for backward compatibility and incremental migration. Future phases will introduce `MOXIE*` aliases and eventually deprecate `GO*` variants.

**Rationale:**
1. **Backward Compatibility** - Existing tools and scripts expect GO* variables
2. **Incremental Migration** - Users can migrate gradually
3. **Build Complexity** - Changing now would break the build system before we test it
4. **Testing** - Need working build to validate changes
5. **Documentation First** - Document strategy before implementation

---

## All GO* Environment Variables

### Core Path Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOROOT` | `MOXIEROOT` | Moxie installation directory | Keep for now |
| `GOPATH` | `MOXIEPATH` | Workspace directory | Keep for now |
| `GOBIN` | `MOXIEBIN` | Binary installation directory | Keep for now |

### Platform Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOOS` | `MOXIEOS` | Target operating system | Keep for now |
| `GOARCH` | `MOXIEARCH` | Target architecture | Keep for now |

### Cache & Module Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOCACHE` | `MOXIECACHE` | Build cache directory | Keep for now |
| `GOMODCACHE` | `MOXIEMODCACHE` | Module cache directory | Keep for now |
| `GOCACHEPROG` | `MOXIECACHEPROG` | Cache program path | Keep for now |

### Architecture-Specific Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOARM` | `MOXIEARM` | ARM architecture version | Keep for now |
| `GOARM64` | `MOXIEARM64` | ARM64 architecture version | Keep for now |
| `GO386` | `MOXIE386` | 386 architecture variant | Keep for now |
| `GOAMD64` | `MOXIEAMD64` | AMD64 architecture level | Keep for now |
| `GOMIPS` | `MOXIEMIPS` | MIPS architecture variant | Keep for now |
| `GOMIPS64` | `MOXIEMIPS64` | MIPS64 architecture variant | Keep for now |
| `GOPPC64` | `MOXIEPPC64` | PPC64 architecture level | Keep for now |
| `GORISCV64` | `MOXIERISCV64` | RISC-V 64 extensions | Keep for now |
| `GOWASM` | `MOXIEWASM` | WebAssembly features | Keep for now |

### Module & Dependency Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOPROXY` | `MOXIEPROXY` | Module proxy URL | Keep for now |
| `GOSUMDB` | `MOXIESUMDB` | Checksum database URL | Keep for now |
| `GOPRIVATE` | `MOXIEPRIVATE` | Private module patterns | Keep for now |
| `GONOPROXY` | `MOXIENOPROXY` | No-proxy patterns | Keep for now |
| `GONOSUMDB` | `MOXIENOSUMDB` | No-sumdb patterns | Keep for now |
| `GOINSECURE` | `MOXIEINSECURE` | Insecure fetching patterns | Keep for now |
| `GOVCS` | `MOXIEVCS` | VCS control | Keep for now |
| `GOAUTH` | `MOXIEAUTH` | Authentication method | Keep for now |

### Security & Compliance Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOFIPS140` | `MOXIEFIPS140` | FIPS 140 mode | Keep for now |

### Experiment & Feature Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `GOEXPERIMENT` | `MOXIEEXPERIMENT` | Experimental features | Keep for now |

### Test Variables

| Current Name | Future Name | Description | Status |
|--------------|-------------|-------------|--------|
| `TESTGO_VERSION` | `TESTMOXIE_VERSION` | Test version override | ✅ Changed |
| `TESTGO_GOHOSTOS` | `TESTMOXIE_HOSTOS` | Test host OS | Keep for now |
| `TESTGO_GOHOSTARCH` | `TESTMOXIE_HOSTARCH` | Test host arch | Keep for now |

---

## Migration Strategy

### Phase 0 (Current)
**Status:** Documentation only
- ✅ Document all GO* variables
- ✅ Define MOXIE* equivalents
- ⏳ Keep GO* names in code
- ⏳ Update comments/documentation to mention future migration

### Phase 1-3 (Type System & Core Features)
**Status:** Pending
- Keep GO* variables unchanged
- Focus on language features
- Maintain stability

### Phase 4-6 (After Core Features)
**Status:** Future
- Add support for MOXIE* variables as aliases
- Both GO* and MOXIE* work simultaneously
- Update documentation to recommend MOXIE*
- Add deprecation warnings for GO* (soft)

### Phase 7 (Standard Library Updates)
**Status:** Future
- Update all stdlib code to check MOXIE* first, then GO* fallback
- Update error messages to suggest MOXIE*
- Comprehensive migration guide

### Phase 8-9 (Testing & Documentation)
**Status:** Future
- Full test coverage for both variable sets
- Migration tools to help users
- Clear deprecation timeline

### Phase 10+ (Future Releases)
**Status:** Future releases
- **Moxie 1.x:** Both supported, MOXIE* recommended
- **Moxie 2.x:** Strong warnings for GO*
- **Moxie 3.x:** GO* deprecated (but still functional with warnings)
- **Moxie 4.x:** GO* removed (breaking change)

---

## Implementation Plan (Future)

### Step 1: Add MOXIE* Support (Aliases)
```go
// In src/cmd/go/internal/cfg/cfg.go

func Getenv(key string) string {
	// Check for MOXIE* variant first
	moxieKey := strings.Replace(key, "GO", "MOXIE", 1)
	if val := os.Getenv(moxieKey); val != "" {
		return val
	}
	// Fall back to GO* variant
	return os.Getenv(key)
}
```

### Step 2: Update Documentation
- All help text mentions both GOROOT and MOXIEROOT
- Examples use MOXIEROOT
- Migration guide explains transition

### Step 3: Add Warnings
```go
func Getenv(key string) string {
	moxieKey := strings.Replace(key, "GO", "MOXIE", 1)
	moxieVal := os.Getenv(moxieKey)
	goVal := os.Getenv(key)

	if moxieVal != "" {
		return moxieVal
	}
	if goVal != "" {
		// Warn in Moxie 2.x+
		if warnDeprecated {
			fmt.Fprintf(os.Stderr, "Warning: %s is deprecated, use %s instead\n", key, moxieKey)
		}
		return goVal
	}
	return ""
}
```

### Step 4: Remove GO* Support (Breaking)
```go
func Getenv(key string) string {
	// Only check MOXIE* variant
	moxieKey := strings.Replace(key, "GO", "MOXIE", 1)
	return os.Getenv(moxieKey)
}
```

---

## Current Code Locations

### Main Configuration
- **File:** `src/cmd/go/internal/cfg/cfg.go`
- **Lines:** 32-33 (GOOS, GOARCH), 455-487 (all GO* variables)

### Environment Reading
- **File:** `src/cmd/go/internal/cfg/cfg.go`
- **Function:** `Getenv()`, `EnvOrAndChanged()`, `envOr()`

### GOROOT Handling
- **File:** `src/cmd/go/internal/cfg/cfg.go`
- **Function:** `SetGOROOT()` (line 238-285)
- **Function:** `findGOROOT()` (line 548+)

### Runtime Variables
- **File:** `src/runtime/extern.go`
- **Constants:** GOOS, GOARCH (line 376+)

---

## Code Comments to Update (Phase 0)

The following comments should be updated to note the future migration:

### In cfg.go
```go
// GOROOT is the root of the Moxie installation.
// Note: In future versions, this will be renamed to MOXIEROOT.
// For now, we keep GOROOT for backward compatibility.
var GOROOT string
```

### In runtime/extern.go
```go
// GOOS is the running program's operating system target:
// one of darwin, freebsd, linux, and so on.
// To view possible combinations of GOOS and GOARCH, run "moxie tool dist list".
// Note: In future versions, this may be renamed to MOXIEOS for consistency.
const GOOS string = goos.GOOS
```

---

## User Communication Strategy

### Phase 0 Documentation
- Note in README: "Moxie currently uses GO* environment variables for compatibility"
- Clear statement: "Future versions will introduce MOXIE* variables"
- Timeline: "Full migration expected by Moxie 4.0"

### Phase 4+ Documentation
- Migration guide with examples
- Shell scripts to help users update
- Clear deprecation timeline
- FAQ about the change

---

## File Updates for Phase 0

### Documentation Comments (Do Now)

#### src/cmd/go/internal/cfg/cfg.go
Add comment block at top of variable section:

```go
// Environment Variable Migration Note:
//
// Moxie uses GO* environment variables (GOROOT, GOPATH, etc.) for backward
// compatibility during the transition from Go. Future versions will introduce
// MOXIE* equivalents (MOXIEROOT, MOXIEPATH, etc.) and eventually deprecate
// the GO* variants.
//
// Timeline:
//   Phase 0-3: GO* only (current)
//   Phase 4+:  Both GO* and MOXIE* supported
//   Moxie 1.x: MOXIE* recommended, GO* still works
//   Moxie 2.x: GO* deprecated with warnings
//   Moxie 3.x: GO* shows strong warnings
//   Moxie 4.x: GO* removed (breaking change)
//
// See ENVIRONMENT-VARIABLES.md for full migration plan.
var (
	GOROOT string
	...
)
```

#### src/runtime/extern.go
Update existing comments:

```go
// GOOS is the running program's operating system target:
// one of darwin, freebsd, linux, and so on.
// To view possible combinations of GOOS and GOARCH, run "moxie tool dist list".
//
// Note: Future versions may rename this to MOXIEOS for consistency with
// Moxie branding. See ENVIRONMENT-VARIABLES.md for migration plan.
const GOOS string = goos.GOOS
```

---

## Testing Strategy (Future Phases)

### When MOXIE* Variables Added

1. **Unit Tests**
   - Test MOXIE* takes precedence over GO*
   - Test GO* still works when MOXIE* not set
   - Test both set (MOXIE* wins)

2. **Integration Tests**
   - Build with MOXIEROOT set
   - Build with GOROOT set (should still work)
   - Build with both (MOXIEROOT should win)

3. **Migration Tests**
   - Script to convert GO* to MOXIE* in env
   - Verify builds work identically

---

## Rationale for Deferring

### Why Not Change Now?

1. **Build System Not Ready**
   - Haven't tested build yet
   - Changing vars would break untested build
   - Need working build first

2. **Incremental Approach**
   - Phase 0: Branding and setup
   - Phase 1-3: Core language features
   - Phase 4+: Peripheral concerns like env vars

3. **Backward Compatibility**
   - Existing tools expect GOROOT
   - Editors/IDEs look for GOROOT
   - Build scripts use GOROOT
   - Users' environments have GOROOT

4. **Test First**
   - Need to verify build works
   - Need to validate current setup
   - Then can introduce changes incrementally

5. **Documentation Value**
   - Document strategy now
   - Implement later with confidence
   - Clear migration path

---

## Summary

### Phase 0.2 Completion Status

- ✅ Identified all 30+ GO* environment variables
- ✅ Documented MOXIE* equivalents
- ✅ Defined migration strategy
- ✅ Planned implementation timeline
- ✅ Added documentation comments (recommended)
- ⏳ Implementation deferred to Phase 4+

### Why This Is The Right Approach

1. **Pragmatic** - Don't break what we haven't tested
2. **Documented** - Clear plan for future
3. **Incremental** - Small, testable steps
4. **Compatible** - Maintains backward compatibility
5. **Clear Timeline** - Users know what to expect

---

## Next Steps

### Immediate (Phase 0)
1. ✅ Document all variables (this file)
2. ⏳ Add migration note comments to source (optional for Phase 0)
3. ⏳ Update README with env var compatibility note
4. ⏳ Move to Phase 0.3 (Build System)

### Phase 4+ (Future)
1. Implement MOXIE* alias support
2. Update documentation
3. Add deprecation warnings
4. Create migration tools

### Moxie 2.x+ (Future Releases)
1. Strengthen warnings
2. Update all examples
3. Create migration utilities

### Moxie 4.x (Major Version)
1. Remove GO* support
2. Only MOXIE* variables
3. Breaking change, well-communicated

---

## Conclusion

**Environment variables will remain as GO* for Phase 0-3.**

This is a strategic decision to maintain stability and backward compatibility while we focus on core language features. The migration path is clear, documented, and will be implemented in future phases when the build system is stable and tested.

**Phase 0.2 Status:** ✅ COMPLETE (Documentation & Strategy)

**Next:** Phase 0.3 - Build System Updates
