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

var imageListBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

func assertRowsIngestIDs[S model.ImageListCursorScalar](
	t *testing.T,
	rows []persist.ImageWithCursorKey[S],
	rowsName string,
	want ...model.IngestIDType,
) {
	t.Helper()

	assert.LenIs(t, rowsName, rows, len(want))
	for i, row := range rows {
		if got, wantID := row.Image.IngestID, want[i]; got != wantID {
			t.Fatalf("%s[%d].Image.IngestID = %d, want %d", rowsName, i, got, wantID)
		}
	}
}

func assertRowsExcludeIngestID[S model.ImageListCursorScalar](
	t *testing.T,
	rows []persist.ImageWithCursorKey[S],
	disallowedID model.IngestIDType,
) {
	t.Helper()

	for i, row := range rows {
		if row.Image.IngestID == disallowedID {
			t.Fatalf("rows[%d].Image.IngestID = %d, must not contain %d", i, row.Image.IngestID, disallowedID)
		}
	}
}

func assertLastRowTimeCursorEquals(
	t *testing.T,
	rows []persist.ImageWithCursorKey[time.Time],
	want time.Time,
) time.Time {
	t.Helper()
	assert.NotEmpty(t, "rows", rows)

	gotCursor := rows[len(rows)-1].PrimaryKey
	assert.EqualFn(t, "nextCursor", gotCursor, want)
	return gotCursor
}

func assertLastRowScoreCursorEquals(
	t *testing.T,
	rows []persist.ImageWithCursorKey[model.ScoreType],
	want model.ScoreType,
) model.ScoreType {
	t.Helper()
	assert.NotEmpty(t, "rows", rows)

	gotCursor := rows[len(rows)-1].PrimaryKey
	assert.Equal(t, "nextCursor", gotCursor, want)
	return gotCursor
}

func mustAddIngestAndImage(
	t testing.TB,
	ops c.TxSession,
	ingestedAt time.Time,
	capturedAt time.Time,
) model.Ingest {
	t.Helper()
	ingest := ops.MustAddIngest(
		t,
		mb.
			Ingest().
			Ingested(ingestedAt).
			Captured(capturedAt).
			Build(),
	)
	ops.MustAddImage(t, mb.Image(ingest.ID).Ingested(ingestedAt).Build())
	return ingest
}

func TestImageListRepositoryListLatest(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			oldest := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC.Add(-48*time.Hour), imageListBaseTimeUTC.Add(-48*time.Hour))
			oldestTie := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC.Add(-48*time.Hour), imageListBaseTimeUTC.Add(-48*time.Hour))
			middle := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC.Add(-24*time.Hour), imageListBaseTimeUTC.Add(-24*time.Hour))
			latest := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC, imageListBaseTimeUTC)

			rows, err := ops.ImageList().ListLatest(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListLatest() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middle.IngestedAt)
			nextRows, err := ops.ImageList().ListLatest(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListLatest() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", oldestTie.ID, oldest.ID)
		})
	})
}

func TestImageListRepositoryListChronological(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			oldest := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC.Add(-2*time.Hour), imageListBaseTimeUTC.Add(-2*time.Hour))
			middleTie := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC, imageListBaseTimeUTC.Add(-1*time.Hour))
			middle := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC.Add(4*time.Hour), imageListBaseTimeUTC.Add(-1*time.Hour))
			latest := mustAddIngestAndImage(t, ops, imageListBaseTimeUTC.Add(2*24*time.Hour), imageListBaseTimeUTC.Add(2*24*time.Hour))

			rows, err := ops.ImageList().ListChronological(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListChronological() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middle.CapturedAt)
			nextRows, err := ops.ImageList().ListChronological(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListChronological() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middleTie.ID, oldest.ID)
		})
	})
}

func TestImageListRepositoryListRecently(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			oldest := mustAddIngestAndImage(t, ops, baseTime.Add(-2*24*time.Hour), baseTime.Add(-2*24*time.Hour))
			middleTie := mustAddIngestAndImage(t, ops, baseTime, baseTime)
			latest := mustAddIngestAndImage(t, ops, baseTime.Add(2*24*time.Hour), baseTime.Add(2*24*time.Hour))
			middle := mustAddIngestAndImage(t, ops, baseTime.Add(4*24*time.Hour), baseTime.Add(4*24*time.Hour))
			withoutLastViewedAt := mustAddIngestAndImage(t, ops, baseTime.Add(5*24*time.Hour), baseTime.Add(5*24*time.Hour))

			ops.MustAddStats(t, mb.Stats(oldest.ID).ViewedOffset(1, -2*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(middleTie.ID).ViewedOffset(1, -1*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(latest.ID).ViewedOffset(1, 0).Build())
			middleStats := ops.MustAddStats(t, mb.Stats(middle.ID).ViewedOffset(1, -1*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(withoutLastViewedAt.ID).Build())

			rows, err := ops.ImageList().ListRecently(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListRecently() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
			assertRowsExcludeIngestID(t, rows, withoutLastViewedAt.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middleStats.LastViewedAt.MustGet())
			nextRows, err := ops.ImageList().ListRecently(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListRecently() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middleTie.ID, oldest.ID)
			assertRowsExcludeIngestID(t, nextRows, withoutLastViewedAt.ID)
		})
	})
}

func TestImageListRepositoryListFirstLove(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			oldest := mustAddIngestAndImage(t, ops, baseTime.Add(-2*24*time.Hour), baseTime.Add(-2*24*time.Hour))
			middleTie := mustAddIngestAndImage(t, ops, baseTime, baseTime)
			latest := mustAddIngestAndImage(t, ops, baseTime.Add(2*24*time.Hour), baseTime.Add(2*24*time.Hour))
			middle := mustAddIngestAndImage(t, ops, baseTime.Add(4*24*time.Hour), baseTime.Add(4*24*time.Hour))
			withoutFirstLovedAt := mustAddIngestAndImage(t, ops, baseTime.Add(5*24*time.Hour), baseTime.Add(5*24*time.Hour))

			ops.MustAddStats(t, mb.Stats(oldest.ID).LovedOffset(-2*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(middleTie.ID).LovedOffset(-1*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(latest.ID).LovedOffset(0).Build())
			middleStats := ops.MustAddStats(t, mb.Stats(middle.ID).LovedOffset(-1*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(withoutFirstLovedAt.ID).Build())

			rows, err := ops.ImageList().ListFirstLove(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListFirstLove() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
			assertRowsExcludeIngestID(t, rows, withoutFirstLovedAt.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middleStats.FirstLovedAt.MustGet())
			nextRows, err := ops.ImageList().ListFirstLove(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListFirstLove() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middleTie.ID, oldest.ID)
			assertRowsExcludeIngestID(t, nextRows, withoutFirstLovedAt.ID)
		})
	})
}

func TestImageListRepositoryListHallOfFame(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			oldest := mustAddIngestAndImage(t, ops, baseTime.Add(-2*24*time.Hour), baseTime.Add(-2*24*time.Hour))
			middleTie := mustAddIngestAndImage(t, ops, baseTime, baseTime)
			latest := mustAddIngestAndImage(t, ops, baseTime.Add(2*24*time.Hour), baseTime.Add(2*24*time.Hour))
			middle := mustAddIngestAndImage(t, ops, baseTime.Add(4*24*time.Hour), baseTime.Add(4*24*time.Hour))
			withoutHallOfFameAt := mustAddIngestAndImage(t, ops, baseTime.Add(5*24*time.Hour), baseTime.Add(5*24*time.Hour))

			ops.MustAddStats(t, mb.Stats(oldest.ID).HallOfFameOffset(-2*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(middleTie.ID).HallOfFameOffset(-1*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(latest.ID).HallOfFameOffset(0).Build())
			middleStats := ops.MustAddStats(t, mb.Stats(middle.ID).HallOfFameOffset(-1*time.Hour).Build())
			ops.MustAddStats(t, mb.Stats(withoutHallOfFameAt.ID).Build())

			rows, err := ops.ImageList().ListHallOfFame(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListHallOfFame() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
			assertRowsExcludeIngestID(t, rows, withoutHallOfFameAt.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middleStats.HallOfFameAt.MustGet())
			nextRows, err := ops.ImageList().ListHallOfFame(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListHallOfFame() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middleTie.ID, oldest.ID)
			assertRowsExcludeIngestID(t, nextRows, withoutHallOfFameAt.ID)
		})
	})
}

func TestImageListRepositoryListEngaged(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	evaluatedAt := baseTime.Add(2 * time.Hour)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			high := mustAddIngestAndImage(t, ops, baseTime.Add(-5*24*time.Hour), baseTime.Add(-5*24*time.Hour))
			middleTie := mustAddIngestAndImage(t, ops, baseTime.Add(-4*24*time.Hour), baseTime.Add(-4*24*time.Hour))
			lowest := mustAddIngestAndImage(t, ops, baseTime.Add(-3*24*time.Hour), baseTime.Add(-3*24*time.Hour))
			hiddenHighest := mustAddIngestAndImage(t, ops, baseTime.Add(-2*24*time.Hour), baseTime.Add(-2*24*time.Hour))
			middle := mustAddIngestAndImage(t, ops, baseTime.Add(-1*24*time.Hour), baseTime.Add(-1*24*time.Hour))
			low := mustAddIngestAndImage(t, ops, baseTime.Add(0*24*time.Hour), baseTime.Add(0*24*time.Hour))

			ops.MustAddStats(t, mb.Stats(high.ID).Score(180).EvaluateScore(evaluatedAt).Build())
			ops.MustAddStats(t, mb.Stats(middleTie.ID).Score(165).EvaluateScore(evaluatedAt).Build())
			ops.MustAddStats(t, mb.Stats(lowest.ID).Score(150).EvaluateScore(evaluatedAt).Build())
			ops.MustAddStats(t, mb.Stats(hiddenHighest.ID).
				Score(190).
				EvaluateScore(evaluatedAt).
				HallOfFameOffset(0).
				Build())
			middleStats := ops.MustAddStats(t, mb.Stats(middle.ID).Score(165).EvaluateScore(evaluatedAt).Build())
			ops.MustAddStats(t, mb.Stats(low.ID).Score(160).EvaluateScore(evaluatedAt).Build())

			rows, err := ops.ImageList().ListEngaged(t.Context(), persist.EngagedImageListSpec{
				ImageListSpec: persist.ImageListSpec[model.ScoreType]{
					MaxCount: 2,
				},
				ScoreThreshold: 160,
			})
			assert.NilError(t, "ListEngaged() error", err)
			assertRowsIngestIDs(t, rows, "rows", high.ID, middle.ID)
			assertRowsExcludeIngestID(t, rows, hiddenHighest.ID)

			nextCursor := assertLastRowScoreCursorEquals(t, rows, middleStats.Score)
			nextRows, err := ops.ImageList().ListEngaged(t.Context(), persist.EngagedImageListSpec{
				ImageListSpec: persist.ImageListSpec[model.ScoreType]{
					CursorKey: mo.Some(model.ImageListCursorKey[model.ScoreType]{
						Primary:   nextCursor,
						Secondary: rows[len(rows)-1].Image.IngestID,
					}),
					MaxCount: 2,
				},
				ScoreThreshold: 160,
			})
			assert.NilError(t, "ListEngaged() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middleTie.ID, low.ID)
			assertRowsExcludeIngestID(t, nextRows, hiddenHighest.ID)
		})
	})
}
