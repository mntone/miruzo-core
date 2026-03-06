package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

func TestRetryRetriesRecoverableErrors(t *testing.T) {
	originalSleepWithContext := sleepWithContext
	sleepCalls := 0
	sleepWithContext = func(_ context.Context, _ time.Duration) error {
		sleepCalls++
		return nil
	}
	t.Cleanup(func() {
		sleepWithContext = originalSleepWithContext
	})

	callCount := 0
	result, err := Retry(
		context.Background(),
		backoff.FixedPolicy{
			MaxAttempts: 3,
			Delays:      []time.Duration{time.Millisecond, 2 * time.Millisecond},
		},
		func(_ context.Context) (string, error) {
			callCount++
			if callCount < 3 {
				return "", fmt.Errorf("retryable: %w", persist.ErrRecoverableUnavailable)
			}

			return "ok", nil
		},
	)

	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	if result != "ok" {
		t.Fatalf("expected result=ok, got %q", result)
	}
	if callCount != 3 {
		t.Fatalf("expected 3 attempts, got %d", callCount)
	}
	if sleepCalls != 2 {
		t.Fatalf("expected 2 sleep calls, got %d", sleepCalls)
	}
}

func TestRetryDoesNotRetryUnrecoverableError(t *testing.T) {
	originalSleepWithContext := sleepWithContext
	sleepCalls := 0
	sleepWithContext = func(_ context.Context, _ time.Duration) error {
		sleepCalls++
		return nil
	}
	t.Cleanup(func() {
		sleepWithContext = originalSleepWithContext
	})

	sourceError := errors.New("unknown")
	callCount := 0
	_, err := Retry(
		context.Background(),
		backoff.FixedPolicy{
			MaxAttempts: 5,
			Delays:      []time.Duration{time.Millisecond},
		},
		func(_ context.Context) (string, error) {
			callCount++
			return "", sourceError
		},
	)

	if !errors.Is(err, sourceError) {
		t.Fatalf("expected source error, got %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 attempt, got %d", callCount)
	}
	if sleepCalls != 0 {
		t.Fatalf("expected 0 sleep calls, got %d", sleepCalls)
	}
}

func TestRetryStopsAtMaxAttempts(t *testing.T) {
	originalSleepWithContext := sleepWithContext
	sleepCalls := 0
	sleepWithContext = func(_ context.Context, _ time.Duration) error {
		sleepCalls++
		return nil
	}
	t.Cleanup(func() {
		sleepWithContext = originalSleepWithContext
	})

	callCount := 0
	_, err := Retry(
		context.Background(),
		backoff.FixedPolicy{
			MaxAttempts: 2,
			Delays:      []time.Duration{time.Millisecond},
		},
		func(_ context.Context) (string, error) {
			callCount++
			return "", fmt.Errorf("retryable: %w", persist.ErrRecoverableConflict)
		},
	)

	if !errors.Is(err, persist.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if errors.Is(err, persist.ErrRecoverableConflict) {
		t.Fatalf("did not expect ErrRecoverableConflict, got %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 attempts, got %d", callCount)
	}
	if sleepCalls != 1 {
		t.Fatalf("expected 1 sleep call, got %d", sleepCalls)
	}
}

func TestRetryReturnsContextErrorWhenSleepCanceled(t *testing.T) {
	originalSleepWithContext := sleepWithContext
	sleepWithContext = func(_ context.Context, _ time.Duration) error {
		return context.Canceled
	}
	t.Cleanup(func() {
		sleepWithContext = originalSleepWithContext
	})

	callCount := 0
	_, err := Retry(
		context.Background(),
		backoff.FixedPolicy{
			MaxAttempts: 3,
			Delays:      []time.Duration{time.Millisecond},
		},
		func(_ context.Context) (string, error) {
			callCount++
			return "", fmt.Errorf("retryable: %w", persist.ErrRecoverableUnavailable)
		},
	)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 attempt, got %d", callCount)
	}
}
