//go:generate go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=postgres

package migrations

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4/database"
	driver "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mntone/miruzo-core/miruzo/internal/database/migration"
)

//go:embed *.sql
var fs embed.FS

func newSourceDriver() (source.Driver, error) {
	return iofs.New(fs, ".")
}

func newDatabaseDriverFunc(pool *pgxpool.Pool) func() (database.Driver, error) {
	return func() (database.Driver, error) {
		db := stdlib.OpenDBFromPool(pool)
		database, err := driver.WithInstance(db, &driver.Config{})
		if err != nil {
			databaseCloseErr := db.Close()
			return nil, errors.Join(err, databaseCloseErr)
		}

		return database, nil
	}
}

func NewSpec(pool *pgxpool.Pool) migration.Spec {
	return migration.Spec{
		SourceName:        "iofs",
		NewSourceDriver:   newSourceDriver,
		DatabaseName:      "pgx",
		NewDatabaseDriver: newDatabaseDriverFunc(pool),
		CloseDatabase:     true,
	}
}
