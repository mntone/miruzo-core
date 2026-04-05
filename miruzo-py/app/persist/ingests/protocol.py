from datetime import datetime
from typing import Protocol

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

	def append_execution(
		self,
		ingest_id: int,
		execution: ExecutionEntry,
	) -> IngestRecord | None: ...
