package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/mysql/gen"
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
	rowCount, err := repo.queries.IncrementDailyLoveUsed(ctx, int8(dailyLoveLimit))
	if err != nil {
		return 0, dberrors.ToPersist("IncrementDailyLoveUsed", err)
	}

	if rowCount != 1 {
		return 0, persist.ErrQuotaExceeded
	}

	dailyLoveUsed, err := repo.queries.GetDailyLoveUsed(ctx)
	if err != nil {
		return 0, dberrors.ToPersist("IncrementDailyLoveUsed", err)
	}

	return model.QuotaInt(dailyLoveUsed), nil
}

func (repo userRepository) DecrementDailyLoveUsed(ctx context.Context) (model.QuotaInt, error) {
	rowCount, err := repo.queries.DecrementDailyLoveUsed(ctx)
	if err != nil {
		mapError := dberrors.ToPersist("DecrementDailyLoveUsed", err)
		if errors.Is(mapError, persist.ErrCheckViolation) {
			return 0, persist.ErrQuotaUnderflow
		}
		return 0, mapError
	}

	if rowCount != 1 {
		return 0, persist.ErrNoRows
	}

	dailyLoveUsed, err := repo.queries.GetDailyLoveUsed(ctx)
	if err != nil {
		return 0, dberrors.ToPersist("DecrementDailyLoveUsed", err)
	}

	return model.QuotaInt(dailyLoveUsed), nil
}

func (repo userRepository) ResetDailyLoveUsed(ctx context.Context) error {
	rowCount, err := repo.queries.ResetDailyLoveUsed(ctx)
	if err != nil {
		return dberrors.ToPersist("ResetDailyLoveUsed", err)
	}

	if rowCount != 1 {
		exists, err := repo.queries.ExistsUser(ctx)
		if err != nil {
			return dberrors.ToPersist("ResetDailyLoveUsed", err)
		}
		if !exists {
			return persist.ErrNoRows
		}
	}

	return nil
}
