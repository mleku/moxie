// Test basic FFI functionality
// Tests: dlopen, dlsym, dlclose, dlerror

package main

import "fmt"

func main() {
	fmt.Println("=== Phase 6: FFI Basic Test ===")

	// Load libc
	lib := dlopen("libc.so.6", RTLD_LAZY)
	if lib == nil {
		fmt.Println("Error loading library:", dlerror())
		return
	}
	defer dlclose(lib)

	fmt.Println("Library loaded successfully")

	// Look up strlen function
	// strlen has signature: size_t strlen(const char *s)
	// In Moxie: func(*byte) int64
	strlen := dlsym[func(*byte) int64](lib, "strlen")

	// Test strlen with a C string (null-terminated)
	msg := "Hello FFI\x00"
	length := strlen(&msg[0])

	fmt.Println("String:", msg[:len(msg)-1])
	fmt.Println("Length (via strlen):", length)

	if length == 9 {
		fmt.Println("✓ FFI basic test PASSED")
	} else {
		fmt.Println("✗ FFI basic test FAILED")
	}
}
