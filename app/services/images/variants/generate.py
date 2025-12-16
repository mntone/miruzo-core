from pathlib import Path

from PIL import Image as PILImage

from app.config.variant import VariantSpec
from app.services.images.variants.types import ImageFileInfo, OriginalImage, VariantReport


def _select_resample_algorithm(original: ImageFileInfo, target_width: int) -> int:
	"""Choose a resize filter based on scale ratio and source losslessness."""

	ratio = target_width / original.width
	if ratio > 1:
		return PILImage.BICUBIC
	if ratio >= 0.3:
		return PILImage.LANCZOS
	if original.lossless:
		return PILImage.HAMMING
	return PILImage.BOX


def _transform_variant(
	spec: VariantSpec,
	original: OriginalImage,
) -> PILImage.Image:
	"""Resize/copy the original image according to the spec."""

	width = spec.width
	height = max(1, int(round(width * (original.image.height / original.image.width))))
	resample = _select_resample_algorithm(original.info, width)

	variant_image = original.image.copy().resize((width, height), resample)
	return variant_image


def _save_variant(
	spec: VariantSpec,
	output_image: PILImage.Image,
	output_path: Path,
) -> ImageFileInfo | None:
	"""Encode the resized image and return filesystem metadata."""

	kwargs: dict[str, object] = {}
	match spec.format.container:
		case 'jpeg':
			lossless = False
			kwargs.setdefault('optimize', True)
			kwargs.setdefault('progressive', True)
		case 'webp':
			lossless = spec.format.codecs == 'vp8l'
			kwargs.setdefault('method', 6)
			kwargs.setdefault('lossless', lossless)
		case _:
			raise ValueError(f'Unsupported variant spec: {spec.format.container}')

	if spec.quality is not None:
		kwargs['quality'] = spec.quality

	pil_format = spec.format.container.upper()
	try:
		output_image.save(output_path, pil_format, **kwargs)
	except OSError:
		return None

	try:
		stat = output_path.lstat()
	except FileNotFoundError:
		return None

	info = ImageFileInfo(
		file_path=output_path,
		container=spec.format.container,
		codecs=spec.format.codecs,
		bytes=stat.st_size,
		width=output_image.width,
		height=output_image.height,
		lossless=lossless,
	)
	return info


def generate_variant(
	spec: VariantSpec,
	original: OriginalImage,
	*,
	media_root: Path,
	relpath_noext: Path,
) -> VariantReport | None:
	"""Render and persist a single variant, returning its report."""

	variant_image = _transform_variant(spec, original)

	filename = relpath_noext.with_suffix(spec.format.file_extension)
	variant_root = media_root / spec.slotkey.label / filename

	variant_info = _save_variant(spec, variant_image, variant_root)
	if variant_info is None:
		return None

	variant_report = VariantReport(spec, variant_info)

	return variant_report
