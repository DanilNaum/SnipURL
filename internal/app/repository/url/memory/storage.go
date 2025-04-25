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
		mu:   sync.RWMutex{},
		urls: make(map[string]*urlstorage.URLRecord),
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return errors.New("not implemented")
}
func (s *storage) SetURL(ctx context.Context, id, url string) (int, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		userID = ""
	}
	s.mu.Lock()
	defer s.mu.Unlock()

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
	for _, url := range urls {
		_, err := s.SetURL(ctx, url.ShortURL, url.OriginalURL)
		if err != nil {
			if errors.Is(err, urlstorage.ErrIDIsBusy) {

			} else {
				return nil, err
			}
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
			urls = append(urls, url)
		}
	}
	return urls, nil
}

func (s *storage) DeleteURLs(ctx context.Context, ids []string) error {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return errors.New("userID not found in context")
	}
	for _, id := range ids {
		s.mu.Lock()
		url, ok := s.urls[id]
		s.mu.Unlock()
		if !ok {
			continue
		}
		if url.UserID != userID {
			continue
		}
		s.mu.Lock()
		url.Deleted = true
		s.mu.Unlock()
	}
	return nil
}
