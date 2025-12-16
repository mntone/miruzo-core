from pathlib import Path

from app.config.variant import VariantSpec
from app.services.images.variants.types import ImageFileInfo, VariantFile


def build_image_info(
	*,
	width: int,
	container: str = 'png',
	codecs: str | None = None,
	lossless: bool = True,
) -> ImageFileInfo:
	return ImageFileInfo(
		file_path=Path('/tmp/source.png'),
		container=container,
		codecs=codecs,
		bytes=1024,
		width=width,
		height=width,
		lossless=lossless,
	)


def build_variant_file(
	spec: VariantSpec,
	*,
	width: int,
	container: str | None = None,
) -> VariantFile:
	container = container or spec.format.container
	info = ImageFileInfo(
		file_path=Path(f'/tmp/{spec.slotkey.label}.{container}'),
		container=container,
		codecs=spec.format.codecs,
		bytes=2048,
		width=width,
		height=width,
		lossless=False,
	)
	return VariantFile(
		variant_dir=spec.slotkey.label,
		file_info=info,
	)
