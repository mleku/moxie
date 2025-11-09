package main

import "fmt"

// Channel literal syntax with anonymous int64 field:
//   ch := &chan T{}      → make(chan T)      (unbuffered)
//   ch := &chan T{10}    → make(chan T, 10)  (buffered with capacity 10)
//   ch := &chan T{n}     → make(chan T, n)   (buffered with capacity n)
//
// The channel literal has one anonymous int64 field for buffer count.
// This is a direct 1:1 mapping to make(chan T, capacity).

func main() {
	// Unbuffered channel using new literal syntax
	ch1 := &chan int{}
	go func() {
		ch1 <- 42
	}()
	value := <-ch1
	fmt.Printf("Received from unbuffered channel: %d\n", value)

	// Buffered channel with capacity 2
	ch2 := &chan string{2}
	ch2 <- "Hello"
	ch2 <- "World"
	fmt.Printf("Sent to buffered channel: %s, %s\n", <-ch2, <-ch2)

	// Buffered channel with capacity 5
	ch3 := &chan int{5}
	ch3 <- 1
	ch3 <- 2
	ch3 <- 3
	ch3 <- 4
	ch3 <- 5
	sum := 0
	sum += <-ch3
	sum += <-ch3
	sum += <-ch3
	sum += <-ch3
	sum += <-ch3
	fmt.Printf("Sum from buffered channel: %d\n", sum)

	// Explicitly zero-buffered → unbuffered
	ch4 := &chan bool{0}
	go func() {
		ch4 <- true
	}()
	if <-ch4 {
		fmt.Println("Zero-buffered channel works (unbuffered)")
	}

	// Test send-only channel (for demonstration)
	sendCh := &chan<- int{1}
	sendCh <- 100
	fmt.Println("Send-only channel created successfully")

	// Test receive-only channel (for demonstration)
	recvCh := &<-chan string{1}
	go func(ch chan<- string) {
		ch <- "test"
	}(make(chan string, 1))
	// Note: recvCh would need to be properly initialized in real code
	_ = recvCh
	fmt.Println("Receive-only channel created successfully")
}
