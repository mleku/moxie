package main

import (
	"fmt"
	moxie "github.com/mleku/moxie/runtime"
)

func main() {
	// Test with old fmt.Println (shows byte arrays)
	s := "Hello, World!"
	fmt.Println("Using fmt.Println:", s)

	// Test with new moxie.Print (shows strings)
	moxie.Print("Using moxie.Print:", s)

	// Test with multiple strings
	greeting := "Hello"
	name := "Moxie"
	moxie.Print(greeting, name)

	// Test with mixed types
	count := 42
	message := "The answer is"
	moxie.Print(message, count)

	// Test Printf
	moxie.Printf("Formatted: %s has %d characters\n", s, len(*s))
}
