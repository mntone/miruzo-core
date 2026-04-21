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
	once    sync.Once
	pool    *pgxpool.Pool
	initErr error
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
// Priority:
// 1) MIRUZO_TEST_POSTGRES_URL (externally managed DB)
// 2) auto-start container (local dev / VSCode)
func GetPostgresTestPool(t testing.TB, reg *testutil.CleanupRegistry) *pgxpool.Pool {
	t.Helper()

	once.Do(func() {
		ctx := context.Background()
		dsn := os.Getenv(TEST_DSN_ENVNAME)
		if dsn == "" {
			container, err := startPostgresContainer(ctx)
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
				initErr = fmt.Errorf("postgres get dsn: %w", err)
				return
			}
		}

		localPool, err := openPoolFromDSN(ctx, dsn)
		if err != nil {
			initErr = fmt.Errorf("postgres open: %w", err)
			return
		}
		reg.Register(func() error {
			localPool.Close()
			return nil
		})

		err = postgres.NewMigrationRunnerFromPool(localPool).Up(ctx)
		if err != nil {
			initErr = fmt.Errorf("postgres migrate: %w", err)
			return
		}

		pool = localPool
	})
	if initErr != nil {
		t.Fatal(initErr)
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
