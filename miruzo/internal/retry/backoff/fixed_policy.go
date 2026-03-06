package backoff

import (
	"math/rand/v2"
	"time"
)

type FixedPolicy struct {
	MaxAttempts uint32
	Delays      []time.Duration
}

var _ Policy = FixedPolicy{}

func (policy FixedPolicy) GetMaxAttempts() uint32 {
	if policy.MaxAttempts == 0 {
		return 1
	}

	return policy.MaxAttempts
}

func (policy FixedPolicy) NextDelay(attempt uint32, rng *rand.Rand) time.Duration {
	if int(attempt) >= len(policy.Delays) {
		return 0
	}

	return policy.Delays[attempt]
}
