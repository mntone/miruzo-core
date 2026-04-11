package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
)

type SuiteFactory struct {
	pool  *pgxpool.Pool
	close func() error
}

func NewSuiteFactory(ctx context.Context) (*SuiteFactory, error) {
	container, err := startPostgreContainer(ctx)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("get postgres container dsn: %w", err)
	}

	pool, err := openTestPoolFromDSN(ctx, dsn)
	if err != nil {
		containerErr := container.Terminate(ctx)
		return nil, fmt.Errorf(
			"open postgres test pool: %w",
			shared.JoinErrors(err, containerErr),
		)
	}

	closeFn := func() error {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate postgres container: %w", err)
		}

		return nil
	}

	if err := postgres.NewMigrationRunnerFromPool(pool).Up(ctx); err != nil {
		closeErr := closeFn()
		return nil, fmt.Errorf(
			"run postgres migrations: %w",
			shared.JoinErrors(err, closeErr),
		)
	}

	return &SuiteFactory{
		pool:  pool,
		close: closeFn,
	}, nil
}

func (ste *SuiteFactory) Close() error {
	if ste.close == nil {
		return nil
	}

	closeFn := ste.close
	ste.close = nil
	return closeFn()
}

func (ste *SuiteFactory) Reset(ctx context.Context) error {
	_, err := ste.pool.Exec(ctx, "TRUNCATE TABLE actions, stats, images, ingests RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("reset postgres database: %w", err)
	}

	_, err = ste.pool.Exec(ctx, "INSERT INTO users(id) VALUES(1) ON CONFLICT (id) DO UPDATE SET daily_love_used=0;")
	if err != nil {
		return fmt.Errorf("reset postgres database: %w", err)
	}

	return nil
}

func (ste *SuiteFactory) MustReset(
	t *testing.T,
	ctx context.Context,
) {
	t.Helper()

	if err := ste.Reset(ctx); err != nil {
		t.Fatalf("reset postgres test suite: %v", err)
	}
}
