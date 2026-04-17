package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/mysql"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/shared"
)

const (
	TEST_DSN_ENVNAME = "MIRUZO_TEST_MYSQL_URL"
)

var (
	once    sync.Once
	db      *sql.DB
	initErr error
)

func openDBFromDSN(ctx context.Context, dsn string) (*sql.DB, error) {
	cfg := database.ConnectConfig{
		DSN:              dsn,
		MultiStatements:  true,
		ConnectionTuning: shared.NewTestConnectionTuning(),
	}

	return database.Open(ctx, cfg)
}

// GetMySQLTestDB returns a shared *sql.DB for this package test process.
// Priority:
// 1) MIRUZO_TEST_MYSQL_URL (externally managed DB)
// 2) auto-start container (local dev / VSCode)
func GetMySQLTestDB(t testing.TB, reg *testutil.CleanupRegistry) *sql.DB {
	t.Helper()

	once.Do(func() {
		ctx := context.Background()
		dsn := os.Getenv(TEST_DSN_ENVNAME)
		if dsn == "" {
			container, err := startMySQLContainer(ctx)
			if err != nil {
				initErr = err
				return
			}
			reg.Register(func() error {
				closeContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				return container.Terminate(closeContext)
			})

			dsn, err = container.ConnectionString(ctx)
			if err != nil {
				initErr = fmt.Errorf("mysql get dsn: %w", err)
				return
			}
		}

		localDB, err := openDBFromDSN(ctx, dsn)
		if err != nil {
			initErr = fmt.Errorf("mysql open: %w", err)
			return
		}
		reg.Register(func() error {
			localDB.Close()
			return nil
		})

		err = mysql.NewMigrationRunnerFromDB(localDB).Up(ctx)
		if err != nil {
			initErr = fmt.Errorf("mysql migrate: %w", err)
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
	if os.Getenv(TEST_DSN_ENVNAME) != "" {
		return fmt.Sprintf("external(%s)", TEST_DSN_ENVNAME)
	}
	return "auto(testcontainers)"
}

func DebugInfo(t *testing.T) {
	t.Helper()
	t.Logf("mysql source: %s", MustDSNForLog())
	fmt.Println()
}
