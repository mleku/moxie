# The Moxie Programming Language Specification

**Version 0.1.0**

Based on the Go Programming Language Specification with targeted improvements.

---

## Introduction

Moxie is a programming language derived from Go that addresses specific design inconsistencies while preserving Go's philosophy of simplicity and clarity. Moxie maintains full compatibility with Go's naming conventions (PascalCase for exported symbols, camelCase for private symbols) and syntax, with focused changes to improve explicitness and consistency.

### Design Principles

1. **Simplicity** - Easy to learn, easy to read, easy to understand
2. **Explicitness** - No hidden behavior, obvious semantics
3. **Consistency** - Uniform rules, fewer special cases
4. **Performance** - Fast compilation, fast execution, predictable behavior
5. **Pragmatism** - Solve real problems, avoid theoretical purity

### Differences from Go

Moxie differs from Go in exactly four areas:

1. **Explicit reference semantics** - `*[]T`, `*map[K]V`, `*chan T` instead of implicit references
2. **Mutable strings** - `string` is an alias for `*[]byte`
3. **True immutability** - `const` works for all types with MMU protection
4. **Native FFI** - Built-in `dlopen`/`dlsym` instead of CGO

Everything else remains identical to Go.

---

## 1. Notation

The syntax is specified using a Extended Backus-Naur Form (EBNF):

```
Production  = production_name "=" [ Expression ] "." .
Expression  = Term { "|" Term } .
Term        = Factor { Factor } .
Factor      = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```

Productions are expressions constructed from terms and the following operators:

```
|   alternation
()  grouping
[]  option (0 or 1 times)
{}  repetition (0 to n times)
```

---

## 2. Source Code Representation

### 2.1 Characters

Moxie programs are written in Unicode using UTF-8 encoding.

```
newline        = /* Unicode code point U+000A */ .
unicode_char   = /* an arbitrary Unicode code point except newline */ .
unicode_letter = /* a Unicode code point classified as "Letter" */ .
unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
```

### 2.2 Letters and Digits

The underscore character `_` (U+005F) is considered a lowercase letter.

```
letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
```

---

## 3. Lexical Elements

### 3.1 Comments

Comments are identical to Go:

- Line comments start with `//`
- General comments start with `/*` and end with `*/`

### 3.2 Tokens

Tokens form the vocabulary of Moxie.

```
Token = identifier | keyword | operator | punctuation | literal .
```

### 3.3 Semicolons

Like Go, Moxie uses semicolons for statement termination, with automatic insertion rules.

### 3.4 Identifiers

```
identifier = letter { letter | unicode_digit } .
```

**Naming Conventions (MANDATORY):**

- **Exported** (visible outside package): Start with Unicode uppercase letter (PascalCase)
  ```go
  type User struct { }          // Exported type
  func ProcessData() { }        // Exported function
  var DefaultTimeout = 30       // Exported variable
  const MaxConnections = 100    // Exported constant
  ```

- **Unexported** (package-private): Start with Unicode lowercase letter (camelCase)
  ```go
  type internalState struct { } // Unexported type
  func parseRequest() { }       // Unexported function
  var errorCount = 0            // Unexported variable
  const bufferSize = 1024       // Unexported constant
  ```

**No snake_case. No kebab-case. Only PascalCase and camelCase.**

### 3.5 Keywords

Moxie has the same keywords as Go:

```
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

### 3.6 Operators and Punctuation

```
+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
     &^          &^=          ~
```

---

## 4. Types

The Moxie type system is simpler than Go's, with explicit pointers for all reference types.

### 4.1 Type Categories

```
Type      = TypeName | TypeLit | "(" Type ")" .
TypeName  = identifier | QualifiedIdent .
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
            SliceType | MapType | ChannelType .
```

**Type categories:**

1. **Value types** - Stored directly, passed by copy
   - Boolean types: `bool`
   - Numeric types: `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64`, `complex128`
   - String type: `string` (alias for `*[]byte`)
   - Array types: `[N]T`
   - Struct types: `struct { ... }`

2. **Pointer types** - Reference to a value
   - Explicit pointers: `*T`
   - Slice pointers: `*[]T`
   - Map pointers: `*map[K]V`
   - Channel pointers: `*chan T`

3. **Const types** - Immutable, MMU-protected
   - Any type prefixed with `const`

### 4.2 Boolean Types

```
bool
```

### 4.3 Numeric Types

```
int8        signed  8-bit integers (-128 to 127)
int16       signed 16-bit integers (-32768 to 32767)
int32       signed 32-bit integers (-2147483648 to 2147483647)
int64       signed 64-bit integers (-9223372036854775808 to 9223372036854775807)

uint8       unsigned  8-bit integers (0 to 255)
uint16      unsigned 16-bit integers (0 to 65535)
uint32      unsigned 32-bit integers (0 to 4294967295)
uint64      unsigned 64-bit integers (0 to 18446744073709551615)

float32     IEEE-754 32-bit floating-point numbers
float64     IEEE-754 64-bit floating-point numbers

complex64   complex numbers with float32 real and imaginary parts
complex128  complex numbers with float64 real and imaginary parts

byte        alias for uint8
rune        alias for int32
```

**IMPORTANT: No platform-dependent `int` or `uint` types.**

Rationale: Platform-dependent integer types cause:
- Serialization ambiguity
- Cross-platform bugs
- Hidden overflow issues

Use explicit bit-width types instead:
```go
// Good
var count int32         // Always 32 bits on all platforms
var largeIndex int64    // Always 64 bits on all platforms

// Bad (not available in Moxie)
var count int           // Platform-dependent (REMOVED)
var size uint           // Platform-dependent (REMOVED)
```

### 4.4 String Types

```
string
```

**String is an alias for `*[]byte`:**

```go
type string = *[]byte
```

Strings are mutable byte slices:

```go
s := "hello"     // Type: string (which is *[]byte)
s[0] = 'H'      // OK - mutable
s = s + " world" // OK - concatenation
```

**UTF-8 Encoding:**

Strings are UTF-8 encoded byte sequences:
- Indexing returns bytes: `s[i]` is `byte`
- Range iteration yields runes: `for i, r := range s { /* r is rune */ }`
- Length is byte count: `len(s)` returns number of bytes

**Immutability through const:**

```go
const greeting = "Hello"  // Immutable, MMU-protected
// greeting[0] = 'h'       // Compile error
```

### 4.5 Array Types

```
ArrayType   = "[" ArrayLength "]" ElementType .
ArrayLength = Expression .
ElementType = Type .
```

Examples:
```go
[32]byte
[2*N] struct { x, y int32 }
[1000]*float64
[3][5]int32
```

Arrays are value types and passed by copy.

### 4.6 Slice Types

**Slices are ALWAYS pointers:**

```
SliceType = "*" "[" "]" ElementType .
```

```go
*[]byte
*[]int32
*[][]*string
```

**Declaration:**

```go
// Explicit pointer allocation
s := &[]int32{1, 2, 3, 4, 5}

// Empty slice (nil)
var s *[]int32

// Empty slice (allocated)
s := &[]int32{}
```

**Zero value:**

The zero value of a slice type is `nil`.

**Nil slices:**

```go
var s *[]int32       // s == nil
if s == nil {        // OK - explicit nil check
    s = &[]int32{}
}
```

**Slice operations:**

```go
s := &[]int32{0, 1, 2, 3, 4}
s[2]          // Indexing: 2
s[1:4]        // Slicing: &[]int32{1, 2, 3}
len(s)        // Length: int64
cap(s)        // Capacity: int64
```

**No `make()` function - use composite literals:**

```go
// Before (Go)
s := make([]int32, 10, 20)

// After (Moxie)
s := &[]int32{}         // Empty slice
s = grow(s, 20)         // Pre-allocate capacity
for i := 0; i < 10; i++ {
    s = append(s, 0)    // Append zero values
}
```

### 4.7 Struct Types

```
StructType    = "struct" "{" { FieldDecl ";" } "}" .
FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
EmbeddedField = [ "*" ] TypeName .
Tag           = string_lit .
```

Examples:
```go
struct {}

struct {
    x, y int32
}

struct {
    T1        // Embedded field
    *T2       // Embedded pointer field
    P.T3      // Embedded qualified type
    *P.T4     // Embedded qualified pointer type
    x, y int32
}

struct {
    name string `json:"name"`  // Field with tag
}
```

**Exported vs unexported fields:**

```go
type User struct {
    ID       int32   // Exported (PascalCase)
    Name     string  // Exported
    email    string  // Unexported (camelCase)
    verified bool    // Unexported
}
```

### 4.8 Pointer Types

```
PointerType = "*" BaseType .
BaseType    = Type .
```

```go
*int32
*map[string]*User
*[128]byte
```

**All slices, maps, and channels are pointers:**

```go
*[]byte       // Slice (pointer)
*map[K]V      // Map (pointer)
*chan T       // Channel (pointer)
```

### 4.9 Function Types

```
FunctionType   = "func" Signature .
Signature      = Parameters [ Result ] .
Result         = Parameters | Type .
Parameters     = "(" [ ParameterList [ "," ] ] ")" .
ParameterList  = ParameterDecl { "," ParameterDecl } .
ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
```

Examples:
```go
func()
func(x int32) int32
func(a, _ int32, z float32) bool
func(a, b int32, z float64) (bool)
func(prefix string, values ...int32)
func(a, b int32, z float32) (success bool)
func(int32, int64, float32) (float32, *error)
```

### 4.10 Interface Types

```
InterfaceType  = "interface" "{" { MethodSpec ";" } "}" .
MethodSpec     = MethodName Signature | InterfaceTypeName .
MethodName     = identifier .
```

Examples:
```go
interface {
    Read(p *[]byte) (n int64, err *error)
    Write(p *[]byte) (n int64, err *error)
    Close() *error
}
```

**Method naming:**
- Exported methods: PascalCase (`Read`, `Write`, `Close`)
- Unexported methods: camelCase (rare, but allowed)

### 4.11 Map Types

**Maps are ALWAYS pointers:**

```
MapType     = "*" "map" "[" KeyType "]" ElementType .
KeyType     = Type .
ElementType = Type .
```

```go
*map[string]int32
*map[int32]*User
*map[string]*[]byte
```

**Declaration:**

```go
// Explicit allocation
m := &map[string]int32{"one": 1, "two": 2}

// Empty map (nil)
var m *map[string]int32

// Empty map (allocated)
m := &map[string]int32{}
```

**Zero value:**

The zero value of a map type is `nil`.

**Map operations:**

```go
m := &map[string]int32{"one": 1}
m["one"]                // Index: 1
m["two"] = 2            // Assignment
delete(m, "one")        // Deletion
_, ok := m["three"]     // Membership test
len(m)                  // Length: int64
```

**No `make()` function:**

```go
// Before (Go)
m := make(map[string]int32, 100)

// After (Moxie)
m := &map[string]int32{}
```

### 4.12 Channel Types

**Channels are ALWAYS pointers:**

```
ChannelType = "*" "chan" [ "<-" | "<-" ] ElementType .
```

```go
*chan int32
*chan *[]byte
*<-chan int32   // Receive-only
*chan<- int32   // Send-only
```

**Declaration:**

```go
// Unbuffered channel
ch := &chan int32{}

// Buffered channel (pre-allocate buffer)
ch := &chan int32{}
ch = grow(ch, 10)  // Buffer capacity 10
```

**Zero value:**

The zero value of a channel type is `nil`.

**Channel operations:**

```go
ch := &chan int32{}
ch <- x        // Send
x = <-ch       // Receive
x, ok = <-ch   // Receive with test
close(ch)      // Close
len(ch)        // Buffer length
cap(ch)        // Buffer capacity
```

### 4.13 Const Types

**Const applies to any type:**

```
ConstType = "const" Type .
```

```go
const MaxSize = 100                // const int32
const Message = "Hello"            // const string (const *[]byte)
const Config = &map[string]int32{  // const *map[string]int32
    "timeout": 30,
}
```

**Const is enforced by MMU:**

Constant values are placed in read-only memory (`.rodata` section) with MMU page protection. Any attempt to modify causes a hardware exception (SIGSEGV).

```go
const secret = "password123"
// secret[0] = 'P'  // Compile error

// Cannot cast away const
// p := (*[]byte)(secret)  // Compile error

// Must clone to get mutable copy
mutable := clone(secret)
mutable[0] = 'P'  // OK
```

---

## 5. Properties of Types and Values

### 5.1 Type Identity

Two types are either identical or different.

**Named types:**
```go
type Age int32
type Height int32

// Age and Height are different types
// Age and int32 are different types
```

**Unnamed types:**

```go
// Identical
*[]byte
*[]byte

// Different
*[]int32
*[]int64
```

### 5.2 Assignability

A value `x` is assignable to a variable of type `T` ("x is assignable to T") if one of the following conditions applies:

1. `x`'s type is identical to `T`
2. `x`'s type `V` and `T` have identical underlying types and at least one of `V` or `T` is not a named type
3. `T` is an interface type and `x` implements `T`
4. `x` is a bidirectional channel value, `T` is a channel type, `x`'s type `V` and `T` have identical element types, and at least one of `V` or `T` is not a named type
5. `x` is the predeclared identifier `nil` and `T` is a pointer, function, slice, map, channel, or interface type
6. `x` is an untyped constant representable by a value of type `T`

### 5.3 Representability

A constant `x` is representable by a value of type `T` if:

- `x` is in the set of values determined by `T`
- For integer types: `x` can be represented by a value of type `T`
- For floating-point types: `x` can be rounded to type `T` without overflow
- For complex types: both `real(x)` and `imag(x)` are representable

---

## 6. Blocks

```
Block = "{" StatementList "}" .
StatementList = { Statement ";" } .
```

Blocks nest and influence scoping.

---

## 7. Declarations and Scope

```
Declaration   = ConstDecl | TypeDecl | VarDecl .
TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
```

### 7.1 Const Declarations

```
ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .
IdentifierList = identifier { "," identifier } .
ExpressionList = Expression { "," Expression } .
```

```go
const Pi = 3.14159
const MaxSize int32 = 1024

const (
    StatusOK    = 200
    StatusError = 500
)

// Const with pointer types (MMU-protected)
const DefaultConfig = &map[string]int32{
    "timeout": 30,
    "retries": 3,
}
```

**Naming:**
- Exported constants: PascalCase (`Pi`, `MaxSize`)
- Unexported constants: camelCase (`bufferSize`, `maxRetries`)

### 7.2 Type Declarations

```
TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
TypeSpec = AliasDecl | TypeDef .
```

**Type definition:**
```go
type User struct {
    ID   int32
    Name string
}
```

**Type alias:**
```go
type MyInt = int32  // MyInt is identical to int32
type string = *[]byte  // Built-in alias
```

**Naming:**
- Exported types: PascalCase (`User`, `HTTPServer`)
- Unexported types: camelCase (`internalState`, `requestData`)

### 7.3 Variable Declarations

```
VarDecl     = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
VarSpec     = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
```

```go
var count int32
var isReady bool
var name, email string

var (
    users *[]User
    cache *map[string]*User
)

// Short variable declaration
count := 0
users := &[]User{}
```

**Naming:**
- Exported variables: PascalCase (`DefaultTimeout`, `MaxConnections`)
- Unexported variables: camelCase (`errorCount`, `requestCounter`)

### 7.4 Short Variable Declarations

```
ShortVarDecl = IdentifierList ":=" ExpressionList .
```

```go
i := 0
i, j := 0, 10
users := &[]User{}
```

### 7.5 Function Declarations

```
FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
FunctionName = identifier .
FunctionBody = Block .
```

```go
func Add(x, y int32) int32 {
    return x + y
}

func ProcessData(data *[]byte) (*Result, *error) {
    // ...
    return result, nil
}

// Exported function (PascalCase)
func ServeHTTP(w ResponseWriter, r *Request) {
    // ...
}

// Unexported function (camelCase)
func parseRequest(data *[]byte) *Request {
    // ...
}
```

**Naming:**
- Exported functions: PascalCase (`Add`, `ServeHTTP`, `NewUser`)
- Unexported functions: camelCase (`parseRequest`, `validateInput`)

### 7.6 Method Declarations

```
MethodDecl = "func" Receiver MethodName Signature [ FunctionBody ] .
Receiver   = Parameters .
```

```go
type User struct {
    ID   int32
    Name string
}

// Exported method
func (u *User) GetID() int32 {
    return u.ID
}

// Unexported method
func (u *User) validate() bool {
    return u.ID > 0
}
```

**Naming:**
- Exported methods: PascalCase (`GetID`, `String`, `Error`)
- Unexported methods: camelCase (`validate`, `parseHeaders`)

---

## 8. Expressions

### 8.1 Operands

```
Operand     = Literal | OperandName | "(" Expression ")" .
Literal     = BasicLit | CompositeLit | FunctionLit .
BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
OperandName = identifier | QualifiedIdent .
```

### 8.2 Composite Literals

```
CompositeLit  = LiteralType LiteralValue .
LiteralType   = StructType | ArrayType | "[" "..." "]" ElementType |
                SliceType | MapType | TypeName .
LiteralValue  = "{" [ ElementList [ "," ] ] "}" .
ElementList   = KeyedElement { "," KeyedElement } .
KeyedElement  = [ Key ":" ] Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .
```

**Slices (with explicit pointer):**
```go
&[]int32{1, 2, 3, 4}
&[]string{"hello", "world"}
```

**Maps (with explicit pointer):**
```go
&map[string]int32{
    "one": 1,
    "two": 2,
}
```

**Channels (with explicit pointer):**
```go
&chan int32{}  // Unbuffered
```

**Structs:**
```go
User{ID: 1, Name: "Alice"}
&User{ID: 1, Name: "Alice"}  // Pointer to struct
```

### 8.3 Selectors

```
selector.field
selector.method
```

**Auto-dereferencing for pointer types:**

```go
users := &[]User{{ID: 1}}
users[0].ID          // Auto-dereferences, same as (*users)[0].ID

config := &map[string]int32{"timeout": 30}
config["timeout"]    // Auto-dereferences, same as (*config)["timeout"]

ch := &chan int32{}
<-ch                 // Auto-dereferences, same as <-(*ch)
```

### 8.4 Index Expressions

```
a[x]
```

- Array: `a[x]` accesses element at index `x`
- Slice pointer: `a[x]` accesses element (auto-dereferences)
- Map pointer: `a[x]` accesses value (auto-dereferences)
- String: `s[x]` accesses byte at index `x`

**Index type:**
- For arrays/slices/strings: any integer type (auto-converted to int64)
- For maps: must be comparable to key type

### 8.5 Slice Expressions

```
a[low : high]
a[low : high : max]
```

Creates a new slice from an array, slice, or string.

```go
a := [5]int32{1, 2, 3, 4, 5}
s := a[1:4]  // &[]int32{2, 3, 4}

b := &[]byte{0, 1, 2, 3, 4}
s := b[1:3]  // &[]byte{1, 2}
```

### 8.6 Calls

```
f(a1, a2, … an)
```

**Special forms:**

```go
append(s, x)         // Append to slice
clone(x)             // Deep copy
grow(s, n)           // Pre-allocate capacity
clear(x)             // Reset slice/map
free(x)              // Release memory (GC hint)
delete(m, k)         // Remove map key
close(ch)            // Close channel
```

### 8.7 Operators

**Arithmetic:**
```
+    sum                    integers, floats, complex, strings
-    difference             integers, floats, complex
*    product                integers, floats, complex
/    quotient               integers, floats, complex
%    remainder              integers

&    bitwise AND            integers
|    bitwise OR             integers
^    bitwise XOR            integers
&^   bit clear (AND NOT)    integers

<<   left shift             integer << unsigned integer
>>   right shift            integer >> unsigned integer
```

**Comparison:**
```
==    equal
!=    not equal
<     less
<=    less or equal
>     greater
>=    greater or equal
```

**Logical:**
```
&&    conditional AND
||    conditional OR
!     NOT
```

**Address:**
```
&x   address of x
*p   dereference p
```

**Channel:**
```
<-   send/receive
```

**Concatenation (for slices and strings):**

The `+` operator concatenates slices and strings:

```go
s1 := "hello"
s2 := " world"
s3 := s1 + s2  // "hello world"

a := &[]int32{1, 2, 3}
b := &[]int32{4, 5, 6}
c := a + b     // &[]int32{1, 2, 3, 4, 5, 6}
```

**Semantics:**
- Always allocates a new slice
- Does not mutate operands
- Equivalent to: `clone(a)` then `append` all elements from `b`

---

## 9. Statements

### 9.1 Assignment

```
Assignment = ExpressionList assign_op ExpressionList .
assign_op  = [ add_op | mul_op ] "=" .
```

```go
x = 1
*p = f()
a[i] = 23
```

### 9.2 If Statements

```
IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
```

```go
if x > 0 {
    return x
}

if err := process(); err != nil {
    return err
}
```

### 9.3 Switch Statements

```
SwitchStmt = ExprSwitchStmt | TypeSwitchStmt .
```

**Expression switch:**
```go
switch tag {
case 1:
    fmt.Println("one")
case 2, 3:
    fmt.Println("two or three")
default:
    fmt.Println("other")
}
```

**Type switch:**
```go
switch v := x.(type) {
case int32:
    // v is int32
case string:
    // v is string
default:
    // v has same type as x
}
```

### 9.4 For Statements

```
ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
```

**Condition:**
```go
for i < 10 {
    i++
}
```

**For clause:**
```go
for i := 0; i < 10; i++ {
    sum += i
}
```

**Range clause:**
```go
// Array/slice
for i, v := range arr {
    // i: int64, v: element type
}

// Map
for k, v := range m {
    // k: key type, v: value type
}

// String (UTF-8 runes)
for i, r := range s {
    // i: byte index (int64), r: rune (int32)
}

// Channel
for v := range ch {
    // v: element type
}
```

**Index type in range:**
- Always `int64` for arrays/slices/strings
- Key type for maps

### 9.5 Go Statements

```
GoStmt = "go" Expression .
```

```go
go f(x, y, z)
go func() {
    // ...
}()
```

### 9.6 Select Statements

```
SelectStmt = "select" "{" { CommClause } "}" .
CommClause = CommCase ":" StatementList .
CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
```

```go
select {
case x := <-ch1:
    // received from ch1
case ch2 <- y:
    // sent to ch2
default:
    // neither ready
}
```

### 9.7 Return Statements

```
ReturnStmt = "return" [ ExpressionList ] .
```

```go
return
return x
return x, y
```

### 9.8 Defer Statements

```
DeferStmt = "defer" Expression .
```

```go
defer closeFile(f)
defer func() {
    // cleanup
}()
```

---

## 10. Built-in Functions

Moxie provides built-in functions that operate on various types.

### 10.1 Length and Capacity

```
len(v)  int64
cap(v)  int64
```

**For slices and strings:**
- `len(s)` returns the number of elements/bytes (always int64)
- `cap(s)` returns the capacity (always int64)

**For arrays:**
- `len(a)` returns the array length
- `cap(a)` returns the array length

**For maps:**
- `len(m)` returns the number of entries

**For channels:**
- `len(ch)` returns the number of queued elements
- `cap(ch)` returns the buffer capacity

**Return type is always int64 (not platform-dependent).**

### 10.2 Allocation

```
new(T) *T
```

Allocates zeroed storage for a value of type `T` and returns a pointer.

```go
p := new(User)  // *User
i := new(int32) // *int32
```

### 10.3 Appending and Copying

```
append(s *[]T, values ...T) *[]T
```

Appends elements to a slice.

```go
s := &[]int32{1, 2}
s = append(s, 3, 4, 5)  // &[]int32{1, 2, 3, 4, 5}
```

**Copying slice elements:**

```
copy(dst, src *[]T) int64
```

Copies elements from source slice to destination slice. Returns the number of elements copied (the minimum of `len(src)` and `len(dst)`).

```go
src := &[]int32{1, 2, 3}
dst := &[]int32{0, 0, 0, 0, 0}
n := copy(dst, src)  // n = 3, dst = &[]int32{1, 2, 3, 0, 0}
```

**Deep copying (entire slice with new backing array):**

```
clone(v *T) *T
```

Creates a deep copy of any value, allocating a new backing array.

```go
s1 := &[]int32{1, 2, 3}
s2 := clone(s1)  // Deep copy with new backing array
s2[0] = 99       // Does not affect s1

m1 := &map[string]int32{"a": 1}
m2 := clone(m1)  // Deep copy
```

**Difference between `copy()` and `clone()`:**

- `copy(dst, src)`: Copies elements between existing slices (shallow copy of elements)
- `clone(src)`: Allocates new backing array and copies all elements (deep copy)

### 10.4 Growing Slices

```
grow(s *[]T, n int64) *[]T
```

Pre-allocates capacity for a slice without changing length.

```go
s := &[]int32{}
s = grow(s, 1000)  // Pre-allocate for 1000 elements
// len(s) == 0, cap(s) >= 1000
```

### 10.5 Clearing

```
clear(x)
```

Resets a collection:
- For slices: sets length to 0 (capacity unchanged)
- For maps: removes all entries

```go
s := &[]int32{1, 2, 3}
clear(s)  // len(s) == 0, cap(s) unchanged

m := &map[string]int32{"a": 1, "b": 2}
clear(m)  // len(m) == 0
```

### 10.6 Deletion

```
delete(m *map[K]V, key K)
```

Removes an entry from a map.

```go
m := &map[string]int32{"a": 1, "b": 2}
delete(m, "a")
```

### 10.7 Freeing Memory

```
free(p *T)
```

Provides a hint to the garbage collector that memory can be freed immediately.

```go
largeData := &[]byte{}
largeData = grow(largeData, 10_000_000)
// ... use largeData ...
free(largeData)  // Hint to GC
```

**Note:** This is a hint, not a guarantee. The GC may ignore it.

### 10.8 Channel Operations

```
close(ch *chan T)
```

Closes a channel.

```go
ch := &chan int32{}
close(ch)
```

### 10.9 Complex Numbers

```
complex(r, i FloatType) ComplexType
real(c ComplexType) FloatType
imag(c ComplexType) FloatType
```

```go
c := complex(1.0, 2.0)  // complex128
r := real(c)            // float64(1.0)
i := imag(c)            // float64(2.0)
```

### 10.10 Panic and Recover

```
panic(interface{})
recover() interface{}
```

```go
panic("something went wrong")

defer func() {
    if r := recover(); r != nil {
        fmt.Println("recovered:", r)
    }
}()
```

---

## 11. Packages

### 11.1 Package Clause

```
PackageClause = "package" PackageName .
PackageName   = identifier .
```

```go
package main
package users
package http
```

**Package naming:**
- All lowercase (no PascalCase for packages)
- Short, concise names
- No underscores or mixed case

### 11.2 Import Declarations

```
ImportDecl       = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
ImportSpec       = [ "." | PackageName ] ImportPath .
ImportPath       = string_lit .
```

```go
import "fmt"
import "net/http"

import (
    "fmt"
    "os"
    "github.com/user/package"
)

import (
    . "fmt"              // Import into current namespace
    m "math"             // Import with alias
    _ "image/png"        // Import for side effects
)
```

### 11.3 Program Execution

A complete program is created by linking a single package called `main` with its dependencies.

The `main` package must have a function called `main`:

```go
package main

func main() {
    // Program entry point
}
```

---

## 12. Native FFI (Foreign Function Interface)

Moxie provides built-in support for calling C libraries without CGO.

### 12.1 Dynamic Library Loading

```
dlopen(filename string, flags int32) *DLib
dlclose(lib *DLib)
dlerror() string
```

**Flags:**
```go
const (
    RTLD_LAZY   = 0x001  // Lazy binding
    RTLD_NOW    = 0x002  // Immediate binding
    RTLD_GLOBAL = 0x100  // Make symbols available globally
    RTLD_LOCAL  = 0x000  // Symbols not available for symbol resolution
)
```

**Example:**
```go
lib := dlopen("libc.so.6", RTLD_LAZY)
if lib == nil {
    panic(dlerror())
}
defer dlclose(lib)
```

### 12.2 Symbol Lookup

```
dlsym[T any](lib *DLib, name string) T
```

Type-safe symbol lookup using generics:

```go
// Load function with specific signature
printf := dlsym[func(*byte, ...interface{}) int32](lib, "printf")

// Call it
msg := "Hello from Moxie\n"
printf(&msg[0])
```

### 12.3 Memory-based Library Loading

```
dlopen_mem(data *[]byte, flags int32) *DLib
```

Load a shared library from memory (embedded in binary):

```go
import _ "embed"

//go:embed libfoo.so
var libfooData []byte

func main() {
    lib := dlopen_mem(&libfooData, RTLD_NOW)
    defer dlclose(lib)

    foo := dlsym[func() int32](lib, "foo")
    result := foo()
}
```

### 12.4 Type Mapping (C to Moxie)

| C Type | Moxie Type |
|--------|------------|
| `char`, `signed char` | `int8` |
| `unsigned char` | `uint8`, `byte` |
| `short` | `int16` |
| `unsigned short` | `uint16` |
| `int` | `int32` |
| `unsigned int` | `uint32` |
| `long` (64-bit) | `int64` |
| `unsigned long` (64-bit) | `uint64` |
| `long long` | `int64` |
| `unsigned long long` | `uint64` |
| `float` | `float32` |
| `double` | `float64` |
| `void*` | `unsafe.Pointer` |
| `char*` | `*byte` (string compatible) |

### 12.5 Example: Using SQLite

```go
import "unsafe"

func openDatabase(path string) *SQLite3 {
    lib := dlopen("libsqlite3.so.0", RTLD_LAZY)

    // Type-safe function loading
    type SQLite3 = unsafe.Pointer

    sqlite3_open := dlsym[func(*byte, **SQLite3) int32](lib, "sqlite3_open")
    sqlite3_close := dlsym[func(*SQLite3) int32](lib, "sqlite3_close")

    var db *SQLite3
    pathCStr := path + "\x00"
    rc := sqlite3_open(&pathCStr[0], &db)
    if rc != 0 {
        panic("failed to open database")
    }

    return db
}
```

---

## 13. Zero-Copy Type Coercion

Moxie supports zero-copy reinterpretation of numeric slices with automatic endianness handling.

### 13.1 Syntax

```
(*[]TargetType)(sourceSlice)                     // Native endian
(*[]TargetType, LittleEndian)(sourceSlice)       // Little-endian
(*[]TargetType, BigEndian)(sourceSlice)          // Big-endian
```

### 13.2 Byte Order Constants

```go
const (
    NativeEndian = 0  // Platform native (default)
    LittleEndian = 1  // x86, x86-64, ARM64
    BigEndian    = 2  // Network byte order
)
```

### 13.3 Allowed Types

**Source and target must be fixed-width numeric types:**
- Integers: `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `uint64`
- Floats: `float32`, `float64`
- Byte: `byte` (alias for `uint8`)

**Not allowed:**
- Structs (alignment issues)
- Slices of slices (not fixed-width)
- Strings (use `*[]byte` directly)
- Maps/channels

### 13.4 Examples

**Basic reinterpretation:**
```go
bytes := &[]byte{0x01, 0x02, 0x03, 0x04}

// Native endian (platform-dependent)
u32s := (*[]uint32)(bytes)  // len=1

// On little-endian machine: u32s[0] == 0x04030201
// On big-endian machine: u32s[0] == 0x01020304
```

**With explicit endianness:**
```go
bytes := &[]byte{0x01, 0x02, 0x03, 0x04}

// Force little-endian (portable)
u32s_le := (*[]uint32, LittleEndian)(bytes)
// u32s_le[0] == 0x04030201 on ALL platforms

// Force big-endian (network byte order)
u32s_be := (*[]uint32, BigEndian)(bytes)
// u32s_be[0] == 0x01020304 on ALL platforms
```

**Network protocol parsing:**
```go
func parsePacket(data *[]byte) Header {
    // Network protocols use big-endian
    fields := (*[]uint32, BigEndian)(data[0:16])

    return Header{
        Magic:   fields[0],  // Auto byte-swapped
        Version: fields[1],
        Length:  fields[2],
        CRC:     fields[3],
    }
}
```

**Cryptography:**
```go
func xorBlocks(a, b *[]byte) {
    // Process 8 bytes at a time
    a64 := (*[]uint64)(a)
    b64 := (*[]uint64)(b)
    for i := range a64 {
        a64[i] ^= b64[i]
    }
}
```

### 13.5 Length Adjustment

When casting, length and capacity are adjusted based on type sizes:

```go
bytes := &[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}  // len=10

u16s := (*[]uint16)(bytes)  // len=5 (10/2)
u32s := (*[]uint32)(bytes)  // len=2 (10/4, 2 bytes unused)
u64s := (*[]uint64)(bytes)  // len=1 (10/8, 2 bytes unused)
```

### 13.6 Alignment

Casting checks alignment at runtime (debug mode) or relies on undefined behavior (release mode for performance):

```go
bytes := &[]byte{1, 2, 3, 4, 5}

// May panic if &bytes[0] is not 4-byte aligned
u32s := (*[]uint32)(bytes)

// Safe alignment
offset := uintptr(unsafe.Pointer(&bytes[0])) & ^uintptr(3)
aligned := bytes[offset:]
u32s := (*[]uint32)(aligned)  // Guaranteed aligned
```

### 13.7 Byte Order Metadata

Slice headers include a byte order field that is preserved across casts:

```go
bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
u32s_le := (*[]uint32, LittleEndian)(bytes)

// Converting back preserves endianness
bytes2 := (*[]byte)(u32s_le)  // Inherits LittleEndian tag

// Further casts inherit endianness
u16s := (*[]uint16)(u32s_le)  // Also LittleEndian
```

### 13.8 Copy vs Reinterpret

**Reinterpret (zero-copy):**
```go
u32s := (*[]uint32)(bytes)  // Same backing array
```

**Copy (explicit):**
```go
u32s_copy := &(*[]uint32)(bytes)  // & forces allocation
// Or
u32s_copy := clone((*[]uint32)(bytes))
```

---

## 14. Summary of Changes from Go

### 14.1 Removed Features

- ❌ `make()` function (use composite literals with `&`)
- ❌ CGO (use `dlopen`/`dlsym`)
- ❌ Immutable strings (strings are mutable `*[]byte`)
- ❌ Implicit reference types (slices/maps/channels are explicit pointers)
- ❌ Platform-dependent `int` and `uint` types (use explicit `int32`/`int64`/`uint32`/`uint64`)

### 14.2 Added Features

- ✅ Explicit pointers for slices: `*[]T`
- ✅ Explicit pointers for maps: `*map[K]V`
- ✅ Explicit pointers for channels: `*chan T`
- ✅ Mutable strings: `string` = `*[]byte`
- ✅ True `const` with MMU protection for all types
- ✅ `clone()` built-in for deep copying
- ✅ `grow()` built-in for pre-allocating capacity
- ✅ `clear()` built-in for resetting slices/maps
- ✅ `free()` built-in for GC hints
- ✅ `+` operator for slice concatenation
- ✅ Native FFI: `dlopen`, `dlsym`, `dlclose`, `dlerror`
- ✅ Memory-based library loading: `dlopen_mem`
- ✅ Zero-copy type coercion with endianness control
- ✅ Always `int64` return from `len()` and `cap()`

### 14.3 Unchanged Features

- ✅ Same syntax for most code
- ✅ Same package system
- ✅ Same interfaces
- ✅ Same goroutines and channels
- ✅ Same defer/panic/recover
- ✅ Same generics (Go 1.18+)
- ✅ PascalCase/camelCase naming conventions
- ✅ Same standard library (mostly)

---

## 15. Examples

### 15.1 Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Moxie!")
}
```

### 15.2 HTTP Server

```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from Moxie!")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

### 15.3 Working with Slices

```go
package main

import "fmt"

func main() {
    // Create slice with explicit pointer
    numbers := &[]int32{1, 2, 3, 4, 5}

    // Append
    numbers = append(numbers, 6, 7, 8)

    // Clone
    copy := clone(numbers)
    copy[0] = 99  // Does not affect numbers

    // Pre-allocate
    large := &[]int32{}
    large = grow(large, 1000)

    // Concatenate
    a := &[]int32{1, 2, 3}
    b := &[]int32{4, 5, 6}
    c := a + b  // &[]int32{1, 2, 3, 4, 5, 6}

    // Clear
    clear(numbers)  // len=0
}
```

### 15.4 Working with Maps

```go
package main

import "fmt"

func main() {
    // Create map with explicit pointer
    ages := &map[string]int32{
        "Alice": 30,
        "Bob":   25,
    }

    // Access
    age := ages["Alice"]

    // Add
    ages["Charlie"] = 35

    // Delete
    delete(ages, "Bob")

    // Check existence
    age, ok := ages["David"]

    // Clone
    copy := clone(ages)

    // Clear
    clear(ages)
}
```

### 15.5 Working with Strings

```go
package main

import "fmt"

func main() {
    // Strings are mutable
    s := "hello"
    s[0] = 'H'  // s = "Hello"

    // Concatenation
    greeting := "Hello" + ", " + "world!"

    // Immutable string with const
    const template = "Dear %s,"
    // template[0] = 'X'  // Compile error

    // Clone for safety
    mutable := clone(template)
    mutable[0] = 'X'  // OK
}
```

### 15.6 Using FFI

```go
package main

import "fmt"

func main() {
    // Load library
    lib := dlopen("libc.so.6", RTLD_LAZY)
    if lib == nil {
        panic(dlerror())
    }
    defer dlclose(lib)

    // Look up symbol with type safety
    strlen := dlsym[func(*byte) int64](lib, "strlen")

    // Use it
    msg := "Hello from Moxie\x00"
    length := strlen(&msg[0])
    fmt.Printf("Length: %d\n", length)
}
```

### 15.7 Zero-Copy Binary Parsing

```go
package main

import "fmt"

type PacketHeader struct {
    Magic   uint32
    Version uint16
    Length  uint16
    CRC     uint32
}

func parsePacket(data *[]byte) PacketHeader {
    // Zero-copy reinterpretation with network byte order
    fields := (*[]uint32, BigEndian)(data[0:12])

    return PacketHeader{
        Magic:   fields[0],
        Version: uint16(fields[1] >> 16),
        Length:  uint16(fields[1] & 0xFFFF),
        CRC:     fields[2],
    }
}

func main() {
    packet := &[]byte{
        0x12, 0x34, 0x56, 0x78,  // Magic (big-endian)
        0x00, 0x01, 0x02, 0x00,  // Version=1, Length=512
        0xAB, 0xCD, 0xEF, 0x12,  // CRC
    }

    header := parsePacket(packet)
    fmt.Printf("Magic: 0x%08X\n", header.Magic)
    fmt.Printf("Version: %d\n", header.Version)
    fmt.Printf("Length: %d\n", header.Length)
}
```

---

## Appendix A: Keywords

```
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

## Appendix B: Predeclared Identifiers

**Types:**
```
bool byte complex64 complex128
float32 float64
int8 int16 int32 int64
string
uint8 uint16 uint32 uint64
rune
```

**Constants:**
```
true false iota
nil
```

**Functions:**
```
append cap close clone complex clear delete free grow
len new panic print println real imag recover
dlopen dlclose dlsym dlopen_mem dlerror
```

**Byte order constants:**
```
NativeEndian LittleEndian BigEndian
```

**FFI constants:**
```
RTLD_LAZY RTLD_NOW RTLD_GLOBAL RTLD_LOCAL
```

---

## Appendix C: Migration from Go

### C.1 Slices

**Before (Go):**
```go
var s []int32
s = make([]int32, 10, 20)
s = append(s, 1, 2, 3)
```

**After (Moxie):**
```go
var s *[]int32
s = &[]int32{}
s = grow(s, 20)
for i := 0; i < 10; i++ {
    s = append(s, 0)
}
s = append(s, 1, 2, 3)
```

### C.2 Maps

**Before (Go):**
```go
var m map[string]int32
m = make(map[string]int32)
m["key"] = 123
```

**After (Moxie):**
```go
var m *map[string]int32
m = &map[string]int32{}
m["key"] = 123
```

### C.3 Channels

**Before (Go):**
```go
ch := make(chan int32, 10)
ch <- 42
x := <-ch
```

**After (Moxie):**
```go
ch := &chan int32{}
ch = grow(ch, 10)
ch <- 42
x := <-ch
```

### C.4 Strings

**Before (Go):**
```go
s := "hello"
// s[0] = 'H'  // Error
b := []byte(s)  // Copy
b[0] = 'H'
s = string(b)   // Copy back
```

**After (Moxie):**
```go
s := "hello"
s[0] = 'H'  // OK - mutable
```

### C.5 CGO

**Before (Go with CGO):**
```go
/*
#include <stdio.h>
void hello() { printf("Hello\n"); }
*/
import "C"

func main() {
    C.hello()
}
```

**After (Moxie with FFI):**
```go
import "unsafe"

func main() {
    lib := dlopen("libc.so.6", RTLD_LAZY)
    defer dlclose(lib)

    printf := dlsym[func(*byte, ...interface{}) int32](lib, "printf")
    msg := "Hello\n\x00"
    printf(&msg[0])
}
```

---

## Appendix D: Rationale

### D.1 Why Explicit Pointers for Slices/Maps/Channels?

**Problem:** Go has two type categories with different semantics:
- Value types (int, struct, array) - copied on assignment
- Reference types (slice, map, channel) - shared on assignment

This creates cognitive overhead and hides mutation.

**Solution:** Make everything explicit. Use pointers for all reference types.

**Benefits:**
- Obvious sharing semantics
- Consistent nil behavior
- Better concurrency safety
- Simpler mental model

### D.2 Why Mutable Strings?

**Problem:** Go has two types for text data:
- `string` - immutable, but must convert to `[]byte` for mutation (allocates and copies)
- `[]byte` - mutable, but must convert to `string` for string operations (allocates and copies)

**Solution:** Merge them into a single mutable type.

**Benefits:**
- No conversion overhead
- Consistent text handling
- Simpler API
- Better performance

### D.3 Why True Const with MMU Protection?

**Problem:** Go's const only works for primitives and is not truly immutable (no memory protection).

**Solution:** Extend const to all types with hardware (MMU) enforcement.

**Benefits:**
- True immutability
- Security (tamper-proof data)
- Performance (compiler optimizations)
- Clarity (explicit immutability)

### D.4 Why Native FFI Instead of CGO?

**Problem:** CGO has many drawbacks:
- Requires C compiler (breaks cross-compilation)
- Slow (200ns per call)
- Complex build process
- Breaks static binaries

**Solution:** Built-in FFI with `dlopen`/`dlsym`.

**Benefits:**
- Fast (10ns per call)
- Type-safe (generic type checking)
- Pure Go cross-compilation
- Static binaries (with `dlopen_mem`)
- Simpler deployment

### D.5 Why Remove Platform-Dependent int/uint?

**Problem:** `int` and `uint` are 32-bit on 32-bit platforms, 64-bit on 64-bit platforms. This causes:
- Serialization issues
- Cross-platform bugs
- Hidden overflow

**Solution:** Require explicit bit widths.

**Benefits:**
- Portable serialization
- Explicit developer intent
- Consistent behavior
- No hidden bugs

### D.6 Why Zero-Copy Type Coercion?

**Problem:** Binary data parsing requires manual byte manipulation or unsafe pointer casts.

**Solution:** Safe, zero-copy reinterpretation with automatic endianness handling.

**Benefits:**
- Eliminates copy overhead
- Type-safe
- Portable (automatic byte swapping)
- Faster (hardware-accelerated)
- Common use case (network protocols, crypto, multimedia)

---

**End of Specification**
