import os
from collections.abc import Iterator

from app.services.images.variants.generate import generate_variant
from app.services.images.variants.types import (
	OriginalImage,
	VariantCommitResult,
	VariantFile,
	VariantPlan,
	VariantPolicy,
	VariantReport,
)


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
) -> Iterator[VariantCommitResult]:
	"""A result of applying a policy patch to a VariantDiff"""

	# 0. matched
	for cmp in plan.matched:
		report = VariantReport(cmp.expected_spec, cmp.actual_file)
		yield VariantCommitResult.success('reuse', report)

	# 1. missing
	if policy.generate_missing:
		for plan_file in plan.missing:
			report = generate_variant(plan_file, original)
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
				report = generate_variant(cmp.planning_file, original)
				if report is None:
					yield VariantCommitResult.failure('regenerate', 'save_failed')
				else:
					yield VariantCommitResult.success('regenerate', report)

	# 3. orphaned
	if policy.delete_orphaned:
		for file in plan.orphaned:
			report = _delete_variant_file(file)
			yield report
