package jobguard

import (
	"context"
	"time"
)

type JobRunGuard interface {
	TryAcquire(
		ctx context.Context,
		name string,
		startedAt time.Time,
	) (acquired bool, err error)

	Release(
		ctx context.Context,
		name string,
		finishedAt time.Time,
	) error
}
