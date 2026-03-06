package backoff

import (
	"math/rand/v2"
	"time"
)

type JitterType uint8

const (
	JitterNone JitterType = iota
	JitterEqual
	JitterFull
)

type ExponentialPolicy struct {
	MaxAttempts  uint32
	BaseDelay    time.Duration
	MinimumDelay time.Duration
	MaximumDelay time.Duration
	Jitter       JitterType
}

var _ Policy = ExponentialPolicy{}

func (policy ExponentialPolicy) GetMaxAttempts() uint32 {
	if policy.MaxAttempts == 0 {
		return 1
	}

	return policy.MaxAttempts
}

func (policy ExponentialPolicy) NextDelay(attempt uint32, rng *rand.Rand) time.Duration {
	if policy.BaseDelay <= 0 || policy.MaximumDelay <= 0 {
		return 0
	}

	delay := policy.BaseDelay
	for range attempt {
		if delay >= policy.MaximumDelay/2 {
			delay = policy.MaximumDelay
			break
		}
		delay *= 2
	}
	if delay > policy.MaximumDelay {
		delay = policy.MaximumDelay
	}

	minimumDelay := max(time.Duration(0), policy.MinimumDelay)
	if rng == nil {
		return max(minimumDelay, delay)
	}

	switch policy.Jitter {
	case JitterFull:
		return max(
			minimumDelay,
			time.Duration(rng.Int64N(delay.Nanoseconds())),
		)

	case JitterEqual:
		halfDelay := delay / 2
		if halfDelay <= 0 {
			return max(minimumDelay, delay)
		}

		return max(
			minimumDelay,
			halfDelay+time.Duration(rng.Int64N(halfDelay.Nanoseconds())),
		)
	}

	return max(minimumDelay, delay)
}
