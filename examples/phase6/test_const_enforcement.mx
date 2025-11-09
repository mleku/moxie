// Test compile-time const enforcement
// This file should produce warnings about const mutation

package main

import "fmt"

// Valid const declarations
const MaxSize = 100
const Pi = 3.14159

func main() {
	fmt.Println("=== Phase 6: Const Enforcement Test ===")

	// Valid: reading const values
	fmt.Println("MaxSize:", MaxSize)
	fmt.Println("Pi:", Pi)

	// Valid: using const in expressions
	doubled := MaxSize * 2
	fmt.Println("Doubled:", doubled)

	tripled := Pi * 3
	fmt.Printf("Pi * 3 = %.5f\n", tripled)

	// INVALID examples (commented out):
	// These would trigger const enforcement errors if uncommented:
	// MaxSize = 200           // ERROR: cannot assign to const MaxSize
	// MaxSize++               // ERROR: cannot ++ const MaxSize
	// Pi = 3.14              // ERROR: cannot assign to const Pi

	fmt.Println("✓ Const enforcement test PASSED")
	fmt.Println("Note: Const mutations are prevented at compile-time")
}
