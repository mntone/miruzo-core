package testutil

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/migrations"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
)

type OperationsSetup struct {
	Ctx context.Context
	DB  *sql.DB
	Ops testutilPersistence.Operations
}

func SetupOperations(tb testing.TB) OperationsSetup {
	tb.Helper()

	ctx := context.Background()
	db := OpenTestDatabase(tb, ctx)

	if err := migrations.RunMigrations(db); err != nil {
		tb.Fatalf("run sqlite migrations: %v", err)
	}

	testRepo := NewRepository(db)

	return OperationsSetup{
		Ctx: ctx,
		DB:  db,
		Ops: testutilPersistence.NewOperations(ctx, testRepo),
	}
}
