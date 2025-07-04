package psql

import (
	"context"
	"errors"
	"fmt"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlewares"
	"github.com/DanilNaum/SnipURL/pkg/utils/placeholder"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jackc/pgerrcode"
)

const (
	expectedNumberOfURLs = 20
)

var key = middlewares.Key{Key: "userID"}

type storage struct {
	conn *pgxpool.Pool
}

func NewStorage(conn *pgxpool.Pool) *storage {
	return &storage{
		conn: conn,
	}
}

func (s *storage) Ping(ctx context.Context) error {
	if s.conn == nil {
		return errors.New("connection is nil")
	}
	err := s.conn.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) SetURL(ctx context.Context, id, url string) (int, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		userID = ""
	}
	query := `INSERT INTO url (id, url, user_uuid) 
	VALUES ($1, $2, $3)
	RETURNING uuid`

	var uuid int

	err := s.conn.QueryRow(context.Background(), query, id, url, userID).Scan(&uuid)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				query = `SELECT uuid FROM url WHERE url = $1`
				err = s.conn.QueryRow(context.Background(), query, url).Scan(&uuid)
				if err != nil {
					return 0, err
				}

				return uuid, urlstorage.ErrConflict
			}
		}
		return 0, err
	}

	return uuid, nil
}

func (s *storage) GetURL(ctx context.Context, id string) (string, error) {
	query := `SELECT url, deleted FROM url WHERE id = $1`
	var url string
	var deleted bool
	err := s.conn.QueryRow(ctx, query, id).Scan(&url, &deleted)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return "", urlstorage.ErrNotFound
		}
		return "", err
	}

	if deleted {
		return "", urlstorage.ErrDeleted
	}
	return url, nil
}

func (s *storage) SetURLs(ctx context.Context, urls []*urlstorage.URLRecord) (insertedURLs []*urlstorage.URLRecord, err error) {
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

	rows, err := s.conn.Query(ctx,
		query,
		valuesForInsert(userID, urls)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	insertedURLs = make([]*urlstorage.URLRecord, 0, len(urls))
	for rows.Next() {
		var urlRecord urlstorage.URLRecord
		err := rows.Scan(&urlRecord.ID, &urlRecord.ShortURL, &urlRecord.OriginalURL)
		if err != nil {
			return nil, err
		}
		insertedURLs = append(insertedURLs, &urlRecord)
	}

	return insertedURLs, nil
}

func (s *storage) GetURLs(ctx context.Context) ([]*urlstorage.URLRecord, error) {
	userID, ok := ctx.Value(key).(string)
	if !ok {
		return nil, errors.New("error get userID from context")
	}
	query := `SELECT id, url FROM url WHERE user_uuid = $1 AND deleted = false`
	rows, err := s.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := make([]*urlstorage.URLRecord, 0, expectedNumberOfURLs)
	for rows.Next() {
		var urlRecord urlstorage.URLRecord
		err := rows.Scan(&urlRecord.ShortURL, &urlRecord.OriginalURL)
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

func (s *storage) DeleteURLs(userID string, ids []string) error {

	query := `UPDATE url SET deleted = true WHERE id = ANY($1) AND user_uuid = $2`
	_, err := s.conn.Exec(context.TODO(), query, ids, userID)
	if err != nil {
		return err
	}
	return nil
}
