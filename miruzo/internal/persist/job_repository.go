package persist

import (
	"context"
	"time"
)

type JobRepository interface {
	MarkStarted(
		ctx context.Context,
		name string,
		startedAt time.Time,
	) error

	MarkFinished(
		ctx context.Context,
		name string,
		finishedAt time.Time,
	) error
}
