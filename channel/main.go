// Package main demonstrates advanced channel patterns and concurrency in Go.
//
// This module builds on basic concurrency concepts to show real-world patterns:
// - Worker pools for parallel task processing
// - Fan-out/Fan-in for distributing and collecting work
// - Pipelines for data processing stages
// - Semaphores for limiting concurrency
// - Done channels for graceful cancellation
// - Context-based cancellation
//
// Each pattern is in its own file:
//   - directions.go  : Channel directions (send-only, receive-only)
//   - workerpool.go  : Worker pool pattern
//   - fanout.go      : Fan-out/Fan-in pattern
//   - pipeline.go    : Pipeline pattern
//   - semaphore.go   : Semaphore pattern
//   - done.go        : Done channel pattern
//   - context.go     : Context pattern
//   - ratelimit.go   : Rate limiting pattern
//   - ordone.go      : Or-done pattern
package main

import "fmt"

func main() {
	fmt.Println("========================================")
	fmt.Println("ADVANCED CHANNEL PATTERNS IN GO")
	fmt.Println("========================================")

	demoChannelDirections()
	demoWorkerPool()
	demoFanOutFanIn()
	demoPipeline()
	demoSemaphore()
	demoDoneChannel()
	demoContext()
	demoRateLimiting()
	demoOrDone()

	fmt.Println("\n========================================")
	fmt.Println("KEY TAKEAWAYS")
	fmt.Println("========================================")
	fmt.Println("1. Channel directions (chan<-, <-chan) prevent bugs at compile time")
	fmt.Println("2. Worker pools efficiently process jobs with fixed resources")
	fmt.Println("3. Fan-out/Fan-in distributes work and collects results")
	fmt.Println("4. Pipelines chain processing stages elegantly")
	fmt.Println("5. Buffered channels make great semaphores")
	fmt.Println("6. Done channels broadcast cancellation to multiple goroutines")
	fmt.Println("7. context.Context is the standard for cancellation in production")
	fmt.Println("8. Rate limiting controls throughput with time.Ticker")
	fmt.Println("9. Or-done pattern safely combines channels with cancellation")
}
