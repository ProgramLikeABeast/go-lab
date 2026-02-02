package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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
