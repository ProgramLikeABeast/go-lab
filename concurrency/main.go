package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// GOROUTINES: Lightweight threads managed by Go runtime
// =============================================================================

func sayHello(id int) {
	fmt.Printf("Goroutine %d: Hello!\n", id)
}

func demoGoroutines() {
	fmt.Println("\n=== GOROUTINES ===")

	// Start 5 goroutines - they run concurrently
	for i := 1; i <= 5; i++ {
		go sayHello(i) // "go" keyword spawns a goroutine
	}

	// Without this sleep, main() might exit before goroutines finish
	// (we'll fix this properly with WaitGroup below)
	time.Sleep(100 * time.Millisecond)
}

// =============================================================================
// WAITGROUP: Wait for goroutines to complete
// =============================================================================

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement counter when goroutine completes

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Duration(id*50) * time.Millisecond) // Simulate work
	fmt.Printf("Worker %d done\n", id)
}

func demoWaitGroup() {
	fmt.Println("\n=== WAITGROUP ===")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)      // Increment counter before starting goroutine
		go worker(i, &wg)
	}

	wg.Wait() // Block until counter reaches 0
	fmt.Println("All workers completed")
}

// =============================================================================
// CHANNELS: Communication between goroutines (Go's CSP model)
// =============================================================================

func demoChannels() {
	fmt.Println("\n=== CHANNELS (unbuffered) ===")

	// Unbuffered channel: sender blocks until receiver is ready
	ch := make(chan string)

	go func() {
		ch <- "Hello from goroutine!" // Send blocks until someone receives
	}()

	msg := <-ch // Receive blocks until someone sends
	fmt.Println("Received:", msg)

	// Buffered channel: can hold N values before blocking
	fmt.Println("\n=== CHANNELS (buffered) ===")
	buffered := make(chan int, 3) // Buffer size 3

	buffered <- 1 // Won't block
	buffered <- 2 // Won't block
	buffered <- 3 // Won't block
	// buffered <- 4 // Would block! Buffer full

	fmt.Println(<-buffered, <-buffered, <-buffered)
}

// =============================================================================
// CHANNEL PATTERNS: Range and Close
// =============================================================================

func producer(ch chan<- int) { // chan<- means send-only channel
	for i := 1; i <= 5; i++ {
		ch <- i * 10
	}
	close(ch) // Signal no more values; receivers can detect this
}

func demoRangeClose() {
	fmt.Println("\n=== RANGE OVER CHANNEL ===")

	ch := make(chan int)
	go producer(ch)

	// Range automatically exits when channel is closed
	for value := range ch {
		fmt.Printf("Received: %d\n", value)
	}
	fmt.Println("Channel closed, loop exited")
}

// =============================================================================
// SELECT: Multiplex multiple channel operations
// =============================================================================

func demoSelect() {
	fmt.Println("\n=== SELECT ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// Select waits on multiple channels, executes first ready case
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received", msg2)
		}
	}

	// Select with timeout
	fmt.Println("\n=== SELECT WITH TIMEOUT ===")
	ch3 := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch3 <- "slow operation"
	}()

	select {
	case msg := <-ch3:
		fmt.Println("Got:", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout! Operation took too long")
	}

	// Select with default (non-blocking)
	fmt.Println("\n=== SELECT WITH DEFAULT (non-blocking) ===")
	ch4 := make(chan int)

	select {
	case val := <-ch4:
		fmt.Println("Got:", val)
	default:
		fmt.Println("No value available, moving on")
	}
}

// =============================================================================
// MUTEX: Protect shared state
// =============================================================================

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func demoMutex() {
	fmt.Println("\n=== MUTEX ===")

	counter := SafeCounter{}
	var wg sync.WaitGroup

	// 100 goroutines incrementing the same counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d (should be 100)\n", counter.Value())
}

func main() {
	fmt.Println("========================================")
	fmt.Println("GO CONCURRENCY FUNDAMENTALS")
	fmt.Println("========================================")

	demoGoroutines()
	demoWaitGroup()
	demoChannels()
	demoRangeClose()
	demoSelect()
	demoMutex()
}
