package main

import "fmt"

func greet(name string) string {
	return "Hello, " + name
}

func main() {
	var s string
	s = "World"

	message := greet(s)
	fmt.Println(message)
}
