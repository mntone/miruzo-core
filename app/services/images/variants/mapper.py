from collections.abc import Iterable, Sequence

from app.config.variant import VariantLayerSpec
from app.models.types import VariantEntry
from app.services.images.variants.types import OriginalFile, VariantCommitResult, VariantReport


def map_original_info_to_variant_record(file: OriginalFile) -> VariantEntry:
	"""Build the original image's record from a collected file."""
	record = VariantEntry(
		rel=file.file_info.relative_path.__str__(),
		layer_id=0,
		format=file.image_info.container,
		codecs=file.image_info.codecs,
		bytes=file.file_info.bytes,
		width=file.image_info.width,
		height=file.image_info.height,
		quality=None,
	)

	return record


def map_commit_result_to_variant_record(report: VariantReport) -> VariantEntry:
	"""Build a variant record from a successful commit report."""
	spec = report.spec
	file = report.file.file_info  # only generate / regenerate results reach here
	image = report.file.image_info

	record = VariantEntry(
		rel=file.relative_path.__str__(),
		layer_id=spec.layer_id,
		format=spec.format.container,
		codecs=spec.format.codecs,
		bytes=file.bytes,
		width=image.width,
		height=image.height,
		quality=spec.quality,
	)

	return record


def map_commit_results_to_variants(results: Iterable[VariantCommitResult]) -> Iterable[VariantEntry]:
	"""Group successful commit results into layer-ordered variant records."""

	for result in results:
		if result.result != 'success':
			continue

		if result.action not in ('reuse', 'generate', 'regenerate'):
			continue

		assert result.report is not None
		report = result.report

		entry = map_commit_result_to_variant_record(report)
		yield entry


def map_variants_to_layers(
	variants: Sequence[VariantEntry],
	*,
	spec: Sequence[VariantLayerSpec],
) -> Sequence[Sequence[VariantEntry]]:
	"""Group flat variant entries into layer-ordered sequences."""

	layered: dict[int, list[VariantEntry]] = {layer.layer_id: [] for layer in spec}

	for variant in variants:
		layer = layered.get(variant['layer_id'])
		if layer is None:
			continue
		layer.append(variant)

	for layer in layered.values():
		layer.sort(key=lambda entry: entry['width'])

	return [layered[layer.layer_id] for layer in spec if layered[layer.layer_id]]
