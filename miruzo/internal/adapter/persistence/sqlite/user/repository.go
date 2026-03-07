package user

import (
	"context"
	"database/sql"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type userRepository struct {
	queries *gen.Queries
}

func NewRepository(db *sql.DB) *userRepository {
	return &userRepository{
		queries: gen.New(db),
	}
}

func (repo *userRepository) GetSingletonUser(
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
