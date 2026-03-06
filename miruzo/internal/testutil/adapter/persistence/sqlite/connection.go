package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
	db "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
)

func OpenTestDatabase(t testing.TB, ctx context.Context) *sql.DB {
	t.Helper()

	conf := config.DefaultDatabaseConfig()
	conf.DSN = fmt.Sprintf(
		"file:%s?mode=memory&cache=shared&_foreign_keys=on",
		t.Name(),
	)
	conf.MaxOpenConnections = 1

	databaseHandle, err := db.OpenDatabase(ctx, conf)
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}

	t.Cleanup(func() {
		_ = databaseHandle.Close()
	})

	return databaseHandle
}
