package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// SEMAPHORE PATTERN: Limit concurrent operations
// =============================================================================

func demoSemaphore() {
	fmt.Println("\n=== SEMAPHORE PATTERN ===")
	fmt.Println("Limit concurrency using a buffered channel as semaphore")

	const maxConcurrent = 3
	const totalTasks = 10

	// Buffered channel acts as a semaphore
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	fmt.Printf("Running %d tasks with max %d concurrent...\n", totalTasks, maxConcurrent)

	for i := 1; i <= totalTasks; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			// Acquire semaphore (blocks if buffer is full)
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore

			fmt.Printf("  Task %d started\n", taskID)
			time.Sleep(50 * time.Millisecond) // Simulate work
			fmt.Printf("  Task %d finished\n", taskID)
		}(i)
	}

	wg.Wait()
	fmt.Println("All tasks completed")
}
