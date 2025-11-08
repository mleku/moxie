package main

import "fmt"

func main() {
	// Test grow()
	s := &[]int{1, 2, 3}
	fmt.Printf("Before grow: len=%d cap=%d\n", len(*s), cap(*s))
	s = grow(s, 100)
	fmt.Printf("After grow: len=%d cap=%d\n", len(*s), cap(*s))

	// Test clone()
	s2 := clone(s)
	fmt.Printf("Original: %v\n", s)
	fmt.Printf("Clone: %v\n", s2)
	*s2 = append(*s2, 999)
	fmt.Printf("After modifying clone - Original: %v, Clone: %v\n", s, s2)

	// Test free()
	free(s2)
	fmt.Println("Freed s2")
}
