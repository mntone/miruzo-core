package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/migrations"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
)

func joinErrors(primary error, secondary error) error {
	if secondary == nil {
		return primary
	}

	return fmt.Errorf("%w (cleanup failed: %v)", primary, secondary)
}

type Suite struct {
	pool     *pgxpool.Pool
	close    func() error
	testRepo repository
}

func NewSuite(ctx context.Context) (*Suite, error) {
	container, err := startPostgreContainer(ctx)
	if err != nil {
		return nil, err
	}

	pool, err := openTestPoolFromContainer(ctx, container)
	if err != nil {
		containerErr := container.Terminate(ctx)
		return nil, fmt.Errorf(
			"open postgre test pool: %w",
			joinErrors(err, containerErr),
		)
	}

	closeFn := func() error {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate postgre container: %w", err)
		}

		return nil
	}

	if err := migrations.RunMigrations(pool); err != nil {
		closeErr := closeFn()
		return nil, fmt.Errorf(
			"run postgre migrations: %w",
			joinErrors(err, closeErr),
		)
	}

	return &Suite{
		pool:     pool,
		close:    closeFn,
		testRepo: NewRepository(pool),
	}, nil
}

func (ste *Suite) Close() error {
	if ste.close == nil {
		return nil
	}

	closeFn := ste.close
	ste.close = nil
	return closeFn()
}

func (ste *Suite) Reset(ctx context.Context) error {
	_, err := ste.pool.Exec(ctx, "TRUNCATE TABLE stats, images, ingests RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("reset postgre database: %w", err)
	}

	return nil
}

func (ste *Suite) MustReset(
	tb testing.TB,
	ctx context.Context,
) {
	tb.Helper()

	if err := ste.Reset(ctx); err != nil {
		tb.Fatalf("reset postgre test suite: %v", err)
	}
}

func (ste *Suite) NewImageList(ctx context.Context) testutilPersistence.ImageListSetup {
	return testutilPersistence.ImageListSetup{
		Ctx:  ctx,
		Ops:  testutilPersistence.NewOperations(ctx, ste.testRepo),
		Repo: imagelist.NewRepository(ste.pool),
	}
}
