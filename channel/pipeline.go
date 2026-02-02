package main

import "fmt"

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
