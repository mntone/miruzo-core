package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	database "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/shared"
)

func openTestDatabase(t testing.TB, ctx context.Context) *sql.DB {
	t.Helper()

	cfg := database.ConnectConfig{
		DSN:              fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()),
		ConnectionTuning: shared.NewTestConnectionTuning(),
	}
	databaseHandle, err := database.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}

	t.Cleanup(func() {
		_ = databaseHandle.Close()
	})

	return databaseHandle
}
