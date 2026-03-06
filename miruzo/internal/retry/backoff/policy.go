package backoff

import (
	"math/rand/v2"
	"time"
)

type Policy interface {
	GetMaxAttempts() uint32
	NextDelay(attempt uint32, rng *rand.Rand) time.Duration
}
