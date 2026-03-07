package persist

import "context"

type UserRepository interface {
	GetSingletonUser(
		requestContext context.Context,
	) (User, error)
}
