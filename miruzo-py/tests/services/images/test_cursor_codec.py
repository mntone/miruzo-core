import base64
from datetime import datetime, timezone

import pytest

from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursorMode,
	UInt8ImageListCursor,
)
from app.services.images.cursor_codec import (
	CURSOR_VERSION_UINT8,
	CursorDecodeError,
	decode_datetime_image_list_cursor,
	decode_uint8_image_list_cursor,
	encode_image_list_cursor,
)


def test_encode_decode_datetime_cursor_roundtrip() -> None:
	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.FIRST_LOVE,
		value=datetime(2026, 2, 1, 12, 0, 0, 123456, tzinfo=timezone.utc),
		ingest_id=123,
	)

	encoded = encode_image_list_cursor(cursor)
	decoded = decode_datetime_image_list_cursor(
		encoded,
		expected_mode=ImageListCursorMode.FIRST_LOVE,
	)

	assert decoded == cursor


def test_encode_decode_uint8_cursor_roundtrip() -> None:
	cursor = UInt8ImageListCursor(
		mode=ImageListCursorMode.ENGAGED,
		value=200,
		ingest_id=999,
	)

	encoded = encode_image_list_cursor(cursor)
	decoded = decode_uint8_image_list_cursor(
		encoded,
		expected_mode=ImageListCursorMode.ENGAGED,
	)

	assert decoded == cursor


def test_decode_rejects_mode_mismatch() -> None:
	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.LATEST,
		value=datetime(2026, 2, 1, 12, 0, 0, tzinfo=timezone.utc),
		ingest_id=1,
	)
	encoded = encode_image_list_cursor(cursor)

	with pytest.raises(CursorDecodeError):
		decode_datetime_image_list_cursor(
			encoded,
			expected_mode=ImageListCursorMode.CHRONOLOGICAL,
		)


def test_decode_rejects_unsupported_version() -> None:
	cursor = UInt8ImageListCursor(
		mode=ImageListCursorMode.ENGAGED,
		value=120,
		ingest_id=5,
	)
	encoded = encode_image_list_cursor(cursor)

	payload = base64.urlsafe_b64decode(encoded + '=' * (-len(encoded) % 4))
	packed = int.from_bytes(payload[:8], byteorder='big', signed=False)
	tampered_version = 7 if CURSOR_VERSION_UINT8 != 7 else 6
	tampered = (packed & ~(0b111 << 53)) | (tampered_version << 53)
	tampered_payload = tampered.to_bytes(8, byteorder='big', signed=False) + payload[8:]
	tampered_encoded = base64.urlsafe_b64encode(tampered_payload).decode('ascii').rstrip('=')

	with pytest.raises(CursorDecodeError):
		decode_uint8_image_list_cursor(
			tampered_encoded,
			expected_mode=ImageListCursorMode.ENGAGED,
		)
