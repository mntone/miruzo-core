package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
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

func NewActionSuite(t *testing.T) testutilPersistence.ActionSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	ops := newOperations(ctx, db, queries)
	return testutilPersistence.ActionSuite{
		Context:    ctx,
		Operations: ops,
		Repository: ops.Action,
	}
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

func NewIngestSuite(t *testing.T) testutilPersistence.IngestSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.IngestSuite{
		Context:    ctx,
		Operations: newOperations(ctx, db, queries),
	}
}

func NewJobSuite(t *testing.T) testutilPersistence.JobSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.JobSuite{
		Context:    ctx,
		Operations: newOperations(ctx, db, queries),
		Repository: sqlite.NewJobRepository(queries),
	}
}

func NewUserSuite(t *testing.T) testutilPersistence.UserSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.UserSuite{
		Context:    ctx,
		Operations: newOperations(ctx, db, queries),
		Repository: user.NewRepository(queries),
	}
}

func NewSettingsSuite(t *testing.T) testutilPersistence.SettingsSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.SettingsSuite{
		Context:    ctx,
		Operations: newOperations(ctx, db, queries),
		Repository: sqlite.NewSettingsRepository(queries),
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

func NewViewSuite(t *testing.T) testutilPersistence.ViewSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.ViewSuite{
		Context:    ctx,
		Operations: newOperations(ctx, db, queries),
		Repository: sqlite.NewViewRepository(queries),
	}
}
