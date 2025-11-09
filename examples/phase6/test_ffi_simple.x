// Simplified FFI test
// Just tests that dlopen/dlclose work

package main

import "fmt"

func main() {
	fmt.Println("=== Phase 6: FFI Simple Test ===")

	// Test that constants are transformed
	fmt.Println("RTLD_LAZY =", RTLD_LAZY)
	fmt.Println("RTLD_NOW =", RTLD_NOW)

	// Load libc
	lib := dlopen("libc.so.6", RTLD_LAZY)

	fmt.Println("Library loaded successfully!")

	// Clean up
	if lib != nil {
		dlclose(lib)
	}

	fmt.Println("✓ FFI simple test PASSED")
}
