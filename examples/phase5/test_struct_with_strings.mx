package main

import moxie "github.com/mleku/moxie/runtime"

type Person struct {
	Name string
	Age  int
}

func main() {
	// This now works without workarounds! (Phase 5 fix)
	p1 := &Person{Name: "Alice", Age: 30}
	p2 := &Person{Name: "Bob", Age: 25}
	p3 := &Person{Name: "Charlie", Age: 35}

	moxie.Print("Person 1:", p1.Name, "Age:", p1.Age)
	moxie.Print("Person 2:", p2.Name, "Age:", p2.Age)
	moxie.Print("Person 3:", p3.Name, "Age:", p3.Age)

	// Test with slice of structs
	people := &[]Person{
		{Name: "Dave", Age: 40},
		{Name: "Eve", Age: 28},
	}

	moxie.Print("Number of people:", len(*people))
	for i, person := range *people {
		moxie.Print("  ", i, ":", person.Name, "-", person.Age)
	}
}
