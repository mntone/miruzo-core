from pathlib import Path

import pytest
from PIL import Image as PILImage

from tests.services.images.variants.utils import build_png_info

from app.config.variant import VariantFormat, VariantSlotkey, VariantSpec
from app.services.images.variants.generate import _save_variant, generate_variant
from app.services.images.variants.types import OriginalImage


def test_save_variant_writes_jpeg(tmp_path: Path) -> None:
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=VariantFormat(container='jpeg', codecs=None, file_extension='.jpg'),
		quality=80,
	)
	image = PILImage.new('RGB', (50, 40), color='green')
	target = tmp_path / 'foo.jpg'

	file = _save_variant(spec, image, target)

	assert file is not None
	assert target.exists()
	assert file.bytes > 0
	assert file.info.container == 'jpeg'
	assert file.info.width == 50
	assert file.info.height == 40


def test_save_variant_raises_for_unsupported_format(tmp_path: Path) -> None:
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=1, width=200),
		layer_id=1,
		width=200,
		format=VariantFormat(container='gif', codecs=None, file_extension='.gif'),
	)
	image = PILImage.new('RGB', (10, 10))

	with pytest.raises(ValueError, match='Unsupported variant spec'):
		_save_variant(spec, image, tmp_path / 'foo.gif')


def test_generate_variant_writes_relative_path(tmp_path: Path) -> None:
	spec = VariantSpec(
		slotkey=VariantSlotkey(layer_id=3, width=480),
		layer_id=3,
		width=480,
		format=VariantFormat(container='jpeg', codecs=None, file_extension='.jpg'),
	)
	image = PILImage.new('RGB', (80, 60), color='blue')
	src_path = tmp_path / 'source.png'
	src_path.write_bytes(b'source')
	original = OriginalImage(
		image=image,
		info=build_png_info(width=80, height=60),
	)
	variant_root = tmp_path / spec.slotkey.label / 'foo'
	variant_root.mkdir(parents=True)

	report = generate_variant(
		spec,
		original,
		media_root=tmp_path,
		relpath_noext=Path('foo/bar'),
	)

	assert report is not None
	output_path = tmp_path / spec.slotkey.label / 'foo' / 'bar.jpg'
	assert report.file.path == output_path
	assert output_path.exists()
