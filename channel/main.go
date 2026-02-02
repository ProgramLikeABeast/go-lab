// Package main demonstrates advanced channel patterns and concurrency in Go.
//
// This module builds on basic concurrency concepts to show real-world patterns:
// - Worker pools for parallel task processing
// - Fan-out/Fan-in for distributing and collecting work
// - Pipelines for data processing stages
// - Semaphores for limiting concurrency
// - Done channels for graceful cancellation
// - Context-based cancellation
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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

// =============================================================================
// WORKER POOL PATTERN: Process tasks with a fixed number of workers
// =============================================================================

// Job represents a unit of work.
type Job struct {
	ID      int
	Payload int
}

// Result represents the outcome of processing a job.
type Result struct {
	JobID  int
	Output int
}

// worker processes jobs from the jobs channel and sends results to results channel.
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(50)+10) * time.Millisecond)

		// Process the job (square the payload)
		results <- Result{
			JobID:  job.ID,
			Output: job.Payload * job.Payload,
		}
		fmt.Printf("  Worker %d processed job %d: %d² = %d\n",
			id, job.ID, job.Payload, job.Payload*job.Payload)
	}
}

func demoWorkerPool() {
	fmt.Println("\n=== WORKER POOL PATTERN ===")
	fmt.Println("Fixed number of workers process jobs from a shared queue")

	const numWorkers = 3
	const numJobs = 9

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	// Start workers
	fmt.Printf("Starting %d workers...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	fmt.Printf("Sending %d jobs...\n", numJobs)
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Payload: j * 2}
	}
	close(jobs) // Signal no more jobs

	// Wait for all workers to complete
	wg.Wait()
	close(results)

	// Collect results
	fmt.Println("\nAll results:")
	for result := range results {
		fmt.Printf("  Job %d -> %d\n", result.JobID, result.Output)
	}
}

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

// =============================================================================
// PIPELINE PATTERN: Chain processing stages
// =============================================================================

// generator produces numbers and sends them to a channel.
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square reads numbers, squares them, and sends to output channel.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// double reads numbers, doubles them, and sends to output channel.
func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * 2
		}
	}()
	return out
}

// filter only passes numbers greater than threshold.
func filter(in <-chan int, threshold int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n > threshold {
				out <- n
			}
		}
	}()
	return out
}

func demoPipeline() {
	fmt.Println("\n=== PIPELINE PATTERN ===")
	fmt.Println("Chain processing stages: generate -> square -> double -> filter")

	// Build the pipeline: generate -> square -> double -> filter(>20)
	// Each stage runs in its own goroutine!
	nums := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(nums)
	doubled := double(squared)
	filtered := filter(doubled, 20)

	// Consume the pipeline output
	fmt.Println("Results (n² × 2 where result > 20):")
	for result := range filtered {
		fmt.Printf("  %d\n", result)
	}
}

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

// =============================================================================
// Main
// =============================================================================

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
