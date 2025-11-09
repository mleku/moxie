package main

import "fmt"

func main() {
	// Test string mutation - strings are mutable in Moxie
	s := "Hello"
	fmt.Println("Original:", s)

	// Modify a character
	(*s)[0] = 'J'
	fmt.Println("After modification:", s)

	// Test indexing
	fmt.Println("First char:", (*s)[0])
	fmt.Println("Last char:", (*s)[4])

	// Test length
	fmt.Println("Length:", len(*s))

	// Append to string
	*s = append(*s, '!')
	fmt.Println("After append:", s)

	// Test with slice operations
	s2 := "World"
	slice := (*s2)[0:3]
	fmt.Println("Slice [0:3]:", &slice)
}
