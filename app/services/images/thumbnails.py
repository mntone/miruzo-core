from collections.abc import Iterable
from pathlib import Path
from shutil import rmtree

from PIL import Image as PILImage

from app.config.variant import VariantLayerSpec
from app.models.records import VariantRecord
from app.services.images.variants.legacy import VariantReportLegacy as VariantReport
from app.services.images.variants.legacy import generate_variants

__all__ = ['generate_variants', 'VariantReport']


def reset_variant_directories(media_root: Path, layers: Iterable[VariantLayerSpec]) -> None:
	"""Ensure every variant output directory exists and is empty."""
	seen: set[Path] = set()

	for layer in layers:
		for spec in layer.specs:
			dir_path = media_root / spec.slotkey.label
			if dir_path in seen:
				continue

			if dir_path.exists():
				if dir_path.is_file():
					dir_path.unlink()
				else:
					rmtree(dir_path)

			dir_path.mkdir(parents=True, exist_ok=True)
			seen.add(dir_path)


def collect_existing_variants(
	relative_path: Path,
	media_root: Path,
	layers: Iterable[VariantLayerSpec],
) -> list[list[VariantRecord]]:
	"""Inspect existing files and rebuild variant metadata without re-rendering."""

	layer_records: list[list[VariantRecord]] = []
	for layer in layers:
		variant_records: list[VariantRecord] = []
		for spec in layer.specs:
			relative_output = relative_path.with_suffix(spec.format.file_extension)
			variant_dirname = spec.slotkey.label
			media_relative_path = variant_dirname / relative_output

			output_path = media_root / media_relative_path
			if not output_path.exists():
				continue

			try:
				with PILImage.open(output_path) as img:
					width = img.width
					height = img.height
			except Exception:
				continue

			filesize = output_path.stat().st_size

			variant_records.append(
				VariantRecord(
					rel=media_relative_path.__str__(),
					format=spec.format.container,
					codecs=spec.format.codecs,
					size=filesize,
					width=width,
					height=height,
					quality=spec.quality,
				),
			)

		if variant_records:
			layer_records.append(variant_records)

	return layer_records
