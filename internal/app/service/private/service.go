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

func NewInternalService(urlStorage urlStorage) *internalService {
	return &internalService{
		urlStorage: urlStorage,
	}
}

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
