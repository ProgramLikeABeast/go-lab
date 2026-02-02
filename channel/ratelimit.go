package main

import (
	"fmt"
	"time"
)

// =============================================================================
// RATE LIMITING PATTERN: Control throughput
// =============================================================================

func demoRateLimiting() {
	fmt.Println("\n=== RATE LIMITING PATTERN ===")
	fmt.Println("Use time.Ticker to control throughput")

	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)

	// Process at most 5 requests per second (200ms interval)
	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()

	fmt.Println("Processing requests (rate limited to 5/second):")
	start := time.Now()

	for req := range requests {
		<-limiter.C // Wait for rate limiter
		fmt.Printf("  Request %d processed at %v\n", req, time.Since(start).Round(time.Millisecond))
	}
}
