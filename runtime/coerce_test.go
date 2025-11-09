package runtime

import (
	"encoding/binary"
	"testing"
	"unsafe"
)

func TestCoerceBasic(t *testing.T) {
	// Test basic coercion from bytes to uint32
	bytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	result := Coerce[byte, uint32](&bytes)

	if result == nil {
		t.Fatal("Coerce returned nil")
	}

	if len(*result) != 2 {
		t.Errorf("Expected length 2, got %d", len(*result))
	}

	// Verify it's zero-copy (same backing array)
	resultPtr := unsafe.Pointer(unsafe.SliceData(*result))
	bytesPtr := unsafe.Pointer(unsafe.SliceData(bytes))
	if resultPtr != bytesPtr {
		t.Error("Coerce did not create zero-copy slice")
	}
}

func TestCoerceEndianness(t *testing.T) {
	// Test endianness conversion
	bytes := []byte{0x01, 0x02, 0x03, 0x04}

	// Test native endian (no conversion)
	nativeResult := Coerce[byte, uint32](&bytes, NativeEndian)
	if nativeResult == nil {
		t.Fatal("Coerce with NativeEndian returned nil")
	}

	// Test little endian
	bytesLE := []byte{0x01, 0x02, 0x03, 0x04}
	leResult := Coerce[byte, uint32](&bytesLE, LittleEndian)
	if leResult == nil {
		t.Fatal("Coerce with LittleEndian returned nil")
	}

	expectedLE := binary.LittleEndian.Uint32([]byte{0x01, 0x02, 0x03, 0x04})
	if (*leResult)[0] != expectedLE {
		t.Errorf("Little endian conversion failed: expected 0x%08X, got 0x%08X", expectedLE, (*leResult)[0])
	}

	// Test big endian
	bytesBE := []byte{0x01, 0x02, 0x03, 0x04}
	beResult := Coerce[byte, uint32](&bytesBE, BigEndian)
	if beResult == nil {
		t.Fatal("Coerce with BigEndian returned nil")
	}

	expectedBE := binary.BigEndian.Uint32([]byte{0x01, 0x02, 0x03, 0x04})
	if (*beResult)[0] != expectedBE {
		t.Errorf("Big endian conversion failed: expected 0x%08X, got 0x%08X", expectedBE, (*beResult)[0])
	}
}

func TestCoerceUint16(t *testing.T) {
	// Test 16-bit endianness conversion
	bytes := []byte{0x01, 0x02, 0x03, 0x04}

	leResult := Coerce[byte, uint16](&bytes, LittleEndian)
	if leResult == nil {
		t.Fatal("Coerce uint16 LittleEndian returned nil")
	}

	if len(*leResult) != 2 {
		t.Errorf("Expected length 2, got %d", len(*leResult))
	}

	// Verify endianness
	expectedLE0 := binary.LittleEndian.Uint16([]byte{0x01, 0x02})
	expectedLE1 := binary.LittleEndian.Uint16([]byte{0x03, 0x04})

	if (*leResult)[0] != expectedLE0 {
		t.Errorf("uint16[0] LE: expected 0x%04X, got 0x%04X", expectedLE0, (*leResult)[0])
	}
	if (*leResult)[1] != expectedLE1 {
		t.Errorf("uint16[1] LE: expected 0x%04X, got 0x%04X", expectedLE1, (*leResult)[1])
	}
}

func TestCoerceUint64(t *testing.T) {
	// Test 64-bit endianness conversion
	bytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	leResult := Coerce[byte, uint64](&bytes, LittleEndian)
	if leResult == nil {
		t.Fatal("Coerce uint64 LittleEndian returned nil")
	}

	if len(*leResult) != 1 {
		t.Errorf("Expected length 1, got %d", len(*leResult))
	}

	// Verify endianness
	expectedLE := binary.LittleEndian.Uint64(bytes)

	if (*leResult)[0] != expectedLE {
		t.Errorf("uint64 LE: expected 0x%016X, got 0x%016X", expectedLE, (*leResult)[0])
	}
}

func TestCoerceNilSlice(t *testing.T) {
	// Test nil handling
	var nilSlice *[]byte = nil
	result := Coerce[byte, uint32](nilSlice)

	if result != nil {
		t.Error("Coerce with nil slice should return nil")
	}
}

func TestToLittleEndian(t *testing.T) {
	// Test ToLittleEndian helper
	values := []uint32{0x01020304, 0x05060708}
	result := ToLittleEndian(&values)

	if result == nil {
		t.Fatal("ToLittleEndian returned nil")
	}

	// The values should now be in little endian format
	// On a little-endian platform, they should be unchanged
	// On a big-endian platform, they should be byte-swapped
	if isLittleEndian() {
		if (*result)[0] != 0x01020304 {
			t.Errorf("ToLittleEndian on LE platform: expected 0x01020304, got 0x%08X", (*result)[0])
		}
	} else {
		if (*result)[0] != 0x04030201 {
			t.Errorf("ToLittleEndian on BE platform: expected 0x04030201, got 0x%08X", (*result)[0])
		}
	}
}

func TestToBigEndian(t *testing.T) {
	// Test ToBigEndian helper
	values := []uint32{0x01020304, 0x05060708}
	result := ToBigEndian(&values)

	if result == nil {
		t.Fatal("ToBigEndian returned nil")
	}

	// The values should now be in big endian format
	// On a big-endian platform, they should be unchanged
	// On a little-endian platform, they should be byte-swapped
	if !isLittleEndian() {
		if (*result)[0] != 0x01020304 {
			t.Errorf("ToBigEndian on BE platform: expected 0x01020304, got 0x%08X", (*result)[0])
		}
	} else {
		if (*result)[0] != 0x04030201 {
			t.Errorf("ToBigEndian on LE platform: expected 0x04030201, got 0x%08X", (*result)[0])
		}
	}
}

func BenchmarkCoerceNoEndian(b *testing.B) {
	bytes := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = Coerce[byte, uint32](&bytes)
	}
}

func BenchmarkCoerceLittleEndian(b *testing.B) {
	bytes := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = Coerce[byte, uint32](&bytes, LittleEndian)
	}
}

func BenchmarkCoerceBigEndian(b *testing.B) {
	bytes := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = Coerce[byte, uint32](&bytes, BigEndian)
	}
}
