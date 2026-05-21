package concurrencyworkerpool

import (
	"context"
	"sync"
	"time"
)

type NewJob struct {
	ID string
}

type NewResult struct {
	JobID string
	Success bool
}

func GoWorkerPool(ctx context.Context, jobs []NewJob) []NewResult {
	inputChan := make(chan(NewJob), len(jobs))
	outputChan := make(chan(NewResult), len(jobs))

	totalWorkers := 10
	var wg sync.WaitGroup

	for i:=0; i<totalWorkers; i++ {
		wg.Add(1)
		go NewWorker(ctx, inputChan, outputChan, &wg)
	}

	for _, job := range jobs {
		select {
		case inputChan <- job:
		case <-ctx.Done():
			return []NewResult{}
		}
	}

	go func() {
		wg.Wait()
	}()

	var results []NewResult
	for result := range outputChan {
		results = append(results, result)
	}
	return results
}

func NewWorker(ctx context.Context, inputChan <-chan NewJob, outputChan chan<- NewResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for job := range inputChan {
		select {
		case <- ctx.Done():
			return
		default:
		}

		result := MimicWork(job)
		select {
		case <- ctx.Done():
			return
		case outputChan <- result:
		}

	}
}

func MimicWork(job NewJob) NewResult {
	time.Sleep(5 * time.Millisecond)
	return NewResult{
		JobID: job.ID,
		Success: true,
	}
}