from collections.abc import Iterable, Sequence
from dataclasses import dataclass
from pathlib import Path

from PIL import Image as PILImage

from app.config.variant import VariantLayerSpec
from app.models.records import VariantRecord
from app.services.images.variants.collect import (
	collect_variant_directories,
	collect_variant_files,
	normalize_media_relative_paths,
)
from app.services.images.variants.commit import commit_variant_plan
from app.services.images.variants.mapper import map_commit_results_to_variant_layers
from app.services.images.variants.path import build_origin_relative_path
from app.services.images.variants.plan import build_variant_plan, emit_variant_specs
from app.services.images.variants.preprocess import preprocess_original
from app.services.images.variants.types import (
	OriginalImage,
	VariantCommitResult,
	VariantPolicy,
)
from app.services.images.variants.utils import get_image_info

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


def _map_commit_results_to_legacy_reports(
	results: Iterable[VariantCommitResult],
	*,
	original_size: int | None,
) -> list[VariantReportLegacy]:
	legacy_reports: list[VariantReportLegacy] = []

	for result in results:
		if result.result != 'success':
			continue
		assert result.report is not None

		if result.action not in ('generate', 'regenerate'):
			continue

		spec = result.report.spec
		file = result.report.file.file_info  # only generate / regenerate results reach here
		image = result.report.file.image_info

		if original_size and original_size > 0:
			ratio = (file.bytes / original_size) * 100
			delta = file.bytes - original_size
		else:
			ratio = None
			delta = None

		legacy_name = _LEGACY_LAYER_NAMES.get(spec.layer_id, str(spec.layer_id))
		legacy_report = VariantReportLegacy(
			layer_name=legacy_name,
			label=f'w{spec.width}',
			width=image.width,
			height=image.height,
			size_bytes=file.bytes,
			ratio_percent=ratio,
			delta_bytes=delta,
		)
		legacy_reports.append(legacy_report)

	return legacy_reports


def generate_variants(
	image: PILImage.Image,
	relative_path: Path,
	media_root: Path,
	layers: Iterable[VariantLayerSpec],
	original_size: int | None = None,
) -> tuple[Sequence[Sequence[VariantRecord]], list[VariantReportLegacy]]:
	"""Render thumbnails for all layers/specs and return DB-ready metadata."""

	media_root = media_root.resolve()  # todo: add path validation
	origin_relpath = build_origin_relative_path(relative_path)
	original_info = get_image_info(image)

	# collect
	variant_dirnames = collect_variant_directories(media_root)
	media_relpaths = normalize_media_relative_paths(origin_relpath, under=variant_dirnames)
	existing_files = collect_variant_files(media_relpaths, under=media_root)

	# plan
	planned_specs = emit_variant_specs(layers, original_info)
	plan = build_variant_plan(
		planned=planned_specs,
		existing=existing_files,
		rel_to=origin_relpath,
	)

	# preprocess
	preprocessed_image = OriginalImage(
		image=preprocess_original(image, original_info),
		info=original_info,
	)

	# commit
	policy = VariantPolicy(
		regenerate_mismatched=True,
		generate_missing=True,
		delete_orphaned=True,
	)
	results = commit_variant_plan(
		plan=plan,
		policy=policy,
		original=preprocessed_image,
		media_root=media_root,
	)

	# mapping
	results = list(results)
	variants = map_commit_results_to_variant_layers(results, layers)
	legacy_reports = _map_commit_results_to_legacy_reports(
		results,
		original_size=original_size,
	)

	return variants, legacy_reports
