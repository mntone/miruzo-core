package postgres

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	database "github.com/mntone/miruzo-core/miruzo/internal/database/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/shared"
)

const (
	TEST_DSN_ENVNAME = "MIRUZO_TEST_POSTGRES_URL"
)

var (
	// Runtime test resources are process-global and shared within this package.
	dsnOnce     sync.Once
	resolvedDSN string
	dsnErr      error
)

func resolvePostgresTestDSN(
	ctx context.Context,
	reg *testutil.CleanupRegistry,
) (string, error) {
	dsn := os.Getenv(TEST_DSN_ENVNAME)
	if dsn != "" {
		return dsn, nil
	}

	container, err := startPostgresContainer(ctx)
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
		return "", fmt.Errorf("postgres get dsn: %w", err)
	}

	return dsn, nil
}

func GetPostgresTestDSN(t testing.TB, reg *testutil.CleanupRegistry) string {
	t.Helper()

	// Priority:
	// 1) MIRUZO_TEST_POSTGRES_URL (externally managed DB)
	// 2) auto-start container (local dev / VSCode)
	dsnOnce.Do(func() {
		resolvedDSN, dsnErr = resolvePostgresTestDSN(context.Background(), reg)
	})
	if dsnErr != nil {
		t.Fatal(dsnErr)
	}
	return resolvedDSN
}

var (
	// Runtime test resources are process-global and shared within this package.
	poolOnce sync.Once
	pool     *pgxpool.Pool
	poolErr  error
)

func openPoolFromDSN(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := database.NewConnectConfigFromDSN(dsn, database.ConnectOptions{
		ConnectionTuning: shared.NewTestConnectionTuning(),
	})
	if err != nil {
		return nil, err
	}

	return database.Open(ctx, cfg)
}

// GetPostgresTestPool returns a shared pgx pool for this package test process.
// It opens the pool from the DSN resolved by GetPostgresTestDSN and applies
// schema migrations once.
func GetPostgresTestPool(t testing.TB, reg *testutil.CleanupRegistry) *pgxpool.Pool {
	t.Helper()

	dsn := GetPostgresTestDSN(t, reg)
	poolOnce.Do(func() {
		ctx := context.Background()

		localPool, err := openPoolFromDSN(ctx, dsn)
		if err != nil {
			poolErr = fmt.Errorf("postgres open: %w", err)
			return
		}
		reg.Register(func() error {
			localPool.Close()
			return nil
		})

		err = postgres.NewMigrationRunnerFromPool(localPool).Up(ctx)
		if err != nil {
			poolErr = fmt.Errorf("postgres migrate: %w", err)
			return
		}

		pool = localPool
	})
	if poolErr != nil {
		t.Fatal(poolErr)
	}

	return pool
}

func MustDSNForLog() string {
	if os.Getenv(TEST_DSN_ENVNAME) != "" {
		return fmt.Sprintf("external(%s)", TEST_DSN_ENVNAME)
	}
	return "auto(testcontainers)"
}

func DebugInfo(t *testing.T) {
	t.Helper()
	t.Logf("postgres source: %s", MustDSNForLog())
	fmt.Println()
}
