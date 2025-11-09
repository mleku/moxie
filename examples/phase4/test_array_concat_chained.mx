package main

import "fmt"

func main() {
	// Test chained concatenation
	a := &[]int{1, 2}
	b := &[]int{3, 4}
	c := &[]int{5, 6}
	result := a + b + c
	fmt.Println("Chained concat:", result)

	// Test mixed chained concatenation
	nums := &[]int{10}
	more := &[]int{20, 30}
	evenMore := &[]int{40, 50, 60}
	all := nums + more + evenMore
	fmt.Println("Mixed chain:", all)

	// Test that strings still work
	s1 := "hello"
	s2 := " "
	s3 := "world"
	greeting := s1 + s2 + s3
	fmt.Println("String concat still works:", greeting)
}
