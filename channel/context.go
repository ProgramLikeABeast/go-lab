package main

import (
	"context"
	"fmt"
	"time"
)

// =============================================================================
// CONTEXT PATTERN: Production-ready cancellation with timeouts
// =============================================================================

// slowOperation simulates a long-running task that respects context cancellation.
func slowOperation(ctx context.Context, id int) error {
	select {
	case <-time.After(500 * time.Millisecond): // Simulate work
		fmt.Printf("  Operation %d completed\n", id)
		return nil
	case <-ctx.Done():
		fmt.Printf("  Operation %d cancelled: %v\n", id, ctx.Err())
		return ctx.Err()
	}
}

func demoContext() {
	fmt.Println("\n=== CONTEXT PATTERN ===")
	fmt.Println("Use context.Context for cancellation, timeouts, and deadlines")

	// --- Context with timeout ---
	fmt.Println("\n1. Context with timeout (150ms for 500ms operation):")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel1()
	slowOperation(ctx1, 1)

	// --- Context with manual cancellation ---
	fmt.Println("\n2. Context with manual cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  Calling cancel()...")
		cancel2()
	}()

	slowOperation(ctx2, 2)

	// --- Context that completes successfully ---
	fmt.Println("\n3. Context with sufficient timeout (1s for 500ms operation):")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel3()
	slowOperation(ctx3, 3)
}
