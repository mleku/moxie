// Test basic type coercion functionality
// Tests: (*[]T)(slice) conversions

package main

func main() {
	moxie.Print("=== Phase 6: Type Coercion Basic Test ===")

	// Create byte slice
	bytes := &[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	moxie.Print("Original bytes:", bytes)

	// Coerce to uint16 (native endian)
	u16s := (*[]uint16)(bytes)
	moxie.Print("As uint16 (native):", u16s, "len=", len(*u16s))

	// Coerce to uint32 (native endian)
	u32s := (*[]uint32)(bytes)
	moxie.Print("As uint32 (native):", u32s, "len=", len(*u32s))

	// Coerce to uint64 (native endian)
	u64s := (*[]uint64)(bytes)
	moxie.Print("As uint64 (native):", u64s, "len=", len(*u64s))

	// Verify lengths
	if len(*u16s) == 4 && len(*u32s) == 2 && len(*u64s) == 1 {
		moxie.Print("✓ Type coercion basic test PASSED")
	} else {
		moxie.Print("✗ Type coercion basic test FAILED")
	}
}
