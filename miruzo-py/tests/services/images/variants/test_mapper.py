from collections.abc import Sequence
from pathlib import Path

from tests.services.images.utils import build_variant_spec

from app.config.variant import VariantLayerSpec, VariantSpec
from app.models.types import VariantEntry
from app.services.images.variants.mapper import (
	map_commit_results_to_variants,
	map_original_info_to_variant_record,
	map_variants_to_layers,
)
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
		Path(f'{spec.slot.key}/foo{spec.format.file_extension}'),
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
		variant_dir=spec.slot.key,
	)


def test_map_original_info_to_variant_record(tmp_path: Path) -> None:
	relative_path = VariantRelativePath(Path('l0orig/foo.webp'))
	absolute_path = tmp_path / relative_path
	file_info = FileInfo(
		absolute_path=absolute_path,
		relative_path=relative_path,
		bytes=123,
	)
	image_info = ImageInfo(
		container='webp',
		codecs='vp8',
		width=100,
		height=80,
		lossless=False,
	)
	original = OriginalFile(
		file_info=file_info,
		image_info=image_info,
	)

	record = map_original_info_to_variant_record(original)

	assert record == {
		'rel': str(relative_path),
		'layer_id': 0,
		'format': 'webp',
		'codecs': 'vp8',
		'bytes': 123,
		'width': 100,
		'height': 80,
		'quality': None,
	}


def test_map_commit_results_to_variant_layers_filters_actions(tmp_path: Path) -> None:
	spec_primary = build_variant_spec(1, 320, container='webp', codecs='vp8')
	spec_primary2 = build_variant_spec(1, 640, container='webp', codecs='vp8')
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
	file_primary2 = _build_variant_file(
		tmp_path,
		spec_primary2,
		bytes=200,
		width=640,
		height=480,
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
		VariantCommitResult.success('regenerate', VariantReport(spec_primary2, file_primary2)),
		VariantCommitResult.failure('generate', 'save_failed'),
		VariantCommitResult.success('reuse', VariantReport(spec_fallback, file_fallback)),
	]
	layers = [
		VariantLayerSpec(name='primary', layer_id=1, specs=(spec_primary, spec_primary2)),
		VariantLayerSpec(name='fallback', layer_id=9, specs=(spec_fallback,)),
	]

	entries = list(map_commit_results_to_variants(results))
	mapped = map_variants_to_layers(entries, spec=layers)

	assert len(mapped) == 2
	assert mapped[0][0]['format'] == 'webp'
	assert mapped[0][0]['bytes'] == 100
	assert mapped[0][1]['format'] == 'webp'
	assert mapped[0][1]['bytes'] == 200
	assert mapped[1][0]['format'] == 'jpeg'
	assert mapped[1][0]['bytes'] == 200


def test_map_commit_results_to_variants_returns_empty_for_no_success() -> None:
	results = [
		VariantCommitResult.success('delete', None),
		VariantCommitResult.failure('generate', 'save_failed'),
	]

	entries = list(map_commit_results_to_variants(results))

	assert entries == []


def test_map_variants_to_layers_ignores_unknown_layers() -> None:
	primary = VariantLayerSpec(name='primary', layer_id=1, specs=())
	variants: Sequence[VariantEntry] = [
		{
			'rel': 'l1w320/foo.webp',
			'layer_id': 1,
			'format': 'webp',
			'codecs': 'vp8',
			'bytes': 120,
			'width': 320,
			'height': 240,
			'quality': None,
		},
		{
			'rel': 'l2w320/bar.webp',
			'layer_id': 2,
			'format': 'webp',
			'codecs': 'vp8',
			'bytes': 220,
			'width': 320,
			'height': 240,
			'quality': None,
		},
	]

	mapped = map_variants_to_layers(variants, spec=[primary])

	assert mapped == [[variants[0]]]


def test_map_variants_to_layers_preserves_spec_order() -> None:
	primary = VariantLayerSpec(name='primary', layer_id=1, specs=())
	secondary = VariantLayerSpec(name='secondary', layer_id=2, specs=())
	variants: Sequence[VariantEntry] = [
		{
			'rel': 'l2w640/bar.webp',
			'layer_id': 2,
			'format': 'webp',
			'codecs': 'vp8',
			'bytes': 220,
			'width': 640,
			'height': 480,
			'quality': None,
		},
		{
			'rel': 'l1w320/foo.webp',
			'layer_id': 1,
			'format': 'webp',
			'codecs': 'vp8',
			'bytes': 120,
			'width': 320,
			'height': 240,
			'quality': None,
		},
	]

	mapped = map_variants_to_layers(variants, spec=[secondary, primary])

	assert mapped == [[variants[0]], [variants[1]]]


def test_map_variants_to_layers_sorts_widths() -> None:
	layer = VariantLayerSpec(name='primary', layer_id=1, specs=())
	variants: Sequence[VariantEntry] = [
		{
			'rel': 'l1w640/foo.webp',
			'layer_id': 1,
			'format': 'webp',
			'codecs': 'vp8',
			'bytes': 220,
			'width': 640,
			'height': 480,
			'quality': None,
		},
		{
			'rel': 'l1w320/bar.webp',
			'layer_id': 1,
			'format': 'webp',
			'codecs': 'vp8',
			'bytes': 120,
			'width': 320,
			'height': 240,
			'quality': None,
		},
	]

	mapped = map_variants_to_layers(variants, spec=[layer])

	assert [entry['width'] for entry in mapped[0]] == [320, 640]
