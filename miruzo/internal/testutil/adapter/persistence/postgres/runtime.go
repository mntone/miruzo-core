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
	"github.com/mntone/miruzo-core/miruzo/internal/testutil"
)

const (
	TEST_DSN_ENVNAME = "TEST_POSTGRES_URL"
)

var (
	once    sync.Once
	pool    *pgxpool.Pool
	initErr error
)

// GetPostgresTestPool returns a shared pgx pool for this package test process.
// Priority:
// 1) TEST_POSTGRES_URL (externally managed DB)
// 2) auto-start container (local dev / VSCode)
func GetPostgresTestPool(t testing.TB, reg *testutil.CleanupRegistry) *pgxpool.Pool {
	t.Helper()

	once.Do(func() {
		ctx := context.Background()
		dsn := os.Getenv(TEST_DSN_ENVNAME)
		if dsn == "" {
			container, err := startPostgreContainer(ctx)
			if err != nil {
				initErr = err
				return
			}
			reg.Register(func() error {
				closeContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				return container.Terminate(closeContext)
			})

			dsn, initErr = container.ConnectionString(ctx)
			if initErr != nil {
				return
			}
		}

		pool, initErr = openTestPoolFromDSN(ctx, dsn)
		if initErr != nil {
			return
		}
		reg.Register(func() error {
			pool.Close()
			return nil
		})

		initErr = postgres.NewMigrationRunnerFromPool(pool).Up(ctx)
	})
	if initErr != nil {
		t.Fatalf("postgres init failed: %v", initErr)
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
