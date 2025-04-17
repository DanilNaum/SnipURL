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

const (
	expectedNumberOfURLs = 20
	key                  = "userID"
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

func (s *storage) SetURL(ctx context.Context, id, url string) (int, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return 0, errors.New("error get userID from context")
	}

	query := `WITH insertion AS (
	INSERT INTO url (id, url, user_uuid)
	VALUES ($1, $2, $3)
	ON CONFLICT (url) DO NOTHING
	RETURNING uuid, true AS is_inserted
	),
	fallback AS (
	SELECT uuid, false AS is_inserted  FROM url
	WHERE url = $2
	)
	SELECT uuid, is_inserted  FROM insertion
	UNION ALL
	SELECT uuid, is_inserted FROM fallback
	WHERE NOT EXISTS (SELECT 1 FROM insertion)
	LIMIT 1`

	var (
		uuid     int
		inserted bool
	)

	err := s.conn.Master().QueryRow(ctx, query, id, url, userID).Scan(&uuid, &inserted)
	if err != nil {
		return 0, err
	}

	if !inserted {
		return 0, urlstorage.ErrConflict
	}

	return uuid, nil
}

func (s *storage) GetURL(ctx context.Context, id string) (string, error) {
	query := `SELECT url FROM url WHERE id = $1`
	var url string
	err := s.conn.Master().QueryRow(ctx, query, id).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", urlstorage.ErrNotFound
		}
		return "", err
	}
	return url, nil
}

func (s *storage) SetURLs(ctx context.Context, urls []*urlstorage.URLRecord) (insertedUrls []*urlstorage.URLRecord, err error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return nil, errors.New("error get userID from context")
	}
	placeholder := placeholder.MakeDollars(
		placeholder.WithColumnNumAndRowNum(3, len(urls)),
	)
	query := fmt.Sprintf(`INSERT INTO url (id, url, user_uuid) VALUES %s 
  	ON CONFLICT (id) DO NOTHING
  	RETURNING uuid, id, url`, placeholder)

	rows, err := s.conn.Master().Query(ctx,
		query,
		valuesForInsert(userID, urls)...,
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

func (s *storage) GetURLs(ctx context.Context) ([]*urlstorage.URLRecord, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return nil, errors.New("error get userID from context")
	}
	query := `SELECT id, url FROM url WHERE user_uuid = $1`
	rows, err := s.conn.Master().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := make([]*urlstorage.URLRecord, 0, expectedNumberOfURLs)
	for rows.Next() {
		var urlRecord urlstorage.URLRecord
		err := rows.Scan(&urlRecord.ID, &urlRecord.OriginalURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &urlRecord)
	}
	return urls, nil
}

func valuesForInsert(userID string, urlRecords []*urlstorage.URLRecord) []interface{} {
	values := make([]interface{}, 0, len(urlRecords)*3)

	for _, urlRecord := range urlRecords {
		values = append(values, urlRecord.ShortURL, urlRecord.OriginalURL, userID)
	}

	return values
}
