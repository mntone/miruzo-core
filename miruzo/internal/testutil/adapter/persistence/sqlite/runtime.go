package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/shared"
)

var (
	once    sync.Once
	db      *sql.DB
	initErr error
)

func openMemoryDB(ctx context.Context, name string) (*sql.DB, error) {
	cfg := database.ConnectConfig{
		DSN:              fmt.Sprintf("file:%s?mode=memory&cache=shared", name),
		ConnectionTuning: shared.NewTestConnectionTuning(),
	}

	return database.Open(ctx, cfg)
}

func GetSQLiteTestDB(t testing.TB, reg *testutil.CleanupRegistry) *sql.DB {
	t.Helper()

	once.Do(func() {
		ctx := context.Background()
		localDB, err := openMemoryDB(ctx, t.Name())
		if err != nil {
			initErr = fmt.Errorf("sqlite open: %w", err)
			return
		}
		reg.Register(localDB.Close)

		err = sqlite.NewMigrationRunnerFromDB(localDB).Up(ctx)
		if err != nil {
			initErr = fmt.Errorf("sqlite migrate: %w", err)
			return
		}

		db = localDB
	})
	if initErr != nil {
		t.Fatal(initErr)
	}

	return db
}

func MustDSNForLog() string {
	return "auto(memory)"
}

func DebugInfo(t *testing.T) {
	t.Helper()
	t.Logf("sqlite source: %s", MustDSNForLog())
	fmt.Println()
}
