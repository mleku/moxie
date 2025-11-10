# Moxie Grammar Examples

This document shows how various Moxie language constructs map to the ANTLR grammar rules.

## Basic Program Structure

```moxie
package main

import "fmt"

func main() {
    fmt.Println("Hello, Moxie!")
}
```

**Parse Tree:**
```
sourceFile
├── packageClause: "package main"
├── importDecl
│   └── importSpec: "fmt"
└── topLevelDecl
    └── functionDecl: "main"
        └── block
            └── statement
                └── expressionStmt
                    └── callExpr: fmt.Println("Hello, Moxie!")
```

## Explicit Pointer Types

### Slice Literal

```moxie
s := &[]int{1, 2, 3}
```

**Grammar Rules:**
- `shortVarDecl`: `s := expressionList`
- `compositeLit`: `literalType literalValue`
- `literalType`: `sliceType` → `'*' '[' ']' elementType`
- `literalValue`: `'{' elementList '}'`

### Map Literal

```moxie
m := &map[string]int{
    "foo": 1,
    "bar": 2,
}
```

**Grammar Rules:**
- `compositeLit`: `mapType literalValue`
- `mapType`: `'*' 'map' '[' type_ ']' elementType`
- `literalValue`: keyed elements

### Channel Literal

```moxie
ch := &chan int{cap: 10}
```

**Grammar Rules:**
- `compositeLit`: `channelType literalValue`
- `channelType`: `'*' 'chan' elementType`
- `literalValue`: `'{' keyedElement '}'` where key is `cap`

## Const Declarations

### Simple Const

```moxie
const MaxSize = 100
const Pi = 3.14159
```

**Grammar Rules:**
- `constDecl`: `'const' constSpec`
- `constSpec`: `identifierList '=' expressionList`

### Const Pointer Types

```moxie
const Message = "immutable"
const Config = &map[string]int{
    "timeout": 30,
    "retries": 3,
}
```

**Grammar Rules:**
- `constSpec` with composite literal expression

### Const Type Modifier

```moxie
func process(data const string) {
    // data is immutable
}
```

**Grammar Rules:**
- `type_`: `'const' type_` (ConstType production)
- Used in parameter declarations

## Zero-Copy Type Coercion

### Basic Cast (Native Endian)

```moxie
bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
u32s := (*[]uint32)(bytes)
```

**Grammar Rules:**
- `conversion`: `'(' '*' '[' ']' type_ ')' '(' expression ')'`
- Production: `SliceCastExpr`

### Cast with Endianness

```moxie
// Little-endian
u32s := (*[]uint32, LittleEndian)(bytes)

// Big-endian
u32s := (*[]uint32, BigEndian)(bytes)
```

**Grammar Rules:**
- `conversion`: `'(' '*' '[' ']' type_ ',' endianness ')' '(' expression ')'`
- Production: `SliceCastEndianExpr`
- `endianness`: `'NativeEndian' | 'LittleEndian' | 'BigEndian'`

### Explicit Copy Cast

```moxie
// Copy with native endian
u32s := &(*[]uint32)(bytes)

// Copy with endianness
u32s := &(*[]uint32, LittleEndian)(bytes)
```

**Grammar Rules:**
- `conversion`: `'&' '(' '*' '[' ']' type_ ')' '(' expression ')'`
- Production: `SliceCastCopyExpr`
- `conversion`: `'&' '(' '*' '[' ']' type_ ',' endianness ')' '(' expression ')'`
- Production: `SliceCastCopyEndianExpr`

## String Operations

### Concatenation (using | operator)

```moxie
s1 := "Hello"
s2 := " World"
result := s1 | s2  // "Hello World"
```

**Grammar Rules:**
- `expression`: `expression '|' expression` (ConcatenationExpr)
- Works for strings (which are `*[]byte`)

**Cryptographic notation:**
The `|` operator follows standard cryptographic notation where `a | b` means concatenation. This makes Moxie code more familiar to cryptographers:

```moxie
// HMAC construction: HMAC(K, m) = H((K ⊕ opad) | H((K ⊕ ipad) | m))
innerHash := hash((key ^ ipad) | message)
outerHash := hash((key ^ opad) | innerHash)
```

### Mutation

```moxie
s := "hello"
s[0] = 'H'  // OK in Moxie
```

**Grammar Rules:**
- `assignment`: `expressionList assign_op expressionList`
- Left side: `indexExpr` → `primaryExpr '[' expression ']'`

### Array/Slice Concatenation

```moxie
a1 := &[]int{1, 2, 3}
a2 := &[]int{4, 5, 6}
a3 := a1 | a2  // &[]int{1, 2, 3, 4, 5, 6}
```

**Grammar Rules:**
- Same as string concatenation
- `expression`: `expression '|' expression` (ConcatenationExpr)

### Appending Elements (replaces append())

```moxie
// Before (Go):
s := []int{1, 2, 3}
s = append(s, 4, 5, 6)

// After (Moxie):
s := &[]int{1, 2, 3}
s = s | &[]int{4, 5, 6}

// Or with single element:
s = s | &[]int{7}
```

## Built-in Functions

**Note:** `append()` is NOT a built-in in Moxie. Use the `|` concatenation operator instead.

### clone() - Deep Copy

```moxie
s1 := &[]int{1, 2, 3}
s2 := clone(s1)  // Deep copy - new backing array
s2[0] = 99       // Does not affect s1
```

**Grammar Rules:**
- `callExpr`: `primaryExpr arguments`
- `operandName`: `IDENTIFIER` (where IDENTIFIER is "clone")
- `arguments`: `'(' expressionList ')'`

**Semantics:**
- Allocates new backing array
- Copies all elements
- Returns new slice with same length and capacity

### copy() - Overwrite Existing Slice

```moxie
dst := &[]int32{0, 0, 0, 0, 0}
src := &[]int32{1, 2, 3}

n := copy(dst, src)  // Returns 3 (number of elements copied)
// dst is now &[]int32{1, 2, 3, 0, 0}
```

**Grammar Rules:**
- `callExpr`: function call with two arguments
- `operandName`: "copy"

**Semantics:**
- Copies `min(len(dst), len(src))` elements from src to dst
- Overwrites destination elements (does not append)
- Returns number of elements copied (int64)
- Both slices must have same element type (or compatible cast)
- Does NOT allocate - modifies dst in place

**Type Requirements:**
```moxie
// Same type - OK
dst := &[]byte{0, 0, 0, 0}
src := &[]byte{1, 2, 3}
copy(dst, src)  // OK

// Different types - ERROR
dst := &[]int32{0, 0, 0}
src := &[]byte{1, 2, 3}
copy(dst, src)  // ERROR: type mismatch
```

### copy() with Type Casting

For numeric types, you can cast the source inline with optional endianness:

```moxie
// Copy uint32 slice to byte slice (native endian)
dst := &[]byte{0, 0, 0, 0, 0, 0, 0, 0}
src := &[]uint32{0x12345678, 0xABCDEF00}
copy(dst, (*[]byte)(src))  // Cast src to bytes
// dst is now the byte representation of src

// Copy with specific endianness
dst := &[]byte{0, 0, 0, 0}
src := &[]uint32{0x12345678}
copy(dst, (*[]byte, LittleEndian)(src))
// dst is now {0x78, 0x56, 0x34, 0x12}

copy(dst, (*[]byte, BigEndian)(src))
// dst is now {0x12, 0x34, 0x56, 0x78}
```

**Zero-copy casting rules apply:**
- Only fixed-width numeric types
- Length adjusted: `newLen = oldLen * sizeof(oldType) / sizeof(newType)`
- Alignment checks in debug mode

**Common pattern - serialize to network byte order:**
```moxie
// Serialize uint32 values to network byte order (big-endian)
values := &[]uint32{0x12345678, 0xABCDEF00, 0xDEADBEEF}
buffer := &[]byte{}
buffer = grow(buffer, len(values) * 4)

copy(buffer, (*[]byte, BigEndian)(values))
// buffer is now big-endian bytes suitable for network transmission
```

**Common pattern - deserialize from network:**
```moxie
// Parse network data (big-endian) to uint32 values
buffer := &[]byte{0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xEF, 0x00}
values := &[]uint32{0, 0}

// Cast buffer to uint32 with big-endian interpretation, then copy
copy(values, (*[]uint32, BigEndian)(buffer))
// values is now &[]uint32{0x12345678, 0xABCDEF00}
```

### copy() vs clone()

| Operation | `copy(dst, src)` | `clone(src)` |
|-----------|------------------|--------------|
| **Allocation** | No allocation | Allocates new backing array |
| **Destination** | Must provide dst | Returns new slice |
| **Overwrites** | Yes, up to min(len(dst), len(src)) | N/A |
| **Returns** | Number of elements copied | New slice |
| **Use case** | Overwrite existing buffer | Create independent copy |

```moxie
// Use copy() when you have a pre-allocated buffer
buffer := &[]byte{}
buffer = grow(buffer, 1024)
n := copy(buffer, data)  // Reuse buffer

// Use clone() when you need a new independent copy
backup := clone(original)
```

### grow()

```moxie
s := &[]int{}
s = grow(s, 100)  // Pre-allocate capacity for 100 elements
// len(s) is still 0, but cap(s) is now at least 100
```

**Grammar Rules:**
- Same as function call
- `operandName`: "grow"

**Semantics:**
- Pre-allocates capacity without changing length
- Avoids repeated reallocations
- If current capacity >= requested, may return same slice

### free()

```moxie
m := &map[string]int{}
// ... use m ...
free(m)  // Hint to GC to release memory
```

**Grammar Rules:**
- Function call with single argument

**Semantics:**
- Provides hint to garbage collector
- Does not guarantee immediate deallocation
- Slice/map becomes unusable after free() (undefined behavior if accessed)

## FFI Operations

### dlopen

```moxie
lib := dlopen("libc.so.6", RTLD_LAZY)
defer dlclose(lib)
```

**Grammar Rules:**
- `callExpr`: function call
- Arguments: string literal and identifier

### dlsym with Type Parameters

```moxie
printf := dlsym[func(*byte, ...any) int](lib, "printf")
```

**Grammar Rules:**
- `callExpr`: `primaryExpr arguments`
- `primaryExpr`: `operandName` with `typeArgs`
- `typeArgs`: `'[' typeList ']'`
- `typeList`: function type
- `arguments`: `lib`, string literal

## Generics

### Generic Function Declaration

```moxie
func Map[T, U any](slice *[]T, f func(T) U) *[]U {
    result := &[]U{}
    for _, item := range slice {
        result = append(result, f(item))
    }
    return result
}
```

**Grammar Rules:**
- `functionDecl`: `'func' IDENTIFIER typeParameters signature block`
- `typeParameters`: `'[' typeParameterDecl ',' typeParameterDecl ']'`
- `typeParameterDecl`: `IDENTIFIER typeConstraint`

### Generic Type Declaration

```moxie
type Stack[T any] struct {
    items *[]T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}
```

**Grammar Rules:**
- `typeSpec`: `IDENTIFIER typeParameters type_`
- `methodDecl`: receiver with type parameters

## Complex Examples

### Network Protocol Parsing

```moxie
func parsePacket(data *[]byte) Packet {
    // Zero-copy cast with big-endian (network byte order)
    fields := (*[]uint32, BigEndian)(data[0:12])

    return Packet{
        Magic:   fields[0],
        Version: uint16(fields[1] >> 16),
        Length:  uint16(fields[1]),
        Flags:   fields[2],
    }
}
```

**Grammar Rules:**
- Function with slice parameter
- Slice expression: `data[0:12]`
- Zero-copy cast with endianness
- Struct composite literal

### Const Config Pattern

```moxie
const DefaultConfig = &struct{
    Timeout  int32
    MaxConns int32
    Hosts    *[]string
}{
    Timeout:  30,
    MaxConns: 100,
    Hosts:    &[]string{"localhost", "api.example.com"},
}

func NewServer() *Server {
    cfg := clone(DefaultConfig)  // Mutable copy
    return &Server{config: cfg}
}
```

**Grammar Rules:**
- `constDecl` with anonymous struct type
- Nested composite literals
- `clone()` built-in call

### FFI with Callbacks

```moxie
import "unsafe"

type Callback = func(int32) int32

func setupCallback() {
    lib := dlopen("libfoo.so", RTLD_NOW)
    defer dlclose(lib)

    registerCB := dlsym[func(Callback)](lib, "register_callback")

    myCallback := func(x int32) int32 {
        return x * 2
    }

    registerCB(myCallback)
}
```

**Grammar Rules:**
- Type alias: `type Callback = functionType`
- Generic `dlsym` with function type parameter
- Function literal as argument

## Error Cases (Should Not Parse or Are Deprecated)

### Platform-Dependent int

```moxie
// ERROR: 'int' is not a valid type in Moxie
var count int
```

The lexer has no `INT` keyword (only `INT8`, `INT16`, `INT32`, `INT64`).

### append() Function

```moxie
// ERROR: append() is not supported in Moxie
s := append(s, 1, 2, 3)
```

Use `|` concatenation instead: `s = s | &[]int{1, 2, 3}`

### make() Function

```moxie
// ERROR: make() is not supported in Moxie
s := make([]int, 10)
```

Use `&[]int{}` with `grow()` instead: `s := grow(&[]int{}, 10)`

### Implicit Slice Type

```moxie
// ERROR: Slices must be explicit pointers in Moxie
s := []int{1, 2, 3}
```

Must use: `s := &[]int{1, 2, 3}`

(Note: The grammar accepts both for compatibility during migration, but semantically only `*[]T` is valid Moxie)

### Bitwise OR with |

```moxie
// ERROR: | is concatenation, not bitwise OR
flags := FLAG_A | FLAG_B
```

In Moxie, `|` is the concatenation operator. For bitwise operations, use Go's other operators or rethink your approach with bit manipulation.

## Testing with ANTLR TestRig

```bash
# Parse and view tree
grun Moxie sourceFile -tree example.x

# Parse and show GUI
grun Moxie sourceFile -gui example.x

# Show tokens
grun Moxie sourceFile -tokens example.x

# Parse specific rules
grun Moxie expression -tree
# Then type: bytes := (*[]uint32, LittleEndian)(data)
# Press Ctrl-D

grun Moxie type_ -tree
# Then type: const *[]byte
# Press Ctrl-D
```

## Common Grammar Patterns

### Type with Optional Const

```
type_
    : 'const' type_    # ConstType
    | typeLit          # LiteralType
    | typeName         # NamedType
    | ...
```

### Composite Literal with Pointer

```
compositeLit
    : literalType literalValue

literalType
    : sliceType      → '*' '[' ']' elementType
    | mapType        → '*' 'map' '[' type_ ']' elementType
    | channelType    → '*' 'chan' elementType
```

### Expression Precedence

```
expression (lowest precedence)
├── expression '||' expression            ← Logical OR
└── expression '&&' expression            ← Logical AND
    ├── expression rel_op expression      ← Comparison
    └── expression '|' expression         ← CONCATENATION (Moxie-specific)
        └── expression add_op expression  ← Arithmetic +/-/^
            └── expression mul_op expression  ← Arithmetic */%/etc
                └── unaryExpr
                    └── primaryExpr (highest precedence)
```

The `|` operator for concatenation is between comparisons and arithmetic operations, giving it appropriate precedence for composing strings and slices.

## Summary

The Moxie grammar provides:

1. **Complete Go compatibility** - All valid Go syntax (except removed features)
2. **Moxie extensions** - Explicit pointers, const types, zero-copy casts
3. **Clear structure** - Easy to traverse with visitor/listener patterns
4. **Type safety** - Strong typing for FFI operations
5. **Simplicity** - Removes ambiguous features (make, implicit references)

Use this grammar as the foundation for:
- **Lexical analysis** - Tokenization
- **Parsing** - Syntax tree construction
- **Semantic analysis** - Type checking (separate phase)
- **Code generation** - AST traversal for transpilation
- **IDE support** - Syntax highlighting, autocomplete
- **Static analysis** - Linting, error detection
