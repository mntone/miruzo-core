import io

from PIL import Image, ImageCms, ImageOps

from app.services.images.variants.utils import ImageInfo

_DEFAULT_BACKGROUND = (255, 255, 255, 255)


def _remove_alpha(
	image: Image.Image,
	background: tuple[int, int, int, int] = _DEFAULT_BACKGROUND,
) -> Image.Image:
	match image.mode:
		case 'LA':
			fill_color = background[0]
			output_mode = 'L'

		case 'La':
			image = image.convert('LA')
			fill_color = background[0]
			output_mode = 'L'

		case 'RGBA':
			fill_color = background
			output_mode = 'RGB'

		case 'RGBa':
			image = image.convert('RGBA')
			fill_color = background
			output_mode = 'RGB'

		case 'P' | 'PA' | 'RGB':
			if not image.has_transparency_data:
				return image

			# NOTE:
			# P / RGB images may carry transparency as palette alpha or color-key metadata
			# (e.g. info["transparency"]) rather than an actual alpha band.
			# In those cases there is no channel that can be used as a paste mask,
			# so we temporarily convert to RGBA to materialize transparency as an alpha band.
			# This path is intentionally limited to these modes to avoid unnecessary
			# memory overhead for images that already provide a usable alpha channel.
			image = image.convert('RGBA')
			output_image = Image.new('RGBA', image.size, background)
			output_image = Image.alpha_composite(output_image, image).convert('RGB')
			return output_image

		case _:
			return image

	output_image = Image.new(output_mode, image.size, fill_color)
	output_image.paste(image, mask=image.getchannel('A'))
	return output_image


def _convert_to_srgb(image: Image.Image) -> Image.Image:
	if image.mode != 'RGB':
		return image

	icc = image.info.get('icc_profile')
	if not icc:
		return image

	# profileToProfile mutates the image in place, so we return the same object.
	try:
		src_profile = ImageCms.ImageCmsProfile(io.BytesIO(icc))
		dst_profile = ImageCms.createProfile('sRGB')  # pyright: ignore[reportUnknownVariableType, reportUnknownMemberType]
		ImageCms.profileToProfile(
			image,
			src_profile,
			dst_profile,  # pyright: ignore[reportUnknownArgumentType]
			ImageCms.Intent.PERCEPTUAL,
			'RGB',
			True,
		)
	except (ImageCms.PyCMSError, OSError):
		# Keep the original pixels if color conversion fails.
		return image
	else:
		# Drop embedded profiles after successful conversion.
		image.info.pop('icc_profile', None)

	return image


# NOTE:
# Image generation policy is intentionally hard-coded for now.
# This pipeline prioritizes correctness and stability over abstraction.
def preprocess_original(original: Image.Image, info: ImageInfo) -> Image.Image:
	if info.supports_exif:
		normalized_image = ImageOps.exif_transpose(original)
	else:
		normalized_image = original

	opaque_image = _remove_alpha(normalized_image)

	srgb_image = _convert_to_srgb(opaque_image)

	return srgb_image
