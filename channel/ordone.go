package main

import (
	"fmt"
	"time"
)

// =============================================================================
// OR-DONE PATTERN: Wrap channel reads with cancellation
// =============================================================================

// orDone wraps a channel to respect cancellation via done channel.
func orDone(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

func demoOrDone() {
	fmt.Println("\n=== OR-DONE PATTERN ===")
	fmt.Println("Safely read from channels while respecting cancellation")

	// Create a channel that produces values slowly
	source := make(chan int)
	go func() {
		defer close(source)
		for i := 1; i <= 100; i++ {
			source <- i
			time.Sleep(20 * time.Millisecond)
		}
	}()

	done := make(chan struct{})

	// Cancel after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(done)
	}()

	// Use orDone to safely iterate
	fmt.Println("Reading values (will be cancelled after ~100ms):")
	for val := range orDone(done, source) {
		fmt.Printf("  Got: %d\n", val)
	}
	fmt.Println("Channel iteration stopped due to cancellation")
}
