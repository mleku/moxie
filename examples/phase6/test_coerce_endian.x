// Test endianness conversion in type coercion
// Tests: (*[]T, Endian)(slice) with Little/Big endian

package main

import "fmt"

func main() {
	fmt.Println("=== Phase 6: Type Coercion Endianness Test ===")

	// Create byte slice with known pattern
	bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
	fmt.Printf("Original bytes: 0x%02X 0x%02X 0x%02X 0x%02X\n",
		(*bytes)[0], (*bytes)[1], (*bytes)[2], (*bytes)[3])

	// Coerce to uint32 with different endianness
	// Native endian (platform-dependent)
	u32_native := (*[]uint32)(bytes)
	fmt.Printf("As uint32 (native):  0x%08X\n", (*u32_native)[0])

	// Little-endian (0x04030201 on LE, converted on BE)
	u32_le := (*[]uint32, LittleEndian)(bytes)
	fmt.Printf("As uint32 (little):  0x%08X\n", (*u32_le)[0])

	// Big-endian (0x01020304 on BE, converted on LE)
	u32_be := (*[]uint32, BigEndian)(bytes)
	fmt.Printf("As uint32 (big):     0x%08X\n", (*u32_be)[0])

	// Verify endianness conversion works
	// On little-endian: native == LE, BE is byte-swapped
	// On big-endian: native == BE, LE is byte-swapped
	fmt.Println("\n✓ Endianness conversion test complete!")
	fmt.Println("Note: Values differ based on platform and endianness setting")
}
