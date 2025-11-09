package runtime

import (
	"encoding/binary"
	"unsafe"
)

// Endianness constants
const (
	NativeEndian = 0 // Platform native endianness
	LittleEndian = 1 // Little-endian (x86, ARM64)
	BigEndian    = 2 // Big-endian (network byte order)
)

// Coerce reinterprets a numeric slice as a different type with optional endianness conversion
// This is a zero-copy operation that reuses the same backing array
//
// Usage:
//   bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
//   u32s := Coerce[byte, uint32](bytes)                    // Native endian
//   u32s_le := Coerce[byte, uint32](bytes, LittleEndian)   // Little-endian
//   u32s_be := Coerce[byte, uint32](bytes, BigEndian)      // Big-endian (network)
func Coerce[From, To any](src *[]From, endian ...int) *[]To {
	if src == nil {
		return nil
	}

	// Get endianness (default to native)
	byteOrder := NativeEndian
	if len(endian) > 0 {
		byteOrder = endian[0]
	}

	// Get type sizes
	var fromZero From
	var toZero To
	fromSize := int(unsafe.Sizeof(fromZero))
	toSize := int(unsafe.Sizeof(toZero))

	if fromSize == 0 || toSize == 0 {
		// Cannot coerce zero-sized types
		panic("moxie.Coerce: cannot coerce zero-sized types")
	}

	// Calculate new length and capacity
	srcSlice := *src
	srcLen := len(srcSlice)
	srcCap := cap(srcSlice)

	dstLen := (srcLen * fromSize) / toSize
	dstCap := (srcCap * fromSize) / toSize

	// Use modern unsafe.Slice API instead of deprecated reflect.SliceHeader
	// Get pointer to the underlying data
	srcData := unsafe.SliceData(srcSlice)

	// Reinterpret as destination type using unsafe.Slice
	// This creates a new slice header pointing to the same backing array
	result := unsafe.Slice((*To)(unsafe.Pointer(srcData)), dstLen)

	// Set capacity by slicing (Go doesn't allow setting cap directly with unsafe.Slice)
	// The capacity is implicitly the maximum length based on the backing array
	result = result[:dstLen:dstCap]

	// Apply endianness conversion if needed
	if byteOrder != NativeEndian && toSize > 1 {
		swapEndianHardwareAccelerated(&result, byteOrder, toSize)
	}

	return &result
}

// swapEndianHardwareAccelerated performs in-place endianness conversion
// using hardware-accelerated encoding/binary where possible
func swapEndianHardwareAccelerated[T any](slice *[]T, targetEndian int, elemSize int) {
	if slice == nil {
		return
	}

	// Determine native endianness
	nativeIsLittle := isLittleEndian()

	// Determine if we need to swap
	needSwap := false
	if targetEndian == LittleEndian && !nativeIsLittle {
		needSwap = true
	} else if targetEndian == BigEndian && nativeIsLittle {
		needSwap = true
	}

	if !needSwap {
		return
	}

	// Try to use hardware-accelerated encoding/binary for common types
	// This provides SIMD acceleration on supported architectures
	s := *slice

	// Get the byte order interface for the target endianness
	var order binary.ByteOrder
	if targetEndian == LittleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}

	// Use hardware-accelerated conversion for supported sizes
	switch elemSize {
	case 2:
		// Try to handle as []uint16 for hardware acceleration
		swapUint16Slice(s, order)
	case 4:
		// Try to handle as []uint32 for hardware acceleration
		swapUint32Slice(s, order)
	case 8:
		// Try to handle as []uint64 for hardware acceleration
		swapUint64Slice(s, order)
	default:
		// Fallback to manual byte swapping for unsupported sizes
		for i := range s {
			swapBytes(unsafe.Pointer(&s[i]), elemSize)
		}
	}
}

// swapUint16Slice uses hardware-accelerated endianness conversion for 16-bit values
func swapUint16Slice[T any](slice []T, order binary.ByteOrder) {
	if len(slice) == 0 {
		return
	}

	// Reinterpret as byte slice for bulk conversion
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), len(slice)*2)

	// Read and write back with correct byte order (hardware accelerated on many platforms)
	for i := 0; i < len(byteSlice); i += 2 {
		val := order.Uint16(byteSlice[i : i+2])
		// Store back in native endian (which means we've done the swap)
		if isLittleEndian() {
			binary.LittleEndian.PutUint16(byteSlice[i:i+2], val)
		} else {
			binary.BigEndian.PutUint16(byteSlice[i:i+2], val)
		}
	}
}

// swapUint32Slice uses hardware-accelerated endianness conversion for 32-bit values
func swapUint32Slice[T any](slice []T, order binary.ByteOrder) {
	if len(slice) == 0 {
		return
	}

	// Reinterpret as byte slice for bulk conversion
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), len(slice)*4)

	// Read and write back with correct byte order (hardware accelerated on many platforms)
	for i := 0; i < len(byteSlice); i += 4 {
		val := order.Uint32(byteSlice[i : i+4])
		// Store back in native endian (which means we've done the swap)
		if isLittleEndian() {
			binary.LittleEndian.PutUint32(byteSlice[i:i+4], val)
		} else {
			binary.BigEndian.PutUint32(byteSlice[i:i+4], val)
		}
	}
}

// swapUint64Slice uses hardware-accelerated endianness conversion for 64-bit values
func swapUint64Slice[T any](slice []T, order binary.ByteOrder) {
	if len(slice) == 0 {
		return
	}

	// Reinterpret as byte slice for bulk conversion
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), len(slice)*8)

	// Read and write back with correct byte order (hardware accelerated on many platforms)
	for i := 0; i < len(byteSlice); i += 8 {
		val := order.Uint64(byteSlice[i : i+8])
		// Store back in native endian (which means we've done the swap)
		if isLittleEndian() {
			binary.LittleEndian.PutUint64(byteSlice[i:i+8], val)
		} else {
			binary.BigEndian.PutUint64(byteSlice[i:i+8], val)
		}
	}
}

// swapBytes performs in-place byte swapping for different sizes
// This is a fallback for unsupported sizes
func swapBytes(ptr unsafe.Pointer, size int) {
	bytes := (*[16]byte)(ptr)

	switch size {
	case 2:
		bytes[0], bytes[1] = bytes[1], bytes[0]
	case 4:
		bytes[0], bytes[1], bytes[2], bytes[3] = bytes[3], bytes[2], bytes[1], bytes[0]
	case 8:
		bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7] =
			bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]
	case 16:
		// Support for 128-bit types (e.g., complex128, SIMD types)
		bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7],
			bytes[8], bytes[9], bytes[10], bytes[11], bytes[12], bytes[13], bytes[14], bytes[15] =
			bytes[15], bytes[14], bytes[13], bytes[12], bytes[11], bytes[10], bytes[9], bytes[8],
			bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]
	default:
		// For other sizes, do a generic byte reversal
		for i := 0; i < size/2; i++ {
			bytes[i], bytes[size-1-i] = bytes[size-1-i], bytes[i]
		}
	}
}

// isLittleEndian returns true if the platform is little-endian
func isLittleEndian() bool {
	var test uint32 = 0x01020304
	return *(*byte)(unsafe.Pointer(&test)) == 0x04
}

// Helper functions for endianness conversion (used by generated code)

// ToLittleEndian converts a value to little-endian representation
func ToLittleEndian[T any](slice *[]T) *[]T {
	if slice == nil {
		return nil
	}

	var zero T
	size := int(unsafe.Sizeof(zero))

	if !isLittleEndian() && size > 1 {
		swapEndianHardwareAccelerated(slice, LittleEndian, size)
	}

	return slice
}

// ToBigEndian converts a value to big-endian representation
func ToBigEndian[T any](slice *[]T) *[]T {
	if slice == nil {
		return nil
	}

	var zero T
	size := int(unsafe.Sizeof(zero))

	if isLittleEndian() && size > 1 {
		swapEndianHardwareAccelerated(slice, BigEndian, size)
	}

	return slice
}

// PutUint16 writes a uint16 to a byte slice in the specified byte order
func PutUint16(b *[]byte, v uint16, order int) {
	if b == nil || len(*b) < 2 {
		panic("moxie.PutUint16: slice too small")
	}

	s := *b
	if order == LittleEndian {
		binary.LittleEndian.PutUint16(s, v)
	} else {
		binary.BigEndian.PutUint16(s, v)
	}
}

// PutUint32 writes a uint32 to a byte slice in the specified byte order
func PutUint32(b *[]byte, v uint32, order int) {
	if b == nil || len(*b) < 4 {
		panic("moxie.PutUint32: slice too small")
	}

	s := *b
	if order == LittleEndian {
		binary.LittleEndian.PutUint32(s, v)
	} else {
		binary.BigEndian.PutUint32(s, v)
	}
}

// PutUint64 writes a uint64 to a byte slice in the specified byte order
func PutUint64(b *[]byte, v uint64, order int) {
	if b == nil || len(*b) < 8 {
		panic("moxie.PutUint64: slice too small")
	}

	s := *b
	if order == LittleEndian {
		binary.LittleEndian.PutUint64(s, v)
	} else {
		binary.BigEndian.PutUint64(s, v)
	}
}
