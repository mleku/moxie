# The copy() Built-in Function

## Overview

The `copy()` built-in function copies elements from a source slice to a destination slice. It is similar to Go's `copy()` but with enhanced support for type casting with endianness control.

## Signature

```moxie
func copy(dst, src *[]T) int64
```

**Returns:** Number of elements copied (minimum of `len(dst)` and `len(src)`)

## Basic Semantics

### Behavior

1. **Overwrites** destination elements (does not append)
2. **Copies** `min(len(dst), len(src))` elements
3. **No allocation** - operates on existing slices
4. **Returns** count of elements copied
5. **In-place** - modifies `dst` directly

### Type Requirement

Both slices must have the **same element type**, or the source must be cast to match the destination type.

## Basic Examples

### Same Type Copy

```moxie
dst := &[]int32{0, 0, 0, 0, 0}
src := &[]int32{1, 2, 3}

n := copy(dst, src)
// n = 3
// dst = &[]int32{1, 2, 3, 0, 0}
```

### Source Longer Than Destination

```moxie
dst := &[]byte{0, 0, 0}
src := &[]byte{1, 2, 3, 4, 5}

n := copy(dst, src)
// n = 3 (limited by dst length)
// dst = &[]byte{1, 2, 3}
```

### Destination Longer Than Source

```moxie
dst := &[]byte{0, 0, 0, 0, 0}
src := &[]byte{1, 2, 3}

n := copy(dst, src)
// n = 3
// dst = &[]byte{1, 2, 3, 0, 0}
// Remaining elements unchanged
```

### Copy to Offset

```moxie
dst := &[]byte{0, 0, 0, 0, 0}
src := &[]byte{1, 2, 3}

// Copy to offset 2
n := copy(dst[2:], src)
// n = 3
// dst = &[]byte{0, 0, 1, 2, 3}
```

### Copy from Offset

```moxie
dst := &[]byte{0, 0, 0, 0, 0}
src := &[]byte{1, 2, 3, 4, 5}

// Copy from offset 2
n := copy(dst, src[2:])
// n = 3
// dst = &[]byte{3, 4, 5, 0, 0}
```

## Type Casting with copy()

The source slice can be **cast inline** to match the destination type. This enables zero-copy conversions between numeric slice types with optional endianness control.

### Basic Type Cast

```moxie
// uint32 to bytes (native endian)
src := &[]uint32{0x12345678}
dst := &[]byte{0, 0, 0, 0}

copy(dst, (*[]byte)(src))
// On little-endian: dst = {0x78, 0x56, 0x34, 0x12}
// On big-endian:    dst = {0x12, 0x34, 0x56, 0x78}
```

### Cast with Endianness

```moxie
src := &[]uint32{0x12345678, 0xABCDEF00}
dst := &[]byte{0, 0, 0, 0, 0, 0, 0, 0}

// Force little-endian
copy(dst, (*[]byte, LittleEndian)(src))
// dst = {0x78, 0x56, 0x34, 0x12, 0x00, 0xEF, 0xCD, 0xAB}
// (Same on ALL platforms)

// Force big-endian
copy(dst, (*[]byte, BigEndian)(src))
// dst = {0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xEF, 0x00}
// (Same on ALL platforms)
```

### Type Cast Rules

1. **Only numeric types**: int8-64, uint8-64, float32, float64
2. **Length calculation**: `newLen = oldLen * sizeof(oldType) / sizeof(newType)`
3. **Alignment**: Checked in debug mode, undefined in release mode if misaligned
4. **Zero-copy**: Cast is a reinterpretation view, not a copy

## Practical Use Cases

### 1. Network Protocol Serialization

```moxie
// Serialize packet header to network byte order
type PacketHeader struct {
    Magic   uint32
    Version uint16
    Length  uint16
    Flags   uint32
}

func serializeHeader(h *PacketHeader) *[]byte {
    // Flatten struct to uint32 array (assuming packed layout)
    values := &[]uint32{
        h.Magic,
        uint32(h.Version)<<16 | uint32(h.Length),
        h.Flags,
    }

    // Serialize to big-endian bytes
    buffer := &[]byte{}
    buffer = grow(buffer, 12)
    copy(buffer, (*[]byte, BigEndian)(values))

    return buffer
}

func parseHeader(data *[]byte) *PacketHeader {
    // Parse big-endian bytes to uint32
    values := &[]uint32{0, 0, 0}
    copy(values, (*[]uint32, BigEndian)(data[0:12]))

    return &PacketHeader{
        Magic:   values[0],
        Version: uint16(values[1] >> 16),
        Length:  uint16(values[1]),
        Flags:   values[2],
    }
}
```

### 2. Cryptographic Operations

```moxie
// XOR two byte slices (common in stream ciphers)
func xor(dst, a, b *[]byte) {
    // Process 8 bytes at a time using uint64
    len64 := len(a) / 8
    a64 := (*[]uint64)(a)[0:len64]
    b64 := (*[]uint64)(b)[0:len64]
    dst64 := (*[]uint64)(dst)[0:len64]

    for i := range a64 {
        dst64[i] = a64[i] ^ b64[i]
    }

    // Handle remaining bytes
    offset := len64 * 8
    for i := offset; i < len(a); i++ {
        dst[i] = a[i] ^ b[i]
    }
}
```

### 3. Binary Data Parsing

```moxie
// Parse file header with mixed sizes
func parseFileHeader(data *[]byte) FileHeader {
    // Read uint32 magic number (little-endian)
    magic := &[]uint32{0}
    copy(magic, (*[]uint32, LittleEndian)(data[0:4]))

    // Read uint16 version
    version := &[]uint16{0}
    copy(version, (*[]uint16, LittleEndian)(data[4:6]))

    // Read uint64 file size
    fileSize := &[]uint64{0}
    copy(fileSize, (*[]uint64, LittleEndian)(data[6:14]))

    return FileHeader{
        Magic:    magic[0],
        Version:  version[0],
        FileSize: fileSize[0],
    }
}
```

### 4. Zero-Copy Buffer Management

```moxie
// Efficiently copy data into pre-allocated ring buffer
type RingBuffer struct {
    buf  *[]byte
    head int64
    tail int64
    size int64
}

func (rb *RingBuffer) Write(data *[]byte) int64 {
    available := rb.size - rb.tail

    if available >= int64(len(data)) {
        // Fits in one contiguous chunk
        n := copy(rb.buf[rb.tail:], data)
        rb.tail += n
        return n
    } else {
        // Wraps around
        n1 := copy(rb.buf[rb.tail:], data)
        n2 := copy(rb.buf[0:], data[n1:])
        rb.tail = n2
        return n1 + n2
    }
}
```

### 5. Image Processing

```moxie
// Convert RGBA pixels to grayscale in-place
func toGrayscale(pixels *[]byte) {
    // Process 4 bytes (RGBA) at a time
    for i := int64(0); i < len(pixels); i += 4 {
        r := pixels[i]
        g := pixels[i+1]
        b := pixels[i+2]

        // Grayscale = 0.299*R + 0.587*G + 0.114*B
        gray := byte((299*int32(r) + 587*int32(g) + 114*int32(b)) / 1000)

        pixels[i] = gray
        pixels[i+1] = gray
        pixels[i+2] = gray
        // pixels[i+3] (alpha) unchanged
    }
}
```

### 6. Memory-Mapped I/O

```moxie
// Copy data from memory-mapped file to process buffer
func readFromMmap(mmapData *[]byte, offset, length int64) *[]byte {
    buffer := &[]byte{}
    buffer = grow(buffer, length)

    n := copy(buffer, mmapData[offset:offset+length])
    return buffer[0:n]
}

// Write to memory-mapped file with endianness conversion
func writeToMmap(mmapData *[]byte, offset int64, values *[]uint32) {
    // Write as big-endian to memory-mapped region
    copy(mmapData[offset:], (*[]byte, BigEndian)(values))
}
```

## Performance Characteristics

### Copy Speed

| Source Type | Destination Type | Speed | Notes |
|-------------|------------------|-------|-------|
| Same type | Same type | Fast | Optimized memcpy |
| uint64 | byte | Fast | Zero-copy cast + memcpy |
| uint32 | byte (LE) | Medium | May require byte swap on BE systems |
| uint32 | byte (BE) | Medium | May require byte swap on LE systems |
| byte | uint64 | Fast | Zero-copy cast + memcpy |

### Optimization Tips

1. **Align buffer sizes** to element boundaries for better cache performance
2. **Use larger types** (uint64) when processing byte arrays for speed
3. **Pre-allocate buffers** with `grow()` to avoid reallocations
4. **Minimize endian conversions** - use native endian when possible

## Comparison with Other Operations

### copy() vs clone()

| Aspect | `copy(dst, src)` | `clone(src)` |
|--------|------------------|--------------|
| **Allocation** | No | Yes (new backing array) |
| **Destination** | Explicit | Returned |
| **Partial copy** | Yes | No (always full) |
| **Return value** | Count copied | New slice |
| **Use case** | Reuse buffers | Independent copy |

```moxie
// clone() - full copy, new allocation
backup := clone(original)

// copy() - partial copy, existing buffer
buffer := &[]byte{}
buffer = grow(buffer, 1024)
n := copy(buffer, data)  // Reuse buffer
```

### copy() vs | (concatenation)

| Aspect | `copy(dst, src)` | `dst | src` |
|--------|------------------|-------------|
| **Allocation** | No | Yes (new backing array) |
| **Overwrites** | Yes | No |
| **Appends** | No | Yes |
| **Return value** | Count | New slice |
| **Use case** | Overwrite buffer | Combine slices |

```moxie
// Concatenation - creates new slice
result := slice1 | slice2  // New allocation

// Copy - overwrites existing
copy(buffer, data)  // No allocation if buffer has capacity
```

## Type Safety

### Valid Casts

```moxie
// Numeric types - OK
copy(dst_bytes, (*[]byte)(src_uint32))
copy(dst_uint64, (*[]uint64)(src_byte))
copy(dst_float32, (*[]float32)(src_uint32))

// With endianness - OK
copy(dst, (*[]byte, LittleEndian)(src))
copy(dst, (*[]uint32, BigEndian)(src))
```

### Invalid Casts

```moxie
// Non-numeric types - ERROR
type MyStruct struct { x int32 }
src := &[]MyStruct{{1}, {2}}
dst := &[]byte{}
copy(dst, (*[]byte)(src))  // ERROR: cannot cast struct slice

// Incompatible types - ERROR
src := &[]string{"hello"}
dst := &[]byte{}
copy(dst, (*[]byte)(src))  // ERROR: string is not numeric type

// Direct type mismatch - ERROR
src := &[]int32{1, 2, 3}
dst := &[]byte{0, 0, 0}
copy(dst, src)  // ERROR: must cast first
```

## Error Handling

The `copy()` function does not return an error. Instead:

1. **Length mismatch**: Copies `min(len(dst), len(src))` elements
2. **Nil slices**: Copies 0 elements (returns 0)
3. **Zero-length slices**: Copies 0 elements (returns 0)
4. **Misaligned cast**: Undefined behavior (panic in debug mode)

```moxie
// Safe usage
if dst != nil && src != nil && len(dst) > 0 && len(src) > 0 {
    n := copy(dst, src)
    // Process n elements
}

// No error checking needed for valid slices
n := copy(dst, src)  // Always safe, may copy 0
```

## Implementation Notes

### Runtime Implementation

The Moxie compiler should optimize `copy()` to:

1. **Use `memmove()`** for same-type copies
2. **Use SIMD instructions** when available
3. **Inline** for small, constant-size copies
4. **Apply endian conversion** at element granularity

### Pseudo-implementation

```moxie
func copy(dst, src *[]T) int64 {
    if dst == nil || src == nil {
        return 0
    }

    n := min(len(dst), len(src))
    if n == 0 {
        return 0
    }

    // If T is numeric and src is cast with endianness
    if needsByteSwap(src) {
        for i := int64(0); i < n; i++ {
            dst[i] = byteSwap(src[i])
        }
    } else {
        // Fast path: memmove
        memmove(&dst[0], &src[0], n * sizeof(T))
    }

    return n
}
```

## Summary

The `copy()` built-in provides:

✅ **Efficient buffer reuse** - No allocation
✅ **Type-safe casting** - Zero-copy numeric conversions
✅ **Endianness control** - Portable serialization
✅ **Predictable behavior** - Returns count, never errors
✅ **Performance** - Optimized to memmove/SIMD

Use `copy()` when:
- Overwriting existing buffers
- Working with pre-allocated memory
- Converting between numeric slice types
- Serializing/deserializing with specific byte order
- Implementing high-performance data processing

**copy() is essential for zero-allocation, high-performance systems programming in Moxie.**
