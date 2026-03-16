package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/stats"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/user"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	migrations "github.com/mntone/miruzo-core/miruzo/internal/database/postgres/migrations_min"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
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

	pool, err := openTestPoolFromContainer(ctx, container)
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

	if err := migrations.RunMigrations(pool); err != nil {
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

func (ste *SuiteFactory) newOperations(
	ctx context.Context,
	pool *pgxpool.Pool,
	queries *gen.Queries,
) testutilPersistence.Operations {
	return testutilPersistence.NewOperations(
		ctx,
		action.NewRepository(queries),
		newRepository(pool, queries),
	)
}

func (ste *SuiteFactory) NewAction(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.ActionSuite {
	t.Helper()

	queries := gen.New(ste.pool)
	ops := ste.newOperations(ctx, ste.pool, queries)
	return testutilPersistence.ActionSuite{
		Context:    ctx,
		Operations: ops,
		Repository: ops.Action,
	}
}

func (ste *SuiteFactory) NewImageList(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.ImageListSuite {
	t.Helper()

	queries := gen.New(ste.pool)
	return testutilPersistence.ImageListSuite{
		Context:    ctx,
		Operations: ste.newOperations(ctx, ste.pool, queries),
		Repository: imagelist.NewRepository(queries),
	}
}

func (ste *SuiteFactory) NewUser(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.UserSuite {
	t.Helper()

	queries := gen.New(ste.pool)
	return testutilPersistence.UserSuite{
		Context:    ctx,
		Operations: ste.newOperations(ctx, ste.pool, queries),
		Repository: user.NewRepository(queries),
	}
}

func (ste *SuiteFactory) NewSettings(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.SettingsSuite {
	t.Helper()

	queries := gen.New(ste.pool)
	return testutilPersistence.SettingsSuite{
		Context:    ctx,
		Operations: ste.newOperations(ctx, ste.pool, queries),
		Repository: postgres.NewSettingsRepository(queries),
	}
}

func (ste *SuiteFactory) NewStats(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.StatsSuite {
	t.Helper()

	queries := gen.New(ste.pool)
	return testutilPersistence.StatsSuite{
		Context:        ctx,
		Operations:     ste.newOperations(ctx, ste.pool, queries),
		Repository:     stats.NewRepository(queries),
		ViewRepository: postgres.NewViewRepository(queries),
	}
}

func (ste *SuiteFactory) NewView(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.ViewSuite {
	t.Helper()

	queries := gen.New(ste.pool)
	return testutilPersistence.ViewSuite{
		Context:    ctx,
		Operations: ste.newOperations(ctx, ste.pool, queries),
		Repository: postgres.NewViewRepository(queries),
	}
}
