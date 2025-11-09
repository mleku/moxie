# Phase 3: String Mutability Implementation Plan

**Goal**: Implement `string = *[]byte` transformation

## Overview

In Moxie, strings are mutable byte slices. The `string` type is an alias for `*[]byte`.

```go
// Moxie
s := "hello"     // string (which is *[]byte)
s[0] = 'H'      // OK - mutable
s = s + " world" // OK - concatenation

// Transpiled Go
s := &[]byte{'h', 'e', 'l', 'l', 'o'}  // *[]byte
(*s)[0] = 'H'                          // Mutable
*s = append(*s, []byte{' ', 'w', 'o', 'r', 'l', 'd'}...)  // Concatenation
```

## Required Transformations

### 1. String Type Alias
**Goal**: Transform `string` type to `*[]byte` in declarations

**Examples**:
```go
// Moxie -> Go
var s string                 -> var s *[]byte
func foo(s string) string    -> func foo(s *[]byte) *[]byte
type User struct {           -> type User struct {
    Name string              ->     Name *[]byte
}                            -> }
```

**Implementation**: Add type transformation in `syntax.go`

### 2. String Literals
**Goal**: Transform string literals to byte slice literals

**Examples**:
```go
// Moxie -> Go
"hello"                      -> &[]byte{'h', 'e', 'l', 'l', 'o'}
""                           -> &[]byte{}
"Hello\nWorld"               -> &[]byte{'H', 'e', 'l', 'l', 'o', '\n', 'W', 'o', 'r', 'l', 'd'}
```

**Implementation**: AST transformation of `ast.BasicLit` with `Kind == token.STRING`

### 3. String Concatenation
**Goal**: Transform `+` operator on strings to append

**Examples**:
```go
// Moxie -> Go
s1 + s2                      -> &append(*s1, *s2...)
s + " world"                 -> &append(*s, []byte{' ', 'w', 'o', 'r', 'l', 'd'}...)
"Hello" + " " + "World"      -> &append(append([]byte{'H', 'e', 'l', 'l', 'o'}, ' '), []byte{'W', 'o', 'r', 'l', 'd'}...)
```

**Implementation**: Transform `ast.BinaryExpr` with `Op == token.ADD` and string operands

### 4. String Conversions
**Goal**: Handle `string()` conversions

**Examples**:
```go
// Moxie -> Go
string(b)                    -> &b  (if b is []byte)
string(i)                    -> strconv.Itoa(i) converted to *[]byte
```

**Implementation**: Transform `ast.CallExpr` where `Fun` is string type

### 5. String Comparison
**Goal**: Transform string comparison operators

**Examples**:
```go
// Moxie -> Go
s1 == s2                     -> bytes.Equal(*s1, *s2)
s1 != s2                     -> !bytes.Equal(*s1, *s2)
s1 < s2                      -> bytes.Compare(*s1, *s2) < 0
```

**Implementation**: Transform `ast.BinaryExpr` with comparison operators and string operands
**Note**: Need to import `bytes` package

## Implementation Steps

### Step 1: Type Transformation ✅ COMPLETE
- Transform `string` type identifier to `*[]byte` in:
  - Variable declarations ✅
  - Function parameters ✅
  - Function returns ✅
  - Struct fields ✅
  - Type assertions ✅
  - Type conversions ✅

### Step 2: String Literal Transformation ✅ COMPLETE
- Convert `"text"` to `&[]byte{'t', 'e', 'x', 't'}` ✅
- Handle escape sequences (\n, \t, \r, \\, \", \') ✅
- Handle empty strings ✅
- Handle multi-line strings (raw string literals with backticks) ✅

### Step 3: String Concatenation ✅ COMPLETE
- Transform `s1 + s2` to `moxie.Concat(s1, s2)` ✅
- Handle multiple concatenations (s1 + s2 + s3) ✅
- Handle mixed literal and variable concatenations ✅
- Multi-pass transformation for chained operations ✅

### Step 4: String Comparison ✅ COMPLETE
- Transform `==`, `!=` to use `bytes.Equal` ✅
- Transform `<`, `>`, `<=`, `>=` to use `bytes.Compare` ✅
- Inject `bytes` package import ✅

### Step 5: String Conversions ⏸️ DEFERRED
- Handle `string([]byte)` -> direct assignment (deferred to later phase)
- Handle `string(int)` -> use strconv (deferred to later phase)
- Handle other conversions (deferred to later phase)

## Testing Plan

### Test Files

1. **test_string_type.mx** ✅ PASSING
   - Function parameters/returns with string type
   - Basic concatenation

2. **test_string_literals.mx** ✅ PASSING
   - String literal creation
   - Empty strings
   - Strings with escape sequences (\n, \t, \r, \\, \", \')
   - Raw string literals (backticks)

3. **test_string_mutation.mx** ✅ PASSING
   - String indexing
   - String modification (mutable)
   - Append operations
   - Slice operations
   - Length operations

4. **test_string_concat.mx** ✅ PASSING
   - String concatenation with +
   - Multiple concatenations (chained)
   - Mixed literal/variable concat
   - Incremental concatenation in assignments

5. **test_string_comparison.mx** ✅ PASSING
   - Equality tests (==, !=)
   - Ordering tests (<, >, <=, >=)
   - Comparison with identical strings

6. **test_string_edge_cases.mx** ✅ PASSING
   - Empty string operations
   - Concatenation with empty strings
   - Special characters (tabs, newlines, quotes, backslashes)
   - Unicode support
   - Repeated concatenation
   - Comparison with empty strings

## Known Limitations

1. ~~**Raw string literals** (backtick strings)~~ ✅ RESOLVED - now fully supported
2. **String formatting** - fmt.Sprintf output is byte array representation, not string
3. **Standard library** - functions expecting `string` need wrappers or conversion
4. **Const strings** - deferred to Phase 4 (MMU protection)
5. **String conversions** - `string(int)`, `string([]rune)` etc. deferred to later phase

## Implementation Details

### Multi-Pass Transformation
The transformation uses a multi-pass approach (up to 10 passes) to handle chained string concatenations:
- Pass 1: Transform types, literals, comparisons, and first-level concatenations
- Subsequent passes: Continue transforming concatenations until no more changes occur
- This ensures expressions like `s1 + s2 + s3` are fully transformed to nested `moxie.Concat()` calls

### Runtime Functions
- **moxie.Concat(s1, s2 *[]byte) *[]byte**: Concatenates two byte slices and returns pointer
- Uses Go's `append()` internally for efficiency
- Handles nil/empty slices correctly

### Import Injection
- Automatically adds `bytes` package when string comparisons are used
- Automatically adds `moxie` runtime package when concatenation is used

## Dependencies

- Phase 2 complete (syntax transformations) ✅
- Import injection infrastructure ✅
- AST traversal with replacement ✅
- Multi-pass transformation support ✅

## Success Criteria

- [x] All type declarations transform `string` to `*[]byte` ✅
- [x] String literals become `&[]byte{...}` ✅
- [x] String concatenation works with `+` ✅
- [x] String comparison works correctly ✅
- [x] All test files pass ✅
- [x] Documentation updated ✅

## PHASE 3 STATUS: ✅ COMPLETE

All core string mutability features are implemented and tested:
- Type transformation
- Literal transformation (including raw strings and escape sequences)
- Concatenation (including chained operations)
- Comparison (all operators)
- Mutation (indexing, modification, slicing)
- Edge cases (empty strings, Unicode, special characters)
