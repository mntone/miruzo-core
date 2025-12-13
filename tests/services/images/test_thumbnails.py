from __future__ import annotations

from pathlib import Path
from typing import Iterator

import pytest
from PIL import Image as PILImage

from app.core.variant_config import VariantFormat, VariantLayer, VariantSpec
from app.services.images import thumbnails


class DummyImage:
	def __init__(self, fmt: str, info: dict | None = None) -> None:
		self.format = fmt
		self.info = info or {}


@pytest.fixture()
def tmp_image(tmp_path: Path) -> Iterator[PILImage.Image]:
	path = tmp_path / 'source.png'
	image = PILImage.new('RGBA', (400, 300), color=(255, 0, 0, 128))
	image.save(path, format='PNG')

	with PILImage.open(path) as img:
		yield img.copy()


def test_is_lossless_source_detects_lossless_webp() -> None:
	image = DummyImage('WEBP', info={'lossless': True})
	assert thumbnails._is_lossless_source(image) is True


def test_is_lossless_source_detects_lossy_webp() -> None:
	image = DummyImage('WEBP', info={'lossless': False})
	assert thumbnails._is_lossless_source(image) is False


def test_generate_variants_skips_unrequired_upscale(tmp_image: PILImage.Image, tmp_path: Path) -> None:
	static_root = tmp_path / 'static'
	static_root.mkdir()

	layer = VariantLayer(
		name='primary',
		layer_id=1,
		specs=(
			VariantSpec(label='w200', width=200, format=VariantFormat('webp', 'image/webp', '.webp')),
			VariantSpec(
				label='w600',
				width=600,
				format=VariantFormat('webp', 'image/webp', '.webp'),
				required=False,
			),
		),
	)

	variants, reports = thumbnails.generate_variants(
		image=tmp_image,
		relative_path=Path('foo/bar.png'),
		static_root=static_root,
		layers=(layer,),
		original_size=tmp_image.width * tmp_image.height,
	)

	assert len(variants) == 1
	assert len(variants[0]) == 1
	assert all(record['format'] == 'webp' and record['width'] == 200 for record in variants[0])
	assert all(report.label == 'w200' for report in reports)


def test_generate_variants_respects_required_flag(tmp_image: PILImage.Image, tmp_path: Path) -> None:
	static_root = tmp_path / 'static'
	static_root.mkdir()

	layer = VariantLayer(
		name='primary',
		layer_id=1,
		specs=(
			VariantSpec(
				label='w600',
				width=600,
				format=VariantFormat('webp', 'image/webp', '.webp'),
				required=True,
			),
		),
	)

	variants, reports = thumbnails.generate_variants(
		image=tmp_image,
		relative_path=Path('foo/bar.png'),
		static_root=static_root,
		layers=(layer,),
		original_size=tmp_image.width * tmp_image.height,
	)

	assert len(variants) == 1
	assert variants[0][0]['width'] == 600
	assert any(report.label == 'w600' for report in reports)
