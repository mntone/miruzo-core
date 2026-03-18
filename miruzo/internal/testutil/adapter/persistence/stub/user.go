package stub

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type UserRepository struct {
	DailyLoveUsed model.QuotaInt

	GetError       error
	IncrementError error
	DecrementError error
}

func NewStubUserRepository(dailyLoveUsed int32) *UserRepository {
	return &UserRepository{
		DailyLoveUsed: model.QuotaInt(dailyLoveUsed),
	}
}

func (repo UserRepository) GetSingletonUser(
	ctx context.Context,
) (persist.User, error) {
	if repo.GetError != nil {
		return persist.User{}, repo.GetError
	}

	return persist.User{
		ID:            1,
		DailyLoveUsed: repo.DailyLoveUsed,
	}, nil
}

func (repo *UserRepository) IncrementDailyLoveUsed(
	ctx context.Context,
	dailyLoveLimit model.QuotaInt,
) (model.QuotaInt, error) {
	if repo.IncrementError != nil {
		return model.QuotaInt(0), repo.IncrementError
	}

	if repo.DailyLoveUsed >= dailyLoveLimit {
		return model.QuotaInt(0), persist.ErrQuotaExceeded
	}

	repo.DailyLoveUsed += 1
	return repo.DailyLoveUsed, nil
}

func (repo *UserRepository) DecrementDailyLoveUsed(
	ctx context.Context,
) (model.QuotaInt, error) {
	if repo.DecrementError != nil {
		return model.QuotaInt(0), repo.DecrementError
	}

	if repo.DailyLoveUsed <= 0 {
		return model.QuotaInt(0), persist.ErrQuotaUnderflow
	}

	repo.DailyLoveUsed -= 1
	return repo.DailyLoveUsed, nil
}
