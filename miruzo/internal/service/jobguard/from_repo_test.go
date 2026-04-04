package jobguard_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/jobguard"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/stub"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

const dailyDecayJobName = "daily_decay"

func TestTryAcquireReturnsTrueOnSuccess(t *testing.T) {
	now := time.Date(2026, 1, 10, 5, 0, 1, 0, time.UTC)
	repo := stub.NewStubJobRepository()
	guard := jobguard.NewWithJobRepository(repo)

	acquired, err := guard.TryAcquire(context.Background(), dailyDecayJobName, now)
	assert.NilError(t, "TryAcquire() error", err)
	assert.Equal(t, "acquired", acquired, true)

	assert.NotEmpty(t, "MarkStarted() called", repo.MarkStartedArgs)
	assert.Equal(t, "MarkStarted() name", repo.MarkStartedArgs[0].Name, dailyDecayJobName)
	assert.Equal(t, "MarkStarted() finishedAt", repo.MarkStartedArgs[0].StartedAt, now)
}

func TestTryAcquireReturnsFalseOnConflict(t *testing.T) {
	now := time.Date(2026, 1, 10, 5, 0, 1, 0, time.UTC)
	repo := stub.NewStubJobRepository()
	repo.MarkStartedError = persist.ErrConflict
	guard := jobguard.NewWithJobRepository(repo)

	acquired, err := guard.TryAcquire(context.Background(), dailyDecayJobName, now)
	assert.NilError(t, "TryAcquire() error", err)
	assert.Equal(t, "acquired", acquired, false)
}

func TestTryAcquireReturnsFalseWhenAlreadyRunning(t *testing.T) {
	now := time.Date(2026, 1, 10, 5, 0, 1, 0, time.UTC)
	repo := stub.NewStubJobRepository(model.Job{
		Name:      dailyDecayJobName,
		StartedAt: time.Date(2026, 1, 10, 5, 0, 0, 0, time.UTC),
		// finished_at is absent => running
	})
	guard := jobguard.NewWithJobRepository(repo)

	acquired, err := guard.TryAcquire(context.Background(), dailyDecayJobName, now)
	assert.NilError(t, "TryAcquire() error", err)
	assert.Equal(t, "acquired", acquired, false)
}

func TestTryAcquireReturnsErrorOnNonConflict(t *testing.T) {
	now := time.Date(2026, 1, 10, 5, 0, 1, 0, time.UTC)
	expectedErr := errors.New("db unavailable")
	repo := stub.NewStubJobRepository()
	repo.MarkStartedError = expectedErr
	guard := jobguard.NewWithJobRepository(repo)

	acquired, err := guard.TryAcquire(context.Background(), dailyDecayJobName, now)
	assert.ErrorIs(t, "TryAcquire() error", err, expectedErr)
	assert.Equal(t, "acquired", acquired, false)
}

func TestReleaseReturnsNilOnSuccess(t *testing.T) {
	now := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
	repo := stub.NewStubJobRepository(model.Job{
		Name:      dailyDecayJobName,
		StartedAt: time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC),
	})
	guard := jobguard.NewWithJobRepository(repo)

	err := guard.Release(context.Background(), dailyDecayJobName, now)
	assert.NilError(t, "Release() error", err)

	assert.NotEmpty(t, "MarkFinished() called", repo.MarkFinishedArgs)
	assert.Equal(t, "MarkFinished() name", repo.MarkFinishedArgs[0].Name, dailyDecayJobName)
	assert.Equal(t, "MarkFinished() finishedAt", repo.MarkFinishedArgs[0].FinishedAt, now)
}

func TestReleaseReturnsNilOnConflict(t *testing.T) {
	repo := stub.NewStubJobRepository()
	repo.MarkFinishedError = persist.ErrConflict
	guard := jobguard.NewWithJobRepository(repo)

	err := guard.Release(context.Background(), dailyDecayJobName, time.Date(2026, 1, 10, 5, 10, 0, 0, time.UTC))
	assert.NilError(t, "Release() error", err)
}

func TestReleaseReturnsNilWhenAlreadyFinished(t *testing.T) {
	repo := stub.NewStubJobRepository(model.Job{
		Name:       dailyDecayJobName,
		StartedAt:  time.Date(2026, 1, 10, 5, 0, 0, 0, time.UTC),
		FinishedAt: mo.Some(time.Date(2026, 1, 10, 5, 9, 0, 0, time.UTC)),
	})
	guard := jobguard.NewWithJobRepository(repo)

	err := guard.Release(context.Background(), dailyDecayJobName, time.Date(2026, 1, 10, 5, 10, 0, 0, time.UTC))
	assert.NilError(t, "Release() error", err)
}

func TestReleaseReturnsErrorOnNonConflict(t *testing.T) {
	expectedErr := errors.New("write failed")
	repo := stub.NewStubJobRepository()
	repo.MarkFinishedError = expectedErr
	guard := jobguard.NewWithJobRepository(repo)

	err := guard.Release(context.Background(), dailyDecayJobName, time.Date(2026, 1, 10, 5, 10, 0, 0, time.UTC))
	assert.ErrorIs(t, "Release() error", err, expectedErr)
}
