from collections.abc import Iterable, Sequence
from datetime import datetime, timezone
from pathlib import Path
from tempfile import gettempdir

from sqlmodel import Session

from app.config.variant import VariantFormat, VariantSlotkey, VariantSpec
from app.models.enums import ImageStatus
from app.models.records import ImageRecord, VariantRecord

TEST_VARIANT_ROOT = Path(gettempdir()) / 'miruzo-test-variants'


def build_variant(fmt: str, width: int, *, label: str = 'primary') -> VariantRecord:
	filepath = TEST_VARIANT_ROOT / f'{label}-{width}.{fmt}'
	filepath.parent.mkdir(parents=True, exist_ok=True)
	payload = f'{label}-{width}-{fmt}'.encode('utf-8')
	filepath.write_bytes(payload)
	return {
		'filepath': filepath.as_posix(),
		'format': fmt,
		'codecs': None,
		'size': filepath.stat().st_size,
		'width': width,
		'height': round(width * 0.75),
		'quality': None,
	}


def _make_image_record(
	image_id: int,
	formats: Sequence[str],
	*,
	captured_at: datetime | None,
	width_offset: int,
) -> ImageRecord:
	now = captured_at or datetime.now(timezone.utc)
	base_width = 320 + width_offset
	fallback_variant = build_variant('jpeg', base_width)
	return ImageRecord(
		id=image_id,
		fingerprint=f'{image_id:064d}',
		captured_at=now,
		ingested_at=now,
		status=ImageStatus.ACTIVE,
		original=build_variant(formats[0], base_width),
		fallback=fallback_variant,
		variants=[
			[build_variant(fmt, base_width) for fmt in formats],
			[fallback_variant],
		],
	)


def build_image_record(
	image_id: int,
	formats: Sequence[str],
	*,
	captured_at: datetime | None = None,
) -> ImageRecord:
	return _make_image_record(
		image_id=image_id,
		formats=formats,
		captured_at=captured_at,
		width_offset=0,
	)


def add_image_record(
	session: Session,
	idx: int,
	*,
	captured_at: datetime | None = None,
	formats: Iterable[str] | None = None,
) -> ImageRecord:
	timestamp = captured_at or datetime.now(timezone.utc)
	format_list = list(formats or ['webp'])
	record = _make_image_record(
		image_id=idx,
		formats=format_list,
		captured_at=timestamp,
		width_offset=idx,
	)
	session.add(record)
	session.commit()
	session.refresh(record)
	return record


def build_variant_spec(
	layer_id: int,
	width: int,
	*,
	container: str = 'jpeg',
	codecs: str | None = None,
	quality: int | None = None,
	required: bool = False,
) -> VariantSpec:
	return VariantSpec(
		slotkey=VariantSlotkey(layer_id, width),
		layer_id=layer_id,
		width=width,
		format=VariantFormat(container=container, codecs=codecs, file_extension=f'.{container}'),  # pyright: ignore[reportArgumentType]
		quality=quality,
		required=required,
	)
