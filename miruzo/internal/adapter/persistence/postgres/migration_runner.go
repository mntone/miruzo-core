package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/postgres/migrations_min"
)

func NewMigrationRunnerFromPool(pool *pgxpool.Pool) shared.MigrationRunner {
	return shared.MigrationRunner{
		Spec: migrations.NewSpec(pool),
	}
}
