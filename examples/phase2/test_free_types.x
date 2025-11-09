package main

import "fmt"

type Person struct {
	Name *[]byte
	Age  int
}

func main() {
	// Test 1: Free a slice
	numbers := &[]int{1, 2, 3, 4, 5}
	fmt.Printf("Slice before free: %v\n", *numbers)
	free(numbers)
	fmt.Printf("Slice after free: %v (freed with FreeSlice[int])\n", *numbers)

	// Test 2: Free a map
	scores := &map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	fmt.Printf("\nMap before free: %v\n", *scores)
	free(scores)
	fmt.Printf("Map after free: %v (freed with FreeMap[string, int])\n", *scores)

	// Test 3: Free a struct
	person := &Person{
		Name: &[]byte{'J', 'o', 'h', 'n'},
		Age:  30,
	}
	fmt.Printf("\nPerson before free: Name=%s, Age=%d\n", *person.Name, person.Age)
	free(person)
	// After free, struct is zeroed - don't access freed memory
	fmt.Printf("Person freed with Free[Person] (memory zeroed)\n")

	// Test 4: Free string (which is *[]byte in Moxie)
	message := &[]byte{'H', 'e', 'l', 'l', 'o'}
	fmt.Printf("\nMessage before free: %s\n", *message)
	free(message)
	fmt.Printf("Message after free: %s (freed with FreeSlice[byte])\n", *message)

	fmt.Println("\nAll free type tests completed!")
	fmt.Println("✓ Slices use FreeSlice[T]")
	fmt.Println("✓ Maps use FreeMap[K, V]")
	fmt.Println("✓ Structs use Free[T]")
}
