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

func NewConnection(ctx context.Context, cnf *connConfig, log logger) *pgxpool.Pool {
	masterDsn := cnf.getDsn()

	pg, err := pgxpool.Connect(ctx, masterDsn)
	if err != nil {
		log.Errorf("сould not establish db connection %s", err.Error())
		return nil
	}

	log.Info("msg", fmt.Sprintf("Database connection established"))
	return pg
}
