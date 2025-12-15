from PIL import Image as PILImage

_DEFAULT_BACKGROUND = (255, 255, 255, 255)


def _remove_alpha(
	image: PILImage.Image,
	background: tuple[int] = _DEFAULT_BACKGROUND,
) -> PILImage.Image:
	match image.mode:
		case 'LA':
			output_mode = 'L'

		case 'La':
			image = image.convert('LA')
			output_mode = 'L'

		case 'RGBA':
			output_mode = 'RGB'

		case 'RGBa':
			image = image.convert('RGBA')
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
			output_image = PILImage.new('RGBA', image.size, background)
			output_image = PILImage.alpha_composite(output_image, image).convert('RGB')
			return output_image

		case _:
			return image

	output_image = PILImage.new(output_mode, image.size, background)
	output_image.paste(image, mask=image.getchannel('A'))
	return output_image


def preprocess_original(original: PILImage.Image) -> PILImage.Image:
	opaque_image = _remove_alpha(original)

	return opaque_image
