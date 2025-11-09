// Test const mutation detection
// This file SHOULD produce const enforcement errors

package main

import "fmt"

const MaxConnections = 100
const ServerName = "Moxie Server"

func main() {
	fmt.Println("Testing const mutation detection...")

	// These lines SHOULD trigger const enforcement errors:
	MaxConnections = 200  // ERROR: cannot assign to const
	ServerName = "Other" // ERROR: cannot assign to const

	fmt.Println("If you see this, const enforcement failed!")
}
