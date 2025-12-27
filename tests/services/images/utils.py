from collections.abc import Iterable
from datetime import datetime, timezone
from pathlib import Path
from tempfile import gettempdir

from sqlmodel import Session

from app.config.variant import VariantFormat, VariantSlotkey, VariantSpec
from app.models.enums import ProcessStatus, VisibilityStatus
from app.models.records import ImageRecord, IngestRecord
from app.models.types import VariantEntry

TEST_VARIANT_ROOT = Path(gettempdir()) / 'miruzo-test-variants'


def build_variant(fmt: str, width: int, *, layer_id: int = 1, label: str = 'primary') -> VariantEntry:
	filepath = TEST_VARIANT_ROOT / f'{label}-{width}.{fmt}'
	filepath.parent.mkdir(parents=True, exist_ok=True)
	payload = f'{label}-{width}-{fmt}'.encode('utf-8')
	filepath.write_bytes(payload)
	return {
		'rel': filepath.as_posix(),
		'layer_id': layer_id,
		'format': fmt,
		'codecs': None,
		'bytes': filepath.stat().st_size,
		'width': width,
		'height': round(width * 0.75),
		'quality': None,
	}


def _make_ingest_record(
	ingest_id: int,
	*,
	relative_path: str = '/foo/bar.webp',
	process: ProcessStatus = ProcessStatus.PROCESSING,
	visibility: VisibilityStatus = VisibilityStatus.PRIVATE,
	captured_at: datetime | None = None,
) -> IngestRecord:
	timestamp = captured_at or datetime.now(timezone.utc)
	return IngestRecord(
		id=ingest_id,
		process=process,
		visibility=visibility,
		relative_path=relative_path,
		fingerprint=f'{ingest_id:064d}',
		ingested_at=timestamp,
		captured_at=timestamp,
	)


def _make_image_record(
	ingest_id: int,
	*,
	captured_at: datetime | None = None,
	widths: Iterable[int] = [320, 480, 640],
) -> ImageRecord:
	timestamp = captured_at or datetime.now(timezone.utc)
	return ImageRecord(
		ingest_id=ingest_id,
		captured_at=timestamp,
		original=build_variant('webp', 960),
		fallback=None,
		variants=[
			*[
				build_variant('webp', width, layer_id=1)
				for width in widths
			],
			build_variant('jpeg', 320, layer_id=9, label='fallback'),
		],
	)


def build_image_record(ingest_id: int) -> ImageRecord:
	return _make_image_record(
		ingest_id=ingest_id,
	)


def add_ingest_record(
	session: Session,
	idx: int,
	*,
	relative_path: str = '/foo/bar.webp',
	process: ProcessStatus = ProcessStatus.PROCESSING,
	visibility: VisibilityStatus = VisibilityStatus.PRIVATE,
	captured_at: datetime | None = None,
) -> IngestRecord:
	record = _make_ingest_record(
		ingest_id=idx,
		relative_path=relative_path,
		process=process,
		visibility=visibility,
		captured_at=captured_at,
	)
	session.add(record)
	session.commit()
	session.refresh(record)
	return record


def add_image_record(
	session: Session,
	idx: int,
	*,
	relative_path: str = '/foo/bar.webp',
	process: ProcessStatus = ProcessStatus.PROCESSING,
	visibility: VisibilityStatus = VisibilityStatus.PRIVATE,
	captured_at: datetime | None = None,
) -> ImageRecord:
	ingest = _make_ingest_record(
		ingest_id=idx,
		relative_path=relative_path,
		process=process,
		visibility=visibility,
		captured_at=captured_at,
	)
	session.add(ingest)

	image = _make_image_record(
		ingest_id=idx,
		captured_at=captured_at,
	)
	session.add(image)

	session.commit()
	session.refresh(image)
	return image


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
