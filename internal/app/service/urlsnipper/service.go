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
	SetURL(ctx context.Context, id, url string) (int, error)
	GetURL(ctx context.Context, id string) (string, error)
}

//go:generate moq -out mock_hasher_moq_test.go . hasher
type hasher interface {
	Hash(s string) string
}

//go:generate moq -out mock_dumper_moq_test.go . dumper
type dumper interface {
	Add(record *dump.Record) error
	ReadAll() (chan dump.Record, error)
}

type urlSnipperService struct {
	storage urlStorage
	hasher  hasher
	dumper  dumper
}

func NewURLSnipperService(storage urlStorage, hasher hasher, dumper dumper) *urlSnipperService {
	return &urlSnipperService{
		storage: storage,
		hasher:  hasher,
		dumper:  dumper,
	}
}

func (s *urlSnipperService) SetURL(ctx context.Context, url string) (string, error) {

	urlCopy := url
	for i := 0; i < _maxAttempts; i++ {
		id := s.hasher.Hash(urlCopy)
		l, err := s.storage.SetURL(ctx, id, url)
		if err == nil {
			rec := &dump.Record{
				UUID:        l,
				ShortURL:    id,
				OriginalURL: url,
			}

			err = s.dumper.Add(rec)
			if err != nil {
				// TODO: add logger
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

func (s *urlSnipperService) RestoreStorage() error {
	records, err := s.dumper.ReadAll()
	if err != nil {
		return err
	}

	for record := range records {
		_, err := s.storage.SetURL(context.Background(), record.ShortURL, record.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}
