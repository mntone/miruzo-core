package list

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type imageListCursorMode uint8

const (
	imageListCursorModeLatest imageListCursorMode = iota + 1
	imageListCursorModeChronological
	imageListCursorModeRecently
	imageListCursorModeFirstLove
	imageListCursorModeHallOfFame
	imageListCursorModeEngaged
)

const (
	cursorVersionDatetime uint8 = 1
	cursorVersionUint8    uint8 = 2

	cursorPackModeShift             = 56
	cursorPackVersionShift          = 53
	cursorPackVersionMask    uint64 = 0b111
	cursorPackIngestIDMask   uint64 = (1 << cursorPackVersionShift) - 1
	imageListCursorMinLength        = 9
)

var (
	errCursorDecode         = errors.New("cursor decode error")
	errEmptyCursor          = errors.New("cursor must not be empty")
	errInvalidPayloadLength = errors.New("cursor payload length is invalid")
	errInvalidMode          = errors.New("cursor mode is invalid")
	errInvalidVersion       = errors.New("cursor version is not supported")
	errModeMismatched       = errors.New("cursor mode is mismatched for the endpoint")
	errVersionMismatched    = errors.New("cursor version is mismatched for the endpoint mode")
	errIngestOutOfRange     = errors.New("cursor ingest_id is out of range")
	errValueOutOfRange      = errors.New("int cursor value must be uint8")
)

func expectedVersionForMode(mode imageListCursorMode) uint8 {
	if mode == imageListCursorModeEngaged {
		return cursorVersionUint8
	}
	return cursorVersionDatetime
}

func isSupportedMode(mode imageListCursorMode) bool {
	return mode >= imageListCursorModeLatest && mode <= imageListCursorModeEngaged
}

// --- encode ---

func packCursorHeader(mode imageListCursorMode, ingestID model.IngestIDType, version uint8) uint64 {
	return (uint64(mode) << cursorPackModeShift) |
		(uint64(version) << cursorPackVersionShift) |
		uint64(ingestID)
}

func encodeTimeImageListCursor(
	mode imageListCursorMode,
	cursor model.ImageListCursorKey[time.Time],
) (string, error) {
	packed := packCursorHeader(mode, cursor.Secondary, cursorVersionDatetime)
	micros := cursor.Primary.UnixMicro()

	var payload [16]byte
	binary.BigEndian.PutUint64(payload[:8], packed)
	binary.BigEndian.PutUint64(payload[8:], uint64(micros))
	return base64.RawURLEncoding.EncodeToString(payload[:]), nil
}

func encodeUint8ImageListCursor(
	mode imageListCursorMode,
	cursor model.ImageListCursorKey[model.ScoreType],
) (string, error) {
	if cursor.Primary < 0 || cursor.Primary > 0xFF {
		return "", errValueOutOfRange
	}

	packed := packCursorHeader(mode, cursor.Secondary, cursorVersionUint8)

	var payload [9]byte
	binary.BigEndian.PutUint64(payload[:8], packed)
	payload[8] = uint8(cursor.Primary)
	return base64.RawURLEncoding.EncodeToString(payload[:]), nil
}

// --- decode ---

func unpackCursorHeader(packed uint64) (imageListCursorMode, uint8, model.IngestIDType, error) {
	mode := imageListCursorMode((packed >> cursorPackModeShift) & 0xFF)
	version := uint8((packed >> cursorPackVersionShift) & cursorPackVersionMask)
	ingestID := model.IngestIDType(packed & cursorPackIngestIDMask)

	if ingestID < model.MinIngestID || ingestID > model.MaxIngestID {
		return 0, 0, 0, errIngestOutOfRange
	}
	if !isSupportedMode(mode) {
		return 0, 0, 0, errInvalidMode
	}
	return mode, version, ingestID, nil
}

func decodePayload(value string) ([]byte, error) {
	if value == "" {
		return nil, errEmptyCursor
	}

	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, errors.Join(errCursorDecode, err)
	}

	return payload, nil
}

func decodeCursor(
	cursor string,
	expectedMode imageListCursorMode,
) (uint8, model.IngestIDType, []byte, error) {
	payload, err := decodePayload(cursor)
	if err != nil {
		return 0, 0, nil, err
	}
	if len(payload) < imageListCursorMinLength {
		return 0, 0, nil, errInvalidPayloadLength
	}

	packed := binary.BigEndian.Uint64(payload[:8])
	mode, version, ingestID, err := unpackCursorHeader(packed)
	if err != nil {
		return 0, 0, nil, err
	}
	if mode != expectedMode {
		return 0, 0, nil, errModeMismatched
	}
	if version != cursorVersionDatetime && version != cursorVersionUint8 {
		return 0, 0, nil, errInvalidVersion
	}
	if version != expectedVersionForMode(expectedMode) {
		return 0, 0, nil, errVersionMismatched
	}

	return version, ingestID, payload, nil
}

func decodeTimeImageListCursor(
	cursor string,
	expectedMode imageListCursorMode,
) (model.ImageListCursorKey[time.Time], error) {
	version, ingestID, payload, err := decodeCursor(cursor, expectedMode)
	if err != nil {
		return model.ImageListCursorKey[time.Time]{}, err
	}
	if version == cursorVersionUint8 || len(payload) != 16 {
		return model.ImageListCursorKey[time.Time]{}, errInvalidPayloadLength
	}

	micros := int64(binary.BigEndian.Uint64(payload[8:]))
	return model.ImageListCursorKey[time.Time]{
		Primary:   time.UnixMicro(micros).UTC(),
		Secondary: ingestID,
	}, nil
}

func decodeUint8ImageListCursor(
	cursor string,
	expectedMode imageListCursorMode,
) (model.ImageListCursorKey[model.ScoreType], error) {
	version, ingestID, payload, err := decodeCursor(cursor, expectedMode)
	if err != nil {
		return model.ImageListCursorKey[model.ScoreType]{}, err
	}
	if version == cursorVersionDatetime || len(payload) != 9 {
		return model.ImageListCursorKey[model.ScoreType]{}, errInvalidPayloadLength
	}

	return model.ImageListCursorKey[model.ScoreType]{
		Primary:   model.ScoreType(payload[8]),
		Secondary: ingestID,
	}, nil
}
