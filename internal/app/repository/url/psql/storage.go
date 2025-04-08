package psql

import (
	"context"
	"errors"
	"fmt"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	"github.com/DanilNaum/SnipURL/pkg/utils/placeholder"
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

func (s *storage) Ping(ctx context.Context) error {
	if s.conn == nil {
		return errors.New("connection is nil")
	}
	err := s.conn.Master().Ping(ctx)
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

func (s *storage) SetURLs(_ context.Context, urls []*urlstorage.URLRecord) (insertedUrls []*urlstorage.URLRecord, err error) {

	placeholder := placeholder.MakeDollars(
		placeholder.WithColumnNumAndRowNum(2, len(urls)),
	)
	query := fmt.Sprintf(`INSERT INTO url (id, url) VALUES %s 
  	ON CONFLICT (id) DO NOTHING
  	RETURNING uuid, id, url`, placeholder)

	rows, err := s.conn.Master().Query(context.Background(),
		query,
		valuesForInsert(urls)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	insertedUrls = make([]*urlstorage.URLRecord, 0, len(urls))
	for rows.Next() {
		var urlRecord urlstorage.URLRecord
		err := rows.Scan(&urlRecord.ID, &urlRecord.ShortURL, &urlRecord.OriginalURL)
		if err != nil {
			return nil, err
		}
		insertedUrls = append(insertedUrls, &urlRecord)
	}

	return insertedUrls, nil
}

func valuesForInsert(urlRecords []*urlstorage.URLRecord) []interface{} {
	values := make([]interface{}, 0, len(urlRecords)*2)

	for _, urlRecord := range urlRecords {
		values = append(values, urlRecord.ShortURL, urlRecord.OriginalURL)
	}

	return values
}
