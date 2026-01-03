from collections.abc import Iterable, Iterator

from app.config.variant import VariantLayerSpec, VariantSpec
from app.services.images.variants.path import VariantBasePath, build_variant_relative_path
from app.services.images.variants.types import (
	ImageInfo,
	VariantComparison,
	VariantDiff,
	VariantFile,
	VariantPlan,
	VariantPlanFile,
	VariantRegeneratePlan,
)


def _should_emit_variant(spec: VariantSpec, original: ImageInfo) -> bool:
	return spec.required or spec.width < original.width


def emit_variant_specs(layers: Iterable[VariantLayerSpec], original: ImageInfo) -> Iterator[VariantSpec]:
	for layer in layers:
		for spec in layer.specs:
			if _should_emit_variant(spec, original):
				yield spec


def _is_content_matched(cmp: VariantComparison) -> bool:
	"""Check whether width, container, and codec attributes align."""

	if cmp.expected_spec.width != cmp.actual_file.image_info.width:
		return False

	if cmp.expected_spec.format.container != cmp.actual_file.image_info.container:
		return False

	if (cmp.expected_spec.format.codecs is not None) and (
		cmp.expected_spec.format.codecs != cmp.actual_file.image_info.codecs
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

			if spec.slot == file.slot:
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
		info = comparison.actual_file.image_info

		if info.container != fmt.container or info.codecs != fmt.codecs:
			del remaining_mismatched[c]
			orphaned.append(comparison.actual_file)

	return VariantDiff(
		matched=diff.matched,
		mismatched=remaining_mismatched,
		missing=diff.missing,
		orphaned=orphaned,
	)


def _prepare_variant_plan(
	diff: VariantDiff,
	relative_path: VariantBasePath,
) -> VariantPlan:
	mismatched_plans: list[VariantRegeneratePlan] = []
	for cmp in diff.mismatched:
		regen_plan = VariantRegeneratePlan(
			actual_file=cmp.actual_file,
			planning_file=VariantPlanFile(
				path=cmp.actual_file.file_info.relative_path,
				spec=cmp.expected_spec,
			),
		)
		mismatched_plans.append(regen_plan)

	missing_plan_files: list[VariantPlanFile] = []
	for spec in diff.missing:
		variant_relpath = build_variant_relative_path(relative_path, under=spec.slot).with_suffix(
			spec.format.file_extension,
		)
		plan_file = VariantPlanFile(
			path=variant_relpath,
			spec=spec,
		)
		missing_plan_files.append(plan_file)

	plan = VariantPlan(
		matched=diff.matched,
		mismatched=mismatched_plans,
		missing=missing_plan_files,
		orphaned=diff.orphaned,
	)
	return plan


def build_variant_plan(
	*,
	planned: Iterable[VariantSpec],
	existing: Iterable[VariantFile],
	rel_to: VariantBasePath,
) -> VariantPlan:
	# Normalize argument name for internal use
	relative_path = rel_to

	diff = _compare_variant_specs(planned, existing)

	diff = _classify_variant_diff(diff)

	plan = _prepare_variant_plan(diff, relative_path)

	return plan
