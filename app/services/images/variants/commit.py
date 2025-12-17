import os
from collections.abc import Iterator
from pathlib import Path

from app.services.images.variants.generate import generate_variant
from app.services.images.variants.types import (
	OriginalImage,
	VariantCommitResult,
	VariantFile,
	VariantPlan,
	VariantPolicy,
	VariantReport,
)


def prepare_variant_directories(
	diff: VariantPlan,
	*,
	media_root: Path,
	relpath_noext: Path,
) -> None:
	media_root_abs = media_root.resolve()

	group_roots: set[Path] = set()
	for spec in diff.missing:
		variant_root = media_root_abs / spec.slotkey.label
		if not variant_root.is_dir():
			raise RuntimeError(f'variant directory missing: {variant_root}')

		group_root = (variant_root / relpath_noext.parent).resolve()
		if group_root == variant_root:
			continue

		if not group_root.is_relative_to(variant_root):
			raise ValueError('Path escapes variant root')

		group_roots.add(group_root)

	for group_root in group_roots:
		group_root.mkdir(parents=True, exist_ok=True)


def _delete_variant_file(
	file: VariantFile,
) -> VariantCommitResult:
	try:
		os.remove(file.path)

	except FileNotFoundError:
		return VariantCommitResult.failure('delete', 'file_already_missing')

	except PermissionError:
		return VariantCommitResult.failure('delete', 'permission_denied')

	except OSError:
		return VariantCommitResult.failure('delete', 'os_error')

	return VariantCommitResult.success('delete', None)


def commit_variant_plan(
	plan: VariantPlan,
	policy: VariantPolicy,
	original: OriginalImage,
	*,
	media_root: Path,
	relpath_noext: Path,
) -> Iterator[VariantCommitResult]:
	"""A result of applying a policy patch to a VariantDiff"""

	# 0. matched
	for cmp in plan.matched:
		yield VariantCommitResult.success('reuse', VariantReport(cmp.spec, cmp.file))

	# 1. missing
	if policy.generate_missing:
		for cmp in plan.missing:
			report = generate_variant(
				cmp,
				original,
				media_root=media_root,
				relpath_noext=relpath_noext,
			)
			if report is None:
				yield VariantCommitResult.failure('generate', 'save_failed')
			else:
				yield VariantCommitResult.success('generate', report)

	# 2. mismatched
	if policy.regenerate_mismatched:
		for cmp in plan.mismatched:
			report = _delete_variant_file(cmp.file)
			if report.result == 'failure':
				yield report
			else:
				report = generate_variant(
					cmp.spec,
					original,
					media_root=media_root,
					relpath_noext=relpath_noext,
				)
				if report is None:
					yield VariantCommitResult.failure('regenerate', 'save_failed')
				else:
					yield VariantCommitResult.success('regenerate', report)

	# 3. orphaned
	if policy.delete_orphaned:
		for file in plan.orphaned:
			report = _delete_variant_file(file)
			yield report
