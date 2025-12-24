from collections.abc import Iterator, Sequence
from pathlib import Path

from app.config.variant import VariantLayerSpec
from app.services.images.variants.collect import (
	collect_variant_directories,
	collect_variant_files,
	normalize_media_relative_paths,
)
from app.services.images.variants.path import map_origin_to_variant_basepath
from app.services.images.variants.pipeline_execution import VariantPipelineExecutionSession
from app.services.images.variants.plan import build_variant_plan, emit_variant_specs
from app.services.images.variants.types import OriginalFile, VariantCommitResult, VariantPolicy


class VariantPipeline:
	def __init__(
		self,
		*,
		media_root: Path,
		policy: VariantPolicy,
		spec: Sequence[VariantLayerSpec],
	) -> None:
		self._media_root = media_root
		self._policy = policy
		self._spec = spec

	@property
	def media_root(self) -> Path:
		return self._media_root

	@property
	def policy(self) -> VariantPolicy:
		return self._policy

	@property
	def spec(self) -> Sequence[VariantLayerSpec]:
		return self._spec

	def run(
		self,
		origin_relative_path: Path,
		file: OriginalFile,
		session: VariantPipelineExecutionSession,
	) -> Iterator[VariantCommitResult]:
		variant_basepath = map_origin_to_variant_basepath(origin_relative_path)

		# collect
		with session.phase('collect'):
			variant_dirnames = collect_variant_directories(self._media_root)
			media_relpaths = normalize_media_relative_paths(variant_basepath, under=variant_dirnames)
			existing_files = collect_variant_files(media_relpaths, under=self._media_root)

		# plan
		with session.phase('plan'):
			planned_specs = emit_variant_specs(self._spec, file.image_info)
			plan = build_variant_plan(
				planned=planned_specs,
				existing=existing_files,
				rel_to=variant_basepath,
			)

		# dispatch
		results = session.execute(
			media_root=self._media_root,
			file=file,
			plan=plan,
			policy=self._policy,
		)

		return results
