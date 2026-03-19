package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
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

func (repo repository) Get(
	ctx context.Context,
) (persist.User, error) {
	user, err := repo.queries.GetUser(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return persist.User{}, persist.ErrNotFound
		}

		return persist.User{}, shared.MapPostgreError("Get", err)
	}

	return persist.User{
		ID:            user.ID,
		DailyLoveUsed: model.QuotaInt(user.DailyLoveUsed),
	}, nil
}

func (repo repository) IncrementDailyLoveUsed(
	ctx context.Context,
	dailyLoveLimit model.QuotaInt,
) (model.QuotaInt, error) {
	dailyLoveUsed, err := repo.queries.IncrementDailyLoveUsed(ctx, int32(dailyLoveLimit))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, persist.ErrQuotaExceeded
		}

		return 0, shared.MapPostgreError("IncrementDailyLoveUsed", err)
	}

	return model.QuotaInt(dailyLoveUsed), nil
}

func (repo repository) DecrementDailyLoveUsed(ctx context.Context) (model.QuotaInt, error) {
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

	return model.QuotaInt(dailyLoveUsed), nil
}
