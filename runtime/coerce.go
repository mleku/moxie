package runtime

import (
	"encoding/binary"
	"reflect"
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

	// Create new slice header pointing to same backing array
	srcHeader := (*reflect.SliceHeader)(unsafe.Pointer(src))
	dstHeader := reflect.SliceHeader{
		Data: srcHeader.Data,
		Len:  dstLen,
		Cap:  dstCap,
	}

	result := (*[]To)(unsafe.Pointer(&dstHeader))

	// Apply endianness conversion if needed
	if byteOrder != NativeEndian && toSize > 1 {
		swapEndian(result, byteOrder, toSize)
	}

	return result
}

// swapEndian performs in-place endianness conversion
func swapEndian[T any](slice *[]T, targetEndian int, elemSize int) {
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

	// Perform byte swapping based on element size
	s := *slice
	for i := range s {
		swapBytes(unsafe.Pointer(&s[i]), elemSize)
	}
}

// swapBytes performs in-place byte swapping for different sizes
func swapBytes(ptr unsafe.Pointer, size int) {
	bytes := (*[8]byte)(ptr)

	switch size {
	case 2:
		bytes[0], bytes[1] = bytes[1], bytes[0]
	case 4:
		bytes[0], bytes[1], bytes[2], bytes[3] = bytes[3], bytes[2], bytes[1], bytes[0]
	case 8:
		bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7] =
			bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]
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
		swapEndian(slice, LittleEndian, size)
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
		swapEndian(slice, BigEndian, size)
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
