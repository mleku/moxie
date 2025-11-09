package main

import "fmt"

func main() {
	// Channel pointer syntax with literal notation
	// &chan T{} creates an unbuffered channel
	ch1 := &chan int{}
	go func() {
		ch1 <- 42
	}()
	value := <-ch1
	fmt.Printf("Received from unbuffered channel: %d\n", value)

	// Buffered channel with capacity 2
	// &chan T{n} creates a buffered channel with capacity n
	ch2 := &chan string{2}
	ch2 <- "Hello"
	ch2 <- "World"
	fmt.Printf("Sent to buffered channel: %s, %s\n", <-ch2, <-ch2)
}

