//go:generate go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=sqlite

package migrations

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4/database"
	driver "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/mntone/miruzo-core/miruzo/internal/database/migration"
)

//go:embed *.sql
var fs embed.FS

func newSourceDriver() (source.Driver, error) {
	return iofs.New(fs, ".")
}

func newDatabaseDriverFunc(db *sql.DB) func() (database.Driver, error) {
	return func() (database.Driver, error) {
		return driver.WithInstance(db, &driver.Config{})
	}
}

func NewSpec(db *sql.DB) migration.Spec {
	return migration.Spec{
		SourceName:        "iofs",
		NewSourceDriver:   newSourceDriver,
		DatabaseName:      "sqlite3",
		NewDatabaseDriver: newDatabaseDriverFunc(db),
	}
}
