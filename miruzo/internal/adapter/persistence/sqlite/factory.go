package sqlite

import (
	"context"
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/action"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/imagelist"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/user"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	db "github.com/mntone/miruzo-core/miruzo/internal/database/sqlite"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type repositoryFactory struct {
	database *sql.DB
}

func newRepositoryFactory(db *sql.DB) repositoryFactory {
	return repositoryFactory{
		database: db,
	}
}

func NewRepositoryFactory(
	ctx context.Context,
	conf config.DatabaseConfig,
) (repositoryFactory, error) {
	db, err := db.OpenDatabase(ctx, conf)
	if err != nil {
		return repositoryFactory{}, err
	}

	return newRepositoryFactory(db), nil
}

func (factory repositoryFactory) Close() error {
	return factory.database.Close()
}

func (factory repositoryFactory) NewAction() persist.ActionRepository {
	return action.NewRepository(factory.database)
}

func (factory repositoryFactory) NewImageList() persist.ImageListRepository {
	return imagelist.NewRepository(factory.database)
}

func (factory repositoryFactory) NewUser() persist.UserRepository {
	return user.NewRepository(factory.database)
}
