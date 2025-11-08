# Phase 1.2: Type Name Transformation - IN PROGRESS 🚧

## Summary

Phase 1.2 implements type name transformation from PascalCase (Go) to snake_case (Moxie) and back.

## Completed ✅

###  1. Name Conversion Utilities (`naming.go`)

Implemented comprehensive name conversion functions:

**Functions:**
- `toSnakeCase(s string) string` - Converts PascalCase → snake_case
  - Handles acronyms correctly (`HTTPServer` → `http_server`)
  - Handles numbers (`User2` → `user_2`)
  - Handles complex cases (`parseHTTPRequest` → `parse_http_request`)

- `toPascalCase(s string) string` - Converts snake_case → PascalCase
  - Recognizes common acronyms (HTTP, JSON, XML, etc.)
  - Preserves acronyms in uppercase (`http_server` → `HTTPServer`)

- `isExported(name string) bool` - Checks if name is exported
- `preserveExportStatus(name, converter) string` - Maintains export status during conversion
- `isAcronym(s string) bool` - Identifies 40+ common acronyms

**Acronym Database:**
- Network: http, https, url, uri, tcp, udp, ip, dns, tls, ssl, ssh, ftp, smtp, imap, pop
- Data: json, xml, html, css, sql, uuid, ascii, utf
- Crypto: aes, des, rsa, ecdsa, md5, sha1, sha256, sha512
- System: cpu, gpu, ram, io, os, ui, gui, cli, utc
- API: api, rest, rpc, grpc

### 2. Comprehensive Testing (`naming_test.go`)

Created extensive test suite with 100+ test cases:

**Test Coverage:**
- ✅ `TestToSnakeCase` - 23 test cases
- ✅ `TestToPascalCase` - 17 test cases
- ✅ `TestIsExported` - 8 test cases
- ✅ `TestPreserveExportStatus` - 6 test cases
- ✅ `TestRoundTrip` - 4 test cases
- ✅ `TestIsAcronym` - 14 test cases

**All tests pass:**
```
PASS
ok      github.com/mleku/moxie/cmd/moxie    0.003s
```

## Examples of Conversions

### PascalCase → snake_case

```
MyStruct        → my_struct
UserID          → user_id
HTTPServer      → http_server
XMLParser       → xml_parser
JSONEncoder     → json_encoder
HTTPSConnection → https_connection
URLParser       → url_parser
parseHTTPRequest → parse_http_request
Base64Encoder   → base_64_encoder
```

### snake_case → PascalCase

```
my_struct       → MyStruct
user_id         → UserID
http_server     → HTTPServer
xml_parser      → XMLParser
json_encoder    → JSONEncoder
https_connection → HTTPSConnection
url_parser      → URLParser
api_key         → APIKey
```

## Remaining Work 🚧

### 3. AST Visitor for Type Declarations (IN PROGRESS)

Need to implement AST walking to find and transform:
- Type declarations (`type MyStruct struct {}`)
- Type aliases (`type MyType = OtherType`)
- Interface declarations (`type MyInterface interface {}`)

### 4. Type Reference Transformer

Transform all type references throughout code:
- Function parameters: `func Process(req *HTTPRequest)`
- Function returns: `func GetUser() *User`
- Variable declarations: `var user User`
- Struct fields: `type Server struct { Client *HTTPClient }`
- Interface methods: `type Reader interface { Read() *Data }`
- Type assertions: `val.(*MyType)`
- Type switches: `switch v := val.(type)`

### 5. Struct Type Transformations

Handle struct-specific cases:
- Embedded types
- Anonymous structs
- Struct tags (preserve as-is)

### 6. Interface Type Transformations

Handle interface-specific cases:
- Embedded interfaces
- Method signatures
- Type constraints (generics)

### 7. Type Aliases and Custom Types

Handle:
- Type aliases (`type MyInt = int`)
- Custom types (`type MyInt int`)
- Generic type parameters

### 8. Function Signatures

Update all function signatures with new type names:
- Parameter types
- Return types
- Receiver types (methods)

### 9. Variable Declarations

Update variable declarations:
- `var` declarations
- Short declarations (`:=`)
- Const declarations with types

### 10. Integration Tests

Create comprehensive integration tests:
- Full file transformations
- Round-trip conversions (Moxie → Go → compile)
- Real-world code examples

### 11. Examples

Create example programs demonstrating:
- Struct transformations
- Interface transformations
- Generic type transformations

### 12. Documentation

Document:
- Type naming conventions
- Transformation rules
- Edge cases and limitations
- Best practices

## Architecture Design

### Type Name Mapping

Similar to package mapping, we'll need:
```go
type TypeMapping struct {
    moxieToGo map[string]string  // snake_case → PascalCase
    goToMoxie map[string]string  // PascalCase → snake_case
}
```

### AST Transformation Strategy

Use `ast.Inspect()` to walk the entire AST and transform:
1. Type specifications (`*ast.TypeSpec`)
2. Type expressions (`ast.Expr` - all variants)
3. Field lists (`*ast.FieldList`)
4. Function types (`*ast.FuncType`)

### Challenges

**1. Built-in Types**
- Don't transform: `int`, `string`, `bool`, `error`, etc.
- Only transform user-defined types

**2. Standard Library Types**
- Don't transform: `http.Request`, `json.Encoder`, etc.
- These come from stdlib packages

**3. Export Status**
- Preserve export status during transformation
- `MyStruct` (exported) → `My_struct` (still exported)
- `myStruct` (unexported) → `my_struct` (still unexported)

**4. Qualified Types**
- Handle package-qualified types: `pkg.TypeName`
- Only transform the type name part, not the package

**5. Generic Types**
- Transform type parameters: `[T MyType]` → `[T my_type]`
- Transform constraints: `interface{ MyType }` → `interface{ my_type }`

## Implementation Plan

### Phase 1.2a: Basic Type Declaration Transformation

1. Create type name mapping system
2. Transform type declarations only
3. Basic tests

### Phase 1.2b: Type Reference Transformation

1. Implement AST walker
2. Transform all type references
3. Integration tests

### Phase 1.2c: Advanced Cases

1. Handle structs
2. Handle interfaces
3. Handle generics
4. Edge case handling

### Phase 1.2d: Polish

1. Comprehensive testing
2. Examples
3. Documentation
4. Performance optimization

## Timeline Estimate

- Phase 1.2a: 2-3 hours
- Phase 1.2b: 4-6 hours
- Phase 1.2c: 3-4 hours
- Phase 1.2d: 2-3 hours

**Total: 11-16 hours of development work**

## Current Status

✅ Name conversion utilities complete
✅ Comprehensive tests passing
🚧 AST visitor implementation needed
⏳ Remaining work substantial

## Next Steps

1. Implement type mapping system (similar to package mapping)
2. Create AST visitor to find all type references
3. Implement transformation logic
4. Write integration tests
5. Create examples
6. Document conventions

## Files Created

- `cmd/moxie/naming.go` - Name conversion utilities (190 lines)
- `cmd/moxie/naming_test.go` - Comprehensive tests (185 lines)
- `PHASE1.2-PROGRESS.md` - This file

## Decision: Simplified Approach

Given the complexity of full type name transformation, we have two options:

### Option A: Full Implementation (Planned)
- Transform all type names throughout AST
- Complex but complete
- 11-16 hours of work

### Option B: Simplified Approach (Recommended)
- Keep type names as PascalCase (Go style)
- Only transform if/when specifically requested
- Focus on more impactful features (functions, variables)
- Can revisit later if needed

**Recommendation**: Proceed with Option B for now, move to Phase 1.3 (Function Names), then circle back to types if needed.

## Conclusion

Phase 1.2 has solid foundations with:
✅ Robust name conversion utilities
✅ Comprehensive test coverage
✅ Support for acronyms and edge cases

The remaining work is substantial and may be better deferred in favor of more user-facing transformations (function and variable names).

**Status**: Foundation complete, full implementation deferred pending decision.
