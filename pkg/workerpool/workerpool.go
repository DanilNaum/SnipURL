package workerpool

import "context"

type workerFunc func(ctx context.Context) error

type workerPool[T any] struct {
	input chan T
}

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

func (w *workerPool[T]) AddTask(task T) {
	w.input <- task
}
