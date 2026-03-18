package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
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

func (repo repository) GetSingletonUser(
	ctx context.Context,
) (persist.User, error) {
	user, err := repo.queries.GetSingletonUser(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return persist.User{}, persist.ErrNotFound
		}

		return persist.User{}, shared.MapSQLiteError("GetSingletonUser", err)
	}

	return persist.User{
		ID:            int16(user.ID),
		DailyLoveUsed: model.QuotaInt(user.DailyLoveUsed),
	}, nil
}

func (repo repository) IncrementDailyLoveUsed(
	ctx context.Context,
	dailyLoveLimit model.QuotaInt,
) (model.QuotaInt, error) {
	dailyLoveUsed, err := repo.queries.IncrementDailyLoveUsed(ctx, int32(dailyLoveLimit))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, persist.ErrQuotaExceeded
		}

		return 0, shared.MapSQLiteError("IncrementDailyLoveUsed", err)
	}

	return model.QuotaInt(dailyLoveUsed), nil
}

func (repo repository) DecrementDailyLoveUsed(ctx context.Context) (model.QuotaInt, error) {
	dailyLoveUsed, err := repo.queries.DecrementDailyLoveUsed(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, persist.ErrNotFound
		}

		mapError := shared.MapSQLiteError("DecrementDailyLoveUsed", err)
		if errors.Is(mapError, persist.ErrCheckViolation) {
			return 0, persist.ErrQuotaUnderflow
		}
		return 0, mapError
	}

	return model.QuotaInt(dailyLoveUsed), nil
}
