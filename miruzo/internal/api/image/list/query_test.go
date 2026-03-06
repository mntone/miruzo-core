package list

import (
	"net/url"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/api/common"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

func TestBindQueryWithTimeCursorUsesDefaultLimitAndEmptyCursor(t *testing.T) {
	got, errs := bindQueryWithTimeCursor(url.Values{})
	assert.Empty(t, "bindQueryWithTimeCursor()", errs)
	assert.Equal(t, "query.Limit", got.Limit, defaultLimit)
	assert.IsAbsent(t, "query.Cursor", got.Cursor)
}

func TestBindQueryWithTimeCursorParsesQueryValues(t *testing.T) {
	wantCursor := time.Date(2026, 3, 6, 12, 34, 56, 0, time.UTC)

	got, errs := bindQueryWithTimeCursor(url.Values{
		"limit":  []string{"42"},
		"cursor": []string{"2026-03-06T12:34:56.000000Z"},
	})
	assert.Empty(t, "bindQueryWithTimeCursor()", errs)
	assert.Equal(t, "query.Limit", got.Limit, 42)
	assert.IsPresent(t, "query.Cursor", got.Cursor)
	assert.EqualFn(t, "query.Cursor", got.Cursor.MustGet(), wantCursor)
}

func TestBindQueryWithTimeCursorReturnsErrorForInvalidLimit(t *testing.T) {
	_, errs := bindQueryWithTimeCursor(url.Values{
		"limit": []string{"abc"},
	})
	assert.LenIs(t, "errors", errs, 1)
	assert.Equal(t, "errors[0].Path", errs[0].Path, "query.limit")
	assert.Equal(t, "errors[0].Type", errs[0].Type, "invalid_type")
}

func TestBindQueryWithTimeCursorReturnsErrorForInvalidCursor(t *testing.T) {
	_, errs := bindQueryWithTimeCursor(url.Values{
		"limit":  []string{"10"},
		"cursor": []string{"invalid"},
	})
	assert.LenIs(t, "errors", errs, 1)
	assert.Equal(t, "errors[0].Path", errs[0].Path, "query.cursor")
	assert.Equal(t, "errors[0].Type", errs[0].Type, "invalid_type")
}

func TestBindQueryWithInt16CursorUsesDefaultLimitAndEmptyCursor(t *testing.T) {
	got, errs := bindQueryWithInt16Cursor(url.Values{})
	assert.Empty(t, "bindQueryWithInt16Cursor()", errs)
	assert.Equal(t, "query.Limit", got.Limit, defaultLimit)
	assert.IsAbsent(t, "query.Cursor", got.Cursor)
}

func TestBindQueryWithInt16CursorParsesQueryValues(t *testing.T) {
	got, errs := bindQueryWithInt16Cursor(url.Values{
		"limit":  []string{"24"},
		"cursor": []string{"170"},
	})
	assert.Empty(t, "bindQueryWithInt16Cursor()", errs)
	assert.Equal(t, "query.Limit", got.Limit, uint16(24))
	assert.IsPresent(t, "query.Cursor", got.Cursor)
	assert.Equal(t, "query.Cursor", got.Cursor.MustGet(), int16(170))
}

func TestBindQueryWithInt16CursorReturnsErrorForInvalidCursor(t *testing.T) {
	_, errs := bindQueryWithInt16Cursor(url.Values{
		"limit":  []string{"10"},
		"cursor": []string{"999999"},
	})
	assert.LenIs(t, "errors", errs, 1)
	assert.Equal(t, "errors[0].Path", errs[0].Path, "query.cursor")
	assert.Equal(t, "errors[0].Type", errs[0].Type, "invalid_type")
}

func TestBuildTimeParamsFromQueryBuildsParams(t *testing.T) {
	wantCursor := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)
	params, errs := buildTimeParamsFromQuery(query[time.Time]{
		PaginationQuery: common.PaginationQuery[time.Time]{
			Limit:  20,
			Cursor: mo.Some(wantCursor),
		},
	})
	assert.Empty(t, "buildTimeParamsFromQuery()", errs)
	assert.NotNil(t, "params", params)
	assert.Equal(t, "params.Limit", params.Limit, uint16(20))
	assert.IsPresent(t, "params.Cursor", params.Cursor)
	assert.EqualFn(t, "params.Cursor", params.Cursor.MustGet(), wantCursor)
	if params.ExcludeFormats != nil {
		t.Fatalf("params.ExcludeFormats = %v, want nil", params.ExcludeFormats)
	}
}

func TestBuildTimeParamsFromQueryReturnsErrorForOutOfRangeLimit(t *testing.T) {
	params, errs := buildTimeParamsFromQuery(query[time.Time]{
		PaginationQuery: common.PaginationQuery[time.Time]{
			Limit: 0,
		},
	})
	assert.Nil(t, "params", params)
	assert.LenIs(t, "errors", errs, 1)
	assert.Equal(t, "errors[0].Path", errs[0].Path, "query.limit")
	assert.Equal(t, "errors[0].Type", errs[0].Type, "out_of_range")
}

func TestBuildInt16ParamsFromQueryBuildsParams(t *testing.T) {
	params, errs := buildInt16ParamsFromQuery(query[int16]{
		PaginationQuery: common.PaginationQuery[int16]{
			Limit:  30,
			Cursor: mo.Some[int16](160),
		},
	})
	assert.Empty(t, "buildInt16ParamsFromQuery()", errs)
	assert.NotNil(t, "params", params)
	assert.Equal(t, "params.Limit", params.Limit, uint16(30))
	assert.IsPresent(t, "params.Cursor", params.Cursor)
	assert.Equal(t, "params.Cursor", params.Cursor.MustGet(), int16(160))
	if params.ExcludeFormats != nil {
		t.Fatalf("params.ExcludeFormats = %v, want nil", params.ExcludeFormats)
	}
}

func TestBuildInt16ParamsFromQueryReturnsErrorForOutOfRangeLimit(t *testing.T) {
	params, errs := buildInt16ParamsFromQuery(query[int16]{
		PaginationQuery: common.PaginationQuery[int16]{
			Limit: limitMaximum + 1,
		},
	})
	assert.Nil(t, "params", params)
	assert.LenIs(t, "errors", errs, 1)
	assert.Equal(t, "errors[0].Path", errs[0].Path, "query.limit")
	assert.Equal(t, "errors[0].Type", errs[0].Type, "out_of_range")
}
