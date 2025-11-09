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
	fmt.Printf("Cloned\n")

	// Test free()
	free(s2)
	fmt.Println("Freed s2")
}
