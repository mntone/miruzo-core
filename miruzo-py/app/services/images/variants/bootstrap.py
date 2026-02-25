from PIL import Image as PILImage

from app.config.constants import MAX_IMAGE_PIXELS

_pillow_configured: bool = False


def configure_pillow() -> None:
	"""Configure Pillow globals for safe image decoding limits."""

	global _pillow_configured
	if _pillow_configured:
		return

	PILImage.MAX_IMAGE_PIXELS = MAX_IMAGE_PIXELS
	_pillow_configured = True
