import re
from dataclasses import dataclass

from PIL import Image as PILImage
from PIL import TiffImagePlugin

from app.config.variant import VariantSlotkey

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
			lossless = image.info.get('lossless')
			return 'webp', 'vp8l' if lossless else 'vp8', lossless
		case 'BMP':
			return 'bmp', None, True
		case 'DIB':
			return 'dib', None, True
		case 'TIFF':
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
