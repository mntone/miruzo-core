package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

// --- Schema ---

func TestJobRepositoryMarkStartedAcceptsValidName(t *testing.T) {
	tests := []struct {
		name    string
		jobName string
	}{
		{
			name:    "AcceptsMinLength",
			jobName: "job_name",
		},
		{
			name:    "AcceptsMaxLength",
			jobName: "job_name_1234567",
		},
		{
			name:    "AcceptsDigitsAndUnderscore",
			jobName: "job_2026_04",
		},
	}

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		t.Parallel()
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				t.Run(tt.name, func(t *testing.T) {
					startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
					err := ops.Job().MarkStarted(t.Context(), tt.jobName, startedAt)
					assert.NilError(t, "MarkStarted() error", err)
				})
			})
		}
	})
}

func TestJobRepositoryMarkStartedReturnsCheckViolation(t *testing.T) {
	tests := []struct {
		name    string
		jobName string
	}{
		{
			name:    "RejectsTooShortName",
			jobName: "short",
		},
		{
			name:    "RejectsTooLongName",
			jobName: "too_long_job_name_1",
		},
		{
			name:    "RejectsUppercaseChars",
			jobName: "Job_Name",
		},
		{
			name:    "RejectsDashChars",
			jobName: "job-name",
		},
		{
			name:    "RejectsDotChars",
			jobName: "job.name",
		},
	}

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		t.Parallel()
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				t.Run(tt.name, func(t *testing.T) {
					startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
					err := ops.Job().MarkStarted(t.Context(), tt.jobName, startedAt)
					assert.ErrorIs(t, "MarkStarted() error", err, persist.ErrCheckViolation)
				})
			})
		}
	})
}

// --- Mark* ---

func TestJobRepositoryMarks(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		t.Parallel()
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
			err := ops.Job().MarkStarted(t.Context(), "test_job", startedAt)
			assert.NilError(t, "MarkStarted() error", err)

			finishedAt := startedAt.Add(2 * time.Second)
			err = ops.Job().MarkFinished(t.Context(), "test_job", finishedAt)
			assert.NilError(t, "MarkFinished() error", err)

			// In MySQL, this path can report affected_rows=2.
			startedAt2 := time.Date(2026, 1, 11, 5, 5, 0, 0, time.UTC)
			err = ops.Job().MarkStarted(t.Context(), "test_job", startedAt2)
			assert.NilError(t, "MarkStarted() error", err)
		})
	})
}

func TestJobRepositoryMarkStartedReturnsConflict(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		t.Parallel()
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
			err := ops.Job().MarkStarted(t.Context(), "strt_already", startedAt)
			assert.NilError(t, "err", err)

			err = ops.Job().MarkStarted(t.Context(), "strt_already", startedAt)
			assert.ErrorIs(t, "err", err, persist.ErrConflict)

			errString := err.Error()
			assert.Contains(t, "err", errString, "operation=MarkStarted")
			assert.Contains(t, "err", errString, "affected_rows=0")
			assert.Contains(t, "err", errString, "name=strt_already")
		})
	})
}

func TestJobRepositoryMarkFinishedReturnsConflictWhenNoRows(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		t.Parallel()
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			finishedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC).Add(2 * time.Second)
			err := ops.Job().MarkFinished(t.Context(), "fin_norows", finishedAt)
			assert.ErrorIs(t, "err", err, persist.ErrConflict)

			errString := err.Error()
			assert.Contains(t, "err", errString, "operation=MarkFinished")
			assert.Contains(t, "err", errString, "affected_rows=0")
			assert.Contains(t, "err", errString, "name=fin_norows")
		})
	})
}

func TestJobRepositoryMarkFinishedReturnsConflictWhenAlreadyFinished(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		t.Parallel()
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			at := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
			ctx := t.Context()

			err := ops.Job().MarkStarted(ctx, "fin_already", at)
			assert.NilError(t, "err", err)

			err = ops.Job().MarkFinished(ctx, "fin_already", at.Add(1*time.Second))
			assert.NilError(t, "err", err)

			err = ops.Job().MarkFinished(ctx, "fin_already", at.Add(2*time.Second))
			assert.ErrorIs(t, "err", err, persist.ErrConflict)

			errString := err.Error()
			assert.Contains(t, "err", errString, "operation=MarkFinished")
			assert.Contains(t, "err", errString, "affected_rows=0")
			assert.Contains(t, "err", errString, "name=fin_already")
		})
	})
}
