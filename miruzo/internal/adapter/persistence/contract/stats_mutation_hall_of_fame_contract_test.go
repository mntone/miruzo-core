package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

const hallOfFameScoreThreshold = model.ScoreType(180)

// --- ApplyHallOfFameGranted ---

func TestStatsRepositoryApplyHallOfFameGrantedUpdates(t *testing.T) {
	ingest := mb.Ingest().Build()
	stats := mb.
		Stats(ingest.ID).
		Score(hallOfFameScoreThreshold).
		Build()
	hallOfFameAt := mb.GetDefaultBaseTime().Add(20 * time.Minute)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			ops.MustAddStats(t, stats)

			err := ops.Stats().ApplyHallOfFameGranted(
				t.Context(),
				ingest.ID,
				hallOfFameAt,
				hallOfFameScoreThreshold,
			)
			assert.NilError(t, "ApplyHallOfFameGranted() error", err)

			dbstats := ops.MustGetStats(t, ingest.ID)
			assert.IsPresent(t, "stats.HallOfFameAt", dbstats.HallOfFameAt)
			assert.Equal(t, "stats.HallOfFameAt", dbstats.HallOfFameAt.MustGet(), hallOfFameAt)
		})
	})
}

func TestStatsRepositoryApplyHallOfFameGrantedReturnsConflict(t *testing.T) {
	ingest := mb.Ingest().Build()
	hallOfFameAt := mb.GetDefaultBaseTime().Add(20 * time.Minute)

	tests := []struct {
		name            string
		score           model.ScoreType
		grantedAtOffset mo.Option[time.Duration]
	}{
		{
			name:  "ScoreBelowThreshold",
			score: hallOfFameScoreThreshold - 10,
		},
		{
			name:            "AlreadyGranted",
			score:           hallOfFameScoreThreshold,
			grantedAtOffset: mo.Some(2 * time.Hour),
		},
	}

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
					ops.MustAddIngest(t, ingest)
					ops.MustAddStats(t, mb.
						Stats(ingest.ID).
						Score(tt.score).
						HallOfFameOffset(tt.grantedAtOffset).
						Build())

					err := ops.Stats().ApplyHallOfFameGranted(
						t.Context(),
						ingest.ID,
						hallOfFameAt,
						hallOfFameScoreThreshold,
					)
					assert.ErrorIs(t, "ApplyHallOfFameGranted() error", err, persist.ErrConflict)
				})
			})
		}
	})
}

func TestStatsRepositoryApplyHallOfFameGrantedReturnsConflictWithoutStats(t *testing.T) {
	ingest := mb.Ingest().Build()
	hallOfFameAt := mb.GetDefaultBaseTime().Add(20 * time.Minute)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)

			err := ops.Stats().ApplyHallOfFameGranted(
				t.Context(),
				ingest.ID,
				hallOfFameAt,
				hallOfFameScoreThreshold,
			)
			assert.ErrorIs(t, "ApplyHallOfFameGranted() error", err, persist.ErrConflict)
		})
	})
}

// --- ApplyHallOfFameRevoked ---

func TestStatsRepositoryApplyHallOfFameRevokedUpdates(t *testing.T) {
	ingest := mb.Ingest().Build()
	stats := mb.
		Stats(ingest.ID).
		Score(hallOfFameScoreThreshold - 10).
		HallOfFameOffset(20 * time.Minute).
		Build()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			ops.MustAddStats(t, stats)

			err := ops.Stats().ApplyHallOfFameRevoked(t.Context(), ingest.ID)
			assert.NilError(t, "ApplyHallOfFameRevoked() error", err)

			dbstats := ops.MustGetStats(t, ingest.ID)
			assert.IsAbsent(t, "stats.HallOfFameAt", dbstats.HallOfFameAt)
		})
	})
}

func TestStatsRepositoryApplyHallOfFameRevokedReturnsConflict(t *testing.T) {
	ingest := mb.Ingest().Build()
	stats := mb.
		Stats(ingest.ID).
		Score(hallOfFameScoreThreshold - 10).
		Build()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			ops.MustAddStats(t, stats)

			err := ops.Stats().ApplyHallOfFameRevoked(t.Context(), ingest.ID)
			assert.ErrorIs(t, "ApplyHallOfFameRevoked() error", err, persist.ErrConflict)
		})
	})
}

func TestStatsRepositoryApplyHallOfFameRevokedReturnsConflictWithoutStats(t *testing.T) {
	ingest := mb.Ingest().Build()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			err := ops.Stats().ApplyHallOfFameRevoked(t.Context(), ingest.ID)
			assert.ErrorIs(t, "ApplyHallOfFameRevoked() error", err, persist.ErrConflict)
		})
	})
}
