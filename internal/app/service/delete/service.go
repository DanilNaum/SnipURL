package deleteurl

import (
	"context"
)

//go:generate moq -out mock_url_storage_moq_test.go . urlStorage
type urlStorage interface {
	DeleteURLs(userID string, ids []string) error
}

// This const allows to configure delete worker number and batch size
const (
	workerNum = 10
)

type data struct {
	userID string
	ids    []string
}

type deleteService struct {
	input   chan *data
	storage urlStorage
}

func NewDeleteService(ctx context.Context, storage urlStorage) *deleteService {
	input := make(chan *data, workerNum)
	d := &deleteService{storage: storage, input: input}
	for i := 0; i < workerNum; i++ {
		go d.deleteWorker(ctx)
	}
	go func() {
		<-ctx.Done()
		close(input)
	}()

	return d
}

func (d *deleteService) Delete(userID string, input []string) {
	d.input <- &data{userID: userID, ids: input}
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
