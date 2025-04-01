package psql

import (
	"context"
	"errors"

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
