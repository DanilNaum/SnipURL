package memory

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
	ErrIDIsBusy = errors.New("id is busy")
)

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

func (s *storage) SetURL(_ context.Context, id, url string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldURL, ok := s.urls[id]; ok && oldURL != url {
		return 0, ErrIDIsBusy
	}

	s.urls[id] = url
	return len(s.urls), nil
}

func (s *storage) GetURL(_ context.Context, id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[id]
	if !ok {
		return "", ErrNotFound
	}

	return url, nil

}
