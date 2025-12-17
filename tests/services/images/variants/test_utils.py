import pytest
from PIL import TiffImagePlugin

from app.config.variant import VariantSlotkey
from app.services.images.variants.utils import _get_image_format, parse_variant_slotkey


class DummyImage:
	def __init__(self, fmt: str, *, info: dict | None = None) -> None:
		self.format = fmt
		self.info = info or {}


def _build_tiff_dummy(compression: int) -> TiffImagePlugin.TiffImageFile:
	image = object.__new__(TiffImagePlugin.TiffImageFile)
	image.format = 'TIFF'
	image.info = {}
	image.tag_v2 = {TiffImagePlugin.COMPRESSION: compression}
	return image


def test_get_image_format_handles_webp_lossless_flag() -> None:
	image = DummyImage('WEBP', info={'lossless': True})
	assert _get_image_format(image) == ('webp', 'vp8l', True, True)


def test_get_image_format_handles_webp_lossy_flag() -> None:
	image = DummyImage('WEBP', info={'lossless': False})
	assert _get_image_format(image) == ('webp', 'vp8', False, True)


def test_get_image_format_detects_png_as_lossless() -> None:
	image = DummyImage('PNG')
	assert _get_image_format(image) == ('png', None, True, False)


def test_get_image_format_detects_jpeg_as_lossy() -> None:
	image = DummyImage('JPEG')
	assert _get_image_format(image) == ('jpeg', None, False, True)


def test_get_image_format_detects_tiff_compression_lossless() -> None:
	image = _build_tiff_dummy(5)
	assert _get_image_format(image) == ('tiff', None, True, True)


def test_get_image_format_defaults_to_lowercased_container() -> None:
	image = DummyImage('ABC')
	assert _get_image_format(image) == ('abc', None, False, False)


def test_parse_variant_slotkey_parses_valid_label() -> None:
	slotkey = parse_variant_slotkey('l12w640')
	assert slotkey == VariantSlotkey(layer_id=12, width=640)


@pytest.mark.parametrize('label', ['lw200', 'l2w', 'l-1w200', 'l2w2x0', 'l2wfoo'])
def test_parse_variant_slotkey_raises_for_invalid_labels(label: str) -> None:
	with pytest.raises(ValueError):
		parse_variant_slotkey(label)
