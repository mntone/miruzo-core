from pathlib import Path

import pytest
from PIL import Image as PILImage

from tests.services.images.utils import build_variant_spec
from tests.services.images.variants.utils import build_png_info

from app.services.images.variants.generate import _save_variant, generate_variant
from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.types import OriginalImage, VariantPlanFile


def test_save_variant_writes_jpeg(tmp_path: Path) -> None:
	spec = build_variant_spec(1, 200, quality=80)
	image = PILImage.new('RGB', (50, 40), color='green')
	file = Path('foo.jpg')
	target = tmp_path / file

	file = _save_variant(
		spec,
		image,
		media_root=tmp_path,
		variant_relpath=VariantRelativePath(file),
		durable_write=False,
	)

	assert file is not None
	assert target.exists()
	assert file.file_info.bytes > 0
	assert file.image_info.container == 'jpeg'
	assert file.image_info.width == 50
	assert file.image_info.height == 40


def test_save_variant_raises_for_unsupported_format(tmp_path: Path) -> None:
	spec = build_variant_spec(1, 200, container='gif', codecs=None)
	image = PILImage.new('RGB', (10, 10))

	with pytest.raises(ValueError, match='Unsupported variant spec'):
		_save_variant(
			spec,
			image,
			media_root=tmp_path,
			variant_relpath=VariantRelativePath(Path('foo.gif')),
			durable_write=False,
		)


def test_generate_variant_writes_relative_path(tmp_path: Path) -> None:
	spec = build_variant_spec(3, 480)
	image = PILImage.new('RGB', (80, 60), color='blue')
	file = Path('source.png')
	src_path = tmp_path / file
	src_path.write_bytes(b'source')
	original = OriginalImage(
		image=image,
		info=build_png_info(width=80, height=60),
	)
	group_path = Path(spec.slot.key) / 'foo' / f'bar{spec.format.file_extension}'
	(tmp_path / group_path.parent).mkdir(parents=True)
	plan_file = VariantPlanFile(VariantRelativePath(group_path), spec)

	report = generate_variant(tmp_path, plan_file, original, durable_write=False)

	assert report is not None
	output_path = tmp_path / group_path
	assert report.file.file_info.absolute_path == output_path
	assert output_path.exists()
