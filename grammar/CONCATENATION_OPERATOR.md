# Concatenation Operator Change: | (Vertical Bar)

## Rationale

Moxie adopts the **vertical bar** (`|`) as the concatenation operator for strings and slices, following standard **cryptographic notation** where `a | b` denotes the concatenation of `a` and `b`.

### Why This Matters

In cryptographic literature, protocols, and papers, concatenation is universally written as `|`:

- **HMAC**: `HMAC(K, m) = H((K ⊕ opad) | H((K ⊕ ipad) | m))`
- **Key derivation**: `KDF(secret, info) = HMAC(secret, counter | info)`
- **Message authentication**: `MAC = H(key | message)`
- **TLS PRF**: `PRF(secret, label, seed) = P_hash(secret, label | seed)`

By using `|` for concatenation, Moxie code becomes **more readable** for developers implementing cryptographic systems, as the code directly mirrors the mathematical notation.

## Language Changes

### 1. Removed: `append()` Built-in Function

The `append()` function is **eliminated** from Moxie. All concatenation operations use the `|` operator.

**Before (Go):**
```go
s := []int{1, 2, 3}
s = append(s, 4, 5, 6)
s = append(s, moreItems...)
```

**After (Moxie):**
```moxie
s := &[]int{1, 2, 3}
s = s | &[]int{4, 5, 6}
s = s | moreItems
```

### 2. Changed: String Concatenation

String concatenation now uses `|` instead of `+`.

**Before (Go):**
```go
greeting := "Hello, " + name + "!"
```

**After (Moxie):**
```moxie
greeting := "Hello, " | name | "!"
```

### 3. Removed: Bitwise OR with `|`

The `|` operator is **no longer available** for bitwise OR operations.

**Before (Go):**
```go
flags := FLAG_READ | FLAG_WRITE | FLAG_EXECUTE
mask := 0xFF00 | 0x00FF
```

**After (Moxie):**
Bitwise OR is not directly available. Alternative approaches:

```moxie
// Option 1: Use addition for non-overlapping flags
flags := FLAG_READ + FLAG_WRITE + FLAG_EXECUTE

// Option 2: Use bit manipulation functions (if available)
flags := bitor(FLAG_READ, FLAG_WRITE, FLAG_EXECUTE)

// Option 3: Manual bit operations
result := a + b - (a & b)  // Equivalent to a | b when bits don't overlap
```

**Note:** For most use cases involving flags, using distinct bit positions and addition works fine:
```moxie
const (
    FLAG_READ    = 1 << 0  // 0x01
    FLAG_WRITE   = 1 << 1  // 0x02
    FLAG_EXECUTE = 1 << 2  // 0x04
)

flags := FLAG_READ + FLAG_WRITE  // Works because bits don't overlap
```

## Operator Precedence

The `|` concatenation operator has its own precedence level, positioned between relational operators and arithmetic operators:

```
Precedence (lowest to highest):
1. ||               (logical OR)
2. &&               (logical AND)
3. == != < <= > >=  (relational)
4. |                (concatenation) ← NEW
5. + - ^            (arithmetic add/subtract/XOR)
6. * / % << >> & &^ (arithmetic multiply/divide/bitwise)
7. unary + - ! ^ * & <-
```

This precedence allows natural composition:

```moxie
// Concatenate if condition is true
result := condition && (prefix | suffix)

// Concatenate with arithmetic
data := header | body | padding[0:len(body)%16]
```

## Examples

### Cryptographic Code

Moxie's `|` operator makes cryptographic implementations directly match the specification:

```moxie
// HMAC-SHA256 implementation mirrors the spec
func hmacSHA256(key, message *[]byte) *[]byte {
    const blockSize = 64
    const ipad = 0x36
    const opad = 0x5C

    // Pad or hash key if needed
    k := key
    if len(key) > blockSize {
        k = sha256(key)
    }
    if len(key) < blockSize {
        k = k | &[]byte{}  // Implicit zero-padding in grow
        k = grow(k, blockSize)
    }

    // XOR key with padding
    keyIpad := xor(k, ipad)
    keyOpad := xor(k, opad)

    // HMAC(K, m) = H((K ⊕ opad) | H((K ⊕ ipad) | m))
    innerHash := sha256(keyIpad | message)
    outerHash := sha256(keyOpad | innerHash)

    return outerHash
}
```

Compare with the mathematical notation:
```
HMAC(K, m) = H((K ⊕ opad) | H((K ⊕ ipad) | m))
```

The code **exactly matches** the formula!

### String Building

```moxie
func buildURL(scheme, host, path string, params *map[string]string) string {
    url := scheme | "://" | host | path

    if len(params) > 0 {
        url = url | "?"
        first := true
        for key, val := range params {
            if !first {
                url = url | "&"
            }
            url = url | key | "=" | val
            first = false
        }
    }

    return url
}
```

### Binary Protocol Construction

```moxie
func buildPacket(payload *[]byte) *[]byte {
    const MAGIC = 0x12345678

    header := &[]byte{
        byte(MAGIC >> 24),
        byte(MAGIC >> 16),
        byte(MAGIC >> 8),
        byte(MAGIC),
        byte(len(payload) >> 8),
        byte(len(payload)),
    }

    checksum := crc32(payload)
    trailer := &[]byte{
        byte(checksum >> 24),
        byte(checksum >> 16),
        byte(checksum >> 8),
        byte(checksum),
    }

    // Packet = header | payload | trailer
    return header | payload | trailer
}
```

### TLS Handshake

```moxie
// TLS 1.3 key derivation (simplified)
func deriveSecret(secret, label *[]byte, messages *[]byte) *[]byte {
    transcriptHash := sha256(messages)

    // HKDF-Expand-Label(secret, label, context, length)
    hkdfLabel := &[]byte{
        byte(32), byte(0),           // length
        byte(len(label)),
    } | label |
        &[]byte{byte(len(transcriptHash))} |
        transcriptHash

    return hmac(secret, hkdfLabel)[0:32]
}
```

### Message Authentication Code

```moxie
func authenticate(key, nonce, data *[]byte) (*[]byte, *[]byte) {
    // Construct authenticated message: nonce | ciphertext | tag
    ciphertext := encrypt(key, nonce, data)
    tag := poly1305(key, nonce | ciphertext)

    return nonce | ciphertext | tag, tag
}

func verify(key, message *[]byte) (*[]byte, bool) {
    // Parse: message = nonce | ciphertext | tag
    const nonceLen = 12
    const tagLen = 16

    nonce := message[0:nonceLen]
    ciphertext := message[nonceLen:len(message)-tagLen]
    receivedTag := message[len(message)-tagLen:]

    // Verify tag
    expectedTag := poly1305(key, nonce | ciphertext)
    if !constantTimeEqual(receivedTag, expectedTag) {
        return nil, false
    }

    plaintext := decrypt(key, nonce, ciphertext)
    return plaintext, true
}
```

## Migration from Go

### Automated Migration

A migration tool can automatically rewrite:

1. `append(s, items...)` → `s | items`
2. `string1 + string2` → `string1 | string2`
3. Detect bitwise `|` usage and warn (requires manual fix)

### Manual Migration - Bitwise OR

For code using `|` as bitwise OR, consider:

**Flag enums:**
```go
// Before
flags := FLAG_A | FLAG_B | FLAG_C

// After (if flags are mutually exclusive bits)
flags := FLAG_A + FLAG_B + FLAG_C
```

**Bitmasks:**
```go
// Before
mask := 0xFF00 | 0x00FF

// After
mask := 0xFF00 + 0x00FF  // If bits don't overlap
// or
mask := 0xFFFF           // Directly compute the constant
```

**Generic bitwise OR:**
If you truly need bitwise OR for overlapping bits, implement a helper:
```moxie
func bitor(a, b uint32) uint32 {
    return a + b - (a & b)
}

func bitor64(a, b uint64) uint64 {
    return a + b - (a & b)
}
```

Or use boolean algebra:
```
a | b = a + b - (a & b)
a | b = (a ^ b) + (a & b)
a | b = a + (b & ^a)
```

## Benefits

### 1. Cryptographic Readability

Implementations directly match specifications and academic papers.

### 2. Unified Concatenation

One operator (`|`) for all concatenation:
- Strings: `s1 | s2`
- Byte slices: `bytes1 | bytes2`
- Any slices: `slice1 | slice2`

No more confusion between `+` and `append()`.

### 3. Explicit Operations

The loss of bitwise `|` forces explicit thinking about bit operations, potentially catching bugs where bitwise OR was used incorrectly.

### 4. Familiar to Domain

Cryptographic developers, network protocol implementers, and security researchers will find Moxie's notation natural and intuitive.

## Trade-offs

### Lost: Bitwise OR Operator

The main trade-off is losing the convenient `|` operator for bitwise OR. However:

1. **Flags still work** - Non-overlapping bits can use `+`
2. **Rare use case** - Bitwise OR is less common than concatenation
3. **Workarounds exist** - Helper functions, boolean algebra
4. **Explicitness** - Forces clarity in bit manipulation

### Different from Go

This is a **breaking change** from Go. Code using:
- `+` for string concatenation must change to `|`
- `append()` for slices must change to `|`
- `|` for bitwise OR must be rewritten

However, these changes are **mechanical** and can be automated for most cases.

## Comparison with Other Languages

| Language | Concatenation | Bitwise OR | Notes |
|----------|--------------|------------|-------|
| **Go** | `+` (string), `append()` (slice) | `\|` | Two ways to concatenate |
| **Moxie** | `\|` (all) | *removed* | One way, cryptographic notation |
| **Python** | `+` | `\|` | Two operators |
| **Rust** | Methods | `\|` | No operator for concatenation |
| **C** | `strcat()` | `\|` | Manual memory management |
| **Cryptographic papers** | `\|` | ⊕ (XOR) | Mathematical notation |

Moxie **aligns with mathematical notation** used in cryptography.

## Conclusion

The `|` concatenation operator makes Moxie particularly well-suited for:

- **Cryptographic implementations** - Code matches specifications
- **Network protocols** - Clear byte stream construction
- **Security systems** - Transparent data composition
- **Binary formats** - Explicit packet/message building

While losing bitwise OR is a trade-off, the benefits for the target domain (systems programming with cryptographic focus) outweigh the costs.

The uniformity of using `|` for all concatenation—combined with the elimination of `append()`—simplifies the language and makes data composition operations explicit and consistent.

**Moxie's `|` operator: Unifying concatenation through cryptographic notation.**
