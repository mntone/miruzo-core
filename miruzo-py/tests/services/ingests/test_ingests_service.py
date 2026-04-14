import logging
from datetime import datetime, timezone
from pathlib import Path
from typing import cast

import pytest

from tests.stubs.clock import FixedClockProvider

from app.config.environments import env
from app.models.enums import ExecutionStatus, IngestMode
from app.models.ingest import Execution
from app.persist.ingests.protocol import IngestAppendExecutionInput, IngestCreateInput
from app.services.ingests.service import IngestService
from app.services.ingests.utils.fingerprint import compute_fingerprint


class _StubIngestRepository:
	def __init__(self, *, fail: bool = False) -> None:
		self.fail = fail
		self.created: IngestCreateInput | None = None
		self.appended: IngestAppendExecutionInput | None = None

	def create(self, entry: IngestCreateInput) -> int:
		self.created = entry
		if self.fail:
			raise RuntimeError('boom')
		return 1

	def append_execution(self, entry: IngestAppendExecutionInput) -> None:
		self.appended = entry


def _setup_roots(tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> Path:
	assets_root = tmp_path / 'assets'
	media_root = tmp_path / 'media'
	assets_root.mkdir()
	media_root.mkdir()

	monkeypatch.setattr(env, 'gataku_assets_root', assets_root)
	monkeypatch.setattr(env, 'media_root', media_root)
	monkeypatch.setattr(env, 'gataku_symlink_dirname', 'gataku')

	return assets_root


def _new_ingest_service_fixture(
	now: datetime,
	*,
	fail: bool = False,
) -> tuple[IngestService, _StubIngestRepository]:
	repo = _StubIngestRepository(fail=fail)
	service = IngestService(
		repository=repo,
		clock=FixedClockProvider(now),
	)
	return service, repo


def test_create_ingest_copy_creates_file(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin_relative = Path('foo') / 'bar.webp'
	origin = assets_root / origin_relative
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	now = datetime(2026, 1, 10, 9, tzinfo=timezone.utc)
	service, repo = _new_ingest_service_fixture(now)

	ingest = service.create_ingest(
		origin_path=origin_relative,
		fingerprint=None,
		captured_at=now,
		ingest_mode=IngestMode.COPY,
	)

	assert ingest.relative_path == 'l0orig/foo/bar.webp'
	assert (tmp_path / 'media' / 'l0orig' / 'foo' / 'bar.webp').exists()
	assert repo.created is not None


def test_create_ingest_symlink_does_not_copy(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin_relative = Path('foo') / 'bar.webp'
	origin = assets_root / origin_relative
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	now = datetime(2026, 1, 10, 9, tzinfo=timezone.utc)
	service, _ = _new_ingest_service_fixture(now)

	ingest = service.create_ingest(
		origin_path=origin_relative,
		fingerprint=None,
		captured_at=now,
		ingest_mode=IngestMode.SYMLINK,
	)

	assert ingest.relative_path == 'gataku/foo/bar.webp'
	assert not (tmp_path / 'media' / 'gataku' / 'foo' / 'bar.webp').exists()


def test_create_ingest_raises_for_unknown_mode(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin_relative = Path('foo') / 'bar.webp'
	origin = assets_root / origin_relative
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	now = datetime(2026, 1, 10, 9, tzinfo=timezone.utc)
	service, _ = _new_ingest_service_fixture(now)

	with pytest.raises(ValueError, match='Unsupported ingest mode'):
		service.create_ingest(
			origin_path=origin_relative,
			fingerprint=None,
			captured_at=now,
			ingest_mode=cast(IngestMode, 999),
		)


def test_create_ingest_copy_cleans_up_on_failure(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin_relative = Path('foo') / 'bar.webp'
	origin = assets_root / origin_relative
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	now = datetime(2026, 1, 10, 9, tzinfo=timezone.utc)
	service, _ = _new_ingest_service_fixture(now, fail=True)
	output_path = tmp_path / 'media' / 'l0orig' / 'foo' / 'bar.webp'

	with pytest.raises(RuntimeError, match='boom'):
		service.create_ingest(
			origin_path=origin_relative,
			fingerprint=None,
			captured_at=datetime.now(timezone.utc),
			ingest_mode=IngestMode.COPY,
		)

	assert not output_path.exists()


def test_create_ingest_recomputes_invalid_fingerprint(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
	caplog: pytest.LogCaptureFixture,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin_relative = Path('foo') / 'bar.webp'
	origin = assets_root / origin_relative
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	now = datetime(2026, 1, 10, 9, tzinfo=timezone.utc)
	service, repo = _new_ingest_service_fixture(now)

	caplog.set_level(logging.WARNING)
	ingest = service.create_ingest(
		origin_path=origin_relative,
		fingerprint='not-a-hash',
		captured_at=now,
		ingest_mode=IngestMode.COPY,
	)

	assert ingest.relative_path == 'l0orig/foo/bar.webp'
	output_path = tmp_path / 'media' / 'l0orig' / 'foo' / 'bar.webp'
	assert repo.created is not None
	assert repo.created.fingerprint == compute_fingerprint(output_path)
	assert 'invalid fingerprint detected' in caplog.text


def test_append_execution_uses_clock() -> None:
	now = datetime(2026, 1, 10, 9, tzinfo=timezone.utc)
	service, repo = _new_ingest_service_fixture(now)
	execution = Execution(
		status=ExecutionStatus.SUCCESS,
		error_type=None,
		error_message=None,
		executed_at=now,
		inspect=None,
		collect=None,
		plan=None,
		execute=None,
		store=None,
		overall=None,
	)

	service.append_execution(ingest_id=10, execution=execution)
	assert repo.appended is not None
	assert repo.appended.ingest_id == 10
	assert repo.appended.updated_at == now
	assert repo.appended.execution == execution
