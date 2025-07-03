package workerpool

import (
	"context"
	"sync"
	"time"
)

// Example task type
type intTask int

func ExampleNewWorkerPool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const workerNum = 3
	input := make(chan intTask, 10)
	var (
		mu      sync.Mutex
		results []intTask
		wg      sync.WaitGroup
	)

	// Worker function: reads from input and appends to results
	workerFunc := func(ctx context.Context) error {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case task, ok := <-input:
				if !ok {
					return nil
				}
				mu.Lock()
				results = append(results, task)
				mu.Unlock()
			}
		}
	}

	// Create worker pool
	pool := NewWorkerPool[intTask](ctx, workerNum, input, workerFunc)

	// Add tasks
	for i := 1; i <= 5; i++ {
		pool.AddTask(intTask(i))
	}

	// Allow some time for workers to process
	time.Sleep(200 * time.Millisecond)
	cancel() // Cancel context to stop workers

	// Wait for all workers to finish
	wg.Wait()

	// Check that all tasks were processed
	mu.Lock()
	defer mu.Unlock()

}
