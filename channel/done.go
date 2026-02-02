package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// DONE CHANNEL PATTERN: Graceful cancellation
// =============================================================================

// infiniteWorker runs until done channel is closed.
func infiniteWorker(id int, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(30 * time.Millisecond)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-done:
			fmt.Printf("  Worker %d stopping after %d iterations\n", id, count)
			return
		case <-ticker.C:
			count++
			// Simulate doing periodic work
		}
	}
}

func demoDoneChannel() {
	fmt.Println("\n=== DONE CHANNEL PATTERN ===")
	fmt.Println("Use a done channel to signal goroutines to stop")

	done := make(chan struct{})
	var wg sync.WaitGroup

	// Start workers
	fmt.Println("Starting 3 infinite workers...")
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go infiniteWorker(i, done, &wg)
	}

	// Let them run for a bit
	time.Sleep(200 * time.Millisecond)

	// Signal all workers to stop by closing the done channel
	fmt.Println("Sending stop signal...")
	close(done) // Closing broadcasts to ALL receivers!

	wg.Wait()
	fmt.Println("All workers stopped gracefully")
}
