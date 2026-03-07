package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type userRepository struct {
	queries *gen.Queries
}

func NewRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{
		queries: gen.New(pool),
	}
}

func (repo *userRepository) GetSingletonUser(
	ctx context.Context,
) (persist.User, error) {
	user, err := repo.queries.GetSingletonUser(ctx)
	if err != nil {
		return persist.User{}, shared.MapPostgreError("GetSingletonUser", err)
	}

	return persist.User{
		ID:            user.ID,
		DailyLoveUsed: uint16(user.DailyLoveUsed),
	}, nil
}
