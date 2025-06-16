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

func NewDeleteService(ctx context.Context, storage urlStorage) *deleteService {
	input := make(chan *data, workerNum)

	d := &deleteService{storage: storage, input: input}

	workerPool := workerpool.NewWorkerPool(ctx, workerNum, input, d.deleteWorker)

	d.workerPool = workerPool

	return d
}

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
