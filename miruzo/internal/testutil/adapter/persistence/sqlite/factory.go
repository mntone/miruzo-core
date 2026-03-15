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
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
)

func setupDatabase(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	db := OpenTestDatabase(t, ctx)
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("run sqlite migrations: %v", err)
	}
	return db
}

func NewActionSuite(t *testing.T) persistence.ActionSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return persistence.ActionSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(queries)),
		Repository: action.NewRepository(queries),
	}
}

func NewImageListSuite(t *testing.T) persistence.ImageListSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return persistence.ImageListSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(queries)),
		Repository: imagelist.NewRepository(queries),
	}
}

func NewUserSuite(t *testing.T) persistence.UserSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return persistence.UserSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(queries)),
		Repository: user.NewRepository(queries),
	}
}

func NewSettingsSuite(t *testing.T) persistence.SettingsSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return persistence.SettingsSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(queries)),
		Repository: sqlite.NewSettingsRepository(queries),
	}
}

func NewStatsSuite(t *testing.T) persistence.StatsSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return persistence.StatsSuite{
		Context:        ctx,
		Operations:     persistence.NewOperations(ctx, newRepository(queries)),
		Repository:     stats.NewRepository(queries),
		ViewRepository: sqlite.NewViewRepository(queries),
	}
}

func NewViewSuite(t *testing.T) persistence.ViewSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	queries := gen.New(db)
	return persistence.ViewSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(queries)),
		Repository: sqlite.NewViewRepository(queries),
	}
}
