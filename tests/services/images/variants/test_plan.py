from tests.services.images.utils import build_variant_spec
from tests.services.images.variants.utils import build_image_info, build_variant_file

from app.core.variant_config import VariantLayer
from app.services.images.variants.plan import (
	_should_emit_variant,
	compare_variant_specs,
	normalize_variant_plan,
	plan_variant_specs,
)
from app.services.images.variants.types import VariantComparison, VariantPlan


def test_should_emit_variant_respects_required_flag() -> None:
	spec_required = build_variant_spec(layer_id=1, width=1000, required=True)
	spec_optional = build_variant_spec(layer_id=1, width=1000, required=False)
	original = build_image_info(width=800)

	assert _should_emit_variant(spec_required, original) is True
	assert _should_emit_variant(spec_optional, original) is False


def test_plan_variant_specs_keeps_order_and_filters_by_width() -> None:
	layer = VariantLayer(
		name='primary',
		layer_id=1,
		specs=(
			build_variant_spec(1, 320, required=True),
			build_variant_spec(1, 640),
			build_variant_spec(1, 960),
		),
	)
	original = build_image_info(width=700)

	result = plan_variant_specs([layer], original)

	assert [spec.width for spec in result] == [320, 640]


def test_compare_variant_specs_classifies_matches_mismatches_and_orphans() -> None:
	spec_match = build_variant_spec(1, 320)
	spec_mismatch = build_variant_spec(1, 640)
	spec_missing = build_variant_spec(1, 960)

	file_match = build_variant_file(spec_match, width=320)
	file_mismatch = build_variant_file(spec_mismatch, width=800)
	orphan_spec = build_variant_spec(2, 320)
	file_orphan = build_variant_file(orphan_spec, width=320)

	diff = compare_variant_specs(
		[spec_match, spec_mismatch, spec_missing],
		[file_match, file_mismatch, file_orphan],
	)

	assert [cmp.spec.slotkey for cmp in diff.matched] == [spec_match.slotkey]
	assert [cmp.spec.slotkey for cmp in diff.mismatched] == [spec_mismatch.slotkey]
	assert [spec.slotkey for spec in diff.missing] == [spec_missing.slotkey]
	assert [file.slotkey for file in diff.orphaned] == [file_orphan.slotkey]


def test_compare_variant_specs_reports_multiple_files_per_spec() -> None:
	spec = build_variant_spec(1, 320)
	file_webp = build_variant_file(spec, width=320)
	file_jpeg = build_variant_file(spec, width=320, container='jpeg')

	diff = compare_variant_specs([spec], [file_webp, file_jpeg])

	assert [cmp.spec.slotkey for cmp in diff.matched] == [spec.slotkey]
	assert [cmp.spec.slotkey for cmp in diff.mismatched] == [spec.slotkey]
	assert diff.missing == []
	assert diff.orphaned == []


def test_normalize_variant_diff_demotes_format_mismatch() -> None:
	spec = build_variant_spec(1, 320, container='webp', codecs='vp8')
	file_incorrect_format = build_variant_file(spec, width=320, container='jpeg')
	diff = VariantPlan(
		matched=[],
		mismatched=[VariantComparison(spec=spec, file=file_incorrect_format)],
		missing=[],
		orphaned=[],
	)

	normalized = normalize_variant_plan(diff)

	assert normalized.mismatched == []
	assert normalized.orphaned == [file_incorrect_format]


def test_normalize_variant_diff_keeps_valid_mismatch() -> None:
	spec = build_variant_spec(1, 640, container='webp', codecs='vp8')
	file_wrong_width = build_variant_file(spec, width=800)
	diff = VariantPlan(
		matched=[],
		mismatched=[VariantComparison(spec=spec, file=file_wrong_width)],
		missing=[],
		orphaned=[],
	)

	normalized = normalize_variant_plan(diff)

	assert normalized.mismatched == [VariantComparison(spec=spec, file=file_wrong_width)]
	assert normalized.orphaned == []
