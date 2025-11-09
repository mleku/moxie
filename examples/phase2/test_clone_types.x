package main

import "fmt"

type Person struct {
	Name *[]byte
	Age  int
}

func main() {
	// Test 1: Clone a slice
	numbers := &[]int{1, 2, 3, 4, 5}
	numbersClone := clone(numbers)
	(*numbersClone)[0] = 999
	fmt.Printf("Original slice: %v\n", *numbers)
	fmt.Printf("Cloned slice: %v\n", *numbersClone)

	// Test 2: Clone a map
	scores := &map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	scoresClone := clone(scores)
	(*scoresClone)["Alice"] = 100
	fmt.Printf("Original map: %v\n", *scores)
	fmt.Printf("Cloned map: %v\n", *scoresClone)

	// Test 3: Clone a struct
	person := &Person{
		Name: &[]byte{'J', 'o', 'h', 'n'},
		Age:  30,
	}
	personClone := clone(person)
	personClone.Age = 35
	(*personClone.Name)[0] = 'M' // Change first letter to M

	fmt.Printf("Original person: Name=%s, Age=%d\n", *person.Name, person.Age)
	fmt.Printf("Cloned person: Name=%s, Age=%d\n", *personClone.Name, personClone.Age)

	// Test 4: Clone string (which is *[]byte in Moxie)
	message := &[]byte{'H', 'e', 'l', 'l', 'o'}
	messageClone := clone(message)
	(*messageClone)[0] = 'Y'
	fmt.Printf("Original message: %s\n", *message)
	fmt.Printf("Cloned message: %s\n", *messageClone)

	// Test 5: Clone nested structures
	nested := &[][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	nestedClone := clone(nested)
	(*nestedClone)[0][0] = 999
	fmt.Printf("Original nested: %v\n", *nested)
	fmt.Printf("Cloned nested: %v\n", *nestedClone)

	fmt.Println("All clone type tests passed!")
}
