from pathlib import Path

from app.config.variant import VariantSpec
from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.types import FileInfo, ImageInfo, VariantFile


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
	relative_path = VariantRelativePath(
		Path(f'{spec.slot.key}/foo/{spec.slot.key}.{container}'),
	)
	absolute_path = Path('/tmp').resolve() / relative_path
	file_info = FileInfo(
		absolute_path=absolute_path,
		relative_path=relative_path,
		bytes=2048,
	)
	return VariantFile(
		file_info=file_info,
		image_info=info,
		variant_dir=spec.slot.key,
	)
