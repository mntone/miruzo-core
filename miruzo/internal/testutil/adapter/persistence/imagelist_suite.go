package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

type ImageListSetup struct {
	Ctx  context.Context
	Ops  Operations
	Repo persist.ImageListRepository
}

var suiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

func assertRowsIngestIDs[C persist.ImageListCursor](
	t *testing.T,
	rows []persist.ImageWithCursor[C],
	rowsName string,
	want ...persist.IngestID,
) {
	t.Helper()

	assert.LenIs(t, rowsName, rows, len(want))
	for i, row := range rows {
		if got, wantID := row.Image.IngestID, want[i]; got != wantID {
			t.Fatalf("%s[%d].Image.IngestID = %d, want %d", rowsName, i, got, wantID)
		}
	}
}

func assertLastRowTimeCursorEquals(
	t *testing.T,
	rows []persist.ImageWithCursor[time.Time],
	want time.Time,
) time.Time {
	t.Helper()
	assert.NotEmpty(t, "rows", rows)

	gotCursor := rows[len(rows)-1].Cursor
	assert.EqualFn(t, "nextCursor", gotCursor, want)
	return gotCursor
}

func assertLastRowInt16CursorEquals(
	t *testing.T,
	rows []persist.ImageWithCursor[int16],
	want int16,
) int16 {
	t.Helper()
	assert.NotEmpty(t, "rows", rows)

	gotCursor := rows[len(rows)-1].Cursor
	assert.Equal(t, "nextCursor", gotCursor, want)
	return gotCursor
}

func assertRowsExcludeIngestID[C persist.ImageListCursor](
	t *testing.T,
	rows []persist.ImageWithCursor[C],
	disallowedID persist.IngestID,
) {
	t.Helper()

	for i, row := range rows {
		if row.Image.IngestID == disallowedID {
			t.Fatalf("rows[%d].Image.IngestID = %d, must not contain %d", i, row.Image.IngestID, disallowedID)
		}
	}
}

func runListLatest(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Run("ListLatest", func(t *testing.T) {
		setup := setupFn(t)

		latest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC))
		middle := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC.Add(-1*24*time.Hour)))
		oldest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

		rows, err := setup.Repo.ListLatest(setup.Ctx, persist.ImageListSpec[time.Time]{
			Limit: 2,
		})
		assert.NilError(t, "ListLatest() error", err)
		assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)

		nextCursor := assertLastRowTimeCursorEquals(t, rows, middle.IngestedAt)
		nextRows, err := setup.Repo.ListLatest(setup.Ctx, persist.ImageListSpec[time.Time]{
			Cursor: mo.Some(nextCursor),
			Limit:  2,
		})
		assert.NilError(t, "ListLatest() error", err)
		assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
	})
}

func runListChronological(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Run("ListChronological", func(t *testing.T) {
		setup := setupFn(t)

		latest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
		middle := setup.Ops.MustAddIngestAndImage(t, NewIngestFixtureWithCapturedAt(
			4,
			suiteBaseTimeUTC.Add(4*time.Hour),
			suiteBaseTimeUTC.Add(-1*time.Hour),
		))
		oldest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*time.Hour)))

		rows, err := setup.Repo.ListChronological(setup.Ctx, persist.ImageListSpec[time.Time]{
			Limit: 2,
		})
		assert.NilError(t, "ListChronological() error", err)
		assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)

		nextCursor := assertLastRowTimeCursorEquals(t, rows, middle.CapturedAt)
		nextRows, err := setup.Repo.ListChronological(setup.Ctx, persist.ImageListSpec[time.Time]{
			Cursor: mo.Some(nextCursor),
			Limit:  2,
		})
		assert.NilError(t, "ListChronological() error", err)
		assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
	})
}

func runListRecently(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Run("ListRecently", func(t *testing.T) {
		setup := setupFn(t)

		withoutLastViewedAt := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(24*time.Hour)))
		latest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
		middle := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(4*24*time.Hour)))
		oldest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

		setup.Ops.MustAddStat(t, NewStatFixture(withoutLastViewedAt.ID))
		setup.Ops.MustAddStat(t, NewStatFixtureWithLastViewedAt(latest.ID, 1, suiteBaseTimeUTC))
		middleStat := setup.Ops.MustAddStat(t, NewStatFixtureWithLastViewedAt(middle.ID, 1, suiteBaseTimeUTC.Add(-1*time.Hour)))
		setup.Ops.MustAddStat(t, NewStatFixtureWithLastViewedAt(oldest.ID, 1, suiteBaseTimeUTC.Add(-2*time.Hour)))

		rows, err := setup.Repo.ListRecently(setup.Ctx, persist.ImageListSpec[time.Time]{
			Limit: 2,
		})
		assert.NilError(t, "ListRecently() error", err)
		assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
		assertRowsExcludeIngestID(t, rows, withoutLastViewedAt.ID)

		nextCursor := assertLastRowTimeCursorEquals(t, rows, middleStat.LastViewedAt.MustGet())
		nextRows, err := setup.Repo.ListRecently(setup.Ctx, persist.ImageListSpec[time.Time]{
			Cursor: mo.Some(nextCursor),
			Limit:  2,
		})
		assert.NilError(t, "ListRecently() error", err)
		assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
		assertRowsExcludeIngestID(t, nextRows, withoutLastViewedAt.ID)
	})
}

func runListFirstLove(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Run("ListFirstLove", func(t *testing.T) {
		setup := setupFn(t)

		withoutFirstLovedAt := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(24*time.Hour)))
		latest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
		middle := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(4*24*time.Hour)))
		oldest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

		setup.Ops.MustAddStat(t, NewStatFixture(withoutFirstLovedAt.ID))
		setup.Ops.MustAddStat(t, NewStatFixtureWithLastLovedAt(latest.ID, suiteBaseTimeUTC))
		middleStat := setup.Ops.MustAddStat(t, NewStatFixtureWithLastLovedAt(middle.ID, suiteBaseTimeUTC.Add(-1*time.Hour)))
		setup.Ops.MustAddStat(t, NewStatFixtureWithLastLovedAt(oldest.ID, suiteBaseTimeUTC.Add(-2*time.Hour)))

		rows, err := setup.Repo.ListFirstLove(setup.Ctx, persist.ImageListSpec[time.Time]{
			Limit: 2,
		})
		assert.NilError(t, "ListFirstLove() error", err)
		assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
		assertRowsExcludeIngestID(t, rows, withoutFirstLovedAt.ID)

		nextCursor := assertLastRowTimeCursorEquals(t, rows, middleStat.FirstLovedAt.MustGet())
		nextRows, err := setup.Repo.ListFirstLove(setup.Ctx, persist.ImageListSpec[time.Time]{
			Cursor: mo.Some(nextCursor),
			Limit:  2,
		})
		assert.NilError(t, "ListFirstLove() error", err)
		assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
		assertRowsExcludeIngestID(t, nextRows, withoutFirstLovedAt.ID)
	})
}

func runListHallOfFame(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Run("ListHallOfFame", func(t *testing.T) {
		setup := setupFn(t)

		withoutHallOfFameAt := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(24*time.Hour)))
		latest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
		middle := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(4*24*time.Hour)))
		oldest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

		setup.Ops.MustAddStat(t, NewStatFixture(withoutHallOfFameAt.ID))
		setup.Ops.MustAddStat(t, NewStatFixtureWithHallOfFameAt(latest.ID, suiteBaseTimeUTC))
		middleStat := setup.Ops.MustAddStat(t, NewStatFixtureWithHallOfFameAt(middle.ID, suiteBaseTimeUTC.Add(-1*time.Hour)))
		setup.Ops.MustAddStat(t, NewStatFixtureWithHallOfFameAt(oldest.ID, suiteBaseTimeUTC.Add(-2*time.Hour)))

		rows, err := setup.Repo.ListHallOfFame(setup.Ctx, persist.ImageListSpec[time.Time]{
			Limit: 2,
		})
		assert.NilError(t, "ListHallOfFame() error", err)
		assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
		assertRowsExcludeIngestID(t, rows, withoutHallOfFameAt.ID)

		nextCursor := assertLastRowTimeCursorEquals(t, rows, middleStat.HallOfFameAt.MustGet())
		nextRows, err := setup.Repo.ListHallOfFame(setup.Ctx, persist.ImageListSpec[time.Time]{
			Cursor: mo.Some(nextCursor),
			Limit:  2,
		})
		assert.NilError(t, "ListHallOfFame() error", err)
		assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
		assertRowsExcludeIngestID(t, nextRows, withoutHallOfFameAt.ID)
	})
}

func runListEngaged(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Run("ListEngaged", func(t *testing.T) {
		setup := setupFn(t)

		hiddenHighest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(-3*24*time.Hour)))
		high := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-5*24*time.Hour)))
		middle := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(-2*24*time.Hour)))
		low := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(6, suiteBaseTimeUTC.Add(-1*24*time.Hour)))
		lowest := setup.Ops.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC.Add(-4*24*time.Hour)))

		hiddenHighestStat := NewStatFixtureWithScore(hiddenHighest.ID, 190, suiteBaseTimeUTC)
		hiddenHighestStat.HallOfFameAt = mo.Some(suiteBaseTimeUTC)
		setup.Ops.MustAddStat(t, hiddenHighestStat)
		setup.Ops.MustAddStat(t, NewStatFixtureWithScore(high.ID, 180, suiteBaseTimeUTC))
		middleStat := setup.Ops.MustAddStat(t, NewStatFixtureWithScore(middle.ID, 165, suiteBaseTimeUTC))
		setup.Ops.MustAddStat(t, NewStatFixtureWithScore(low.ID, 160, suiteBaseTimeUTC))
		setup.Ops.MustAddStat(t, NewStatFixtureWithScore(lowest.ID, 150, suiteBaseTimeUTC))

		rows, err := setup.Repo.ListEngaged(setup.Ctx, persist.EngagedImageListSpec{
			ImageListSpec: persist.ImageListSpec[int16]{
				Limit: 2,
			},
			ScoreThreshold: 160,
		})
		assert.NilError(t, "ListEngaged() error", err)
		assertRowsIngestIDs(t, rows, "rows", high.ID, middle.ID)
		assertRowsExcludeIngestID(t, rows, hiddenHighest.ID)

		nextCursor := assertLastRowInt16CursorEquals(t, rows, middleStat.Score)
		nextRows, err := setup.Repo.ListEngaged(setup.Ctx, persist.EngagedImageListSpec{
			ImageListSpec: persist.ImageListSpec[int16]{
				Cursor: mo.Some(nextCursor),
				Limit:  2,
			},
			ScoreThreshold: 160,
		})
		assert.NilError(t, "ListEngaged() error", err)
		assertRowsIngestIDs(t, nextRows, "nextRows", low.ID)
		assertRowsExcludeIngestID(t, nextRows, hiddenHighest.ID)
	})
}

func RunImageListSuite(t *testing.T, setupFn func(testing.TB) ImageListSetup) {
	t.Helper()

	runListLatest(t, setupFn)
	runListChronological(t, setupFn)
	runListRecently(t, setupFn)
	runListFirstLove(t, setupFn)
	runListHallOfFame(t, setupFn)
	runListEngaged(t, setupFn)
}
