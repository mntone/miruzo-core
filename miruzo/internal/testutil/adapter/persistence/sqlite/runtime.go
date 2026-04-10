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

func GetSQLiteTestDB(t testing.TB, reg *testutil.CleanupRegistry) *sql.DB {
	t.Helper()

	once.Do(func() {
		cfg := database.ConnectConfig{
			DSN:              fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()),
			ConnectionTuning: shared.NewTestConnectionTuning(),
		}

		ctx := context.Background()
		db, initErr = database.Open(ctx, cfg)
		if initErr != nil {
			return
		}
		reg.Register(db.Close)

		initErr = sqlite.NewMigrationRunnerFromDB(db).Up(ctx)
	})
	if initErr != nil {
		t.Fatalf("sqlite init failed: %v", initErr)
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
