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
	expectedNumberOfURLs = 20
)

var key = middlewares.Key{Key: "userID"}

type dumper interface {
	ReadAll() (chan dump.URLRecord, error)
}

type storage struct {
	mu   sync.RWMutex
	urls map[string]*urlstorage.URLRecord
}

func NewStorage() *storage {
	return &storage{
		urls: make(map[string]*urlstorage.URLRecord),
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return errors.New("not implemented")
}
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
