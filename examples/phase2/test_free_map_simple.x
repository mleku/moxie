package main

import "fmt"

func getMap() *map[int]int {
	return &map[int]int{1: 100}
}

func freeMapParam(m *map[int]int) {
	fmt.Printf("  In function before free: %v\n", *m)
	free(m)
	fmt.Printf("  In function after free: %v\n", *m)
}

func main() {
	// Test 1: Map declared with var
	var varMap *map[int]int
	varMap = &map[int]int{1: 10, 2: 20}
	fmt.Printf("Test 1 - varMap before free: %v\n", *varMap)
	free(varMap)
	fmt.Printf("Test 1 - varMap after free: %v\n\n", *varMap)

	// Test 2: Map from function return
	returnedMap := getMap()
	fmt.Printf("Test 2 - returnedMap before free: %v\n", *returnedMap)
	free(returnedMap)
	fmt.Printf("Test 2 - returnedMap after free: %v\n\n", *returnedMap)

	// Test 3: Map as function parameter
	paramMap := &map[int]int{1: 100, 2: 200}
	fmt.Printf("Test 3 - paramMap before freeMapParam:\n")
	freeMapParam(paramMap)
	fmt.Printf("Test 3 - paramMap after freeMapParam: %v\n\n", *paramMap)

	// Test 4: Multiple maps with different types
	intIntMap := &map[int]int{1: 1}
	intBoolMap := &map[int]bool{1: true, 2: false}
	fmt.Printf("Test 4 - intIntMap: %v, intBoolMap: %v\n", *intIntMap, *intBoolMap)
	free(intIntMap)
	free(intBoolMap)
	fmt.Printf("Test 4 - After free: intIntMap=%v, intBoolMap=%v\n\n", *intIntMap, *intBoolMap)

	fmt.Println("✓ All map free tests completed!")
}
