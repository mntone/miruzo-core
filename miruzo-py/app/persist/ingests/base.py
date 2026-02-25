from collections.abc import Sequence
from datetime import datetime, timezone

from sqlmodel import Session

from app.config.constants import EXECUTION_MAXIMUM
from app.models.enums import ExecutionStatus, ProcessStatus, VisibilityStatus
from app.models.records import ExecutionEntry, IngestRecord


class BaseIngestRepository:
	def __init__(self, session: Session) -> None:
		self._session = session

	def create_ingest(
		self,
		*,
		relative_path: str,
		fingerprint: str,
		ingested_at: datetime,
		captured_at: datetime,
	) -> IngestRecord:
		"""Insert a new ingest record."""
		ingest = IngestRecord(
			relative_path=relative_path,
			fingerprint=fingerprint,
			ingested_at=ingested_at,
			captured_at=captured_at,
			updated_at=ingested_at,
			executions=None,
		)

		self._session.add(ingest)
		self._session.flush()
		self._session.refresh(ingest)

		return ingest

	def get_ingest(self, ingest_id: int) -> IngestRecord | None:
		"""Fetch an ingest record by its identifier."""
		ingest = self._session.get(IngestRecord, ingest_id)

		return ingest

	def append_execution(
		self,
		ingest_id: int,
		execution: ExecutionEntry,
	) -> IngestRecord | None:
		"""Append an execution entry and return the updated record."""
		ingest = self._session.get(IngestRecord, ingest_id)
		if ingest is None:
			return None

		executions: Sequence[ExecutionEntry]
		if ingest.executions is None:
			executions = []
		elif execution['status'] == ExecutionStatus.SUCCESS:
			executions = [e for e in ingest.executions if e['status'] != ExecutionStatus.SUCCESS]
		else:
			executions = list(ingest.executions)
		executions.append(execution)

		if execution['status'] == ExecutionStatus.SUCCESS:
			ingest.process = ProcessStatus.FINISHED

		ingest.executions = executions[-EXECUTION_MAXIMUM:]
		ingest.updated_at = datetime.now(timezone.utc)

		self._session.flush()
		self._session.refresh(ingest)

		return ingest

	def set_visibility(
		self,
		ingest_id: int,
		visibility: VisibilityStatus,
	) -> IngestRecord | None:
		"""Update the ingest visibility flag."""
		ingest = self._session.get(IngestRecord, ingest_id)
		if ingest is None:
			return None

		ingest.visibility = visibility

		self._session.flush()
		self._session.refresh(ingest)

		return ingest
