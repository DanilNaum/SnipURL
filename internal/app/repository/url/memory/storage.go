package memory

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
	ErrIdIsBusy = errors.New("id is busy")
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

func (s *storage) SetUrl(_ context.Context, id, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldUrl, ok := s.urls[id]; ok && oldUrl != url {
		return ErrIdIsBusy
	}

	s.urls[id] = url
	return nil
}

func (s *storage) GetUrl(_ context.Context, id string) (string, error) {
	s.mu.RLock()
	s.mu.RUnlock()

	url, ok := s.urls[id]
	if !ok {
		return "", ErrNotFound
	}

	return url, nil

}
