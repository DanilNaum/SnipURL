package deleteurl

import (
	"context"

	"github.com/DanilNaum/SnipURL/pkg/workerpool"
)

//go:generate moq -out mock_url_storage_moq_test.go . urlStorage
type urlStorage interface {
	DeleteURLs(userID string, ids []string) error
}

// This const allows to configure delete worker number and batch size
const (
	workerNum = 10
)

type workerPool interface {
	AddTask(task *data)
}

type data struct {
	userID string
	ids    []string
}

type deleteService struct {
	input      chan *data
	storage    urlStorage
	workerPool workerPool
}

// NewDeleteService creates a new delete service with a worker pool for asynchronous URL deletion.
// It initializes a channel for tasks, sets up a worker pool with the specified number of workers,
// and returns a configured delete service instance.
//
// Parameters:
//   - ctx: the context for managing worker pool lifecycle
//   - storage: the URL storage interface for performing deletion operations
//
// Returns:
//   - *deleteService: a configured delete service ready to process deletion tasks
func NewDeleteService(ctx context.Context, storage urlStorage) *deleteService {
	input := make(chan *data, workerNum)

	d := &deleteService{storage: storage, input: input}

	workerPool := workerpool.NewWorkerPool(ctx, workerNum, input, d.deleteWorker)

	d.workerPool = workerPool

	return d
}

// Delete adds a task to delete URLs with the specified IDs for a given user to the worker pool.
// The deletion is performed asynchronously by worker goroutines.
//
// Parameters:
//   - userID: the identifier of the user who owns the URLs
//   - input: a slice of URL IDs to be deleted
func (d *deleteService) Delete(userID string, input []string) {
	d.workerPool.AddTask(&data{userID: userID, ids: input})
}

func (d *deleteService) deleteWorker(ctx context.Context) error {
	for {
		select {
		case data, ok := <-d.input:
			if !ok {
				return nil
			}
			err := d.storage.DeleteURLs(data.userID, data.ids)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}
