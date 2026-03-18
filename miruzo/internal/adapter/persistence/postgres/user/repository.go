package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
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
		return persist.User{}, shared.MapPostgreError("GetSingletonUser", err)
	}

	return persist.User{
		ID:            user.ID,
		DailyLoveUsed: user.DailyLoveUsed,
	}, nil
}

func (repo repository) IncrementDailyLoveUsed(
	ctx context.Context,
	dailyLoveLimit int16,
) (int16, error) {
	dailyLoveUsed, err := repo.queries.IncrementDailyLoveUsed(ctx, dailyLoveLimit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, persist.ErrQuotaExceeded
		}

		return 0, shared.MapPostgreError("IncrementDailyLoveUsed", err)
	}

	return dailyLoveUsed, nil
}

func (repo repository) DecrementDailyLoveUsed(ctx context.Context) (int16, error) {
	dailyLoveUsed, err := repo.queries.DecrementDailyLoveUsed(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, persist.ErrNotFound
		}

		mapError := shared.MapPostgreError("DecrementDailyLoveUsed", err)
		if errors.Is(mapError, persist.ErrCheckViolation) {
			return 0, persist.ErrQuotaUnderflow
		}
		return 0, mapError
	}

	return dailyLoveUsed, nil
}
