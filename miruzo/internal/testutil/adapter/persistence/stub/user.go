package stub

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type userStorage struct {
	DailyLoveUsed model.QuotaInt
}

type UserRepository struct {
	userStorage

	GetError                 error
	IncrementError           error
	IncrementDailyLoveLimits []model.QuotaInt
	DecrementError           error
	ResetError               error
}

func NewStubUserRepository(dailyLoveUsed int32) *UserRepository {
	return &UserRepository{
		userStorage: userStorage{
			DailyLoveUsed: model.QuotaInt(dailyLoveUsed),
		},
	}
}

func (repo UserRepository) snapshot() userStorage {
	return userStorage{
		DailyLoveUsed: repo.DailyLoveUsed,
	}
}

func (repo UserRepository) Get(
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
	repo.IncrementDailyLoveLimits = append(repo.IncrementDailyLoveLimits, dailyLoveLimit)

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

func (repo *UserRepository) ResetDailyLoveUsed(
	ctx context.Context,
) error {
	if repo.ResetError != nil {
		return repo.ResetError
	}

	repo.DailyLoveUsed = 0
	return nil
}
