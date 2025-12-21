from datetime import datetime
from pathlib import Path
from typing import final

from app.models.enums import IngestMode, VisibilityStatus
from app.models.records import ExecutionEntry, IngestRecord
from app.services.ingests.repository.base import IngestRepository
from app.services.ingests.utils.file import copy_origin_file, delete_origin_file
from app.services.ingests.utils.fingerprint import compute_fingerprint
from app.services.ingests.utils.path import (
	map_origin_to_paths,
	map_origin_to_symlink_paths,
	validate_origin_path,
)


@final
class IngestService:
	def __init__(self, repository: IngestRepository) -> None:
		self._repository = repository

	def create_ingest(
		self,
		*,
		origin_path: Path,
		fingerprint: str | None,
		captured_at: datetime | None,
		ingest_mode: IngestMode,
	) -> IngestRecord:
		"""Create an ingest record and optionally persist the original asset."""
		origin_path = validate_origin_path(origin_path)

		match ingest_mode:
			case IngestMode.SYMLINK:
				relative_path, output_path = map_origin_to_symlink_paths(origin_path)
			case IngestMode.COPY:
				relative_path, output_path = map_origin_to_paths(origin_path)
				copy_origin_file(origin_path, output_path)
			case _:
				raise ValueError(f'Unsupported ingest mode: {ingest_mode}')

		if fingerprint is None:
			fingerprint = compute_fingerprint(origin_path)

		try:
			ingest = self._repository.create_ingest(
				relative_path=relative_path,
				fingerprint=fingerprint,
				captured_at=captured_at,
			)
		except Exception:
			if ingest_mode == IngestMode.COPY:
				delete_origin_file(output_path)
			raise

		return ingest

	def get_ingest(self, ingest_id: int) -> IngestRecord | None:
		"""Fetch an ingest record by its identifier."""
		ingest = self._repository.get_ingest(ingest_id)

		return ingest

	def append_execution(self, ingest_id: int, execution: ExecutionEntry) -> IngestRecord | None:
		"""Append an execution entry to the ingest record."""
		ingest = self._repository.append_execution(ingest_id, execution)

		return ingest

	def set_visibility(self, ingest_id: int, visibility: VisibilityStatus) -> IngestRecord | None:
		"""Update visibility for a given ingest record."""
		ingest = self._repository.set_visibility(ingest_id, visibility)

		return ingest
