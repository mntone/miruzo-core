package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
)

func wrapCloseError(operation string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", operation, err)
}

type Spec struct {
	SourceName      string
	NewSourceDriver func() (source.Driver, error)

	DatabaseName      string
	NewDatabaseDriver func() (database.Driver, error)
	KeepDatabaseOpen  bool
}

func (spec Spec) NewInstance() (*migrate.Migrate, func() error, error) {
	source, err := spec.NewSourceDriver()
	if err != nil {
		return nil, nil, fmt.Errorf("create migration source: %w", err)
	}

	database, err := spec.NewDatabaseDriver()
	if err != nil {
		sourceCloseErr := source.Close()
		return nil, nil, errors.Join(
			fmt.Errorf("create migration driver: %w", err),
			wrapCloseError("close migration source", sourceCloseErr),
		)
	}

	m, err := migrate.NewWithInstance(
		spec.SourceName, source,
		spec.DatabaseName, database,
	)
	if err != nil {
		sourceCloseErr := source.Close()

		if spec.KeepDatabaseOpen {
			return nil, nil, errors.Join(
				fmt.Errorf("create migration instance: %w", err),
				wrapCloseError("close migration source", sourceCloseErr),
			)
		}

		databaseCloseErr := database.Close()
		return nil, nil, errors.Join(
			fmt.Errorf("create migration instance: %w", err),
			wrapCloseError("close migration source", sourceCloseErr),
			wrapCloseError("close migration database", databaseCloseErr),
		)
	}

	var close func() error
	if spec.KeepDatabaseOpen {
		close = func() error {
			sourceCloseErr := source.Close()
			return wrapCloseError("close migration source", sourceCloseErr)
		}
	} else {
		close = func() error {
			sourceCloseErr := source.Close()
			databaseCloseErr := database.Close()
			return errors.Join(
				wrapCloseError("close migration source", sourceCloseErr),
				wrapCloseError("close migration database", databaseCloseErr),
			)
		}
	}

	return m, close, nil
}
