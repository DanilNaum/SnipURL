package urlsnipper

import (
	"context"
	"errors"
	"fmt"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlewares"
	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
)

var (
	ErrFailedToGenerateID = fmt.Errorf("failed to generate id")
	ErrFailedToGetURL     = fmt.Errorf("failed to get url")
	ErrConflict           = fmt.Errorf("conflict")
	ErrDeleted            = fmt.Errorf("deleted")
)

const (
	_maxAttempts = 10
	batchSize    = 10
)

//go:generate moq -out mock_url_storage_moq_test.go . urlStorage
type urlStorage interface {
	SetURL(ctx context.Context, id, url string) (length int, err error)
	GetURL(ctx context.Context, id string) (string, error)
	SetURLs(ctx context.Context, urls []*urlstorage.URLRecord) (insertedURLs []*urlstorage.URLRecord, err error)
	GetURLs(ctx context.Context) ([]*urlstorage.URLRecord, error)
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

//go:generate moq -out mock_logger_moq_test.go . logger
type logger interface {
	Errorf(string, ...interface{})
}

//go:generate moq -out mock_delete_service_moq_test.go . deleteService
type deleteService interface {
	Delete(userID string, input []string)
}

type urlSnipperService struct {
	storage       urlStorage
	hasher        hasher
	dumper        dumper
	logger        logger
	deleteService deleteService
}

func NewURLSnipperService(storage urlStorage, hasher hasher, dumper dumper, deleteService deleteService, logger logger) *urlSnipperService {
	return &urlSnipperService{
		storage:       storage,
		hasher:        hasher,
		dumper:        dumper,
		deleteService: deleteService,
		logger:        logger,
	}
}

func (s *urlSnipperService) SetURL(ctx context.Context, url string) (string, error) {

	urlCopy := url
	for i := 0; i < _maxAttempts; i++ {
		id := s.hasher.Hash(urlCopy)
		length, err := s.storage.SetURL(ctx, id, url)
		if errors.Is(err, urlstorage.ErrConflict) {
			return id, ErrConflict
		}
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
		switch {
		case errors.Is(err, urlstorage.ErrDeleted):
			return "", ErrDeleted
		default:
			return "", fmt.Errorf("%w: %w", ErrFailedToGetURL, err)
		}
	}

	return url, nil
}

func (s *urlSnipperService) SetURLs(ctx context.Context, urls []*SetURLsInput) (map[string]*SetURLsOutput, error) {
	output := make(map[string]*SetURLsOutput, len(urls))

	toInsert := make([]*urlstorage.URLRecord, 0, len(urls))

	for _, url := range urls {
		id := s.hasher.Hash(url.OriginalURL)

		toInsert = append(toInsert, &urlstorage.URLRecord{
			ShortURL:    id,
			OriginalURL: url.OriginalURL,
		})

		output[url.CorrelationID] = &SetURLsOutput{
			CorrelationID: url.CorrelationID,
			ShortURLID:    id,
		}

	}

	inserted, err := s.storage.SetURLs(ctx, toInsert)
	if err != nil {
		return nil, err
	}
	for _, record := range inserted {
		rec := &dump.URLRecord{
			UUID:        record.ID,
			ShortURL:    record.ShortURL,
			OriginalURL: record.OriginalURL,
		}
		err := s.dumper.Add(rec)
		if err != nil {
			s.logger.Errorf("failed to dump record: %v", err)
		}
	}

	if len(inserted) != len(toInsert) {
		return nil, ErrFailedToGenerateID
	}

	return output, nil
}

func (s *urlSnipperService) GetURLs(ctx context.Context) ([]*URL, error) {
	urls, err := s.storage.GetURLs(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*URL, 0, len(urls))
	for _, url := range urls {
		output = append(output, &URL{
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
		})
	}
	return output, nil

}

var key = middlewares.Key{Key: "userID"}

func (s *urlSnipperService) DeleteURLs(ctx context.Context, ids []string) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return
	}

	go func() {
		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			s.deleteService.Delete(userID, ids[i:end])
		}
	}()

}
