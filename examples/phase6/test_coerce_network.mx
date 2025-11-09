// Test zero-copy network protocol parsing
// Demonstrates practical use of type coercion for binary protocols

package main

type PacketHeader struct {
	Magic   uint32
	Version uint16
	Length  uint16
	CRC     uint32
}

func parseHeader(data *[]byte) PacketHeader {
	// Network protocols use big-endian byte order
	// Note: Endianness syntax requires parser extension
	// For now, using native endian (demonstration purposes)
	fields := (*[]uint32)(data)

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
	moxie.Print("=== Phase 6: Network Protocol Parsing Test ===")

	// Simulated network packet (12 bytes)
	packet := &[]byte{
		0x12, 0x34, 0x56, 0x78, // Magic: 0x12345678
		0x00, 0x01, 0x02, 0x00, // Version: 1, Length: 512 (0x0200)
		0xAB, 0xCD, 0xEF, 0x12, // CRC: 0xABCDEF12
	}

	moxie.Print("Packet bytes:", packet)

	// Parse header using zero-copy coercion
	header := parseHeader(packet)

	moxie.Printf("Magic:   0x%08X\n", header.Magic)
	moxie.Printf("Version: %d\n", header.Version)
	moxie.Printf("Length:  %d\n", header.Length)
	moxie.Printf("CRC:     0x%08X\n", header.CRC)

	// Note: Values will be different on little-endian vs big-endian platforms
	// This demonstrates zero-copy parsing works, but endianness requires parser extension
	moxie.Print("✓ Zero-copy type coercion demonstrated")
	moxie.Print("Note: Proper endianness conversion requires parser extension")
}
