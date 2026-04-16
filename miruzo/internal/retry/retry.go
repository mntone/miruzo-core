package retry

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

func waitWithContext(requestContext context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-requestContext.Done():
		return requestContext.Err()
	case <-timer.C:
		return nil
	}
}

var sleepWithContext = waitWithContext

type retryFunc[T any] func(requestContext context.Context) (T, error)

func Retry[T any](
	requestContext context.Context,
	policy backoff.Policy,
	retry retryFunc[T],
) (T, error) {
	var maxAttempts uint32
	if policy != nil {
		maxAttempts = policy.GetMaxAttempts()
	} else {
		maxAttempts = 1
	}

	var rng *rand.Rand
	if policy != nil {
		rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}

	var zero T
	var lastError error
	for attempt := range maxAttempts {
		if err := requestContext.Err(); err != nil {
			return zero, err
		}

		result, err := retry(requestContext)
		if err == nil {
			return result, nil
		}

		lastError = err
		if !persist.IsRecoverable(err) {
			return zero, err
		}
		if attempt+1 >= maxAttempts {
			return zero, err
		}

		var delay time.Duration
		if policy != nil {
			delay = policy.NextDelay(attempt, rng)
		}

		if err := sleepWithContext(requestContext, delay); err != nil {
			return zero, err
		}
	}

	return zero, lastError
}
