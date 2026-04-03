package memory

import (
	"context"
	"errors"
	"sync"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlewares"
	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
)

const (
	expectedNumberOfURLs = 1000
)

var key = middlewares.Key{Key: "userID"}

type dumper interface {
	ReadAll() (chan dump.URLRecord, error)
}

type storage struct {
	mu   sync.RWMutex
	urls map[string]*urlstorage.URLRecord
}

// NewStorage creates and returns a new in-memory storage for URL records.
// It initializes an empty map to store URL records with thread-safe access.
func NewStorage() *storage {
	return &storage{
		urls: make(map[string]*urlstorage.URLRecord),
	}
}

// Ping checks the availability of the storage service.
// Currently returns a "not implemented" error.
// Implements the urlstorage.Repository interface.
func (s *storage) Ping(ctx context.Context) error {
	return errors.New("not implemented")
}

// SetURL adds a new URL record to the in-memory storage with thread-safe synchronization.
// It locks the mutex, calls the internal setURL method, and returns the total number of URLs or an error.
func (s *storage) SetURL(ctx context.Context, id, url string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.setURL(ctx, id, url)

}

func (s *storage) setURL(ctx context.Context, id, url string) (int, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		userID = ""
	}

	if oldURL, ok := s.urls[id]; ok && oldURL.OriginalURL != url {
		return 0, urlstorage.ErrIDIsBusy
	} else if ok {
		return 0, urlstorage.ErrConflict
	}

	s.urls[id] = &urlstorage.URLRecord{
		ShortURL:    id,
		OriginalURL: url,
		UserID:      userID,
		Deleted:     false,
	}
	return len(s.urls), nil
}

// GetURL retrieves the original URL for a given short URL ID.
// It uses a read lock to ensure thread-safe access to the in-memory storage.
// Returns the original URL if found, or an error if the URL is not found or has been deleted.
func (s *storage) GetURL(_ context.Context, id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[id]
	if !ok {
		return "", urlstorage.ErrNotFound
	}
	if url.Deleted {
		return "", urlstorage.ErrDeleted
	}
	return url.OriginalURL, nil

}

// RestoreStorage populates the in-memory storage with URL records from a dumper.
// It reads all records and attempts to set them in the storage,
// returning an error if any record fails to be added.
func (s *storage) RestoreStorage(dumper dumper) error {
	records, err := dumper.ReadAll()
	if err != nil {
		return err
	}

	for record := range records {
		_, err := s.SetURL(context.Background(), record.ShortURL, record.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetURLs adds multiple URL records to the storage.
// It attempts to insert each URL, skipping URLs that would cause ID conflicts.
// Returns a slice of successfully inserted URLs and any error encountered during insertion.
func (s *storage) SetURLs(ctx context.Context, urls []*urlstorage.URLRecord) (insertedURLs []*urlstorage.URLRecord, err error) {
	inserted := make([]*urlstorage.URLRecord, 0, len(urls))
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, url := range urls {
		_, err := s.setURL(ctx, url.ShortURL, url.OriginalURL)
		if err != nil {
			if errors.Is(err, urlstorage.ErrIDIsBusy) {
				continue
			}
			return nil, err

		}
		url.ID = len(s.urls)
		inserted = append(inserted, url)
	}
	return inserted, nil
}

// GetURLs retrieves all non-deleted URL records for a specific user.
// Requires a user ID in the context. Returns an error if no user ID is found.
// Returns a slice of URL records belonging to the user.
func (s *storage) GetURLs(ctx context.Context) ([]*urlstorage.URLRecord, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return nil, errors.New("userID not found in context")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	urls := make([]*urlstorage.URLRecord, 0, expectedNumberOfURLs)
	for _, url := range s.urls {
		if url.UserID == userID {
			if !url.Deleted {
				urls = append(urls, url)
			}
		}
	}
	return urls, nil
}

// DeleteURLs marks specified URL records as deleted for a given user.
// Only deletes URLs that belong to the specified user.
// Silently skips URLs that do not exist or belong to a different user.
func (s *storage) DeleteURLs(userID string, ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		url, ok := s.urls[id]

		if !ok {
			continue
		}
		if url.UserID != userID {
			continue
		}

		url.Deleted = true
	}
	return nil
}

// GetState returns statistics about the current state of the storage.
// It counts the number of non-deleted URLs and unique users with non-empty user IDs.
func (s *storage) GetState(ctx context.Context) (*urlstorage.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	urlsCount := 0
	users := make(map[string]bool)

	for _, urlRecord := range s.urls {
		if !urlRecord.Deleted && urlRecord.UserID != "" {
			urlsCount++
			users[urlRecord.UserID] = true
		}
	}

	return &urlstorage.State{
		UrlsNum:  urlsCount,
		UsersNum: len(users),
	}, nil
}
