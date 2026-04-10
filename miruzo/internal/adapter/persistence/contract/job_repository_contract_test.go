package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestJobRepositoryMarks(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
			err := ops.Job().MarkStarted(t.Context(), "test_job", startedAt)
			assert.NilError(t, "MarkStarted() error", err)

			finishedAt := startedAt.Add(2 * time.Second)
			err = ops.Job().MarkFinished(t.Context(), "test_job", finishedAt)
			assert.NilError(t, "MarkFinished() error", err)
		})
	})
}

func TestJobRepositoryMarkStartedReturnsConflict(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			startedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC)
			err := ops.Job().MarkStarted(t.Context(), "test_job", startedAt)
			assert.NilError(t, "MarkStarted() error", err)

			err = ops.Job().MarkStarted(t.Context(), "test_job", startedAt)
			assert.ErrorIs(t, "MarkStarted() error", err, persist.ErrConflict)
		})
	})
}

func TestJobRepositoryMarkFinishedReturnsConflict(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			finishedAt := time.Date(2026, 1, 10, 5, 5, 0, 0, time.UTC).Add(2 * time.Second)
			err := ops.Job().MarkFinished(t.Context(), "test_job", finishedAt)
			assert.ErrorIs(t, "MarkFinished() error", err, persist.ErrConflict)
		})
	})
}
