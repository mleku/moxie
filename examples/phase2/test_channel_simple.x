package main

import "fmt"

func main() {
	// Test all channel literal syntax variants

	// 1. Unbuffered bidirectional channel
	ch1 := &chan int{}
	go func() { ch1 <- 1 }()
	fmt.Printf("Unbuffered: %d\n", <-ch1)

	// 2. Buffered bidirectional channel
	ch2 := &chan string{2}
	ch2 <- "A"
	ch2 <- "B"
	fmt.Printf("Buffered: %s %s\n", <-ch2, <-ch2)

	// 3. Send-only channel
	sendCh := &chan<- int{1}
	go func() { sendCh <- 99 }()
	fmt.Println("Send-only channel created")

	// 4. Receive-only channel (created from bidirectional)
	biCh := &chan bool{1}
	var recvCh <-chan bool = biCh
	go func() { biCh <- true }()
	fmt.Printf("Receive-only: %v\n", <-recvCh)

	fmt.Println("All channel literal tests passed!")
}
