from datetime import datetime
from logging import getLogger
from pathlib import Path
from typing import final

from app.domain.clock.protocol import ClockProvider
from app.models.enums import IngestMode
from app.models.ingest import Execution, Ingest
from app.persist.ingests.protocol import IngestAppendExecutionInput, IngestCreateInput, IngestRepository
from app.services.ingests.utils.file import copy_origin_file, delete_origin_file
from app.services.ingests.utils.fingerprint import compute_fingerprint, normalize_fingerprint
from app.services.ingests.utils.path import (
	map_relative_to_output_path,
	map_relative_to_pathstr,
	map_relative_to_symlink_pathstr,
	resolve_origin_absolute_path,
)

log = getLogger(__name__)


@final
class IngestService:
	def __init__(
		self,
		*,
		repository: IngestRepository,
		clock: ClockProvider,
	) -> None:
		self._repository = repository
		self._clock = clock

	def create_ingest(
		self,
		*,
		origin_path: Path,
		fingerprint: str | None,
		captured_at: datetime,
		ingest_mode: IngestMode,
	) -> Ingest:
		"""Create an ingest record and optionally persist the original asset."""

		origin_absolute_path = resolve_origin_absolute_path(origin_path)

		match ingest_mode:
			case IngestMode.SYMLINK:
				output_path = origin_absolute_path
				relative_path = map_relative_to_symlink_pathstr(origin_path)
			case IngestMode.COPY:
				output_path = map_relative_to_output_path(origin_path)
				copy_origin_file(origin_absolute_path, output_path)
				relative_path = map_relative_to_pathstr(origin_path)
			case _:
				raise ValueError(f'Unsupported ingest mode: {ingest_mode}')

		if fingerprint is None:
			fingerprint = compute_fingerprint(output_path)
		else:
			fingerprint = normalize_fingerprint(fingerprint)
			if fingerprint is None:
				log.warning('invalid fingerprint detected; recomputing for %s', origin_path)
				fingerprint = compute_fingerprint(output_path)

		now = self._clock.now()
		try:
			ingest_id = self._repository.create(
				IngestCreateInput(
					relative_path=relative_path,
					fingerprint=fingerprint,
					ingested_at=now,
					captured_at=captured_at,
				),
			)
		except Exception:
			if ingest_mode == IngestMode.COPY:
				delete_origin_file(output_path)
			raise

		return Ingest(
			id=ingest_id,
			relative_path=relative_path,
			fingerprint=fingerprint,
			ingested_at=now,
			captured_at=captured_at,
			updated_at=now,
			executions=[],
		)

	def append_execution(self, ingest_id: int, execution: Execution) -> None:
		"""Append an execution entry to the ingest record."""

		now = self._clock.now()
		self._repository.append_execution(
			IngestAppendExecutionInput(
				ingest_id=ingest_id,
				updated_at=now,
				execution=execution,
			),
		)
