package postgre

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	db "github.com/mntone/miruzo-core/miruzo/internal/database/postgre"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type repositoryFactory struct {
	pool *pgxpool.Pool
}

func newRepositoryFactory(pool *pgxpool.Pool) repositoryFactory {
	return repositoryFactory{
		pool: pool,
	}
}

func NewRepositoryFactory(
	ctx context.Context,
	conf config.DatabaseConfig,
) (repositoryFactory, error) {
	pool, err := db.OpenDatabase(ctx, conf)
	if err != nil {
		return repositoryFactory{}, err
	}

	return newRepositoryFactory(pool), nil
}

func (factory repositoryFactory) Close() error {
	factory.pool.Close()
	return nil
}

func (factory repositoryFactory) NewImageList() persist.ImageListRepository {
	return imagelist.NewRepository(factory.pool)
}
