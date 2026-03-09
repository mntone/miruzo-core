package user

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(queries *gen.Queries) repository {
	return repository{
		queries: queries,
	}
}

func (repo repository) GetSingletonUser(
	ctx context.Context,
) (persist.User, error) {
	user, err := repo.queries.GetSingletonUser(ctx)
	if err != nil {
		return persist.User{}, shared.MapSQLiteError("GetSingletonUser", err)
	}

	return persist.User{
		ID:            int16(user.ID),
		DailyLoveUsed: uint16(user.DailyLoveUsed),
	}, nil
}
