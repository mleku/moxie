package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	// Test with custom struct slices
	people1 := &[]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	people2 := &[]Person{
		{Name: "Charlie", Age: 35},
	}
	allPeople := people1 + people2
	fmt.Println("Struct concat:", allPeople)

	// Test with pointer slices
	p1 := &Person{Name: "Dave", Age: 40}
	p2 := &Person{Name: "Eve", Age: 28}
	p3 := &Person{Name: "Frank", Age: 33}

	ptrs1 := &[]*Person{p1, p2}
	ptrs2 := &[]*Person{p3}
	allPtrs := ptrs1 + ptrs2
	fmt.Println("Pointer slice concat length:", len(*allPtrs))

	// Test with float slices
	floats1 := &[]float64{1.1, 2.2, 3.3}
	floats2 := &[]float64{4.4, 5.5}
	allFloats := floats1 + floats2
	fmt.Println("Float concat:", allFloats)
}
