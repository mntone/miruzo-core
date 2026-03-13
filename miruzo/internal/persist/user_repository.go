package persist

import "context"

type UserRepository interface {
	GetSingletonUser(
		requestContext context.Context,
	) (User, error)

	IncrementDailyLoveUsed(
		requestContext context.Context,
		dailyLoveLimit int16,
	) (int16, error)
}
