from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Sequence

from PIL import Image as PILImage

from app.config.environments import env
from app.config.variant import VariantLayer
from app.models.records import VariantRecord
from app.services.images.variants.collect import collect_variant_directories, collect_variant_files
from app.services.images.variants.commit import commit_variant_plan, prepare_variant_directories
from app.services.images.variants.plan import (
	compare_variant_specs,
	normalize_variant_plan,
	plan_variant_specs,
)
from app.services.images.variants.preprocess import preprocess_original
from app.services.images.variants.security import validate_relative_path
from app.services.images.variants.types import (
	ImageFileInfo,
	OriginalImage,
	VariantCommitResult,
	VariantPolicy,
	VariantReport,
)
from app.services.images.variants.utils import get_image_format, validate_variant_slotkey

_LEGACY_LAYER_NAMES = {
	0: 'original',
	1: 'primary',
	2: 'secondary',
	3: 'tertiary',
	4: 'quaternary',
	5: 'quinary',
	6: 'senary',
	7: 'septenary',
	8: 'octonary',
	9: 'fallback',
}


@dataclass(frozen=True, slots=True)
class VariantReportLegacy:
	layer_name: str
	label: str
	width: int
	height: int
	size_bytes: int
	ratio_percent: float | None
	delta_bytes: int | None


def _map_records(
	results: Iterable[VariantCommitResult],
	*,
	relpath_noext: Path,
	public_prefix: str,
	original_size: int | None,
) -> tuple[list[tuple[int, VariantRecord]], list[VariantReportLegacy]]:
	records: list[tuple[int, VariantRecord]] = []
	legacy_reports: list[VariantReportLegacy] = []

	for result in results:
		if result.result != 'success':
			continue

		if result.action not in ('generate', 'regenerate'):
			continue

		spec = result.report.spec
		info = result.report.info  # only generate / regenerate results reach here

		filename = relpath_noext.with_suffix(spec.format.file_extension)
		public_path = f'{public_prefix}/{spec.slotkey.label}/{filename}'

		record = VariantRecord(
			filepath=public_path,
			format=spec.format.container,
			codecs=spec.format.codecs,
			size=info.bytes,
			width=info.width,
			height=info.height,
			lossless=info.lossless,
		)
		records.append((spec.layer_id, record))

		if original_size and original_size > 0:
			ratio = (info.bytes / original_size) * 100
			delta = info.bytes - original_size
		else:
			ratio = None
			delta = None

		legacy_name = _LEGACY_LAYER_NAMES.get(spec.layer_id, str(spec.layer_id))
		legacy_report = VariantReportLegacy(
			layer_name=legacy_name,
			label=f'w{spec.width}',
			width=info.width,
			height=info.height,
			size_bytes=info.bytes,
			ratio_percent=ratio,
			delta_bytes=delta,
		)
		legacy_reports.append(legacy_report)

	return records, legacy_reports


def _group_records_by_layer(
	records: Sequence[tuple[int, VariantRecord]],
	layers: Iterable[VariantLayer],
) -> list[list[VariantRecord]]:
	grouped: dict[int, list[VariantRecord]] = {layer.layer_id: [] for layer in layers}

	for layer_id, record in records:
		if layer_id in grouped:
			grouped[layer_id].append(record)

	return [grouped[layer.layer_id] for layer in layers if grouped[layer.layer_id]]


def _get_image_info(original: PILImage.Image) -> ImageFileInfo | None:
	file_path = Path(getattr(original, 'filename', None))
	try:
		stat = file_path.lstat()
	except FileNotFoundError:
		return None

	container, codecs, lossless = get_image_format(original)

	info = ImageFileInfo(
		file_path=file_path,
		container=container,
		codecs=codecs,
		bytes=stat.st_size,
		width=original.width,
		height=original.height,
		lossless=lossless,
	)
	return info


def generate_variants(
	image: PILImage.Image,
	relative_path: Path,
	media_root: Path,
	layers: Iterable[VariantLayer],
	public_prefix: str | None = None,
	original_size: int | None = None,
) -> tuple[list[list[VariantRecord]], list[VariantReport]]:
	"""Render thumbnails for all layers/specs and return DB-ready metadata."""

	# 0. security
	relpath_noext = relative_path.with_suffix('')
	validate_relative_path(relpath_noext)

	# 1. get info
	original_info = _get_image_info(image)
	if original_info is None:
		return [[]], []

	# 2. collect
	variant_dirs = collect_variant_directories(media_root=media_root)

	valid_variant_dirs = [d for d in variant_dirs if validate_variant_slotkey(d)]

	existing = collect_variant_files(
		media_root=media_root,
		variant_dirs=valid_variant_dirs,
		relative_path_noext=relpath_noext,
	)

	# 3. plan
	planned_specs = plan_variant_specs(layers, original_info)

	plan = compare_variant_specs(planned_specs, list(existing))

	normalized_plan = normalize_variant_plan(plan)

	# 4. preprocess
	preprocessed_image = OriginalImage(
		image=preprocess_original(image),
		info=original_info,
	)

	# 5. commit
	policy = VariantPolicy(
		regenerate_mismatched=True,
		generate_missing=True,
		delete_orphaned=True,
	)

	prepare_variant_directories(
		normalized_plan,
		media_root=media_root,
		relpath_noext=relpath_noext,
	)

	results = commit_variant_plan(
		normalized_plan,
		policy,
		preprocessed_image,
		media_root=media_root,
		relpath_noext=relpath_noext,
	)

	# 6. mapping
	public_prefix = public_prefix or env.public_media_root
	records, legacy_reports = _map_records(
		results,
		public_prefix=public_prefix,
		relpath_noext=relpath_noext,
		original_size=original_size,
	)
	by_layers = _group_records_by_layer(records, layers)

	return by_layers, legacy_reports
