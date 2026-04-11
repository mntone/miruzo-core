package contract_test

import (
	"slices"
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

func assertRowsIngestIDs[S model.ImageListCursorScalar](
	t *testing.T,
	rows []persist.ImageWithCursorKey[S],
	rowsName string,
	wantIDs ...model.IngestIDType,
) {
	t.Helper()

	assert.LenIs(t, rowsName, rows, len(wantIDs))

	gotIDs := make([]model.IngestIDType, len(rows))
	for i, row := range rows {
		gotIDs[i] = row.Image.IngestID
	}
	if !slices.Equal(gotIDs, wantIDs) {
		t.Fatalf("%s = %v, want %v", rowsName, gotIDs, wantIDs)
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

type imageListSeed struct {
	Ingested       time.Duration
	Score          mo.Option[model.ScoreType]
	ScoreEvaluated mo.Option[time.Duration]
	Captured       mo.Option[time.Duration]
	HallOfFame     mo.Option[time.Duration]
	Loved          mo.Option[time.Duration]
	Viewed         mo.Option[time.Duration]
	ViewCount      int64
}

type imageListFixture struct {
	model.Ingest
	Image persist.Image
	Stats model.Stats
}

func buildImageListFixtures(t testing.TB, seeds []imageListSeed) []imageListFixture {
	t.Helper()

	fixtures := make([]imageListFixture, len(seeds))
	for i, seed := range seeds {
		ingest := mb.
			Ingest().
			IngestedOffset(seed.Ingested).
			CapturedOffset(seed.Captured).
			Build()
		fixtures[i].Ingest = ingest
		fixtures[i].Image = mb.
			Image(ingest.ID).
			Ingested(ingest.IngestedAt).
			Build()
		fixtures[i].Stats = mb.
			Stats(ingest.ID).
			ScoreOption(seed.Score).
			EvaluateScoreOffset(seed.ScoreEvaluated).
			HallOfFameOffset(seed.HallOfFame).
			LovedOffset(seed.Loved).
			ViewedOffset(seed.ViewCount, seed.Viewed).
			Build()
	}
	return fixtures
}

func mustInsertImageListFixtures(t testing.TB, ops c.TxSession, fixtures []imageListFixture) {
	t.Helper()

	for _, fixture := range fixtures {
		ops.MustAddIngest(t, fixture.Ingest)
		ops.MustAddImage(t, fixture.Image)
		ops.MustAddStats(t, fixture.Stats)
	}
}

func TestImageListRepositoryListLatest(t *testing.T) {
	fixtures := buildImageListFixtures(t, []imageListSeed{
		{
			Ingested: -2 * time.Hour,
		},
		{
			Ingested: -2 * time.Hour,
		},
		{
			Ingested: -1 * time.Hour,
		},
		{
			Ingested: 0,
		},
	})

	oldest := fixtures[0]    // ID: k
	oldestTie := fixtures[1] // ID: k + 1
	middle := fixtures[2]    // ID: k + 2
	latest := fixtures[3]    // ID: k + 3

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			mustInsertImageListFixtures(t, ops, fixtures)

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
	fixtures := buildImageListFixtures(t, []imageListSeed{
		{
			Ingested: 1 * time.Hour,
			Captured: mo.Some(-2 * time.Hour),
		},
		{
			Ingested: 2 * time.Hour,
			Captured: mo.Some(-1 * time.Hour),
		},
		{
			Ingested: 3 * time.Hour,
			Captured: mo.Some(0 * time.Hour),
		},
		{
			Ingested: 4 * time.Hour,
			Captured: mo.Some(-1 * time.Hour),
		},
	})

	oldest := fixtures[0]    // ID: k
	middleTie := fixtures[1] // ID: k + 1
	latest := fixtures[2]    // ID: k + 2
	middle := fixtures[3]    // ID: k + 3

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			mustInsertImageListFixtures(t, ops, fixtures)

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
	fixtures := buildImageListFixtures(t, []imageListSeed{
		{
			Ingested:  1 * time.Hour,
			Viewed:    mo.Some(4 * time.Hour),
			ViewCount: 15,
		},
		{
			Ingested:  2 * time.Hour,
			Viewed:    mo.Some(5 * time.Hour),
			ViewCount: 3,
		},
		{
			Ingested:  3 * time.Hour,
			Viewed:    mo.Some(6 * time.Hour),
			ViewCount: 32,
		},
		{
			Ingested:  4 * time.Hour,
			Viewed:    mo.Some(5 * time.Hour),
			ViewCount: 42,
		},
		{
			Ingested: 5 * time.Hour,
		},
	})

	oldest := fixtures[0]              // ID: k
	middle := fixtures[1]              // ID: k + 1
	latest := fixtures[2]              // ID: k + 2
	middleTie := fixtures[3]           // ID: k + 3
	withoutLastViewedAt := fixtures[4] // ID: k + 4

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			mustInsertImageListFixtures(t, ops, fixtures)

			rows, err := ops.ImageList().ListRecently(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListRecently() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middleTie.ID)
			assertRowsExcludeIngestID(t, rows, withoutLastViewedAt.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middleTie.Stats.LastViewedAt.MustGet())
			nextRows, err := ops.ImageList().ListRecently(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListRecently() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middle.ID, oldest.ID)
			assertRowsExcludeIngestID(t, nextRows, withoutLastViewedAt.ID)
		})
	})
}

func TestImageListRepositoryListFirstLove(t *testing.T) {
	fixtures := buildImageListFixtures(t, []imageListSeed{
		{
			Ingested: 1 * time.Hour,
			Loved:    mo.Some(4 * time.Hour),
		},
		{
			Ingested: 2 * time.Hour,
			Loved:    mo.Some(5 * time.Hour),
		},
		{
			Ingested: 3 * time.Hour,
			Loved:    mo.Some(6 * time.Hour),
		},
		{
			Ingested: 4 * time.Hour,
			Loved:    mo.Some(5 * time.Hour),
		},
		{
			Ingested: 5 * time.Hour,
		},
	})

	oldest := fixtures[0]              // ID: k
	middle := fixtures[1]              // ID: k + 1
	latest := fixtures[2]              // ID: k + 2
	middleTie := fixtures[3]           // ID: k + 3
	withoutFirstLovedAt := fixtures[4] // ID: k + 4

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			mustInsertImageListFixtures(t, ops, fixtures)

			rows, err := ops.ImageList().ListFirstLove(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListFirstLove() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middleTie.ID)
			assertRowsExcludeIngestID(t, rows, withoutFirstLovedAt.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middleTie.Stats.FirstLovedAt.MustGet())
			nextRows, err := ops.ImageList().ListFirstLove(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListFirstLove() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middle.ID, oldest.ID)
			assertRowsExcludeIngestID(t, nextRows, withoutFirstLovedAt.ID)
		})
	})
}

func TestImageListRepositoryListHallOfFame(t *testing.T) {
	fixtures := buildImageListFixtures(t, []imageListSeed{
		{
			Ingested:   1 * time.Hour,
			HallOfFame: mo.Some(4 * time.Hour),
		},
		{
			Ingested:   2 * time.Hour,
			HallOfFame: mo.Some(5 * time.Hour),
		},
		{
			Ingested:   3 * time.Hour,
			HallOfFame: mo.Some(6 * time.Hour),
		},
		{
			Ingested:   4 * time.Hour,
			HallOfFame: mo.Some(5 * time.Hour),
		},
		{
			Ingested: 5 * time.Hour,
		},
	})

	oldest := fixtures[0]              // ID: k
	middle := fixtures[1]              // ID: k + 1
	latest := fixtures[2]              // ID: k + 2
	middleTie := fixtures[3]           // ID: k + 3
	withoutHallOfFameAt := fixtures[4] // ID: k + 4

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			mustInsertImageListFixtures(t, ops, fixtures)

			rows, err := ops.ImageList().ListHallOfFame(t.Context(), persist.ImageListSpec[time.Time]{
				MaxCount: 2,
			})
			assert.NilError(t, "ListHallOfFame() error", err)
			assertRowsIngestIDs(t, rows, "rows", latest.ID, middleTie.ID)
			assertRowsExcludeIngestID(t, rows, withoutHallOfFameAt.ID)

			nextCursorAt := assertLastRowTimeCursorEquals(t, rows, middleTie.Stats.HallOfFameAt.MustGet())
			nextRows, err := ops.ImageList().ListHallOfFame(t.Context(), persist.ImageListSpec[time.Time]{
				CursorKey: mo.Some(model.ImageListCursorKey[time.Time]{
					Primary:   nextCursorAt,
					Secondary: rows[len(rows)-1].Image.IngestID,
				}),
				MaxCount: 2,
			})
			assert.NilError(t, "ListHallOfFame() error", err)
			assertRowsIngestIDs(t, nextRows, "nextRows", middle.ID, oldest.ID)
			assertRowsExcludeIngestID(t, nextRows, withoutHallOfFameAt.ID)
		})
	})
}

func TestImageListRepositoryListEngaged(t *testing.T) {
	evaluatedOffset := mo.Some(12 * time.Hour)
	fixtures := buildImageListFixtures(t, []imageListSeed{
		{
			Ingested:       1 * time.Hour,
			Score:          mo.Some(model.ScoreType(180)),
			ScoreEvaluated: evaluatedOffset,
		},
		{
			Ingested:       2 * time.Hour,
			Score:          mo.Some(model.ScoreType(165)),
			ScoreEvaluated: evaluatedOffset,
		},
		{
			Ingested:       3 * time.Hour,
			Score:          mo.Some(model.ScoreType(150)), // below threshold
			ScoreEvaluated: evaluatedOffset,
		},
		{
			Ingested:       4 * time.Hour,
			Score:          mo.Some(model.ScoreType(190)),
			HallOfFame:     mo.Some(7 * time.Hour), // excluded by hall_of_fame_at
			ScoreEvaluated: evaluatedOffset,
		},
		{
			Ingested:       5 * time.Hour,
			Score:          mo.Some(model.ScoreType(165)),
			ScoreEvaluated: evaluatedOffset,
		},
		{
			Ingested:       6 * time.Hour,
			Score:          mo.Some(model.ScoreType(160)),
			ScoreEvaluated: evaluatedOffset,
		},
	})

	high := fixtures[0]          // ID: k
	middle := fixtures[1]        // ID: k + 1
	hiddenLowest := fixtures[2]  // ID: k + 2
	hiddenHighest := fixtures[3] // ID: k + 3
	middleTie := fixtures[4]     // ID: k + 4
	low := fixtures[5]           // ID: k + 5

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			mustInsertImageListFixtures(t, ops, fixtures)

			rows, err := ops.ImageList().ListEngaged(t.Context(), persist.EngagedImageListSpec{
				ImageListSpec: persist.ImageListSpec[model.ScoreType]{
					MaxCount: 2,
				},
				ScoreThreshold: 160,
			})
			assert.NilError(t, "ListEngaged() error", err)
			assertRowsIngestIDs(t, rows, "rows", high.ID, middleTie.ID)
			assertRowsExcludeIngestID(t, rows, hiddenHighest.ID)
			assertRowsExcludeIngestID(t, rows, hiddenLowest.ID)

			nextCursor := assertLastRowScoreCursorEquals(t, rows, middleTie.Stats.Score)
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
			assertRowsIngestIDs(t, nextRows, "nextRows", middle.ID, low.ID)
			assertRowsExcludeIngestID(t, nextRows, hiddenHighest.ID)
			assertRowsExcludeIngestID(t, nextRows, hiddenLowest.ID)
		})
	})
}
