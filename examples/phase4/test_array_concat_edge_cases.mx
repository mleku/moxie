package main

import "fmt"

func main() {
	// Test with empty slices
	empty1 := &[]int{}
	empty2 := &[]int{}
	bothEmpty := empty1 + empty2
	fmt.Println("Both empty:", len(*bothEmpty))

	// Test single element
	single := &[]int{42}
	nums := &[]int{1, 2, 3}
	withSingle := single + nums
	fmt.Println("Single + multiple:", withSingle)

	// Test large slices
	large1 := &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	large2 := &[]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	large := large1 + large2
	fmt.Println("Large concat length:", len(*large))

	// Test float slices
	floats1 := &[]float64{1.1, 2.2}
	floats2 := &[]float64{3.3, 4.4}
	floats := floats1 + floats2
	fmt.Println("Float concat:", floats)
}
