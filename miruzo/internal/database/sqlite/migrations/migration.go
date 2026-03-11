//go:generate go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=sqlite

package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	driver "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func RunMigrations(db *sql.DB) (err error) {
	sourceDriver, err := iofs.New(fs, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}
	defer func() {
		closeErr := sourceDriver.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close migration source: %w", closeErr))
		}
	}()

	databaseDriver, err := driver.WithInstance(db, &driver.Config{})
	if err != nil {
		return fmt.Errorf("create sqlite migration driver: %w", err)
	}

	migrateInstance, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"sqlite3",
		databaseDriver,
	)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}

	err = migrateInstance.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
