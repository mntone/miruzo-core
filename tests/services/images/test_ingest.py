from collections.abc import Iterator, Sequence
from datetime import datetime
from pathlib import Path
from typing import cast

import pytest
from PIL import Image as PILImage

from tests.services.images.utils import build_variant_spec
from tests.services.images.variants.utils import build_variant_file

from app.config.variant import VariantLayerSpec
from app.models.enums import ExecutionStatus, IngestMode
from app.models.records import ImageRecord, IngestRecord
from app.models.types import ExecutionEntry
from app.services.images.ingest import ImageIngestService
from app.services.images.variants.types import (
	OriginalFile,
	VariantCommitResult,
	VariantPolicy,
	VariantReport,
)


class DummyIngestCore:
	def __init__(self, record: IngestRecord) -> None:
		self.record = record
		self.created_args: dict[str, object] | None = None
		self.appended: tuple[int, ExecutionEntry] | None = None

	def create_ingest(
		self,
		*,
		origin_path: Path,
		fingerprint: str | None,
		captured_at: datetime | None,
		ingest_mode: IngestMode,
	) -> IngestRecord:
		self.created_args = {
			'origin_path': origin_path,
			'fingerprint': fingerprint,
			'captured_at': captured_at,
			'ingest_mode': ingest_mode,
		}
		return self.record

	def append_execution(self, ingest_id: int, entry: ExecutionEntry) -> IngestRecord:
		self.appended = (ingest_id, entry)
		return self.record


class DummyPersist:
	def __init__(self) -> None:
		self.recorded: ImageRecord | None = None

	def record(self, image: ImageRecord) -> None:
		self.recorded = image


class DummyPipeline:
	def __init__(
		self,
		media_root: Path,
		spec: Sequence[VariantLayerSpec],
		results: Sequence[VariantCommitResult],
	) -> None:
		self.media_root = media_root
		self.spec = spec
		self._results = results
		self.run_args: dict[str, object] | None = None

	def run(
		self,
		origin_relative_path: Path,
		file: OriginalFile,
		session: object,
	) -> Iterator[VariantCommitResult]:
		self.run_args = {
			'origin_relative_path': origin_relative_path,
			'file': file,
			'session': session,
		}
		return iter(self._results)


class FailingPipeline:
	def __init__(self, media_root: Path, spec: Sequence[VariantLayerSpec]) -> None:
		self.media_root = media_root
		self.spec = spec

	def run(
		self,
		origin_relative_path: Path,  # noqa: ARG002
		file: OriginalFile,  # noqa: ARG002
		session: object,  # noqa: ARG002
	) -> Iterator[VariantCommitResult]:
		raise ValueError('boom')


def test_image_ingest_service_records_image(tmp_path: Path) -> None:
	origin_relpath = Path('l0orig/sample.png')
	origin_path = tmp_path / origin_relpath
	origin_path.parent.mkdir(parents=True, exist_ok=True)
	PILImage.new('RGB', (10, 8), color='blue').save(origin_path)

	ingest_record = IngestRecord(
		id=5,
		relative_path=str(origin_relpath),
		fingerprint='f' * 64,
		captured_at=None,
	)

	spec = build_variant_spec(1, 320, container='webp', codecs='vp8')
	layer = VariantLayerSpec(name='primary', layer_id=1, specs=(spec,))
	variant_file = build_variant_file(spec, width=320)
	results = [VariantCommitResult.success('generate', VariantReport(spec, variant_file))]

	policy = VariantPolicy(
		regenerate_mismatched=False,
		generate_missing=True,
		delete_orphaned=False,
	)

	service = ImageIngestService(
		image_repo=object(),  # pyright: ignore[reportArgumentType]
		ingest_repo=object(),  # pyright: ignore[reportArgumentType]
		policy=policy,
	)
	service._ingest_core = DummyIngestCore(ingest_record)  # pyright: ignore[reportAttributeAccessIssue]
	service._persist = DummyPersist()  # pyright: ignore[reportAttributeAccessIssue]
	service._pipeline = DummyPipeline(tmp_path, [layer], results)  # pyright: ignore[reportAttributeAccessIssue]

	ingest = service.ingest(
		origin_path=origin_path,
		fingerprint=None,
		captured_at=None,
		ingest_mode=IngestMode.COPY,
	)

	assert ingest is ingest_record
	assert service._persist.recorded is not None  # pyright: ignore[reportAttributeAccessIssue, reportUnknownMemberType]

	image = cast(ImageRecord | None, service._persist.recorded)  # pyright: ignore[reportAttributeAccessIssue, reportUnknownMemberType]
	assert image is not None
	assert image.ingest_id == ingest_record.id
	assert image.original['rel'] == origin_relpath.__str__()
	assert image.original['format'] == 'png'
	assert image.original['width'] == 10
	assert image.original['height'] == 8
	assert image.original['size'] == origin_path.stat().st_size
	assert len(image.variants) == 1
	assert image.variants[0][0]['format'] == 'webp'

	appended = cast(tuple[int, ExecutionEntry] | None, service._ingest_core.appended)  # pyright: ignore[reportAttributeAccessIssue, reportUnknownMemberType]
	assert appended is not None
	ingest_id, entry = appended
	assert ingest_id == ingest_record.id
	assert entry['status'] == ExecutionStatus.SUCCESS


def test_image_ingest_service_records_failure_entry(tmp_path: Path) -> None:
	origin_relpath = Path('l0orig/sample.png')
	origin_path = tmp_path / origin_relpath
	origin_path.parent.mkdir(parents=True, exist_ok=True)
	PILImage.new('RGB', (10, 8), color='blue').save(origin_path)

	ingest_record = IngestRecord(
		id=7,
		relative_path=str(origin_relpath),
		fingerprint='f' * 64,
		captured_at=None,
	)

	policy = VariantPolicy(
		regenerate_mismatched=False,
		generate_missing=True,
		delete_orphaned=False,
	)

	service = ImageIngestService(
		image_repo=object(),  # pyright: ignore[reportArgumentType]
		ingest_repo=object(),  # pyright: ignore[reportArgumentType]
		policy=policy,
	)
	service._ingest_core = DummyIngestCore(ingest_record)  # pyright: ignore[reportAttributeAccessIssue]
	service._persist = DummyPersist()  # pyright: ignore[reportAttributeAccessIssue]
	service._pipeline = FailingPipeline(tmp_path, [])  # pyright: ignore[reportAttributeAccessIssue]

	with pytest.raises(ValueError, match='boom'):
		service.ingest(
			origin_path=origin_path,
			fingerprint=None,
			captured_at=None,
			ingest_mode=IngestMode.COPY,
		)

	appended = cast(tuple[int, ExecutionEntry] | None, service._ingest_core.appended)  # pyright: ignore[reportAttributeAccessIssue, reportUnknownMemberType]
	assert appended is not None
	ingest_id, entry = appended
	assert ingest_id == ingest_record.id
	assert entry['status'] == ExecutionStatus.UNKNOWN_ERROR
