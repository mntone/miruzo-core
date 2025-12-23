from collections.abc import Iterable, Sequence

from app.config.variant import VariantLayerSpec
from app.models.records import VariantRecord
from app.services.images.variants.types import OriginalFile, VariantCommitResult, VariantReport


def map_original_info_to_variant_record(file: OriginalFile) -> VariantRecord:
	"""Build the original image's record from a collected file."""
	record = VariantRecord(
		rel=file.file_info.relative_path.__str__(),
		format=file.image_info.container,
		codecs=file.image_info.codecs,
		size=file.file_info.bytes,
		width=file.image_info.width,
		height=file.image_info.height,
		quality=None,
	)

	return record


def map_commit_result_to_variant_record(report: VariantReport) -> VariantRecord:
	"""Build a variant record from a successful commit report."""
	spec = report.spec
	file = report.file.file_info  # only generate / regenerate results reach here
	image = report.file.image_info

	record = VariantRecord(
		rel=file.relative_path.__str__(),
		format=spec.format.container,
		codecs=spec.format.codecs,
		size=file.bytes,
		width=image.width,
		height=image.height,
		quality=spec.quality,
	)

	return record


def map_commit_results_to_variant_layers(
	results: Iterable[VariantCommitResult],
	layers: Iterable[VariantLayerSpec],
) -> Sequence[Sequence[VariantRecord]]:
	"""Group successful commit results into layer-ordered variant records."""
	variants: dict[int, list[VariantRecord]] = {layer.layer_id: [] for layer in layers}

	for result in results:
		if result.result != 'success':
			continue

		if result.action not in ('generate', 'regenerate'):
			continue

		assert result.report is not None
		report = result.report

		layer_id = report.spec.layer_id
		record = map_commit_result_to_variant_record(report)

		variants[layer_id].append(record)

	return [variants[layer.layer_id] for layer in layers if variants[layer.layer_id]]
