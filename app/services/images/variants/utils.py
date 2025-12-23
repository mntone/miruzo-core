import re
from dataclasses import dataclass
from pathlib import Path

from PIL import ExifTags, TiffImagePlugin
from PIL import Image as PILImage

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

	@property
	def supports_exif(self) -> bool:
		return self.container in ('jpeg', 'webp', 'tiff')


def _get_image_format(image: PILImage.Image) -> tuple[str, str | None, bool, bool]:
	container = (image.format or '').upper()
	match container:
		case 'GIF':
			return 'gif', None, True, False
		case 'PNG':
			return 'png', None, True, False
		case 'JPEG':
			return 'jpeg', None, False, True
		case 'WEBP':
			lossless = bool(image.info.get('lossless'))
			return 'webp', 'vp8l' if lossless else 'vp8', lossless, True
		case 'BMP':
			return 'bmp', None, True, False
		case 'DIB':
			return 'dib', None, True, False
		case 'TIFF':
			assert isinstance(image, TiffImagePlugin.TiffImageFile), f'TIFF format, but got {type(image)}'
			lossless = image.tag_v2[TiffImagePlugin.COMPRESSION] in _TIFF_LOSSLESS_COMPRESSIONS
			return 'tiff', None, lossless, True
		case _:
			return container.lower(), None, False, False


def get_image_info(image: PILImage.Image) -> ImageInfo:
	container, codecs, lossless, supports_exif = _get_image_format(image)
	width = image.width
	height = image.height

	if supports_exif:
		exif = image.getexif()
		orientation = exif.get(ExifTags.Base.Orientation)
		if orientation in (5, 6, 7, 8):
			width, height = height, width

	info = ImageInfo(
		container=container,
		codecs=codecs,
		width=width,
		height=height,
		lossless=lossless,
	)
	return info


def get_image_info_from_file(path: Path) -> ImageInfo:
	with PILImage.open(path) as image:
		info = get_image_info(image)
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
