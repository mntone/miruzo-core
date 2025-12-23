import os
from collections.abc import Iterator
from pathlib import Path

from app.services.images.variants.generate import generate_variant
from app.services.images.variants.path import build_absolute_path
from app.services.images.variants.types import (
	OriginalImage,
	VariantCommitResult,
	VariantFile,
	VariantPlan,
	VariantPlanFile,
	VariantPolicy,
	VariantReport,
)


def _delete_variant_file(
	file: VariantFile,
) -> VariantCommitResult:
	try:
		os.remove(file.file_info.absolute_path)

	except FileNotFoundError:
		return VariantCommitResult.failure('delete', 'file_already_missing')

	except PermissionError:
		return VariantCommitResult.failure('delete', 'permission_denied')

	except OSError:
		return VariantCommitResult.failure('delete', 'os_error')

	return VariantCommitResult.success('delete', None)


def _prepare_variant(media_root: Path, plan_file: VariantPlanFile) -> None:
	absolute_path = build_absolute_path(plan_file.path, under=media_root)
	absolute_path.parent.mkdir(parents=True, exist_ok=True)


def commit_variant_plan(
	*,
	plan: VariantPlan,
	policy: VariantPolicy,
	original: OriginalImage,
	media_root: Path,
) -> Iterator[VariantCommitResult]:
	"""A result of applying a policy patch to a VariantDiff"""

	if not media_root.is_dir():
		raise RuntimeError(f'media_root does not exist or is not a directory: {media_root}')

	# 0. matched
	for cmp in plan.matched:
		report = VariantReport(cmp.expected_spec, cmp.actual_file)
		yield VariantCommitResult.success('reuse', report)

	# 1. missing
	if policy.generate_missing:
		for plan_file in plan.missing:
			_prepare_variant(media_root, plan_file)
			report = generate_variant(media_root, plan_file, original)
			if report is None:
				yield VariantCommitResult.failure('generate', 'save_failed')
			else:
				yield VariantCommitResult.success('generate', report)

	# 2. mismatched
	if policy.regenerate_mismatched:
		for cmp in plan.mismatched:
			report = _delete_variant_file(cmp.actual_file)
			if report.result == 'failure':
				yield report
			else:
				_prepare_variant(media_root, cmp.planning_file)
				report = generate_variant(media_root, cmp.planning_file, original)
				if report is None:
					yield VariantCommitResult.failure('regenerate', 'save_failed')
				else:
					yield VariantCommitResult.success('regenerate', report)

	# 3. orphaned
	if policy.delete_orphaned:
		for file in plan.orphaned:
			report = _delete_variant_file(file)
			yield report
