from pathlib import Path

from tests.services.images.utils import build_variant_spec
from tests.services.images.variants.utils import build_png_info, build_variant_file

from app.config.variant import VariantLayerSpec
from app.services.images.variants.path import map_origin_to_variant_basepath
from app.services.images.variants.plan import (
	_classify_variant_diff,
	_compare_variant_specs,
	_prepare_variant_plan,
	_should_emit_variant,
	build_variant_plan,
	emit_variant_specs,
)
from app.services.images.variants.types import VariantComparison, VariantDiff


def test_should_emit_variant_respects_required_flag() -> None:
	spec_required = build_variant_spec(layer_id=1, width=1000, required=True)
	spec_optional = build_variant_spec(layer_id=1, width=1000, required=False)
	original = build_png_info(width=800)

	assert _should_emit_variant(spec_required, original) is True
	assert _should_emit_variant(spec_optional, original) is False


def test_emit_variant_specs_keeps_order_and_filters_by_width() -> None:
	layer = VariantLayerSpec(
		name='primary',
		layer_id=1,
		specs=(
			build_variant_spec(1, 320, required=True),
			build_variant_spec(1, 640),
			build_variant_spec(1, 960),
		),
	)
	original = build_png_info(width=700)

	result = emit_variant_specs([layer], original)

	assert [spec.width for spec in result] == [320, 640]


def test_compare_variant_specs_classifies_matches_mismatches_and_orphans() -> None:
	spec_match = build_variant_spec(1, 320)
	spec_mismatch = build_variant_spec(1, 640)
	spec_missing = build_variant_spec(1, 960)

	file_match = build_variant_file(spec_match, width=320)
	file_mismatch = build_variant_file(spec_mismatch, width=800)
	orphan_spec = build_variant_spec(2, 320)
	file_orphan = build_variant_file(orphan_spec, width=320)

	diff = _compare_variant_specs(
		[spec_match, spec_mismatch, spec_missing],
		[file_match, file_mismatch, file_orphan],
	)

	assert [cmp.expected_spec.slotkey for cmp in diff.matched] == [spec_match.slotkey]
	assert [cmp.expected_spec.slotkey for cmp in diff.mismatched] == [spec_mismatch.slotkey]
	assert [spec.slotkey for spec in diff.missing] == [spec_missing.slotkey]
	assert [file.slotkey for file in diff.orphaned] == [file_orphan.slotkey]


def test_compare_variant_specs_reports_multiple_files_per_spec() -> None:
	spec = build_variant_spec(1, 320)
	file_webp = build_variant_file(spec, width=320, container='webp')
	file_jpeg = build_variant_file(spec, width=320)

	diff = _compare_variant_specs([spec], [file_webp, file_jpeg])

	assert [cmp.expected_spec.slotkey for cmp in diff.matched] == [spec.slotkey]
	assert [cmp.expected_spec.slotkey for cmp in diff.mismatched] == [spec.slotkey]
	assert diff.missing == []
	assert diff.orphaned == []


def test_classify_variant_diff_demotes_format_mismatch() -> None:
	spec = build_variant_spec(1, 320, container='webp', codecs='vp8')
	file_incorrect_format = build_variant_file(spec, width=320, container='jpeg')
	diff = VariantDiff(
		matched=[],
		mismatched=[VariantComparison(spec, file_incorrect_format)],
		missing=[],
		orphaned=[],
	)

	classified = _classify_variant_diff(diff)

	assert classified.mismatched == []
	assert classified.orphaned == [file_incorrect_format]


def test_classify_variant_diff_keeps_valid_mismatch() -> None:
	spec = build_variant_spec(1, 640, container='webp', codecs='vp8')
	file_wrong_width = build_variant_file(spec, width=800)
	diff = VariantDiff(
		matched=[],
		mismatched=[VariantComparison(spec, file_wrong_width)],
		missing=[],
		orphaned=[],
	)

	classified = _classify_variant_diff(diff)

	assert classified.mismatched == [VariantComparison(spec, file_wrong_width)]
	assert classified.orphaned == []


def test_prepare_variant_plan_builds_plan_files() -> None:
	spec = build_variant_spec(1, 320)
	diff = VariantDiff(
		matched=[],
		mismatched=[],
		missing=[spec],
		orphaned=[],
	)
	variant_basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))

	plan = _prepare_variant_plan(diff, variant_basepath)

	assert len(plan.missing) == 1
	assert plan.missing[0].spec == spec
	assert plan.missing[0].path == Path('l1w320/foo/bar.jpeg')


def test_prepare_variant_plan_reuses_existing_relative_paths() -> None:
	spec = build_variant_spec(1, 320)
	file = build_variant_file(spec, width=320)
	diff = VariantDiff(
		matched=[],
		mismatched=[VariantComparison(spec, file)],
		missing=[],
		orphaned=[],
	)
	variant_basepath = map_origin_to_variant_basepath(Path('l0orig/foo/bar.webp'))

	plan = _prepare_variant_plan(diff, variant_basepath)

	assert len(plan.mismatched) == 1
	assert plan.mismatched[0].planning_file.path == file.file_info.relative_path


def test_build_variant_plan_produces_relative_paths() -> None:
	spec = build_variant_spec(1, 480)
	variant_basepath = map_origin_to_variant_basepath(Path('l0orig/foo/example.webp'))

	plan = build_variant_plan(
		planned=[spec],
		existing=[],
		rel_to=variant_basepath,
	)

	assert len(plan.missing) == 1
	assert plan.missing[0].path == Path('l1w480/foo/example.jpeg')
