package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/migrations"
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

func NewImageListSuite(t *testing.T) persistence.ImageListSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	return persistence.ImageListSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(db)),
		Repository: imagelist.NewRepository(db),
	}
}

func NewUserSuite(t *testing.T) persistence.UserSuite {
	t.Helper()

	ctx := context.Background()
	db := setupDatabase(t, ctx)
	return persistence.UserSuite{
		Context:    ctx,
		Operations: persistence.NewOperations(ctx, newRepository(db)),
		Repository: user.NewRepository(db),
	}
}
