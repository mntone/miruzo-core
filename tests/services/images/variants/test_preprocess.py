import io
from pathlib import Path

import pytest
from PIL import ExifTags, Image

from tests.services.images.variants.utils import build_jpeg_info, build_png_info

from app.services.images.variants.preprocess import preprocess_original


@pytest.mark.parametrize(
	'mode,color,expected_mode,expected_pixel',
	[
		('RGBA', (0, 0, 0, 0), 'RGB', (255, 255, 255)),
		('RGBa', (0, 0, 0, 0), 'RGB', (255, 255, 255)),
		('LA', (0, 0), 'L', 255),
		('La', (0, 0), 'L', 255),
	],
)
def test_preprocess_original_strips_alpha_channel(
	mode: str,
	color: int | tuple[int, ...],
	expected_mode: str,
	expected_pixel: tuple[int, int, int] | int,
) -> None:
	base_mode = mode
	if mode == 'RGBa':
		base_mode = 'RGBA'
	if mode == 'La':
		base_mode = 'LA'

	image = Image.new(base_mode, (1, 1))
	image.putpixel((0, 0), color)
	if base_mode != mode:
		image = image.convert(mode)

	result = preprocess_original(image, build_png_info(width=1))

	assert result.mode == expected_mode
	assert result.getpixel((0, 0)) == expected_pixel


@pytest.mark.parametrize('orientation', [5, 6, 7, 8])
def test_preprocess_original_rotates_based_on_exif_orientation(orientation: int) -> None:
	base = Image.new('RGB', (2, 1), 'red')
	exif = Image.Exif()
	exif[ExifTags.Base.Orientation] = orientation  # orientation tag
	buffer = io.BytesIO()
	base.save(buffer, format='JPEG', exif=exif)
	buffer.seek(0)
	oriented = Image.open(buffer)

	info = build_jpeg_info(width=2, height=1)
	result = preprocess_original(oriented, info)

	assert result.size == (1, 2)


# Cyan/Magenta are excluded:
# relative colorimetric conversion changes RGB values
# due to different primaries between sRGB and Display P3.
@pytest.mark.parametrize(
	'color,expected_pixel',
	[
		((238, 18, 0), (255, 0, 0)),  # red
		((0, 252, 0), (0, 255, 0)),  # green
		((0, 6, 255), (0, 0, 255)),  # blue
		# ((0, 245, 250), (0, 255, 255)),  # cyan
		# ((250, 5, 245), (255, 0, 255)),  # magenta
		((255, 255, 0), (255, 255, 0)),  # yellow
		((0, 0, 0), (0, 0, 0)),  # black
		((128, 128, 128), (128, 128, 128)),  # gray
		((255, 255, 255), (255, 255, 255)),  # white
	],
)
def test_preprocess_original_converts_icc_profile(
	color: tuple[int, int, int],
	expected_pixel: tuple[int, int, int],
) -> None:
	image = Image.new('RGB', (1, 1), color)
	icc_path = Path(__file__).parent / 'icc' / 'display_p3.icc'
	image.info['icc_profile'] = icc_path.read_bytes()

	result = preprocess_original(image, build_png_info(width=1))

	assert 'icc_profile' not in result.info
	assert result.getpixel((0, 0)) == expected_pixel
