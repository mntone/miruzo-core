from pathlib import Path

import pytest

from app.config.variant import DEFAULT_VARIANT_LAYERS
from app.services.ingests.bootstrap import (
	_ensure_media_root,
	_ensure_original_root,
	_ensure_variant_roots,
)


def _make_symlink(target: Path, link: Path) -> None:
	try:
		link.symlink_to(target, target_is_directory=True)
	except (OSError, NotImplementedError) as exc:
		pytest.skip(f'symlink not supported: {exc}')


def _collect_variant_dirnames() -> set[str]:
	return {spec.slot.key for layer in DEFAULT_VARIANT_LAYERS for spec in layer.specs}


def test_ensure_media_root_creates_directory(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'

	_ensure_media_root(media_root)

	assert media_root.is_dir()


def test_ensure_media_root_rejects_files(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.write_text('not a dir')

	with pytest.raises(RuntimeError, match='media_root must be a directory'):
		_ensure_media_root(media_root)


def test_ensure_media_root_rejects_symlinks(tmp_path: Path) -> None:
	target = tmp_path / 'media_target'
	target.mkdir()
	media_root = tmp_path / 'media'

	_make_symlink(target, media_root)

	with pytest.raises(RuntimeError, match='media_root must not be a symlink'):
		_ensure_media_root(media_root)


def test_ensure_original_root_creates_directory(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()

	_ensure_original_root(media_root)

	assert (media_root / 'l0orig').is_dir()


def test_ensure_original_root_rejects_files(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()
	(media_root / 'l0orig').write_text('not a dir')

	with pytest.raises(RuntimeError, match='Original directory must be a directory: l0orig'):
		_ensure_original_root(media_root)


def test_ensure_original_root_rejects_symlinks(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()
	target = tmp_path / 'original_target'
	target.mkdir()

	_make_symlink(target, media_root / 'l0orig')

	with pytest.raises(RuntimeError, match='Original directory must not be a symlink: l0orig'):
		_ensure_original_root(media_root)


def test_ensure_variant_roots_creates_directories(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()
	dirnames = _collect_variant_dirnames()

	_ensure_variant_roots(media_root, DEFAULT_VARIANT_LAYERS)

	for dirname in dirnames:
		assert (media_root / dirname).is_dir()


def test_ensure_variant_roots_rejects_files(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()
	dirname = next(iter(_collect_variant_dirnames()))
	(media_root / dirname).write_text('not a dir')

	with pytest.raises(RuntimeError, match='Variant directory must be a directory'):
		_ensure_variant_roots(media_root, DEFAULT_VARIANT_LAYERS)


def test_ensure_variant_roots_rejects_symlinks(tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()
	dirname = next(iter(_collect_variant_dirnames()))
	target = tmp_path / 'variant_target'
	target.mkdir()

	_make_symlink(target, media_root / dirname)

	with pytest.raises(RuntimeError, match='Variant directory must not be a symlink'):
		_ensure_variant_roots(media_root, DEFAULT_VARIANT_LAYERS)
