from collections.abc import Iterator, Sequence
from contextlib import contextmanager
from pathlib import Path

import pytest

from tests.services.images.utils import build_variant_spec

from app.config.variant import VariantLayerSpec, VariantSpec
from app.services.images.variants.path import (
	VariantBasePath,
	VariantRelativePath,
	map_origin_to_variant_basepath,
)
from app.services.images.variants.pipeline import VariantPipeline
from app.services.images.variants.types import (
	FileInfo,
	ImageInfo,
	OriginalFile,
	VariantCommitResult,
	VariantFile,
	VariantPlan,
	VariantPolicy,
)


class DummySession:
	def __init__(self) -> None:
		self.phases: list[str] = []
		self.execute_args: dict[str, object] = {}

	@contextmanager
	def phase(self, name: str) -> Iterator[None]:
		self.phases.append(name)
		yield None

	def execute(
		self,
		*,
		media_root: Path,
		file: OriginalFile,
		plan: VariantPlan,
		policy: VariantPolicy,
	) -> Sequence[VariantCommitResult]:
		self.execute_args = {
			'media_root': media_root,
			'file': file,
			'plan': plan,
			'policy': policy,
		}
		return []


def test_pipeline_run_builds_plan_and_executes(
	tmp_path: Path,
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	spec = build_variant_spec(1, 320, container='webp', codecs='vp8')
	layers = [VariantLayerSpec(name='primary', layer_id=1, specs=(spec,))]
	policy = VariantPolicy(
		regenerate_mismatched=False,
		generate_missing=True,
		delete_orphaned=False,
	)
	pipeline = VariantPipeline(media_root=tmp_path, policy=policy, spec=layers)

	origin_relative_path = Path('l0orig/foo/bar.webp')
	variant_basepath = map_origin_to_variant_basepath(origin_relative_path)
	relative_path = VariantRelativePath(Path('l0orig/foo/bar.webp'))
	absolute_path = tmp_path / relative_path
	file_info = FileInfo(
		absolute_path=absolute_path,
		relative_path=relative_path,
		bytes=123,
	)
	image_info = ImageInfo(
		container='webp',
		codecs='vp8',
		width=320,
		height=240,
		lossless=False,
	)
	original = OriginalFile(file_info=file_info, image_info=image_info)

	expected_media_relpaths = [VariantRelativePath(Path('l1w320/foo/bar'))]
	expected_plan = VariantPlan(matched=[], mismatched=[], missing=[], orphaned=[])

	def fake_collect_variant_directories(media_root: Path) -> list[str]:
		assert media_root == tmp_path
		return ['l1w320']

	def fake_normalize_media_relative_paths(
		relative_path: Path,
		*,
		under: list[str],
	) -> list[VariantRelativePath]:
		assert relative_path == variant_basepath
		assert under == ['l1w320']
		return expected_media_relpaths

	def fake_collect_variant_files(
		media_relpaths: list[VariantRelativePath],
		*,
		under: Path,
	) -> list[object]:
		assert media_relpaths == expected_media_relpaths
		assert under == tmp_path
		return []

	def fake_emit_variant_specs(
		layers_arg: list[VariantLayerSpec],
		image_info_arg: ImageInfo,
	) -> list[VariantSpec]:
		assert layers_arg == layers
		assert image_info_arg == image_info
		return [spec]

	def fake_build_variant_plan(
		*,
		planned: list[VariantSpec],
		existing: list[VariantFile],
		rel_to: VariantBasePath,
	) -> VariantPlan:
		assert list(planned) == [spec]
		assert list(existing) == []
		assert rel_to == variant_basepath
		return expected_plan

	monkeypatch.setattr(
		'app.services.images.variants.pipeline.collect_variant_directories',
		fake_collect_variant_directories,
	)
	monkeypatch.setattr(
		'app.services.images.variants.pipeline.normalize_media_relative_paths',
		fake_normalize_media_relative_paths,
	)
	monkeypatch.setattr(
		'app.services.images.variants.pipeline.collect_variant_files',
		fake_collect_variant_files,
	)
	monkeypatch.setattr(
		'app.services.images.variants.pipeline.emit_variant_specs',
		fake_emit_variant_specs,
	)
	monkeypatch.setattr(
		'app.services.images.variants.pipeline.build_variant_plan',
		fake_build_variant_plan,
	)

	session = DummySession()

	results = list(pipeline.run(origin_relative_path, original, session))  # pyright: ignore[reportArgumentType]

	assert results == []
	assert session.phases == ['collect', 'plan', 'execute']
	assert session.execute_args['media_root'] == tmp_path
	assert session.execute_args['file'] == original
	assert session.execute_args['plan'] == expected_plan
	assert session.execute_args['policy'] == policy
