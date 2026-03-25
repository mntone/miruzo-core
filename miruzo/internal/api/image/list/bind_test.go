package list

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

// --- time cursor ---

func TestBindParamsForTimeCursorUsesDefaultLimitAndEmptyCursor(t *testing.T) {
	got, err := bindParamsForTimeCursor(url.Values{}, imageListCursorModeLatest)
	assert.NilArray(t, "bindParamsForTimeCursor() error", err)
	assert.Equal(t, "params.Limit", got.Limit, defaultLimit)
	assert.IsAbsent(t, "params.Cursor", got.Cursor)
}

func TestBindParamsForTimeCursorParsesQueryValues(t *testing.T) {
	wantCursor := time.Date(2026, 3, 6, 12, 34, 56, 0, time.UTC)
	encodedCursor, encodeErr := encodeTimeImageListCursor(imageListCursorModeLatest, model.ImageListCursorKey[time.Time]{
		Primary:   wantCursor,
		Secondary: 123,
	})
	assert.NilError(t, "encodeDatetimeImageListCursor()", encodeErr)

	got, err := bindParamsForTimeCursor(url.Values{
		"limit":  []string{"42"},
		"cursor": []string{encodedCursor},
	}, imageListCursorModeLatest)
	assert.NilArray(t, "bindParamsForTimeCursor() error", err)
	assert.Equal(t, "params.Limit", got.Limit, 42)
	assert.IsPresent(t, "params.Cursor", got.Cursor)
	assert.EqualFn(t, "params.Cursor.Value", got.Cursor.MustGet().Primary, wantCursor)
	assert.Equal(t, "params.Cursor.ID", got.Cursor.MustGet().Secondary, int64(123))
}

func TestBindParamsForTimeCursorReturnsError(t *testing.T) {
	tests := []struct {
		name         string
		values       url.Values
		errorType    string
		errorPath    string
		errorMessage string
	}{
		{
			name: "Unsupported",
			values: url.Values{
				"Limit": []string{"11"},
			},
			errorType:    "unsupported",
			errorPath:    "query.Limit",
			errorMessage: "is not supported",
		},
		{
			name: "Duplicate",
			values: url.Values{
				"limit": []string{"24", "18"},
			},
			errorType:    "duplicate",
			errorPath:    "query.limit",
			errorMessage: "must not be specified multiple times",
		},
		{
			name: "InvalidCursor",
			values: url.Values{
				"cursor": []string{"invalid"},
			},
			errorType:    "invalid",
			errorPath:    "query.cursor",
			errorMessage: "must be a valid cursor",
		},
		{
			name: "InvalidLimit",
			values: url.Values{
				"limit": []string{"abc"},
			},
			errorType:    "invalid",
			errorPath:    "query.limit",
			errorMessage: "must be an integer",
		},
		{
			name: "LimitTooSmall",
			values: url.Values{
				"limit": []string{"0"},
			},
			errorType:    "invalid",
			errorPath:    "query.limit",
			errorMessage: "must be between 1 and 200",
		},
		{
			name: "LimitTooLarge",
			values: url.Values{
				"limit": []string{fmt.Sprint(limitMaximum + 1)},
			},
			errorType:    "invalid",
			errorPath:    "query.limit",
			errorMessage: "must be between 1 and 200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bindParamsForTimeCursor(tt.values, imageListCursorModeLatest)
			assert.LenIs(t, "bindParamsForTimeCursor() error", err, 1)
			assert.Equal(t, "err.Type", err[0].Type, tt.errorType)
			assert.Equal(t, "err.Path", err[0].Path, tt.errorPath)
			assert.Equal(t, "err.Message", err[0].Message, tt.errorMessage)
		})
	}
}

// --- score cursor ---

func TestBindParamsForScoreCursorUsesDefaultLimitAndEmptyCursor(t *testing.T) {
	got, err := bindParamsForScoreCursor(url.Values{}, imageListCursorModeEngaged)
	assert.NilArray(t, "bindParamsForScoreCursor() error", err)
	assert.Equal(t, "params.Limit", got.Limit, defaultLimit)
	assert.IsAbsent(t, "params.Cursor", got.Cursor)
}

func TestBindParamsForScoreCursorParsesQueryValues(t *testing.T) {
	encodedCursor, encodeErr := encodeUint8ImageListCursor(imageListCursorModeEngaged, model.ImageListCursorKey[model.ScoreType]{
		Primary:   170,
		Secondary: 123,
	})
	assert.NilError(t, "encodeUint8ImageListCursor()", encodeErr)

	got, err := bindParamsForScoreCursor(url.Values{
		"limit":  []string{"42"},
		"cursor": []string{encodedCursor},
	}, imageListCursorModeEngaged)
	assert.NilArray(t, "bindParamsForScoreCursor() error", err)
	assert.Equal(t, "params.Limit", got.Limit, 42)
	assert.IsPresent(t, "params.Cursor", got.Cursor)
	assert.Equal(t, "params.Cursor.Value", got.Cursor.MustGet().Primary, model.ScoreType(170))
	assert.Equal(t, "params.Cursor.ID", got.Cursor.MustGet().Secondary, int64(123))
}

func TestBindParamsForScoreCursorReturnsErrorForInvalidCursor(t *testing.T) {
	_, err := bindParamsForScoreCursor(url.Values{
		"limit":  []string{"42"},
		"cursor": []string{"999999"},
	}, imageListCursorModeEngaged)
	assert.LenIs(t, "bindParamsForScoreCursor() error", err, 1)
	assert.Equal(t, "err.Type", err[0].Type, "invalid")
	assert.Equal(t, "err.Path", err[0].Path, "query.cursor")
	assert.Equal(t, "err.Message", err[0].Message, "must be a valid cursor")
}
