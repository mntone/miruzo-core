from pathlib import Path

import pytest

from tests.services.images.variants.utils import build_jpeg_info

from app.services.images.variants.collect import (
	collect_variant_directories,
	collect_variant_files,
	normalize_media_relative_paths,
)
from app.services.images.variants.path import VariantRelativePath, map_origin_to_variant_basepath
from app.services.images.variants.types import FileInfo, VariantFile


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
	variant_root = tmp_path / 'l1w200'
	output_dir = variant_root / 'foo'
	output_dir.mkdir(parents=True)

	target = output_dir / 'bar.webp'
	target.write_bytes(b'data')

	called = {}

	def fake_load_variant_file(
		absolute_path: Path,
		relative_path: VariantRelativePath,
		variant_dirname: str,
	) -> VariantFile | None:
		called['absolute'] = absolute_path
		called['relative'] = relative_path
		info = build_jpeg_info(width=100, height=80)
		file_info = FileInfo(
			absolute_path=absolute_path,
			relative_path=relative_path,
			bytes=absolute_path.stat().st_size,
		)
		return VariantFile(
			file_info=file_info,
			image_info=info,
			variant_dir=variant_dirname,
		)

	monkeypatch.setattr(
		'app.services.images.variants.collect._load_variant_file',
		fake_load_variant_file,
	)

	basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))
	media_relpaths = list(normalize_media_relative_paths(basepath, under=['l1w200']))

	result = list(collect_variant_files(media_relpaths, under=tmp_path))

	assert len(result) == 1
	assert result[0].variant_dir == 'l1w200'
	assert called['absolute'] == target
	assert called['relative'] == Path('l1w200/foo/bar')


def test_normalize_media_relative_paths_filters_invalid() -> None:
	basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))

	valid = ['l1w200', 'l2w640']
	invalid = ['foo', 'l-1w200', 'l1wxyz']

	paths = list(normalize_media_relative_paths(basepath, under=valid + invalid))

	assert [str(path) for path in paths] == [f'{name}/foo/bar' for name in valid]
