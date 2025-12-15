from pathlib import Path

import pytest

from tests.services.images.utils import build_variant_spec

from app.core.variant_config import WEBP_FORMAT, VariantSlotkey, VariantSpec
from app.services.images.variants.commit import _delete_variant_file, prepare_variant_directories
from app.services.images.variants.types import ImageFileInfo, VariantFile, VariantPlan


def test_prepare_variant_directories_creates_missing_groups(tmp_path: Path) -> None:
	variant_root = tmp_path / 'l1w200'
	variant_root.mkdir()

	spec = build_variant_spec(1, width=200)
	diff = VariantPlan(matched=[], mismatched=[], missing=[spec], orphaned=[])

	prepare_variant_directories(
		diff,
		media_root=tmp_path,
		relpath_noext=Path('foo/bar'),
	)

	assert (variant_root / 'foo').is_dir()


def test_prepare_variant_directories_rejects_missing_root(tmp_path: Path) -> None:
	spec = build_variant_spec(1, width=200)
	diff = VariantPlan(matched=[], mismatched=[], missing=[spec], orphaned=[])

	with pytest.raises(RuntimeError):
		prepare_variant_directories(
			diff,
			media_root=tmp_path,
			relpath_noext=Path('foo/bar'),
		)


def test_prepare_variant_directories_prevents_escape(tmp_path: Path) -> None:
	variant_root = tmp_path / 'l1w200'
	variant_root.mkdir()

	spec = build_variant_spec(1, width=200)
	diff = VariantPlan(matched=[], mismatched=[], missing=[spec], orphaned=[])

	with pytest.raises(ValueError):
		prepare_variant_directories(
			diff,
			media_root=tmp_path,
			relpath_noext=Path('../escape'),
		)


def test_prepare_variant_directories_skips_root_path(tmp_path: Path) -> None:
	variant_root = tmp_path / 'l1w200'
	variant_root.mkdir()

	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=WEBP_FORMAT,
	)
	diff = VariantPlan(matched=[], mismatched=[], missing=[spec], orphaned=[])

	prepare_variant_directories(
		diff,
		media_root=tmp_path,
		relpath_noext=Path('image'),
	)

	assert list(variant_root.iterdir()) == []


def test_prepare_variant_directories_accepts_relative_media_root(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	static_root = tmp_path / 'static'
	variant_root = static_root / 'l1w200'
	variant_root.mkdir(parents=True)
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=WEBP_FORMAT,
	)
	diff = VariantPlan(matched=[], mismatched=[], missing=[spec], orphaned=[])

	monkeypatch.chdir(tmp_path)
	prepare_variant_directories(
		diff,
		media_root=Path('static'),
		relpath_noext=Path('foo/bar'),
	)

	assert (variant_root / 'foo').is_dir()


def _build_variant_file(file_path: Path) -> VariantFile:
	info = ImageFileInfo(
		file_path=file_path,
		container='webp',
		codecs='vp8',
		bytes=0,
		width=100,
		height=80,
		lossless=False,
	)
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=WEBP_FORMAT,
	)
	return VariantFile(variant_dir=spec.slotkey.label, relative_path=Path('foo/bar'), file_info=info)


def test_delete_variant_file_succeeds_when_file_exists(tmp_path: Path) -> None:
	file_path = tmp_path / 'foo.webp'
	file_path.write_bytes(b'data')
	variant_file = _build_variant_file(file_path)

	result = _delete_variant_file(variant_file)

	assert result.action == 'delete'
	assert result.result == 'success'
	assert not file_path.exists()


def test_delete_variant_file_returns_missing_when_file_absent(tmp_path: Path) -> None:
	file_path = tmp_path / 'missing.webp'
	variant_file = _build_variant_file(file_path)

	result = _delete_variant_file(variant_file)

	assert result.result == 'failure'
	assert result.reason == 'file_already_missing'
