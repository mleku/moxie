// Test string conversion functions
// Tests string(int), string(rune), string(*[]rune), and []rune(string)

package main

import "fmt"

func main() {
	fmt.Println("=== Phase 5: String Conversions Test ===")

	// Test 1: Int to string
	n := 42
	s1 := string(n)
	fmt.Print("Int to string: ")
	fmt.Println(*s1)  // Should print: 42

	// Test 2: Rune to string
	r := 'A'
	s2 := string(r)
	fmt.Print("Rune to string: ")
	fmt.Println(*s2)  // Should print: A

	// Test 3: Rune literal to string
	s3 := string('🎉')
	fmt.Print("Emoji rune to string: ")
	fmt.Println(*s3)  // Should print: 🎉

	// Test 4: Rune slice to string
	runes := &[]rune{'H', 'e', 'l', 'l', 'o'}
	s4 := string(*runes)
	fmt.Print("Rune slice to string: ")
	fmt.Println(*s4)  // Should print: Hello

	// Test 5: String to rune slice
	text := "Hello, 世界"
	runeSlice := *[]rune(text)
	fmt.Print("String to rune slice length: ")
	fmt.Println(len(*runeSlice))  // Should print: 9 (5 ASCII + 1 comma + 1 space + 2 Chinese chars)

	// Test 6: String to rune slice and back
	original := "Test 123"
	asRunes := *[]rune(original)
	backToString := string(*asRunes)
	fmt.Print("Round trip: ")
	fmt.Println(*backToString)  // Should print: Test 123

	// Test 7: Using moxie.Print for better output
	fmt.Println("\nUsing moxie.Print:")
	moxie.Print("Int 42 as string:", s1)
	moxie.Print("Rune 'A' as string:", s2)
	moxie.Print("Emoji:", s3)
	moxie.Print("From runes:", s4)
	moxie.Print("Original:", original)
	moxie.Print("Round trip:", backToString)

	// Test 8: Multiple conversions
	numbers := &[]int{1, 2, 3, 4, 5}
	fmt.Println("\nConverting numbers to strings:")
	for i := 0; i < len(*numbers); i++ {
		numStr := string((*numbers)[i])
		moxie.Print("  Number:", (*numbers)[i], "->", numStr)
	}

	fmt.Println("\n✓ String conversions test PASSED")
}
