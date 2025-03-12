package urlsnipper

import (
	"context"
	"fmt"
)

var (
	ErrFailedToGenerateID = fmt.Errorf("failed to generate id")
	ErrFailedToGetURL     = fmt.Errorf("failed to get url")
)

const _maxAttempts = 10

//go:generate moq -out mock_url_storage_moq_test.go . urlStorage
type urlStorage interface {
	SetURL(ctx context.Context, id, url string) error
	GetURL(ctx context.Context, id string) (string, error)
}

//go:generate moq -out mock_hasher_moq_test.go . hasher
type hasher interface {
	Hash(s string) string
}

type urlSnipperService struct {
	storage urlStorage
	hasher  hasher
}

func NewURLSnipperService(storage urlStorage, hasher hasher) *urlSnipperService {
	return &urlSnipperService{
		storage: storage,
		hasher:  hasher,
	}
}

func (s *urlSnipperService) SetURL(ctx context.Context, url string) (string, error) {

	urlCopy := url
	for i := 0; i < _maxAttempts; i++ {
		id := s.hasher.Hash(urlCopy)
		err := s.storage.SetURL(ctx, id, url)
		if err == nil {
			return id, nil
		}

		urlCopy = fmt.Sprint(urlCopy, id)

	}
	return "", ErrFailedToGenerateID
}

func (s *urlSnipperService) GetURL(ctx context.Context, id string) (string, error) {
	url, err := s.storage.GetURL(ctx, id)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetURL, err)
	}
	return url, nil
}
