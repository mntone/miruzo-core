import re
from dataclasses import dataclass
from pathlib import Path

from PIL import Image as PILImage
from PIL import TiffImagePlugin

from app.config.variant import VariantSlotkey
from app.services.images.variants.path import VariantDirectoryPath

_TIFF_LOSSLESS_COMPRESSIONS = {
	1,  # No compression
	5,  # LZW
	8,  # Deflate (ZIP)
	32773,  # PackBits
}


@dataclass(frozen=True, slots=True)
class ImageInfo:
	container: str
	codecs: str | None
	width: int
	height: int
	lossless: bool


def _get_image_format(image: PILImage.Image) -> tuple[str, str | None, bool]:
	container = (image.format or '').upper()
	match container:
		case 'GIF':
			return 'gif', None, True
		case 'PNG':
			return 'png', None, True
		case 'JPEG':
			return 'jpeg', None, False
		case 'WEBP':
			lossless = bool(image.info.get('lossless'))
			return 'webp', 'vp8l' if lossless else 'vp8', lossless
		case 'BMP':
			return 'bmp', None, True
		case 'DIB':
			return 'dib', None, True
		case 'TIFF':
			assert isinstance(image, TiffImagePlugin.TiffImageFile), f'TIFF format, but got {type(image)}'
			lossless = image.tag_v2[TiffImagePlugin.COMPRESSION] in _TIFF_LOSSLESS_COMPRESSIONS
			return 'tiff', None, lossless
		case _:
			return container.lower(), None, False


def get_image_info(image: PILImage.Image) -> ImageInfo:
	container, codecs, lossless = _get_image_format(image)
	width = image.width
	height = image.height

	info = ImageInfo(
		container=container,
		codecs=codecs,
		width=width,
		height=height,
		lossless=lossless,
	)
	return info


def inspect_variant_subdir(
	relative_dirpath: Path,
	*,
	under: VariantDirectoryPath,
) -> Path | None:
	"""
	Inspect relative_dirpath under variant_root.

	Returns:
		Path: directory to mkdir
		None: mkdir should be skipped

	Raises:
		ValueError: if path escapes variant_root
	"""

	# Normalize argument name for internal use
	variant_root = under

	if relative_dirpath == Path('.'):
		return None

	variant_root = variant_root
	group_root = (variant_root / relative_dirpath).resolve()

	if not group_root.is_relative_to(variant_root):
		raise ValueError(f'Path escapes variant root: {group_root}')

	return group_root


def parse_variant_slotkey(label: str) -> VariantSlotkey:
	"""
	Parse slotkey and return structured representation.
	Raises ValueError if invalid.
	Returns structured, normalized representation.
	"""

	match = re.fullmatch(r'l(?P<layer>\d+)w(?P<width>\d+)', label)
	if not match:
		raise ValueError(f'Malformed variant slotkey: {label}')

	return VariantSlotkey(
		layer_id=int(match['layer']),
		width=int(match['width']),
	)
