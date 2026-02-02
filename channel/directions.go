package main

import "fmt"

// =============================================================================
// CHANNEL DIRECTIONS: Enforcing send-only or receive-only at compile time
// =============================================================================

// sendOnly demonstrates a send-only channel parameter.
// The function can ONLY send to this channel, not receive.
func sendOnly(ch chan<- int, value int) {
	ch <- value
	// _ = <-ch // Compile error: cannot receive from send-only channel
}

// receiveOnly demonstrates a receive-only channel parameter.
// The function can ONLY receive from this channel, not send.
func receiveOnly(ch <-chan int) int {
	return <-ch
	// ch <- 42 // Compile error: cannot send to receive-only channel
}

func demoChannelDirections() {
	fmt.Println("\n=== CHANNEL DIRECTIONS ===")
	fmt.Println("Use chan<- for send-only, <-chan for receive-only")

	ch := make(chan int, 1) // Bidirectional channel

	// When passed to functions, it becomes directional
	sendOnly(ch, 42)
	value := receiveOnly(ch)
	fmt.Printf("Sent and received: %d\n", value)

	fmt.Println("Channel directions prevent bugs at compile time!")
}
