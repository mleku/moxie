// Test zero-copy network protocol parsing
// Demonstrates practical use of type coercion for binary protocols

package main

import "fmt"

type PacketHeader struct {
	Magic   uint32
	Version uint16
	Length  uint16
	CRC     uint32
}

func parseHeader(data *[]byte) PacketHeader {
	// Network protocols use big-endian byte order
	// Use endianness tuple syntax to specify BigEndian conversion
	fields := (*[]uint32, BigEndian)(data)

	// Extract fields from the uint32 array
	magic := (*fields)[0]
	verLen := (*fields)[1]
	crc := (*fields)[2]

	return PacketHeader{
		Magic:   magic,
		Version: uint16(verLen >> 16),
		Length:  uint16(verLen & 0xFFFF),
		CRC:     crc,
	}
}

func main() {
	fmt.Println("=== Phase 6: Network Protocol Parsing Test ===")

	// Simulated network packet (12 bytes)
	packet := &[]byte{
		0x12, 0x34, 0x56, 0x78, // Magic: 0x12345678
		0x00, 0x01, 0x02, 0x00, // Version: 1, Length: 512 (0x0200)
		0xAB, 0xCD, 0xEF, 0x12, // CRC: 0xABCDEF12
	}

	fmt.Printf("Packet bytes: %v\n", *packet)

	// Parse header using zero-copy coercion
	header := parseHeader(packet)

	fmt.Printf("Magic:   0x%08X\n", header.Magic)
	fmt.Printf("Version: %d\n", header.Version)
	fmt.Printf("Length:  %d\n", header.Length)
	fmt.Printf("CRC:     0x%08X\n", header.CRC)

	// Note: BigEndian conversion ensures consistent interpretation regardless of platform
	// This demonstrates zero-copy parsing with automatic endianness handling
	fmt.Println("\n✓ Zero-copy type coercion with endianness conversion")
	fmt.Println("Note: Big-endian conversion applied for network byte order")
}
