package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type userRepository struct {
	queries *gen.Queries
}

func (repo userRepository) Get(
	ctx context.Context,
) (persist.User, error) {
	user, err := repo.queries.GetUser(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return persist.User{}, persist.ErrNoRows
		}

		return persist.User{}, dberrors.ToPersist("Get", err)
	}

	return persist.User{
		ID:            int16(user.ID),
		DailyLoveUsed: model.QuotaInt(user.DailyLoveUsed),
	}, nil
}

func (repo userRepository) IncrementDailyLoveUsed(
	ctx context.Context,
	dailyLoveLimit model.QuotaInt,
) (model.QuotaInt, error) {
	dailyLoveUsed, err := repo.queries.IncrementDailyLoveUsed(ctx, int32(dailyLoveLimit))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, persist.ErrQuotaExceeded
		}

		return 0, dberrors.ToPersist("IncrementDailyLoveUsed", err)
	}

	return model.QuotaInt(dailyLoveUsed), nil
}

func (repo userRepository) DecrementDailyLoveUsed(ctx context.Context) (model.QuotaInt, error) {
	dailyLoveUsed, err := repo.queries.DecrementDailyLoveUsed(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, persist.ErrNoRows
		}

		mapError := dberrors.ToPersist("DecrementDailyLoveUsed", err)
		if errors.Is(mapError, persist.ErrCheckViolation) {
			return 0, persist.ErrQuotaUnderflow
		}
		return 0, mapError
	}

	return model.QuotaInt(dailyLoveUsed), nil
}

func (repo userRepository) ResetDailyLoveUsed(ctx context.Context) error {
	rowCount, err := repo.queries.ResetDailyLoveUsed(ctx)
	if err != nil {
		return dberrors.ToPersist("ResetDailyLoveUsed", err)
	}

	if rowCount == 0 {
		return persist.ErrNoRows
	}

	return nil
}
