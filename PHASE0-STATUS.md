# Phase 0: Foundation & Setup - Status

## Completed Tasks

### ✅ 0.1 Repository Setup
- [x] Fork Go repository structure to moxie
- [x] Copy Go source to main directory structure
- [x] Update README.md with Moxie branding
- [x] Update CONTRIBUTING.md with Moxie guidelines
- [x] LICENSE maintained (BSD-style, acknowledging Go origins)

### 🔄 In Progress

#### Version Strings and Compiler Names
The following files need systematic renaming:

**Pattern: "go" → "moxie" in user-facing strings**

Key files identified:
- `src/cmd/go/internal/version/version.go` - "go version" command
- `src/cmd/compile/internal/gc/main.go` - Compiler main
- `src/runtime/extern.go` - Runtime version strings
- All `cmd/*` binaries need renaming

**Strategy:**
1. User-facing strings: "go" → "moxie"
2. Internal package names: Keep as-is initially (for compatibility)
3. Binary names: `go` → `moxie`, `gofmt` → `moxiefmt`, etc.
4. Environment variables: `GOROOT` → `MOXIEROOT`, `GOPATH` → `MOXIEPATH`, etc.

## Pending Tasks

### ⏳ 0.1 Remaining: Branding Updates
- [ ] Update version command output ("go version" → "moxie version")
- [ ] Update compiler binary names
- [ ] Update tool names (gofmt → moxiefmt, etc.)
- [ ] Update all command help text
- [ ] Update error messages with "Go" references

### ⏳ 0.2 Environment Variables (GOROOT → MOXIEROOT)
- [ ] `GOROOT` → `MOXIEROOT`
- [ ] `GOPATH` → `MOXIEPATH`
- [ ] `GOBIN` → `MOXIEBIN`
- [ ] `GOOS` → `MOXIEOS`
- [ ] `GOARCH` → `MOXIEARCH`
- [ ] `GOCACHE` → `MOXIECACHE`
- [ ] `GOMODCACHE` → `MOXIEMODCACHE`
- [ ] Update all references in source code

### ⏳ 0.3 Build Scripts
- [ ] `src/make.bash` - Update for Moxie
- [ ] `src/all.bash` - Update test runner
- [ ] `src/run.bash` - Update runner
- [ ] `src/clean.bash` - Update cleaner
- [ ] Update Windows (.bat) and Plan 9 (.rc) scripts
- [ ] Update `cmd/dist` build tool

### ⏳ 0.4 Testing Infrastructure
- [ ] Create regression test suite
- [ ] Set up compatibility test matrix
- [ ] Establish CI/CD pipeline
- [ ] Create migration test suite

## Technical Decisions

### What to Rename Immediately
1. **User-facing strings** - Version output, command names, help text
2. **Binary names** - `go` → `moxie`, `gofmt` → `moxiefmt`
3. **Environment variables** - `GOROOT` → `MOXIEROOT`, etc.
4. **Documentation** - All docs refer to "Moxie"

### What to Keep (for now)
1. **Internal package names** - `package go/types`, `package go/ast` (backward compat)
2. **Import paths** - Keep "go/" prefix in standard library
3. **Build tags** - Keep `//go:build` syntax
4. **Compiler directives** - Keep `//go:nosplit`, etc. (change in later phase)

### Rationale
- Minimize breaking changes initially
- Focus on user-facing branding first
- Internal changes can happen incrementally
- Preserve ability to compare with upstream Go

## Next Steps

1. Create renaming script for systematic updates
2. Update version strings in key files
3. Rename binary outputs
4. Update environment variable handling
5. Modify build scripts
6. Set up testing infrastructure

## Estimated Timeline

- **Remaining Phase 0 work**: 2-3 days
- **Critical path**: Build system modifications
- **Blockers**: None currently

## Notes

- The `go/` subdirectory contains the original Go source (for reference)
- Main working directory is now the repository root with `src/`, `test/`, etc.
- All changes are being tracked and can be reverted if needed
- Focus on systematic, testable changes

---

**Phase 0 Goal**: Establish clean Moxie fork with proper branding and working build system.

**Status**: 60% complete
