package main

import "fmt"

// Example Moxie file demonstrating syntax highlighting
const MaxConnections = 100

type Server struct {
	Name    *[]byte  // Mutable string in Moxie
	Port    int
	Active  bool
}

func main() {
	// Explicit pointer syntax for slices
	numbers := &[]int{1, 2, 3, 4, 5}

	// Moxie built-in functions
	grow(numbers, 10)
	backup := clone(numbers)
	defer free(backup)

	// String operations (strings are mutable in Moxie)
	message := &[]byte("Hello, Moxie!")
	message[0] = 'h'  // This is legal in Moxie!

	// Map with explicit pointer
	config := &map[string]int{
		"timeout":     30,
		"maxRetries":  3,
		"bufferSize":  1024,
	}

	// Traditional Go features
	for i, num := range *numbers {
		fmt.Printf("Index %d: %d\n", i, num)
	}

	// FFI example
	/*
	   Multi-line comment
	   showing Moxie's FFI capabilities
	*/
	lib := dlopen("libc.so.6", RTLD_LAZY)
	if lib != nil {
		defer dlclose(lib)
		strlen := dlsym[func(*byte) int64](lib, "strlen")
		length := strlen(&message[0])
		fmt.Printf("String length: %d\n", length)
	}

	// Type coercion with endianness
	bytes := &[]byte{0x01, 0x02, 0x03, 0x04}
	u32s := (*[]uint32)(bytes)  // Zero-copy coercion
	fmt.Printf("As uint32: %v\n", *u32s)
}

// Function with multiple return values
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Generic function example
func filter[T any](items *[]T, predicate func(T) bool) *[]T {
	result := &[]T{}
	for _, item := range *items {
		if predicate(item) {
			*result = append(*result, item)
		}
	}
	return result
}
