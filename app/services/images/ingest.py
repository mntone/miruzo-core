from datetime import datetime
from pathlib import Path

from app.config.environments import env
from app.models.enums import IngestMode
from app.models.records import ImageRecord, IngestRecord
from app.services.images.persist import ImagePersistService
from app.services.images.repository.base import ImageRepository
from app.services.images.variants.executors.local import LocalVariantExecutor
from app.services.images.variants.mapper import (
	map_commit_results_to_variant_layers,
	map_original_info_to_variant_record,
)
from app.services.images.variants.path import VariantRelativePath
from app.services.images.variants.pipeline import VariantPipeline
from app.services.images.variants.pipeline_execution import VariantPipelineExecutionSession
from app.services.images.variants.types import FileInfo, OriginalFile, VariantPolicy
from app.services.images.variants.utils import get_image_info_from_file
from app.services.ingests.repository.base import IngestRepository
from app.services.ingests.service import IngestService


class ImageIngestService:
	def __init__(
		self,
		image_repo: ImageRepository,
		ingest_repo: IngestRepository,
		policy: VariantPolicy,
	) -> None:
		self._ingest_core = IngestService(ingest_repo)
		self._persist = ImagePersistService(image_repo)
		self._pipeline = VariantPipeline(
			media_root=env.media_root,
			policy=policy,
			spec=env.variant_layers,
		)

	def ingest(
		self,
		*,
		origin_path: Path,
		fingerprint: str | None,
		captured_at: datetime | None,
		ingest_mode: IngestMode,
	) -> IngestRecord:
		ingest = self._ingest_core.create_ingest(
			origin_path=origin_path,
			fingerprint=fingerprint,
			captured_at=captured_at,
			ingest_mode=ingest_mode,
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
					variants = map_commit_results_to_variant_layers(results, self._pipeline.spec)

					image = ImageRecord(
						ingest_id=ingest.id,
						captured_at=captured_at,
						original=original,
						fallback=None,
						variants=variants,
					)

					self._persist.record(image)
		finally:
			entry = session.to_entry()
			self._ingest_core.append_execution(ingest.id, entry)

		return ingest
