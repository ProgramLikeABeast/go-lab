package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// FAN-OUT / FAN-IN PATTERN: Distribute work, then collect results
// =============================================================================

// fanOutWorker simulates processing and returns result on its own channel.
func fanOutWorker(id int, input int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
		result := input * input
		fmt.Printf("  Fan-out worker %d: %d² = %d\n", id, input, result)
		out <- result
	}()
	return out
}

// fanIn merges multiple channels into one.
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	// Close output channel when all inputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func demoFanOutFanIn() {
	fmt.Println("\n=== FAN-OUT / FAN-IN PATTERN ===")
	fmt.Println("Distribute work to multiple goroutines, merge results")

	inputs := []int{2, 3, 4, 5, 6}

	// Fan-out: start a worker for each input
	fmt.Println("Fan-out: starting workers...")
	var channels []<-chan int
	for i, input := range inputs {
		channels = append(channels, fanOutWorker(i+1, input))
	}

	// Fan-in: merge all results into one channel
	fmt.Println("Fan-in: collecting results...")
	merged := fanIn(channels...)

	// Collect all results
	var sum int
	for result := range merged {
		sum += result
	}
	fmt.Printf("Sum of squares: %d\n", sum)
}
