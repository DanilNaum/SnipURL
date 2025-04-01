package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type logger interface {
	Info(args ...interface{})
	Errorf(template string, args ...interface{})
}

type connection struct {
	masterPool *pgxpool.Pool
}

func (c *connection) Master() *pgxpool.Pool {
	return c.masterPool
}

func (c *connection) Close() {
	c.Master().Close()

}

func NewConnection(ctx context.Context, cnf *connConfig, log logger) *connection {
	masterDsn := cnf.getDsn()

	masterPool := createPool(ctx, masterDsn, "master", log)
	if masterPool == nil {
		return nil
	}

	return &connection{
		masterPool: masterPool,
	}
}

func createPool(ctx context.Context, dsn, tp string, log logger) *pgxpool.Pool {
	pg, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Errorf("сould not establish db %s connection %s", tp, err.Error())
		return nil
	}

	log.Info("msg", fmt.Sprintf("Database connection %s established", tp))
	return pg
}
