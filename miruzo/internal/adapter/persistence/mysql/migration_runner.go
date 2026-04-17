package mysql

import (
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/mysql/migrations_min"
)

func NewMigrationRunnerFromDB(db *sql.DB) shared.MigrationRunner {
	return shared.MigrationRunner{
		Spec: migrations.NewSpec(db),
	}
}
