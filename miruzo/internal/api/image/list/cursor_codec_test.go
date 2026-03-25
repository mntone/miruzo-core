package list

import (
	"encoding/base64"
	"encoding/binary"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestEncodeDecodeTimeCursorRoundtrip(t *testing.T) {
	wantValue := time.Date(2026, 2, 1, 12, 0, 0, 123456000, time.UTC)
	wantIngestID := int64(123)

	encoded, encodeErr := encodeTimeImageListCursor(imageListCursorModeFirstLove, model.ImageListCursorKey[time.Time]{
		Primary:   wantValue,
		Secondary: wantIngestID,
	})
	assert.NilError(t, "encodeTimeImageListCursor()", encodeErr)

	got, decodeErr := decodeTimeImageListCursor(
		encoded,
		imageListCursorModeFirstLove,
	)
	assert.NilError(t, "decodeTimeImageListCursor()", decodeErr)
	assert.EqualFn(t, "decodeTimeImageListCursor().Primary", got.Primary, wantValue)
	assert.Equal(t, "decodeTimeImageListCursor().Secondary", got.Secondary, wantIngestID)
}

func TestEncodeDecodeUint8CursorRoundtrip(t *testing.T) {
	wantValue := uint8(200)
	wantIngestID := int64(999)

	encoded, encodeErr := encodeUint8ImageListCursor(imageListCursorModeEngaged, model.ImageListCursorKey[model.ScoreType]{
		Primary:   200,
		Secondary: wantIngestID,
	})
	assert.NilError(t, "encodeUint8ImageListCursor()", encodeErr)

	got, decodeErr := decodeUint8ImageListCursor(
		encoded,
		imageListCursorModeEngaged,
	)
	assert.NilError(t, "decodeUint8ImageListCursor()", decodeErr)
	assert.Equal(t, "decodeUint8ImageListCursor().Primary", got.Primary, model.ScoreType(wantValue))
	assert.Equal(t, "decodeUint8ImageListCursor().Secondary", got.Secondary, wantIngestID)
}

func TestDecodeRejectsModeMismatch(t *testing.T) {
	encoded, encodeErr := encodeTimeImageListCursor(imageListCursorModeLatest, model.ImageListCursorKey[time.Time]{
		Primary:   time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC),
		Secondary: 1,
	})
	assert.NilError(t, "encodeTimeImageListCursor()", encodeErr)

	_, decodeErr := decodeTimeImageListCursor(encoded, imageListCursorModeChronological)
	assert.Error(t, "decodeTimeImageListCursor()", decodeErr)
}

func TestDecodeRejectsUnsupportedVersion(t *testing.T) {
	encoded, encodeErr := encodeUint8ImageListCursor(imageListCursorModeEngaged, model.ImageListCursorKey[model.ScoreType]{
		Primary:   120,
		Secondary: 5,
	})
	assert.NilError(t, "encodeUint8ImageListCursor()", encodeErr)

	payload, decodeErr := base64.RawURLEncoding.DecodeString(encoded)
	assert.NilError(t, "DecodeString()", decodeErr)

	packed := binary.BigEndian.Uint64(payload[:8])
	tamperedVersion := uint8(7)
	if cursorVersionUint8 == 7 {
		tamperedVersion = 6
	}
	tampered := (packed & ^(uint64(0b111) << cursorPackVersionShift)) |
		(uint64(tamperedVersion) << cursorPackVersionShift)
	binary.BigEndian.PutUint64(payload[:8], tampered)
	tamperedEncoded := base64.RawURLEncoding.EncodeToString(payload)

	_, decodeCursorErr := decodeUint8ImageListCursor(
		tamperedEncoded,
		imageListCursorModeEngaged,
	)
	assert.Error(t, "decodeUint8ImageListCursor()", decodeCursorErr)
}
