import logging
import re
from dataclasses import dataclass
from pathlib import Path
from typing import Literal, Self, TypeAlias

from PIL import Image as PILImage
from PIL import TiffImagePlugin
from PIL import UnidentifiedImageError as PILUnidentifiedImageError

from app.config.variant import VariantSlotkey

log = logging.getLogger(__name__)

_TIFF_LOSSLESS_COMPRESSIONS = {
	1,  # No compression
	5,  # LZW
	8,  # Deflate (ZIP)
	32773,  # PackBits
}

_ImageFileInfoReason: TypeAlias = Literal['success', 'not_found', 'invalid_image']


@dataclass(frozen=True, slots=True)
class ImageFileInfo:
	file_path: Path
	container: str
	codecs: str | None
	bytes: int
	width: int
	height: int
	lossless: bool


@dataclass(frozen=True, slots=True)
class _ImageInfoResult:
	file_info: ImageFileInfo | None
	reason: _ImageFileInfoReason

	@classmethod
	def success(cls, info: ImageFileInfo) -> Self:
		return cls(info, 'success')

	@classmethod
	def failure(cls, reason: _ImageFileInfoReason) -> Self:
		return cls(None, reason)


def get_image_format(image: PILImage.Image) -> tuple[str, str | None, bool]:
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


def _load_image_info_internal(file_path: Path) -> _ImageInfoResult:
	try:
		stat = file_path.lstat()
	except FileNotFoundError:
		return _ImageInfoResult.failure('not_found')

	try:
		with PILImage.open(file_path) as image:
			container, codecs, lossless = get_image_format(image)
			width = image.width
			height = image.height
	except FileNotFoundError:
		return _ImageInfoResult.failure('not_found')
	except PILUnidentifiedImageError:
		return _ImageInfoResult.failure('invalid_image')

	image_info = ImageFileInfo(
		file_path=file_path,
		container=container,
		codecs=codecs,
		bytes=stat.st_size,
		width=width,
		height=height,
		lossless=lossless,
	)
	return _ImageInfoResult.success(image_info)


def load_image_info(file_path: Path) -> ImageFileInfo | None:
	result = _load_image_info_internal(file_path)
	match result.reason:
		case 'success':
			return result.file_info
		case 'not_found':
			log.debug('image not found: %s', file_path)
		case 'invalid_image':
			log.warning('invalid image: %s', file_path)

	return None


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
