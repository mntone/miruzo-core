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
	// Runtime test resources are process-global and shared within this package.
	dsnOnce     sync.Once
	resolvedDSN string
	dsnErr      error
)

func resolveMySQLTestDSN(
	ctx context.Context,
	reg *testutil.CleanupRegistry,
) (string, error) {
	dsn := os.Getenv(TEST_DSN_ENVNAME)
	if dsn != "" {
		return dsn, nil
	}

	container, err := startMySQLContainer(ctx)
	if err != nil {
		return "", err
	}
	reg.Register(func() error {
		closeContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		return container.Terminate(closeContext)
	})

	dsn, err = container.ConnectionString(ctx)
	if err != nil {
		return "", fmt.Errorf("mysql get dsn: %w", err)
	}

	return dsn, nil
}

func GetMySQLTestDSN(t testing.TB, reg *testutil.CleanupRegistry) string {
	t.Helper()

	// Priority:
	// 1) MIRUZO_TEST_MYSQL_URL (externally managed DB)
	// 2) auto-start container (local dev / VSCode)
	dsnOnce.Do(func() {
		resolvedDSN, dsnErr = resolveMySQLTestDSN(context.Background(), reg)
	})
	if dsnErr != nil {
		t.Fatal(dsnErr)
	}
	return resolvedDSN
}

var (
	// Runtime test resources are process-global and shared within this package.
	dbOnce sync.Once
	db     *sql.DB
	dbErr  error
)

func openDBFromDSN(ctx context.Context, dsn string) (*sql.DB, error) {
	cfg, err := database.NewConnectConfigFromDSN(dsn, database.ConnectOptions{
		MultiStatements:  true,
		ConnectionTuning: shared.NewTestConnectionTuning(),
	})
	if err != nil {
		return nil, err
	}

	return database.Open(ctx, cfg)
}

// GetMySQLTestDB returns a shared *sql.DB for this package test process.
// It opens the DB from the DSN resolved by GetMySQLTestDSN and applies schema
// migrations once.
func GetMySQLTestDB(t testing.TB, reg *testutil.CleanupRegistry) *sql.DB {
	t.Helper()

	dsn := GetMySQLTestDSN(t, reg)
	dbOnce.Do(func() {
		ctx := context.Background()

		localDB, err := openDBFromDSN(ctx, dsn)
		if err != nil {
			dbErr = fmt.Errorf("mysql open: %w", err)
			return
		}
		reg.Register(func() error {
			localDB.Close()
			return nil
		})

		err = mysql.NewMigrationRunnerFromDB(localDB).Up(ctx)
		if err != nil {
			dbErr = fmt.Errorf("mysql migrate: %w", err)
			return
		}

		db = localDB
	})
	if dbErr != nil {
		t.Fatal(dbErr)
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
