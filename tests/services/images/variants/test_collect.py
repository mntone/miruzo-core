from pathlib import Path

import pytest

from tests.services.images.variants.utils import build_jpeg_info

from app.services.images.variants.collect import (
	collect_variant_directories,
	collect_variant_files,
	normalize_variant_directories,
)
from app.services.images.variants.path import normalize_relative_path
from app.services.images.variants.types import VariantFile


def test_collect_variant_directories_filters_symlinks_and_suffixes(tmp_path: Path) -> None:
	dirs = [
		tmp_path / 'l1w200',
		tmp_path / 'l2w400',
		tmp_path / 'l3orig',  # excluded because the suffix marks non-variants
	]
	for directory in dirs:
		directory.mkdir()

	# non-directory entry
	(tmp_path / 'l4w600.txt').write_text('ignore me', encoding='utf-8')

	# symlink entry
	symlink_dir = tmp_path / 'l5w800'
	symlink_dir.symlink_to(tmp_path / 'l1w200')

	collected = sorted(collect_variant_directories(tmp_path))

	assert collected == ['l1w200', 'l2w400']


def test_collect_variant_files_yields_existing_variants(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	variant_dir = tmp_path / 'l1w200'
	output_dir = variant_dir / 'foo'
	output_dir.mkdir(parents=True)

	target = output_dir / 'bar.webp'
	target.write_bytes(b'data')

	called = {}

	def fake_load_variant_file(file_path: Path, variant_dirname: str) -> VariantFile | None:
		called['path'] = file_path
		info = build_jpeg_info(width=100, height=80)
		return VariantFile(
			bytes=file_path.stat().st_size,
			info=info,
			path=file_path,
			variant_dir=variant_dirname,
		)

	monkeypatch.setattr(
		'app.services.images.variants.collect._load_variant_file',
		fake_load_variant_file,
	)

	variant_dirs = list(normalize_variant_directories(['l1w200'], under=tmp_path))
	relative = normalize_relative_path(Path('foo/bar.webp'))

	result = list(collect_variant_files(variant_dirs, rel_to=relative))

	assert len(result) == 1
	assert result[0].variant_dir == 'l1w200'
	assert called['path'] == target


def test_normalize_variant_directories_filters_invalid(tmp_path: Path) -> None:
	valid = ['l1w200', 'l2w640']
	invalid = ['foo', 'l-1w200', 'l1wxyz']

	dirs = list(normalize_variant_directories(valid + invalid, under=tmp_path))

	assert [path.name for path in dirs] == valid
