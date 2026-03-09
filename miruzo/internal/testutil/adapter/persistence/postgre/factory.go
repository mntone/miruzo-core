package postgre

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/user"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/migrations"
	testutilPersistence "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence"
)

func joinErrors(primary error, secondary error) error {
	if secondary == nil {
		return primary
	}

	return fmt.Errorf("%w (cleanup failed: %v)", primary, secondary)
}

type SuiteFactory struct {
	pool     *pgxpool.Pool
	close    func() error
	testRepo repository
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

	return &SuiteFactory{
		pool:     pool,
		close:    closeFn,
		testRepo: newRepository(pool),
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
	_, err := ste.pool.Exec(ctx, "TRUNCATE TABLE stats, images, ingests RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("reset postgre database: %w", err)
	}

	_, err = ste.pool.Exec(ctx, "UPDATE users SET daily_love_used=0;")
	if err != nil {
		return fmt.Errorf("reset postgre database: %w", err)
	}

	return nil
}

func (ste *SuiteFactory) MustReset(
	t *testing.T,
	ctx context.Context,
) {
	t.Helper()

	if err := ste.Reset(ctx); err != nil {
		t.Fatalf("reset postgre test suite: %v", err)
	}
}

func (ste *SuiteFactory) NewAction(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.ActionSuite {
	t.Helper()

	return testutilPersistence.ActionSuite{
		Context:    ctx,
		Operations: testutilPersistence.NewOperations(ctx, ste.testRepo),
		Repository: action.NewRepository(gen.New(ste.pool)),
	}
}

func (ste *SuiteFactory) NewImageList(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.ImageListSuite {
	t.Helper()

	return testutilPersistence.ImageListSuite{
		Context:    ctx,
		Operations: testutilPersistence.NewOperations(ctx, ste.testRepo),
		Repository: imagelist.NewRepository(gen.New(ste.pool)),
	}
}

func (ste *SuiteFactory) NewUser(
	t *testing.T,
	ctx context.Context,
) testutilPersistence.UserSuite {
	t.Helper()

	return testutilPersistence.UserSuite{
		Context:    ctx,
		Operations: testutilPersistence.NewOperations(ctx, ste.testRepo),
		Repository: user.NewRepository(gen.New(ste.pool)),
	}
}
