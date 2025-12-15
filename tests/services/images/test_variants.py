from __future__ import annotations

from typing import Iterable

from tests.services.images.utils import build_variant

from app.models.records import VariantRecord
from app.services.images.variants.api import compute_allowed_formats, normalize_variants_for_format


def _flatten_non_fallback(layers: Iterable[list[VariantRecord]]) -> list[VariantRecord]:
	return [variant for layer in layers[:-1] for variant in layer]


def test_normalizebuild_variants_filters_and_preserves_fallback() -> None:
	primary = [
		build_variant('webp', 320, label='primary'),
		build_variant('webp', 480, label='primary'),
		build_variant('jxl', 640, label='primary'),
		build_variant('jxl', 960, label='primary'),
	]
	secondary = [
		build_variant('webp', 640, label='secondary'),
		build_variant('webp', 960, label='secondary'),
	]
	fallback = [build_variant('jpeg', 320, label='fallback')]

	result = normalize_variants_for_format(
		[primary, secondary, fallback],
		{'webp'},
		keep_fallback=True,
	)

	allowed = _flatten_non_fallback(result)
	assert [(v['format'], v['width']) for v in allowed] == [
		('webp', 320),
		('webp', 480),
		('webp', 640),
		('webp', 960),
	]
	assert result[-1] == fallback


def test_keep_fallback_false_removes_last_layer_when_filtered() -> None:
	layers = [
		[build_variant('webp', 320, label='primary')],
		[build_variant('jpeg', 320, label='fallback')],
	]

	result = normalize_variants_for_format(layers, {'webp'}, keep_fallback=False)
	assert len(result) == 1
	assert result[0][0]['format'] == 'webp'


def test_empty_layers_are_removed_after_filtering() -> None:
	layers = [
		[build_variant('avif', 320, label='primary')],
		[build_variant('jpeg', 320, label='fallback')],
	]

	result = normalize_variants_for_format(layers, {'webp'}, keep_fallback=True)
	assert len(result) == 1
	assert result[0][0]['format'] == 'jpeg'


def test_compute_allowed_formats_always_keeps_fallback_formats() -> None:
	assert {'jpeg', 'png'}.issubset(compute_allowed_formats(None))
	assert {'jpeg', 'png'}.issubset(compute_allowed_formats([]))


def test_compute_allowed_formats_drops_excluded_items() -> None:
	result = compute_allowed_formats(['webp'])
	assert result == {'gif', 'jpeg', 'png'}


def test_compute_allowed_formats_handles_mixed_case_and_unknowns() -> None:
	result = compute_allowed_formats(['WEBP', 'avif', 'gif'])
	assert result == {'jpeg', 'png'}
