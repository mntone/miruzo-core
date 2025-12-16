from collections.abc import Sequence

from app.config.variant import VariantLayer, VariantSpec
from app.services.images.variants.types import ImageFileInfo, VariantComparison, VariantFile, VariantPlan


def _should_emit_variant(spec: VariantSpec, original: ImageFileInfo) -> bool:
	return spec.required or spec.width < original.width


def plan_variant_specs(layers: Sequence[VariantLayer], original: ImageFileInfo) -> list[VariantSpec]:
	specs: list[VariantSpec] = []

	for layer in layers:
		for spec in layer.specs:
			if not _should_emit_variant(spec, original):
				continue

			specs.append(spec)

	return specs


def _is_content_matched(cmp: VariantComparison) -> bool:
	"""Check whether width, container, and codec attributes align."""

	if cmp.spec.width != cmp.file.file_info.width:
		return False

	if cmp.spec.format.container != cmp.file.file_info.container:
		return False

	if (cmp.spec.format.codecs is not None) and (cmp.spec.format.codecs != cmp.file.file_info.codecs):
		return False

	return True


def compare_variant_specs(
	planned: Sequence[VariantSpec],
	existing: Sequence[VariantFile],
) -> VariantPlan:
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

	return VariantPlan(
		matched=matched,
		mismatched=candidates,
		missing=missing_specs,
		orphaned=remaining_files,
	)


def normalize_variant_plan(diff: VariantPlan) -> VariantPlan:
	"""Move format-incompatible mismatches into orphaned files."""

	remaining_mismatched: list[VariantComparison] = list(diff.mismatched)
	orphaned = list(diff.orphaned)

	for c in range(len(remaining_mismatched) - 1, -1, -1):
		comparison = remaining_mismatched[c]
		fmt = comparison.spec.format
		info = comparison.file.file_info

		if info.container != fmt.container or info.codecs != fmt.codecs:
			del remaining_mismatched[c]
			orphaned.append(comparison.file)

	return VariantPlan(
		matched=diff.matched,
		mismatched=remaining_mismatched,
		missing=diff.missing,
		orphaned=orphaned,
	)
