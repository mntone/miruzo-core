package persist

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type UserRepository interface {
	Get(
		ctx context.Context,
	) (User, error)
}

type SessionUserRepository interface {
	UserRepository

	IncrementDailyLoveUsed(
		ctx context.Context,
		dailyLoveLimit model.QuotaInt,
	) (model.QuotaInt, error)

	DecrementDailyLoveUsed(ctx context.Context) (model.QuotaInt, error)

	ResetDailyLoveUsed(ctx context.Context) error
}
