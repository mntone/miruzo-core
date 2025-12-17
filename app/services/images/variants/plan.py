from collections.abc import Iterable, Iterator
from pathlib import Path

from app.config.variant import VariantLayer, VariantSpec
from app.services.images.variants.path import (
	NormalizedRelativePath,
	build_variant_dirpath,
	build_variant_filename,
)
from app.services.images.variants.types import (
	ImageInfo,
	VariantComparison,
	VariantDiff,
	VariantFile,
	VariantPlan,
	VariantPlanFile,
	VariantRegeneratePlan,
)
from app.services.images.variants.utils import inspect_variant_subdir


def _should_emit_variant(spec: VariantSpec, original: ImageInfo) -> bool:
	return spec.required or spec.width < original.width


def emit_variant_specs(layers: Iterable[VariantLayer], original: ImageInfo) -> Iterator[VariantSpec]:
	for layer in layers:
		for spec in layer.specs:
			if _should_emit_variant(spec, original):
				yield spec


def _is_content_matched(cmp: VariantComparison) -> bool:
	"""Check whether width, container, and codec attributes align."""

	if cmp.expected_spec.width != cmp.actual_file.info.width:
		return False

	if cmp.expected_spec.format.container != cmp.actual_file.info.container:
		return False

	if (cmp.expected_spec.format.codecs is not None) and (
		cmp.expected_spec.format.codecs != cmp.actual_file.info.codecs
	):
		return False

	return True


def _compare_variant_specs(
	planned: Iterable[VariantSpec],
	existing: Iterable[VariantFile],
) -> VariantDiff:
	"""Classify planned specs against on-disk variants."""

	candidates: list[VariantComparison] = []
	missing_specs: list[VariantSpec] = []
	remaining_files: list[VariantFile] = list(existing)

	for spec in planned:
		found_candidate = False

		for f in range(len(remaining_files) - 1, -1, -1):
			file = remaining_files[f]

			if spec.slotkey == file.slotkey:
				candidates.append(VariantComparison(spec, file))
				del remaining_files[f]
				found_candidate = True

		if not found_candidate:
			missing_specs.append(spec)

	matched: list[VariantComparison] = []
	for i in range(len(candidates) - 1, -1, -1):
		comparison = candidates[i]

		if _is_content_matched(comparison):
			del candidates[i]
			matched.append(comparison)

	return VariantDiff(
		matched=matched,
		mismatched=candidates,
		missing=missing_specs,
		orphaned=remaining_files,
	)


def _classify_variant_diff(diff: VariantDiff) -> VariantDiff:
	"""Move format-incompatible mismatches into orphaned files."""

	remaining_mismatched: list[VariantComparison] = list(diff.mismatched)
	orphaned = list(diff.orphaned)

	for c in range(len(remaining_mismatched) - 1, -1, -1):
		comparison = remaining_mismatched[c]
		fmt = comparison.expected_spec.format
		info = comparison.actual_file.info

		if info.container != fmt.container or info.codecs != fmt.codecs:
			del remaining_mismatched[c]
			orphaned.append(comparison.actual_file)

	return VariantDiff(
		matched=diff.matched,
		mismatched=remaining_mismatched,
		missing=diff.missing,
		orphaned=orphaned,
	)


def _prepare_variant_directories(paths: Iterable[Path]) -> None:
	for path in paths:
		path.mkdir(parents=True, exist_ok=True)


def _prepare_variant_plan(
	diff: VariantDiff,
	media_root: Path,
	relative_path: NormalizedRelativePath,
) -> tuple[VariantPlan, set[Path]]:
	target_variant_dirpaths: set[Path] = set()

	mismatched_plans: list[VariantRegeneratePlan] = []
	for cmp in diff.mismatched:
		plan_path = cmp.actual_file.path.with_suffix(cmp.expected_spec.format.file_extension)
		regen_plan = VariantRegeneratePlan(
			actual_file=cmp.actual_file,
			planning_file=VariantPlanFile(plan_path, cmp.expected_spec),
		)
		mismatched_plans.append(regen_plan)

	missing_plan_files: list[VariantPlanFile] = []
	for spec in diff.missing:
		variant_root = build_variant_dirpath(media_root, spec.slotkey.label)

		group_root = inspect_variant_subdir(relative_path.parent, under=variant_root)
		if group_root is not None:
			target_variant_dirpaths.add(group_root)
		else:
			# When normalized relative dir is ".", we keep files directly under the
			# variant root (no subdir), so mkdir is unnecessary.
			group_root = variant_root

		file_name = build_variant_filename(relative_path, spec)
		file_path = group_root / file_name
		plan_file = VariantPlanFile(file_path, spec)
		missing_plan_files.append(plan_file)

	plan = VariantPlan(
		matched=diff.matched,
		mismatched=mismatched_plans,
		missing=missing_plan_files,
		orphaned=diff.orphaned,
	)

	return plan, target_variant_dirpaths


def build_variant_plan(
	*,
	planned: Iterable[VariantSpec],
	existing: Iterable[VariantFile],
	rel_to: NormalizedRelativePath,
	under: Path,
) -> VariantPlan:
	# Normalize argument name for internal use
	relative_path = rel_to
	media_root = under

	diff = _compare_variant_specs(planned, existing)

	diff = _classify_variant_diff(diff)

	plan, variant_dirpaths = _prepare_variant_plan(diff, media_root, relative_path)

	_prepare_variant_directories(variant_dirpaths)

	return plan
