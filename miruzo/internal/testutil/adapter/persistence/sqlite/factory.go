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
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/migrations_min"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
)

func setupDatabase(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	db := OpenTestDatabase(t, ctx)
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("run sqlite migrations: %v", err)
	}
	return db
}

func newOperations(
	ctx context.Context,
	queries *gen.Queries,
) testutilPersistence.Operations {
	return testutilPersistence.NewOperations(
		ctx,
		newRepository(queries),
	)
}

func NewActionSuite(t *testing.T) testutilPersistence.ActionSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.ActionSuite{
		Context:    ctx,
		Operations: newOperations(ctx, queries),
		Repository: action.NewRepository(queries),
	}
}

func NewImageListSuite(t *testing.T) testutilPersistence.ImageListSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.ImageListSuite{
		Context:    ctx,
		Operations: newOperations(ctx, queries),
		Repository: imagelist.NewRepository(queries),
	}
}

func NewUserSuite(t *testing.T) testutilPersistence.UserSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return testutilPersistence.UserSuite{
		Context:    ctx,
		Operations: newOperations(ctx, queries),
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
		Operations: newOperations(ctx, queries),
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
		Operations:     newOperations(ctx, queries),
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
		Operations: newOperations(ctx, queries),
		Repository: sqlite.NewViewRepository(queries),
	}
}
