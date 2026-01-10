from datetime import datetime
from typing import Protocol

from app.models.enums import VisibilityStatus
from app.models.records import IngestRecord
from app.models.types import ExecutionEntry


class IngestRepository(Protocol):
	def create_ingest(
		self,
		*,
		relative_path: str,
		fingerprint: str,
		ingested_at: datetime,
		captured_at: datetime,
	) -> IngestRecord: ...

	def get_ingest(self, ingest_id: int) -> IngestRecord | None: ...

	def append_execution(
		self,
		ingest_id: int,
		execution: ExecutionEntry,
	) -> IngestRecord | None: ...

	def set_visibility(
		self,
		ingest_id: int,
		visibility: VisibilityStatus,
	) -> IngestRecord | None: ...
