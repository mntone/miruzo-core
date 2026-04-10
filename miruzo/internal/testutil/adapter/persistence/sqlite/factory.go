package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
)

func setupDatabase(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	db := openTestDatabase(t, ctx)
	if err := sqlite.NewMigrationRunnerFromDB(db).Up(ctx); err != nil {
		t.Fatalf("run sqlite migrations: %v", err)
	}
	return db
}

func newOperations(
	ctx context.Context,
	db *sql.DB,
	queries *gen.Queries,
) testutilPersistence.Operations {
	return testutilPersistence.NewOperations(
		ctx,
		action.NewRepository(queries),
		newRepository(db, queries),
	)
}

func NewImageListSuite(t *testing.T) testutilPersistence.ImageListSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.ImageListSuite{
		Context:    ctx,
		Operations: newOperations(ctx, db, queries),
		Repository: imagelist.NewRepository(queries),
	}
}

func NewStatsSuite(t *testing.T) testutilPersistence.StatsSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.StatsSuite{
		Context:        ctx,
		Operations:     newOperations(ctx, db, queries),
		Repository:     stats.NewRepository(queries),
		ViewRepository: sqlite.NewViewRepository(queries),
	}
}
