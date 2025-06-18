package migration

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrator struct {
	migrationFolderPath string
	dbName              string
	dsn                 string
}

type migrationFolderPathOpt func() string

func WithAbsolutePath(path string) migrationFolderPathOpt {
	return func() string {
		return "file:///" + path
	}
}

func WithRelativePath(path string) migrationFolderPathOpt {
	return func() string {
		return "file://" + path
	}
}

func NewMigrator(dsn string, migrationFolderPathOpt migrationFolderPathOpt) *migrator {
	mig := &migrator{
		migrationFolderPath: migrationFolderPathOpt(),
		dbName:              "postgres",
		dsn:                 dsn,
	}
	return mig
}

// Migrate executes database migrations using the configured migration path.
// It opens a connection to the PostgreSQL database, sets up the migration driver,
// and applies all pending migrations in the forward direction.
// Returns an error if any step in the migration process fails.
func (mig *migrator) Migrate() error {
	db, err := sql.Open("postgres", mig.dsn)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		mig.migrationFolderPath,
		mig.dbName, driver)
	if err != nil {
		return err
	}
	m.Up()
	return nil
}
