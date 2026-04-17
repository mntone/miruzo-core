package sqlite

import (
	"database/sql"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/migrations_min"
)

func NewMigrationRunnerFromDB(db *sql.DB) persistshared.MigrationRunner {
	return persistshared.MigrationRunner{
		Spec: migrations.NewSpec(db),
	}
}
