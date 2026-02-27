import base64
from datetime import datetime, timedelta, timezone

from app.config.constants import INGEST_ID_MAXIMUM, INGEST_ID_MINIMUM
from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursor,
	ImageListCursorMode,
	UInt8ImageListCursor,
)

_EPOCH = datetime(1970, 1, 1, tzinfo=timezone.utc)
_PACK_MODE_SHIFT = 56
_PACK_VERSION_SHIFT = 53
_PACK_VERSION_MASK = 0b111
_PACK_INGEST_ID_MASK = (1 << _PACK_VERSION_SHIFT) - 1

CURSOR_VERSION_DATETIME = 1
CURSOR_VERSION_UINT8 = 2


class CursorDecodeError(ValueError):
	pass


def _expected_version_for_mode(mode: ImageListCursorMode) -> int:
	if mode == ImageListCursorMode.ENGAGED:
		return CURSOR_VERSION_UINT8
	else:
		return CURSOR_VERSION_DATETIME


def _pack_format(mode: ImageListCursorMode, ingest_id: int, *, version: int) -> int:
	if ingest_id < INGEST_ID_MINIMUM or ingest_id > INGEST_ID_MAXIMUM:
		raise ValueError('cursor ingest_id is out of range')

	return (mode.value << _PACK_MODE_SHIFT) | (version << _PACK_VERSION_SHIFT) | ingest_id


def _unpack_format(packed: int) -> tuple[ImageListCursorMode, int, int]:
	mode_raw = (packed >> _PACK_MODE_SHIFT) & 0xFF
	version = (packed >> _PACK_VERSION_SHIFT) & _PACK_VERSION_MASK
	ingest_id = packed & _PACK_INGEST_ID_MASK

	if ingest_id < INGEST_ID_MINIMUM or ingest_id > INGEST_ID_MAXIMUM:
		raise CursorDecodeError('cursor ingest_id is out of range')

	try:
		mode = ImageListCursorMode(mode_raw)
	except ValueError as exception:
		raise CursorDecodeError('cursor mode is invalid') from exception

	return mode, version, ingest_id


def _to_epoch_microseconds(value: datetime) -> int:
	if value.tzinfo is None:
		raise ValueError('datetime cursor must include timezone')

	value_utc = value.astimezone(timezone.utc)
	delta = value_utc - _EPOCH

	return (delta.days * 86_400 + delta.seconds) * 1_000_000 + delta.microseconds


def _from_epoch_microseconds(value: int) -> datetime:
	return _EPOCH + timedelta(microseconds=value)


def encode_image_list_cursor(cursor: ImageListCursor) -> str:
	if isinstance(cursor, DatetimeImageListCursor):
		packed = _pack_format(
			cursor.mode,
			cursor.ingest_id,
			version=CURSOR_VERSION_DATETIME,
		)
		value = _to_epoch_microseconds(cursor.value)
		payload = packed.to_bytes(8, byteorder='big') + value.to_bytes(8, byteorder='big', signed=True)
	else:
		packed = _pack_format(
			cursor.mode,
			cursor.ingest_id,
			version=CURSOR_VERSION_UINT8,
		)
		if cursor.value < 0 or cursor.value > 0xFF:
			raise ValueError('int cursor value must be uint8')
		payload = packed.to_bytes(8, byteorder='big') + bytes((cursor.value,))

	encoded = base64.urlsafe_b64encode(payload).decode('ascii')

	return encoded.rstrip('=')


def _decode_payload(value: str) -> bytes:
	if value == '':
		raise CursorDecodeError('cursor must not be empty')

	padding = '=' * (-len(value) % 4)
	try:
		return base64.urlsafe_b64decode(value + padding)
	except Exception as exception:  # noqa: BLE001
		raise CursorDecodeError('cursor is not valid base64url') from exception


def decode_image_list_cursor(
	cursor: str,
	*,
	expected_mode: ImageListCursorMode,
) -> ImageListCursor:
	payload = _decode_payload(cursor)
	if len(payload) < 9:
		raise CursorDecodeError('cursor payload length is invalid')

	packed = int.from_bytes(payload[:8], byteorder='big', signed=False)
	mode, version, ingest_id = _unpack_format(packed)

	if mode != expected_mode:
		raise CursorDecodeError('cursor mode is mismatched for the endpoint')
	if version not in (CURSOR_VERSION_DATETIME, CURSOR_VERSION_UINT8):
		raise CursorDecodeError('cursor version is not supported')
	if version != _expected_version_for_mode(expected_mode):
		raise CursorDecodeError('cursor version is mismatched for the endpoint mode')

	if version == CURSOR_VERSION_UINT8:
		if len(payload) != 9:
			raise CursorDecodeError('cursor payload length is invalid for uint8 cursor')
		value = payload[8]
		return UInt8ImageListCursor(
			mode=mode,
			value=value,
			ingest_id=ingest_id,
		)

	if len(payload) != 16:
		raise CursorDecodeError('cursor payload length is invalid for datetime cursor')

	value = int.from_bytes(payload[8:], byteorder='big', signed=True)
	return DatetimeImageListCursor(
		mode=mode,
		value=_from_epoch_microseconds(value),
		ingest_id=ingest_id,
	)


def decode_datetime_image_list_cursor(
	cursor: str,
	*,
	expected_mode: ImageListCursorMode,
) -> DatetimeImageListCursor:
	decoded = decode_image_list_cursor(cursor, expected_mode=expected_mode)
	if not isinstance(decoded, DatetimeImageListCursor):
		raise CursorDecodeError('cursor payload type is mismatched for datetime mode')
	return decoded


def decode_uint8_image_list_cursor(
	cursor: str,
	*,
	expected_mode: ImageListCursorMode,
) -> UInt8ImageListCursor:
	decoded = decode_image_list_cursor(cursor, expected_mode=expected_mode)
	if not isinstance(decoded, UInt8ImageListCursor):
		raise CursorDecodeError('cursor payload type is mismatched for uint8 mode')
	return decoded
