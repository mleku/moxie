package main

import "fmt"

func main() {
	// Test empty string
	empty := ""
	fmt.Println("Empty string length:", len(*empty))
	fmt.Println("Empty == empty:", empty == empty)

	// Test concatenation with empty
	s1 := "Hello"
	result := s1 + empty
	fmt.Println("Hello + empty:", result)

	result2 := empty + s1
	fmt.Println("empty + Hello:", result2)

	// Test special characters
	special := "Tab:\t Newline:\n Quote:\" Backslash:\\"
	fmt.Println("Special chars:", special)

	// Test Unicode (if supported)
	unicode := "Hello 世界"
	fmt.Println("Unicode:", unicode)
	fmt.Println("Unicode length:", len(*unicode))

	// Test repeated concatenation
	repeated := ""
	repeated = repeated + "a"
	repeated = repeated + "b"
	repeated = repeated + "c"
	fmt.Println("Repeated concat:", repeated)

	// Test comparison with empty
	fmt.Println("empty < non-empty:", empty < s1)
	fmt.Println("non-empty > empty:", s1 > empty)
}
