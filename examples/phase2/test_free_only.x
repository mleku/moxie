package main

import "fmt"

func main() {
	s := &[]int{1, 2, 3}
	fmt.Printf("Created slice: %v\n", *s)
	free(s)
	fmt.Println("Freed s")
}
