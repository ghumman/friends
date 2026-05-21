package concurrencyworkerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID string
	Target string
}

type Result struct {
	JobID string
	Success bool
	ErrorMsg string
	Duration time.Duration
}

func RunWorkerPool(ctx context.Context, jobs []Job, workerCount int) []Result {
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	var wg sync.WaitGroup

	// Start Workers
	for i:=0; i<workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i, jobChan, resultChan, &wg)
	}

	// Send all jobs to the channel
	for _, job := range jobs {
		select {
		case jobChan <- job:
		case <-ctx.Done():
			fmt.Println("Cancelled before sending all jobs")
			break
		}
	}

	close(jobChan)

	// Wait for workers to finish in a separate goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func worker (ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Cancelled while waiting\n", id)
			return
		default:
		}

		start := time.Now()
		err := processJob(job)
		duration := time.Since(start)

		result := Result {
			JobID: job.ID,
			Success: err == nil,
			Duration: duration,
		}

		if err != nil {
			result.ErrorMsg = err.Error()
		}

		select {
		case results <- result:
		case <-ctx.Done():
			return
		}
	}
}

func processJob(job Job) error {
	fmt.Printf("Processing job %s on target %s\n", job.ID, job.Target)
	time.Sleep(100 * time.Millisecond)
	return nil
}