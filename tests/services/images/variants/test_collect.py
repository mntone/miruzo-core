from pathlib import Path

import pytest

from app.services.images.variants import collect
from app.services.images.variants.utils import ImageFileInfo


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

	collected = sorted(collect.collect_variant_directories(media_root=tmp_path))

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

	def fake_load_image_info(file_path: Path) -> ImageFileInfo:
		called['path'] = file_path
		return ImageFileInfo(
			file_path=file_path,
			container='webp',
			codecs='vp8',
			bytes=target.stat().st_size,
			width=100,
			height=80,
			lossless=False,
		)

	monkeypatch.setattr(collect, 'load_image_info', fake_load_image_info)

	result = list(
		collect.collect_variant_files(
			media_root=tmp_path,
			variant_dirs=['l1w200'],
			relative_path_noext=Path('foo/bar'),
		),
	)

	assert len(result) == 1
	assert result[0].variant_dir == 'l1w200'
	assert result[0].relative_path == Path('foo/bar')
	assert called['path'] == target
