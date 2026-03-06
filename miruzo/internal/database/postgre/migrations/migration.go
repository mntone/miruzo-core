package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	driver "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed *.sql
var fs embed.FS

func wrapCloseError(operation string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", operation, err)
}

func RunMigrations(pool *pgxpool.Pool) (err error) {
	sourceDriver, err := iofs.New(fs, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	databaseDriver, err := driver.WithInstance(db, &driver.Config{})
	if err != nil {
		sourceCloseErr := sourceDriver.Close()
		dbCloseErr := db.Close()
		return errors.Join(
			fmt.Errorf("create postgresql migration driver: %w", err),
			wrapCloseError("close migration source", sourceCloseErr),
			wrapCloseError("close migration database", dbCloseErr),
		)
	}

	migrateInstance, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"pgx",
		databaseDriver,
	)
	if err != nil {
		sourceCloseErr := sourceDriver.Close()
		databaseCloseErr := databaseDriver.Close()
		return errors.Join(
			fmt.Errorf("create migration instance: %w", err),
			wrapCloseError("close migration source", sourceCloseErr),
			wrapCloseError("close migration database", databaseCloseErr),
		)
	}
	defer func() {
		sourceCloseErr, databaseCloseErr := migrateInstance.Close()
		err = errors.Join(
			err,
			wrapCloseError("close migration source", sourceCloseErr),
			wrapCloseError("close migration database", databaseCloseErr),
		)
	}()

	err = migrateInstance.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
