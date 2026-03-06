package backoff

import (
	"math/rand/v2"
	"time"
)

type NoRetryPolicy struct{}

var _ Policy = NoRetryPolicy{}

func (policy NoRetryPolicy) GetMaxAttempts() uint32 {
	return 1
}

func (policy NoRetryPolicy) NextDelay(attempt uint32, rng *rand.Rand) time.Duration {
	return 0
}
