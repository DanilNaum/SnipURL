package private

import (
	"context"
	"fmt"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
)

type urlStorage interface {
	GetState(ctx context.Context) (*urlstorage.State, error)
}

type internalService struct {
	urlStorage urlStorage
}

// NewInternalService creates a new instance of internalService with the provided URL storage.
func NewInternalService(urlStorage urlStorage) *internalService {
	return &internalService{
		urlStorage: urlStorage,
	}
}

// GetState retrieves the current system state including URL and user counts.
// Returns a State object containing the number of URLs and users in the system.
// Returns an error if the underlying storage operation fails.
func (i *internalService) GetState(ctx context.Context) (*State, error) {
	state, err := i.urlStorage.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("unexpected err:  %w", err)
	}

	return &State{
		UrlsNum:  state.UrlsNum,
		UsersNum: state.UsersNum,
	}, nil
}
