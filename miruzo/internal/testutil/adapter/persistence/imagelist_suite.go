package persistence

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

var suiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

func assertRowsIngestIDs[C persist.ImageListCursor](
	t *testing.T,
	rows []persist.ImageWithCursor[C],
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
	disallowedID model.IngestIDType,
) {
	t.Helper()

	for i, row := range rows {
		if row.Image.IngestID == disallowedID {
			t.Fatalf("rows[%d].Image.IngestID = %d, must not contain %d", i, row.Image.IngestID, disallowedID)
		}
	}
}

type ImageListSuite SuiteBase[persist.ImageListRepository]

func (ste ImageListSuite) RunTestListLatest(t *testing.T) {
	latest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC))
	middle := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC.Add(-1*24*time.Hour)))
	oldest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

	rows, err := ste.Repository.ListLatest(ste.Context, persist.ImageListSpec[time.Time]{
		Limit: 2,
	})
	assert.NilError(t, "ListLatest() error", err)
	assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)

	nextCursor := assertLastRowTimeCursorEquals(t, rows, middle.IngestedAt)
	nextRows, err := ste.Repository.ListLatest(ste.Context, persist.ImageListSpec[time.Time]{
		Cursor: mo.Some(nextCursor),
		Limit:  2,
	})
	assert.NilError(t, "ListLatest() error", err)
	assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
}

func (ste ImageListSuite) RunTestListChronological(t *testing.T) {
	latest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
	middle := ste.Operations.MustAddIngestAndImage(t, NewIngestFixtureWithCapturedAt(
		4,
		suiteBaseTimeUTC.Add(4*time.Hour),
		suiteBaseTimeUTC.Add(-1*time.Hour),
	))
	oldest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*time.Hour)))

	rows, err := ste.Repository.ListChronological(ste.Context, persist.ImageListSpec[time.Time]{
		Limit: 2,
	})
	assert.NilError(t, "ListChronological() error", err)
	assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)

	nextCursor := assertLastRowTimeCursorEquals(t, rows, middle.CapturedAt)
	nextRows, err := ste.Repository.ListChronological(ste.Context, persist.ImageListSpec[time.Time]{
		Cursor: mo.Some(nextCursor),
		Limit:  2,
	})
	assert.NilError(t, "ListChronological() error", err)
	assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
}

func (ste ImageListSuite) RunTestListRecently(t *testing.T) {
	withoutLastViewedAt := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(24*time.Hour)))
	latest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
	middle := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(4*24*time.Hour)))
	oldest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

	ste.Operations.MustAddStat(t, NewStatFixture(withoutLastViewedAt.ID))
	ste.Operations.MustAddStat(t, NewStatFixtureWithLastViewedAt(latest.ID, 1, suiteBaseTimeUTC))
	middleStat := ste.Operations.MustAddStat(t, NewStatFixtureWithLastViewedAt(middle.ID, 1, suiteBaseTimeUTC.Add(-1*time.Hour)))
	ste.Operations.MustAddStat(t, NewStatFixtureWithLastViewedAt(oldest.ID, 1, suiteBaseTimeUTC.Add(-2*time.Hour)))

	rows, err := ste.Repository.ListRecently(ste.Context, persist.ImageListSpec[time.Time]{
		Limit: 2,
	})
	assert.NilError(t, "ListRecently() error", err)
	assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
	assertRowsExcludeIngestID(t, rows, withoutLastViewedAt.ID)

	nextCursor := assertLastRowTimeCursorEquals(t, rows, middleStat.LastViewedAt.MustGet())
	nextRows, err := ste.Repository.ListRecently(ste.Context, persist.ImageListSpec[time.Time]{
		Cursor: mo.Some(nextCursor),
		Limit:  2,
	})
	assert.NilError(t, "ListRecently() error", err)
	assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
	assertRowsExcludeIngestID(t, nextRows, withoutLastViewedAt.ID)
}

func (ste ImageListSuite) RunTestListFirstLove(t *testing.T) {
	withoutFirstLovedAt := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(24*time.Hour)))
	latest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
	middle := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(4*24*time.Hour)))
	oldest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

	ste.Operations.MustAddStat(t, NewStatFixture(withoutFirstLovedAt.ID))
	ste.Operations.MustAddStat(t, NewStatFixtureWithLastLovedAt(latest.ID, suiteBaseTimeUTC))
	middleStat := ste.Operations.MustAddStat(t, NewStatFixtureWithLastLovedAt(middle.ID, suiteBaseTimeUTC.Add(-1*time.Hour)))
	ste.Operations.MustAddStat(t, NewStatFixtureWithLastLovedAt(oldest.ID, suiteBaseTimeUTC.Add(-2*time.Hour)))

	rows, err := ste.Repository.ListFirstLove(ste.Context, persist.ImageListSpec[time.Time]{
		Limit: 2,
	})
	assert.NilError(t, "ListFirstLove() error", err)
	assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
	assertRowsExcludeIngestID(t, rows, withoutFirstLovedAt.ID)

	nextCursor := assertLastRowTimeCursorEquals(t, rows, middleStat.FirstLovedAt.MustGet())
	nextRows, err := ste.Repository.ListFirstLove(ste.Context, persist.ImageListSpec[time.Time]{
		Cursor: mo.Some(nextCursor),
		Limit:  2,
	})
	assert.NilError(t, "ListFirstLove() error", err)
	assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
	assertRowsExcludeIngestID(t, nextRows, withoutFirstLovedAt.ID)
}

func (ste ImageListSuite) RunTestListHallOfFame(t *testing.T) {
	withoutHallOfFameAt := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(24*time.Hour)))
	latest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC))
	middle := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(4*24*time.Hour)))
	oldest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-2*24*time.Hour)))

	ste.Operations.MustAddStat(t, NewStatFixture(withoutHallOfFameAt.ID))
	ste.Operations.MustAddStat(t, NewStatFixtureWithHallOfFameAt(latest.ID, suiteBaseTimeUTC))
	middleStat := ste.Operations.MustAddStat(t, NewStatFixtureWithHallOfFameAt(middle.ID, suiteBaseTimeUTC.Add(-1*time.Hour)))
	ste.Operations.MustAddStat(t, NewStatFixtureWithHallOfFameAt(oldest.ID, suiteBaseTimeUTC.Add(-2*time.Hour)))

	rows, err := ste.Repository.ListHallOfFame(ste.Context, persist.ImageListSpec[time.Time]{
		Limit: 2,
	})
	assert.NilError(t, "ListHallOfFame() error", err)
	assertRowsIngestIDs(t, rows, "rows", latest.ID, middle.ID)
	assertRowsExcludeIngestID(t, rows, withoutHallOfFameAt.ID)

	nextCursor := assertLastRowTimeCursorEquals(t, rows, middleStat.HallOfFameAt.MustGet())
	nextRows, err := ste.Repository.ListHallOfFame(ste.Context, persist.ImageListSpec[time.Time]{
		Cursor: mo.Some(nextCursor),
		Limit:  2,
	})
	assert.NilError(t, "ListHallOfFame() error", err)
	assertRowsIngestIDs(t, nextRows, "nextRows", oldest.ID)
	assertRowsExcludeIngestID(t, nextRows, withoutHallOfFameAt.ID)
}

func (ste ImageListSuite) RunTestListEngaged(t *testing.T) {
	hiddenHighest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(4, suiteBaseTimeUTC.Add(-3*24*time.Hour)))
	high := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, suiteBaseTimeUTC.Add(-5*24*time.Hour)))
	middle := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(5, suiteBaseTimeUTC.Add(-2*24*time.Hour)))
	low := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(6, suiteBaseTimeUTC.Add(-1*24*time.Hour)))
	lowest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(3, suiteBaseTimeUTC.Add(-4*24*time.Hour)))

	hiddenHighestStat := NewStatFixtureWithScore(hiddenHighest.ID, 190, suiteBaseTimeUTC)
	hiddenHighestStat.HallOfFameAt = mo.Some(suiteBaseTimeUTC)
	ste.Operations.MustAddStat(t, hiddenHighestStat)
	ste.Operations.MustAddStat(t, NewStatFixtureWithScore(high.ID, 180, suiteBaseTimeUTC))
	middleStat := ste.Operations.MustAddStat(t, NewStatFixtureWithScore(middle.ID, 165, suiteBaseTimeUTC))
	ste.Operations.MustAddStat(t, NewStatFixtureWithScore(low.ID, 160, suiteBaseTimeUTC))
	ste.Operations.MustAddStat(t, NewStatFixtureWithScore(lowest.ID, 150, suiteBaseTimeUTC))

	rows, err := ste.Repository.ListEngaged(ste.Context, persist.EngagedImageListSpec{
		ImageListSpec: persist.ImageListSpec[int16]{
			Limit: 2,
		},
		ScoreThreshold: 160,
	})
	assert.NilError(t, "ListEngaged() error", err)
	assertRowsIngestIDs(t, rows, "rows", high.ID, middle.ID)
	assertRowsExcludeIngestID(t, rows, hiddenHighest.ID)

	nextCursor := assertLastRowInt16CursorEquals(t, rows, middleStat.Score)
	nextRows, err := ste.Repository.ListEngaged(ste.Context, persist.EngagedImageListSpec{
		ImageListSpec: persist.ImageListSpec[int16]{
			Cursor: mo.Some(nextCursor),
			Limit:  2,
		},
		ScoreThreshold: 160,
	})
	assert.NilError(t, "ListEngaged() error", err)
	assertRowsIngestIDs(t, nextRows, "nextRows", low.ID)
	assertRowsExcludeIngestID(t, nextRows, hiddenHighest.ID)
}
