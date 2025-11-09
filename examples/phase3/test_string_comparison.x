package main

import "fmt"

func main() {
	s1 := "hello"
	s2 := "hello"
	s3 := "world"

	// Test equality
	fmt.Println("s1 == s2:", s1 == s2) // Should be true
	fmt.Println("s1 == s3:", s1 == s3) // Should be false
	fmt.Println("s1 != s3:", s1 != s3) // Should be true

	// Test ordering
	fmt.Println("s1 < s3:", s1 < s3)   // Should be true (hello < world)
	fmt.Println("s3 > s1:", s3 > s1)   // Should be true (world > hello)
	fmt.Println("s1 <= s2:", s1 <= s2) // Should be true (equal)
	fmt.Println("s1 >= s2:", s1 >= s2) // Should be true (equal)
}
