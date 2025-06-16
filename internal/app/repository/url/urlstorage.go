package url

import "context"

type URLStorage interface {
	Ping(ctx context.Context) error
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, id, url string) (int, error)
	SetURLs(ctx context.Context, urls []*URLRecord) ([]*URLRecord, error)
	GetURLs(ctx context.Context) ([]*URLRecord, error)
	DeleteURLs(userID string, ids []string) error
}
