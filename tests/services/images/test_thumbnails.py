from collections.abc import Iterator
from pathlib import Path

import pytest
from PIL import Image as PILImage

from tests.services.images.utils import build_variant_spec

from app.config.variant import VariantLayer
from app.services.images import thumbnails


@pytest.fixture()
def tmp_image(tmp_path: Path) -> Iterator[PILImage.Image]:
	path = tmp_path / 'source.png'
	image = PILImage.new('RGBA', (400, 300), color=(255, 0, 0, 128))
	image.save(path, format='PNG')

	with PILImage.open(path) as img:
		image_copy = img.copy()
		yield image_copy


def test_generate_variants_skips_unrequired_upscale(tmp_image: PILImage.Image, tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()

	specs = (
		build_variant_spec(1, 200, required=True),
		build_variant_spec(1, 600, required=False),
	)
	for spec in specs:
		(media_root / spec.slotkey.label).mkdir(parents=True)

	layer = VariantLayer(name='primary', layer_id=1, specs=specs)

	variants, reports = thumbnails.generate_variants(
		image=tmp_image,
		relative_path=Path('foo/bar.png'),
		media_root=media_root,
		layers=(layer,),
		original_size=tmp_image.width * tmp_image.height,
	)

	assert len(variants) == 1
	assert len(variants[0]) == 1
	assert all(record['format'] == 'jpeg' and record['width'] == 200 for record in variants[0])
	assert all(report.label == 'w200' for report in reports)


def test_generate_variants_respects_required_flag(tmp_image: PILImage.Image, tmp_path: Path) -> None:
	media_root = tmp_path / 'media'
	media_root.mkdir()

	spec = build_variant_spec(1, 600, required=True)
	(media_root / spec.slotkey.label).mkdir(parents=True)

	layer = VariantLayer(name='primary', layer_id=1, specs=(spec,))

	variants, reports = thumbnails.generate_variants(
		image=tmp_image,
		relative_path=Path('foo/bar.png'),
		media_root=media_root,
		layers=(layer,),
		original_size=tmp_image.width * tmp_image.height,
	)

	assert len(variants) == 1
	assert variants[0][0]['width'] == 600
	assert any(report.label == 'w600' for report in reports)
