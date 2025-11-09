package main

import "fmt"

func main() {
	// Test 1: Slice with explicit pointer (should work)
	s := &[]int{1, 2, 3}
	fmt.Println("Slice:", s)

	// Test 2: Map with explicit pointer (should work)
	m := &map[string]int{"one": 1, "two": 2}
	fmt.Println("Map:", m)

	// Test 3: make() should error
	// s2 := make([]int, 10)

	// Test 4: clear() should work (Go 1.21+)
	clear(m)
	fmt.Println("After clear:", m)
}
