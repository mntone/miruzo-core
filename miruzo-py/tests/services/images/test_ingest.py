from collections.abc import Iterator, Sequence
from datetime import datetime, timezone
from pathlib import Path
from typing import cast

import pytest

from tests.fixtures.image_file import new_image_file_fixture
from tests.fixtures.ingest import make_ingest_fixture
from tests.services.images.utils import build_variant_spec
from tests.services.images.variants.utils import build_variant_file
from tests.stubs.image import StubImageRepository
from tests.stubs.stats import StubStatsRepository

from app.config.variant import VariantLayerSpec
from app.models.enums import ExecutionStatus, IngestMode
from app.models.ingest import Execution, Ingest
from app.persist.uow import Repositories
from app.services.images.ingest import ImageIngestService
from app.services.images.variants.types import (
	OriginalFile,
	VariantCommitResult,
	VariantPolicy,
	VariantReport,
)


class DummyIngestCore:
	def __init__(self, dto: Ingest) -> None:
		self.entry = dto
		self.created_args: dict[str, object] | None = None
		self.appended: tuple[int, Execution] | None = None

	def create_ingest(
		self,
		*,
		origin_path: Path,
		fingerprint: str | None,
		captured_at: datetime,
		ingest_mode: IngestMode,
	) -> Ingest:
		self.created_args = {
			'origin_path': origin_path,
			'fingerprint': fingerprint,
			'captured_at': captured_at,
			'ingest_mode': ingest_mode,
		}
		return self.entry

	def append_execution(self, ingest_id: int, entry: Execution) -> None:
		self.appended = (ingest_id, entry)


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


def new_image_ingest_service_fixture() -> ImageIngestService:
	policy = VariantPolicy(
		durable_write=False,
		regenerate_mismatched=False,
		generate_missing=True,
		delete_orphaned=False,
	)

	service = ImageIngestService(
		repos=Repositories(
			ingest=object(),  # pyright: ignore[reportArgumentType]
			image=StubImageRepository(),
			stats=StubStatsRepository(),
		),
		policy=policy,
		initial_score=100,
	)
	return service


def test_image_ingest_service_records_image(tmp_path: Path) -> None:
	ingest_id = 5
	image_pathes = new_image_file_fixture(tmp_path)
	ingest = make_ingest_fixture(ingest_id)

	spec = build_variant_spec(1, 320, container='webp', codecs='vp8')
	layer = VariantLayerSpec(name='primary', layer_id=1, specs=(spec,))
	variant_file = build_variant_file(spec, width=320)
	results = [VariantCommitResult.success('generate', VariantReport(spec, variant_file))]

	ingest_core = DummyIngestCore(ingest)
	service = new_image_ingest_service_fixture()
	service._ingest_core = ingest_core  # pyright: ignore[reportAttributeAccessIssue]
	service._pipeline = DummyPipeline(tmp_path, [layer], results)  # pyright: ignore[reportAttributeAccessIssue]

	entry = service.ingest(
		origin_path=image_pathes.relpath,
		fingerprint=None,
		captured_at=datetime.now(timezone.utc),
		ingest_mode=IngestMode.COPY,
	)

	assert entry[0] is ingest

	image = entry[1]
	assert image is not None
	assert image.ingest_id == ingest_id
	assert image.original['rel'] == image_pathes.relpath_str
	assert image.original['format'] == 'png'
	assert image.original['width'] == 10
	assert image.original['height'] == 8
	assert image.original['bytes'] == image_pathes.path.stat().st_size
	assert len(image.variants) == 1
	assert image.variants[0]['format'] == 'webp'

	stats = cast(StubStatsRepository, service._stats_repo).create_called_with
	assert stats is not None
	assert stats.ingest_id == ingest_id
	assert stats.initial_score == 100

	appended = ingest_core.appended
	assert appended is not None
	appended_ingest_id, entry = appended
	assert appended_ingest_id == ingest_id
	assert entry.status == ExecutionStatus.SUCCESS


def test_image_ingest_service_records_failure_entry(tmp_path: Path) -> None:
	ingest_id = 7
	image_pathes = new_image_file_fixture(tmp_path)
	ingest = make_ingest_fixture(ingest_id)

	service = new_image_ingest_service_fixture()
	ingest_core = DummyIngestCore(ingest)
	service._ingest_core = ingest_core  # pyright: ignore[reportAttributeAccessIssue]
	service._pipeline = FailingPipeline(tmp_path, [])  # pyright: ignore[reportAttributeAccessIssue]

	with pytest.raises(ValueError, match='boom'):
		service.ingest(
			origin_path=image_pathes.relpath,
			fingerprint=None,
			captured_at=datetime.now(timezone.utc),
			ingest_mode=IngestMode.COPY,
		)

	appended = ingest_core.appended
	assert appended is not None
	appended_ingest_id, entry = appended
	assert appended_ingest_id == ingest_id
	assert entry.status == ExecutionStatus.UNKNOWN_ERROR
