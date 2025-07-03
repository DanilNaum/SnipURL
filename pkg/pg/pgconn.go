package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type logger interface {
	Info(args ...interface{})
	Errorf(template string, args ...interface{})
}

// NewConnection establishes a new database connection pool using the provided context,
// connection configuration, and logger. Returns a pgxpool.Pool instance on success,
// or nil if the connection fails.
func NewConnection(ctx context.Context, cnf *connConfig, log logger) *pgxpool.Pool {
	masterDsn := cnf.getDsn()

	pg, err := pgxpool.Connect(ctx, masterDsn)
	if err != nil {
		log.Errorf("сould not establish db connection %s", err.Error())
		return nil
	}

	log.Info("msg", "Database connection established")
	return pg
}
