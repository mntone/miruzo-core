package sqlite

import (
	"context"
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/migrations_min"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func NewMigrationRunnerFromDB(db *sql.DB) shared.MigrationRunner {
	return shared.MigrationRunner{
		Spec: migrations.NewSpec(db),
	}
}

func NewMigrationRunner(
	ctx context.Context,
	conf config.DatabaseConfig,
) (persist.MigrationRunner, error) {
	db, err := database.OpenDatabase(ctx, conf)
	if err != nil {
		return nil, err
	}

	return shared.MigrationRunner{
		Spec:          migrations.NewSpec(db),
		CloseDatabase: db.Close,
	}, nil
}
