package urlsnipper

import (
	"context"
	"errors"
	"fmt"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlewares"
	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
)

// Predefined error variables for URL-related operations, providing specific error conditions
// during URL generation, retrieval, storage, and deletion processes.
var (
	// ErrFailedToGenerateID indicates a failure in generating a unique URL identifier.
	ErrFailedToGenerateID = fmt.Errorf("failed to generate id")

	// ErrFailedToGetURL indicates a failure in retrieving a URL from storage.
	ErrFailedToGetURL = fmt.Errorf("failed to get url")

	// ErrConflict represents a conflict during URL storage or retrieval.
	ErrConflict = fmt.Errorf("conflict")

	// ErrDeleted indicates that the requested URL has been deleted.
	ErrDeleted = fmt.Errorf("deleted")
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

// NewURLSnipperService creates and returns a new instance of urlSnipperService with the provided dependencies.
// It initializes a URL snipper service with storage, hashing, dumping, deletion, and logging capabilities.
//
// Parameters:
//   - storage: Implementation of URL storage interface
//   - hasher: Hash generator for creating short URL IDs
//   - dumper: URL record dumper
//   - deleteService: Service for handling URL deletions
//   - logger: Logger for recording errors
//
// Returns:
//   - *urlSnipperService: Configured URL snipper service instance
func NewURLSnipperService(storage urlStorage, hasher hasher, dumper dumper, deleteService deleteService, logger logger) *urlSnipperService {
	return &urlSnipperService{
		storage:       storage,
		hasher:        hasher,
		dumper:        dumper,
		deleteService: deleteService,
		logger:        logger,
	}
}

// SetURL creates a short URL from the given original URL. It attempts to generate a unique short URL ID
// by hashing the URL. If a collision occurs (ErrConflict), it returns the conflicting ID. If generation fails after
// _maxAttempts, it returns ErrFailedToGenerateID. On success, it stores the URL mapping and dumps the record.
//
// Parameters:
//   - ctx: The context for the operation
//   - url: The original URL to be shortened
//
// Returns:
//   - string: The generated short URL ID on success, or empty string on failure
//   - error: ErrConflict if ID exists, ErrFailedToGenerateID if generation fails, or nil on success
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

// GetURL retrieves the original URL associated with the given short URL ID.
// If the URL has been deleted, it returns ErrDeleted. For any other errors,
// it wraps them with ErrFailedToGetURL.
//
// Parameters:
//   - ctx: The context for the operation
//   - id: The short URL ID to look up
//
// Returns:
//   - string: The original URL if found, or empty string on failure
//   - error: ErrDeleted if URL was deleted, wrapped error with ErrFailedToGetURL for other errors, or nil on success
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

// SetURLs creates multiple short URLs from the given array of original URLs in batch.
// It generates unique short URL IDs by hashing each original URL, stores the URL mappings,
// and dumps the records. If any operation fails, it returns an error.
//
// Parameters:
//   - ctx: The context for the operation
//   - urls: Array of SetURLsInput containing original URLs and correlation IDs
//
// Returns:
//   - map[string]*SetURLsOutput: Map of correlation IDs to their corresponding short URL outputs
//   - error: ErrFailedToGenerateID if not all URLs were inserted, storage error, or nil on success
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

// GetURLs retrieves all URLs stored in the system.
// It returns a slice of URL objects containing both short and original URLs.
// If there's an error retrieving URLs from storage, it returns the error.
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - []*URL: Slice of URL objects containing short and original URLs
//   - error: Storage error or nil on success
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

// DeleteURLs asynchronously deletes multiple URLs in batches.
// It extracts the user ID from the context and processes URL deletions in batches of size defined by batchSize.
// If the user ID cannot be extracted from the context, the operation is aborted.
//
// Parameters:
//   - ctx: The context containing the user ID
//   - ids: Slice of URL IDs to be deleted
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
