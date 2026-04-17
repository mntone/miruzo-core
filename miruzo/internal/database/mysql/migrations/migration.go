//go:generate go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=mysql

package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4/database"
	driver "github.com/golang-migrate/migrate/v4/database/mysql"
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
		conn, err := db.Conn(context.Background())
		if err != nil {
			return nil, err
		}

		// Use WithConnection(*sql.Conn) instead of WithInstance(*sql.DB).
		// The migration driver must Close() this dedicated conn, otherwise
		// subsequent queries can block and hit timeouts.
		return driver.WithConnection(context.Background(), conn, &driver.Config{})
	}
}

func NewSpec(db *sql.DB) migration.Spec {
	return migration.Spec{
		SourceName:        "iofs",
		NewSourceDriver:   newSourceDriver,
		DatabaseName:      "mysql",
		NewDatabaseDriver: newDatabaseDriverFunc(db),
	}
}
