package persistence

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type JobSuite SuiteBase[persist.JobRepository]

func (ste JobSuite) RunTestMarks(t *testing.T) {
	t.Helper()

	ste.Operations.MustTruncateJobs(t)

	startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
	err := ste.Repository.MarkStarted(ste.Context, "test_job", startedAt)
	assert.NilError(t, "MarkStarted() error", err)

	finishedAt := startedAt.Add(2 * time.Second)
	err = ste.Repository.MarkFinished(ste.Context, "test_job", finishedAt)
	assert.NilError(t, "MarkFinished() error", err)
}

func (ste JobSuite) RunTestMarkStartedReturnsConflict(t *testing.T) {
	t.Helper()

	ste.Operations.MustTruncateJobs(t)

	startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
	err := ste.Repository.MarkStarted(ste.Context, "test_job", startedAt)
	assert.NilError(t, "MarkStarted() error", err)

	err = ste.Repository.MarkStarted(ste.Context, "test_job", startedAt)
	assert.ErrorIs(t, "MarkStarted() error", err, persist.ErrConflict)
}

func (ste JobSuite) RunTestMarkFinishedReturnsConflict(t *testing.T) {
	t.Helper()

	ste.Operations.MustTruncateJobs(t)

	finishedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC).Add(2 * time.Second)
	err := ste.Repository.MarkFinished(ste.Context, "test_job", finishedAt)
	assert.ErrorIs(t, "MarkFinished() error", err, persist.ErrConflict)
}
