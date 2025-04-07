package urlsnipper

import (
	"context"
	"fmt"

	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
)

var (
	ErrFailedToGenerateID = fmt.Errorf("failed to generate id")
	ErrFailedToGetURL     = fmt.Errorf("failed to get url")
)

const _maxAttempts = 10

//go:generate moq -out mock_url_storage_moq_test.go . urlStorage
type urlStorage interface {
	SetURL(ctx context.Context, id, url string) (length int, err error)
	GetURL(ctx context.Context, id string) (string, error)
}

//go:generate moq -out mock_hasher_moq_test.go . hasher
type hasher interface {
	Hash(s string) string
}

//go:generate moq -out mock_dumper_moq_test.go . dumper
type dumper interface {
	Add(record *dump.URLRecord) error
	ReadAll() (chan dump.URLRecord, error)
}

type logger interface {
	Errorf(string, ...interface{})
}

type urlSnipperService struct {
	storage urlStorage
	hasher  hasher
	dumper  dumper
	logger  logger
}

func NewURLSnipperService(storage urlStorage, hasher hasher, dumper dumper, logger logger) *urlSnipperService {
	return &urlSnipperService{
		storage: storage,
		hasher:  hasher,
		dumper:  dumper,
		logger:  logger,
	}
}

func (s *urlSnipperService) SetURL(ctx context.Context, url string) (string, error) {

	urlCopy := url
	for i := 0; i < _maxAttempts; i++ {
		id := s.hasher.Hash(urlCopy)
		length, err := s.storage.SetURL(ctx, id, url)
		if err == nil {
			rec := &dump.URLRecord{
				UUID:        length,
				ShortURL:    id,
				OriginalURL: url,
			}

			err = s.dumper.Add(rec)
			if err != nil {
				s.logger.Errorf("failed to dump record: %v", err)
			}
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
