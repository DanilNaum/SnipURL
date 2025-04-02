package memory

import (
	"context"
	"errors"
	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
	"sync"
)

type dumper interface {
	ReadAll() (chan dump.Record, error)
}

type storage struct {
	mu   sync.RWMutex
	urls map[string]string
}

func NewStorage() *storage {
	return &storage{
		mu:   sync.RWMutex{},
		urls: make(map[string]string),
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return errors.New("not implemented")
}
func (s *storage) SetURL(_ context.Context, id, url string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldURL, ok := s.urls[id]; ok && oldURL != url {
		return -1, urlstorage.ErrIDIsBusy
	}

	s.urls[id] = url
	return len(s.urls), nil
}

func (s *storage) GetURL(_ context.Context, id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[id]
	if !ok {
		return "", urlstorage.ErrNotFound
	}

	return url, nil

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
