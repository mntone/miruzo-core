from datetime import datetime
from pathlib import Path
from typing import cast

import pytest

from app.config.environments import env
from app.models.enums import IngestMode
from app.models.records import IngestRecord
from app.services.ingests.service import IngestService


class _StubRepository:
	def __init__(self, *, fail: bool = False) -> None:
		self.fail = fail
		self.created: dict[str, object] | None = None

	def create_ingest(
		self,
		*,
		relative_path: str,
		fingerprint: str,
		captured_at: datetime | None,
	) -> IngestRecord:
		self.created = {
			'relative_path': relative_path,
			'fingerprint': fingerprint,
			'captured_at': captured_at,
		}
		if self.fail:
			raise RuntimeError('boom')
		return IngestRecord(
			id=1,
			relative_path=relative_path,
			fingerprint=fingerprint,
			captured_at=captured_at,
		)


def _setup_roots(tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> Path:
	assets_root = tmp_path / 'assets'
	media_root = tmp_path / 'media'
	assets_root.mkdir()
	media_root.mkdir()

	monkeypatch.setattr(env, 'gataku_assets_root', assets_root)
	monkeypatch.setattr(env, 'media_root', media_root)
	monkeypatch.setattr(env, 'gataku_symlink_dirname', 'gataku')

	return assets_root


def test_create_ingest_copy_creates_file(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin = assets_root / 'foo' / 'bar.webp'
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	repo = _StubRepository()
	service = IngestService(repo)  # type: ignore[arg-type]

	ingest = service.create_ingest(
		origin_path=origin,
		fingerprint=None,
		captured_at=None,
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
	origin = assets_root / 'foo' / 'bar.webp'
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	repo = _StubRepository()
	service = IngestService(repo)  # type: ignore[arg-type]

	ingest = service.create_ingest(
		origin_path=origin,
		fingerprint=None,
		captured_at=None,
		ingest_mode=IngestMode.SYMLINK,
	)

	assert ingest.relative_path == 'gataku/foo/bar.webp'
	assert not (tmp_path / 'media' / 'gataku' / 'foo' / 'bar.webp').exists()


def test_create_ingest_raises_for_unknown_mode(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin = assets_root / 'foo' / 'bar.webp'
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	service = IngestService(_StubRepository())  # type: ignore[arg-type]

	with pytest.raises(ValueError, match='Unsupported ingest mode'):
		service.create_ingest(
			origin_path=origin,
			fingerprint=None,
			captured_at=None,
			ingest_mode=cast(IngestMode, 999),
		)


def test_create_ingest_copy_cleans_up_on_failure(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	origin = assets_root / 'foo' / 'bar.webp'
	origin.parent.mkdir(parents=True)
	origin.write_bytes(b'data')

	service = IngestService(_StubRepository(fail=True))  # type: ignore[arg-type]
	output_path = tmp_path / 'media' / 'l0orig' / 'foo' / 'bar.webp'

	with pytest.raises(RuntimeError, match='boom'):
		service.create_ingest(
			origin_path=origin,
			fingerprint=None,
			captured_at=None,
			ingest_mode=IngestMode.COPY,
		)

	assert not output_path.exists()
