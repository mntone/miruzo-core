from datetime import datetime
from pathlib import Path
from typing import final

from app.config.environments import env
from app.domain.clock.protocol import ClockProvider
from app.models.enums import ImageKind, IngestMode
from app.models.image import Image
from app.models.ingest import Ingest
from app.persist.stats.protocol import StatsCreateInput
from app.persist.uow import Repositories
from app.services.images.variants.executors.local import LocalVariantExecutor
from app.services.images.variants.mapper import (
	map_commit_results_to_variants,
	map_original_info_to_variant_record,
)
from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.pipeline import VariantPipeline
from app.services.images.variants.pipeline_execution import VariantPipelineExecutionSession
from app.services.images.variants.types import FileInfo, OriginalFile, VariantPolicy
from app.services.images.variants.utils import get_image_info_from_file
from app.services.ingests.service import IngestService


@final
class ImageIngestService:
	def __init__(
		self,
		repos: Repositories,
		*,
		clock: ClockProvider,
		policy: VariantPolicy,
		initial_score: int,
	) -> None:
		self._image_repo = repos.image
		self._stats_repo = repos.stats
		self._ingest_core = IngestService(
			repository=repos.ingest,
			clock=clock,
		)
		self._pipeline = VariantPipeline(
			media_root=env.media_root,
			policy=policy,
			spec=env.variant_layers,
		)
		self._initial_score = initial_score

	def ingest(
		self,
		*,
		origin_path: Path,
		fingerprint: str | None,
		captured_at: datetime,
		ingest_mode: IngestMode,
	) -> tuple[Ingest, Image]:
		ingest = self._ingest_core.create_ingest(
			origin_path=origin_path,
			fingerprint=fingerprint,
			captured_at=captured_at,
			ingest_mode=ingest_mode,
		)

		self._stats_repo.create(
			StatsCreateInput(
				ingest_id=ingest.id,
				initial_score=self._initial_score,
			),
		)

		executor = LocalVariantExecutor()
		session = VariantPipelineExecutionSession(executor)
		try:
			with session:
				with session.phase('inspect'):
					origin_relpath = VariantRelativePath(Path(ingest.relative_path))
					original_fileinfo = FileInfo.from_relative_path(
						origin_relpath,
						under=self._pipeline.media_root,
					)
					original_file = OriginalFile(
						file_info=original_fileinfo,
						image_info=get_image_info_from_file(original_fileinfo.absolute_path),
					)

				results = self._pipeline.run(origin_relpath, original_file, session)

				with session.phase('store'):
					original = map_original_info_to_variant_record(original_file)
					variants = map_commit_results_to_variants(results)

					image = Image(
						ingest_id=ingest.id,
						ingested_at=ingest.ingested_at,
						kind=ImageKind.UNSPECIFIED,
						original=original,
						fallback=None,
						variants=list(variants),
					)
					self._image_repo.create(image)
		finally:
			entry = session.to_dto()
			self._ingest_core.append_execution(ingest.id, entry)

		return ingest, image
