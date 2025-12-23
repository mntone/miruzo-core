from pathlib import Path

from tests.services.images.utils import build_variant_spec

from app.config.variant import VariantLayerSpec, VariantSpec
from app.services.images.variants.mapper import map_commit_results_to_variant_layers
from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.types import (
	FileInfo,
	ImageInfo,
	OriginalFile,
	VariantCommitResult,
	VariantFile,
	VariantReport,
)


def _build_variant_file(
	tmp_path: Path,
	spec: VariantSpec,
	*,
	bytes: int,
	width: int,
	height: int,
	container: str,
	codecs: str | None,
) -> VariantFile:
	relative_path = VariantRelativePath(
		Path(f'{spec.slotkey.label}/foo{spec.format.file_extension}'),
	)
	absolute_path = tmp_path / relative_path
	file_info = FileInfo(
		absolute_path=absolute_path,
		relative_path=relative_path,
		bytes=bytes,
	)
	image_info = ImageInfo(
		container=container,
		codecs=codecs,
		width=width,
		height=height,
		lossless=False,
	)
	return VariantFile(
		file_info=file_info,
		image_info=image_info,
		variant_dir=spec.slotkey.label,
	)


def test_map_commit_results_to_variant_layers_filters_actions(tmp_path: Path) -> None:
	spec_primary = build_variant_spec(1, 320, container='webp', codecs='vp8')
	spec_fallback = build_variant_spec(9, 320, container='jpeg', codecs=None)

	file_primary = _build_variant_file(
		tmp_path,
		spec_primary,
		bytes=100,
		width=320,
		height=240,
		container='webp',
		codecs='vp8',
	)
	file_fallback = _build_variant_file(
		tmp_path,
		spec_fallback,
		bytes=200,
		width=320,
		height=240,
		container='jpeg',
		codecs=None,
	)

	results = [
		VariantCommitResult.success('generate', VariantReport(spec_primary, file_primary)),
		VariantCommitResult.success('regenerate', VariantReport(spec_fallback, file_fallback)),
		VariantCommitResult.failure('generate', 'save_failed'),
		VariantCommitResult.success('reuse', VariantReport(spec_primary, file_primary)),
	]
	layers = [
		VariantLayerSpec(name='primary', layer_id=1, specs=(spec_primary,)),
		VariantLayerSpec(name='fallback', layer_id=9, specs=(spec_fallback,)),
	]

	mapped = map_commit_results_to_variant_layers(results, layers)

	assert len(mapped) == 2
	assert mapped[0][0]['format'] == 'webp'
	assert mapped[0][0]['size'] == 100
	assert mapped[1][0]['format'] == 'jpeg'
	assert mapped[1][0]['size'] == 200


def test_map_commit_results_to_variant_layers_returns_empty_for_no_success(tmp_path: Path) -> None:
	spec_primary = build_variant_spec(1, 320, container='webp', codecs='vp8')
	layer = VariantLayerSpec(name='primary', layer_id=1, specs=(spec_primary,))

	file_primary = _build_variant_file(
		tmp_path,
		spec_primary,
		bytes=100,
		width=320,
		height=240,
		container='webp',
		codecs='vp8',
	)
	results = [
		VariantCommitResult.success('reuse', VariantReport(spec_primary, file_primary)),
		VariantCommitResult.failure('generate', 'save_failed'),
	]

	mapped = map_commit_results_to_variant_layers(results, [layer])

	assert mapped == []
