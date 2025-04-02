package psql

import (
	"context"
	"errors"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type connection interface {
	Master() *pgxpool.Pool
	Close()
}
type storage struct {
	conn connection
}

func NewStorage(conn connection) *storage {
	return &storage{
		conn: conn,
	}
}

func (c *storage) Ping(ctx context.Context) error {
	if c.conn == nil {
		return errors.New("connection is nil")
	}
	err := c.conn.Master().Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) SetURL(_ context.Context, id, url string) (int, error) {
	query := `INSERT INTO url (id, url) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING RETURNING uuid, id, url`
	var (
		idDB, urlDB string
		uuid        int
	)

	err := s.conn.Master().QueryRow(context.Background(), query, id, url).Scan(&uuid, &idDB, &urlDB)
	if err != nil {
		return 0, err
	}

	if url != urlDB {
		return 0, urlstorage.ErrIDIsBusy
	}

	return uuid, nil
}

func (s *storage) GetURL(_ context.Context, id string) (string, error) {
	query := `SELECT url FROM url WHERE id = $1`
	var url string
	err := s.conn.Master().QueryRow(context.Background(), query, id).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", urlstorage.ErrNotFound
		}
		return "", err
	}
	return url, nil
}
