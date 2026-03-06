package backoff_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

func TestExponentialPolicyNextDelayWithoutJitter(t *testing.T) {
	policy := backoff.ExponentialPolicy{
		BaseDelay:    10 * time.Millisecond,
		MaximumDelay: 80 * time.Millisecond,
		Jitter:       backoff.JitterNone,
	}

	tests := []struct {
		attempt uint32
		want    time.Duration
	}{
		{attempt: 0, want: 10 * time.Millisecond},
		{attempt: 1, want: 20 * time.Millisecond},
		{attempt: 2, want: 40 * time.Millisecond},
		{attempt: 3, want: 80 * time.Millisecond},
		{attempt: 4, want: 80 * time.Millisecond},
	}

	for _, tt := range tests {
		got := policy.NextDelay(tt.attempt, nil)
		if got != tt.want {
			t.Fatalf("attempt=%d: expected %s, got %s", tt.attempt, tt.want, got)
		}
	}
}

func TestExponentialPolicyNextDelayWithFullJitterInRange(t *testing.T) {
	policy := backoff.ExponentialPolicy{
		BaseDelay:    100 * time.Millisecond,
		MinimumDelay: 10 * time.Millisecond,
		MaximumDelay: 100 * time.Millisecond,
		Jitter:       backoff.JitterFull,
	}

	rng := rand.New(rand.NewPCG(1, 2))
	got := policy.NextDelay(0, rng)

	if got < 10*time.Millisecond || got >= 100*time.Millisecond {
		t.Fatalf("expected delay in [10ms, 100ms), got %s", got)
	}
}
