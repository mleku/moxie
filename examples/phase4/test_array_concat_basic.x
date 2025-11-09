package main

import "fmt"

func main() {
	// Test integer slice concatenation
	nums1 := &[]int{1, 2, 3}
	nums2 := &[]int{4, 5, 6}
	result := nums1 + nums2
	fmt.Println("Integer concat:", result)

	// Test string slice concatenation (NOT the same as string type)
	words1 := &[]string{"hello", "world"}
	words2 := &[]string{"foo", "bar"}
	combined := words1 + words2
	fmt.Println("String slice concat:", combined)

	// Test boolean slice concatenation
	bool1 := &[]bool{true, false}
	bool2 := &[]bool{false, true}
	bools := bool1 + bool2
	fmt.Println("Bool concat:", bools)

	// Test empty slice concatenation
	empty := &[]int{}
	moreNums := &[]int{7, 8, 9}
	withEmpty := empty + moreNums
	fmt.Println("Empty + nums:", withEmpty)
}
