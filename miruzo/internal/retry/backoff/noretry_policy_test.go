package backoff_test

import (
	"math/rand/v2"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

func TestNoRetryPolicyPolicyNextDelayAlwaysReturnsZero(t *testing.T) {
	policy := backoff.NoRetryPolicy{}
	rng := rand.New(rand.NewPCG(1, 2))

	if got := policy.NextDelay(0, rng); got != 0 {
		t.Fatalf("expected zero delay, got %s", got)
	}
	if got := policy.NextDelay(10, rng); got != 0 {
		t.Fatalf("expected zero delay, got %s", got)
	}
}
