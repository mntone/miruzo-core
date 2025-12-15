from pathlib import Path
from shutil import rmtree
from typing import Iterable

from PIL import Image as PILImage

from app.core.settings import settings
from app.core.variant_config import VariantLayer
from app.models.records import VariantRecord
from app.services.images.variants.legacy import VariantReportLegacy as VariantReport

__all__ = ['VariantReport']


def reset_variant_directories(media_root: Path, layers: Iterable[VariantLayer]) -> None:
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


def generate_variants(
	image: PILImage.Image,
	relative_path: Path,
	media_root: Path,
	layers: Iterable[VariantLayer],
	public_prefix: str | None = None,
	original_size: int | None = None,
) -> tuple[list[list[VariantRecord]], list[VariantReport]]:
	"""
	Backward-compatible shim that delegates to the variants.legacy implementation.
	"""
	from app.services.images.variants import legacy as variant_legacy

	return variant_legacy.generate_variants(
		image=image,
		relative_path=relative_path,
		media_root=media_root,
		layers=tuple(layers),
		public_prefix=public_prefix,
		original_size=original_size,
	)


def collect_existing_variants(
	relative_path: Path,
	media_root: Path,
	layers: Iterable[VariantLayer],
	public_prefix: str | None = None,
) -> list[list[VariantRecord]]:
	"""Inspect existing files and rebuild variant metadata without re-rendering."""

	public_prefix = public_prefix or settings.public_media_root

	layer_records: list[list[VariantRecord]] = []
	for layer in layers:
		variant_records: list[VariantRecord] = []
		for spec in layer.specs:
			relative_output = relative_path.with_suffix(spec.format.file_extension)
			output_dir = spec.slotkey.label
			output_path = media_root / output_dir / relative_output
			if not output_path.exists():
				continue

			try:
				with PILImage.open(output_path) as img:
					width = img.width
					height = img.height
			except Exception:
				continue

			filesize = output_path.stat().st_size

			public_path = f'{public_prefix}/{output_dir}/{relative_output}'
			variant_records.append(
				VariantRecord(
					filepath=public_path,
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
