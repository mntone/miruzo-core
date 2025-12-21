from pathlib import Path

import pytest

from app.config.environments import env
from app.services.ingests.utils.path import (
	map_origin_to_paths,
	map_origin_to_symlink_paths,
	validate_origin_path,
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


def test_validate_origin_path_accepts_files_under_assets(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	target = assets_root / 'foo' / 'bar.png'
	target.parent.mkdir(parents=True)
	target.write_bytes(b'data')

	resolved = validate_origin_path(target)

	assert resolved == target.resolve()


def test_validate_origin_path_rejects_directories(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	target = assets_root / 'foo'
	target.mkdir()

	with pytest.raises(ValueError, match='must be a file'):
		validate_origin_path(target)


def test_validate_origin_path_rejects_outside_assets(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	_setup_roots(tmp_path, monkeypatch)
	target = tmp_path / 'outside.txt'
	target.write_text('x')

	with pytest.raises(ValueError, match='escapes allowed root'):
		validate_origin_path(target)


def test_map_origin_to_paths_builds_copy_paths(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	target = assets_root / 'foo' / 'bar.png'
	target.parent.mkdir(parents=True)
	target.write_bytes(b'data')

	relative_path, output_path = map_origin_to_paths(target)

	assert relative_path == 'l0orig/foo/bar.png'
	assert output_path == tmp_path / 'media' / 'l0orig' / 'foo' / 'bar.png'


def test_map_origin_to_symlink_paths_builds_symlink_paths(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	target = assets_root / 'foo' / 'bar.png'
	target.parent.mkdir(parents=True)
	target.write_bytes(b'data')

	relative_path, output_path = map_origin_to_symlink_paths(target)

	assert relative_path == 'gataku/foo/bar.png'
	assert output_path == tmp_path / 'media' / 'gataku' / 'foo' / 'bar.png'


def test_map_origin_to_symlink_paths_uses_custom_dirname(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	assets_root = _setup_roots(tmp_path, monkeypatch)
	monkeypatch.setattr(env, 'gataku_symlink_dirname', 'shared')
	target = assets_root / 'foo' / 'bar.png'
	target.parent.mkdir(parents=True)
	target.write_bytes(b'data')

	relative_path, output_path = map_origin_to_symlink_paths(target)

	assert relative_path == 'shared/foo/bar.png'
	assert output_path == tmp_path / 'media' / 'shared' / 'foo' / 'bar.png'
