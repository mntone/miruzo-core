package backoff_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

func TestFixedPolicyGetMaxAttempts(t *testing.T) {
	policy := backoff.FixedPolicy{
		MaxAttempts: 3,
		Delays: []time.Duration{
			10 * time.Millisecond,
			20 * time.Millisecond,
		},
	}

	if got := policy.GetMaxAttempts(); got != 3 {
		t.Fatalf("expected max attempts=3, got %d", got)
	}
}

func TestFixedPolicyNextDelayReturnsConfiguredDelay(t *testing.T) {
	policy := backoff.FixedPolicy{
		MaxAttempts: 4,
		Delays: []time.Duration{
			10 * time.Millisecond,
			20 * time.Millisecond,
			30 * time.Millisecond,
		},
	}

	rng := rand.New(rand.NewPCG(1, 2))

	tests := []struct {
		attempt uint32
		want    time.Duration
	}{
		{attempt: 0, want: 10 * time.Millisecond},
		{attempt: 1, want: 20 * time.Millisecond},
		{attempt: 2, want: 30 * time.Millisecond},
		{attempt: 3, want: 0},
		{attempt: 100, want: 0},
	}

	for _, tt := range tests {
		got := policy.NextDelay(tt.attempt, rng)
		if got != tt.want {
			t.Fatalf("attempt=%d: expected %s, got %s", tt.attempt, tt.want, got)
		}
	}
}
