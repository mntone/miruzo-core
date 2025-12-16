from pathlib import Path

from app.config.variant import VariantSpec
from app.services.images.variants.types import ImageInfo, VariantFile


def build_jpeg_info(
	*,
	width: int,
	height: int | None = None,
) -> ImageInfo:
	return ImageInfo(
		container='jpeg',
		codecs=None,
		width=width,
		height=height or width,
		lossless=False,
	)


def build_png_info(
	*,
	width: int,
	height: int | None = None,
) -> ImageInfo:
	return ImageInfo(
		container='png',
		codecs=None,
		width=width,
		height=height or width,
		lossless=True,
	)


def build_webp_info(
	*,
	width: int,
	height: int | None = None,
) -> ImageInfo:
	return ImageInfo(
		container='webp',
		codecs='vp8',
		width=width,
		height=height or width,
		lossless=False,
	)


def build_variant_file(
	spec: VariantSpec,
	*,
	width: int,
	height: int | None = None,
	container: str | None = None,
) -> VariantFile:
	container = container or spec.format.container
	info = ImageInfo(
		container=container,
		codecs=spec.format.codecs,
		width=width,
		height=height or width,
		lossless=False,
	)
	return VariantFile(
		bytes=2048,
		info=info,
		path=Path(f'/tmp/{spec.slotkey.label}.{container}'),
		variant_dir=spec.slotkey.label,
	)
