package workerpool

import "context"

type workerFunc func(ctx context.Context) error

type workerPool[T any] struct {
	input chan T
}

// NewWorkerPool creates a new worker pool with a specified number of workers.
// It takes a context, number of workers, an input channel, and a worker function.
// The worker pool starts workers that process tasks from the input channel.
// When the context is cancelled, the input channel is closed.
//
// Parameters:
//   - ctx: Context for controlling the worker pool's lifecycle
//   - workerNum: Number of concurrent workers to process tasks
//   - input: Channel for sending tasks to workers
//   - workerFunc: Function to be executed by each worker
//
// Returns a pointer to the created worker pool.
func NewWorkerPool[T any](ctx context.Context, workerNum int, input chan T, workerFunc workerFunc) *workerPool[T] {
	workerPool := &workerPool[T]{
		input: input,
	}

	for i := 0; i < workerNum; i++ {
		go workerFunc(ctx)
	}

	go func() {
		<-ctx.Done()
		close(input)
	}()

	return workerPool
}

// AddTask adds a new task to the worker pool's input channel.
// This method allows submitting tasks to be processed by the workers.
func (w *workerPool[T]) AddTask(task T) {
	w.input <- task
}
