from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from shutil import rmtree
from typing import Iterable

from PIL import Image as PILImage

from app.core.variant_config import VariantLayer
from app.models.records import VariantRecord


@dataclass(frozen=True, slots=True)
class VariantReport:
	layer_name: str
	label: str
	width: int
	height: int
	size_bytes: int
	ratio_percent: float | None
	delta_bytes: int | None


def _is_lossless_source(image: PILImage.Image) -> bool:
	fmt = (image.format or '').upper()
	if fmt in {'PNG', 'TIFF', 'BMP'}:
		return True
	if fmt == 'WEBP':
		return bool(image.info.get('lossless'))
	return False


def _select_resample(image: PILImage.Image, target_width: int) -> int:
	ratio = target_width / image.width
	if ratio > 1:
		return PILImage.BICUBIC
	if ratio >= 0.3:
		return PILImage.LANCZOS
	if _is_lossless_source(image):
		return PILImage.HAMMING
	return PILImage.BOX


def _layer_dir(layer: VariantLayer, spec_label: str) -> Path:
	return Path(f'l{layer.layer_id}{spec_label}')


def reset_variant_directories(static_root: Path, layers: Iterable[VariantLayer]) -> None:
	"""Ensure every variant output directory exists and is empty."""
	seen: set[Path] = set()

	for layer in layers:
		for spec in layer.specs:
			dir_path = static_root / _layer_dir(layer, spec.label)
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
	static_root: Path,
	layers: Iterable[VariantLayer],
	public_prefix: str = '/static',
	original_size: int | None = None,
) -> tuple[list[list[VariantRecord]], list[VariantReport]]:
	"""Render thumbnails for all layers/specs and return DB-ready metadata."""
	if image.width <= 0 or image.height <= 0:
		return [], []

	layer_records: list[list[VariantRecord]] = []
	reports: list[VariantReport] = []
	for layer in layers:
		variant_records: list[VariantRecord] = []

		for spec in layer.specs:
			target_width = spec.width
			if not spec.required and target_width >= image.width:
				continue

			target_height = max(1, round(target_width * (image.height / image.width)))

			resized = image.copy()
			if target_width != image.width:
				resample = _select_resample(image, target_width)
				resized = resized.resize((target_width, target_height), resample=resample)

			if spec.format.name.lower() == 'jpeg':
				bands = resized.getbands()
				has_alpha = 'A' in bands or 'alpha' in {b.lower() for b in bands}
				has_transparency = 'transparency' in resized.info
				if has_alpha or has_transparency or resized.mode == 'P':
					resized = resized.convert('RGB')

			relative_output = relative_path.with_suffix(spec.format.extension)
			output_dir = _layer_dir(layer, spec.label)
			output_path = static_root / output_dir / relative_output
			output_path.parent.mkdir(parents=True, exist_ok=True)

			save_kwargs: dict[str, object] = {}
			if spec.quality is not None:
				save_kwargs['quality'] = spec.quality

			if spec.format.name.lower() == 'jpeg':
				save_kwargs.setdefault('optimize', True)
				save_kwargs.setdefault('progressive', True)
			elif spec.format.name.lower() == 'webp':
				save_kwargs.setdefault('method', 6)
				save_kwargs.setdefault('lossless', spec.format.lossless)

			resized.save(output_path, format=spec.format.name.upper(), **save_kwargs)
			filesize = output_path.stat().st_size

			public_path = f'{public_prefix}/{(output_dir / relative_output).as_posix()}'
			variant_records.append(
				VariantRecord(
					filepath=public_path,
					format=spec.format.name,
					codecs=spec.format.codecs,
					size=filesize,
					width=target_width,
					height=target_height,
					quality=spec.quality,
				),
			)

			ratio = None
			delta = None
			if original_size and original_size > 0:
				ratio = (filesize / original_size) * 100
				delta = filesize - original_size

			reports.append(
				VariantReport(
					layer_name=layer.name,
					label=spec.label,
					width=target_width,
					height=target_height,
					size_bytes=filesize,
					ratio_percent=ratio,
					delta_bytes=delta,
				),
			)

		layer_records.append(variant_records)

	return layer_records, reports


def collect_existing_variants(
	relative_path: Path,
	static_root: Path,
	layers: Iterable[VariantLayer],
	public_prefix: str = '/static',
) -> list[list[VariantRecord]]:
	"""Inspect existing files and rebuild variant metadata without re-rendering."""
	layer_records: list[list[VariantRecord]] = []

	for layer in layers:
		variant_records: list[VariantRecord] = []
		for spec in layer.specs:
			relative_output = relative_path.with_suffix(spec.format.extension)
			output_dir = _layer_dir(layer, spec.label)
			output_path = static_root / output_dir / relative_output
			if not output_path.exists():
				continue

			try:
				with PILImage.open(output_path) as img:
					width = img.width
					height = img.height
			except Exception:
				continue

			filesize = output_path.stat().st_size

			public_path = f'{public_prefix}/{(output_dir / relative_output).as_posix()}'
			variant_records.append(
				VariantRecord(
					filepath=public_path,
					format=spec.format.name,
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
