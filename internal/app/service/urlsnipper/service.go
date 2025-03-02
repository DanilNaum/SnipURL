package urlsnipper

import (
	"context"
	"fmt"
)

var (
	ErrFailedToGenerateID = fmt.Errorf("failed to generate id")
	ErrFailedToGetUrl     = fmt.Errorf("failed to get url")
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

	//	I am currently using a simple ID generation algorithm that ensures
	//	that each ID is unique, but does not guarantee how many attempts it
	//	will take to create it.
	//	To avoid an endless loop, the number of attempts is limited.

	//	It is possible to improve this process by adding a random "salt" and
	//	storing it in a separate table. In case of unsuccessful operations,
	//	you can first check for a link in the table, which stores links that
	//	cannot be encoded without a "salt". However, even this does not
	//	guarantee a limit on the number of attempts to generate IDs,
	//	since there are a finite number of possible IDs for a given length.

	//	One solution is to increase the hash length after several unsuccessful attempts.
	//	This will ensure the stable operation of the application, but it may lead to an
	//	increase in the length of the original link, rather than reducing it.
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
		return "", fmt.Errorf("%w: %w", ErrFailedToGetUrl, err)
	}
	return url, nil
}
