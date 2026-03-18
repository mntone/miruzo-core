package persist

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type UserRepository interface {
	GetSingletonUser(
		requestContext context.Context,
	) (User, error)

	IncrementDailyLoveUsed(
		requestContext context.Context,
		dailyLoveLimit model.QuotaInt,
	) (model.QuotaInt, error)

	DecrementDailyLoveUsed(ctx context.Context) (model.QuotaInt, error)
}
