from datetime import datetime
from typing import Annotated, Protocol, final

from pydantic import BaseModel, Field

from app.config import constants as c
from app.models.ingest import Execution
from app.models.types import RelativePathType


@final
class IngestCreateInput(BaseModel):
	relative_path: RelativePathType
	fingerprint: Annotated[str, Field(min_length=64, max_length=64)]
	ingested_at: datetime
	captured_at: datetime


@final
class IngestAppendExecutionInput(BaseModel):
	ingest_id: Annotated[
		int,
		Field(ge=c.INGEST_ID_MINIMUM, le=c.INGEST_ID_MAXIMUM),
	]
	updated_at: datetime
	execution: Execution


class IngestRepository(Protocol):
	def create(self, entry: IngestCreateInput) -> int:
		"""Insert a new ingest row."""
		...

	def append_execution(self, entry: IngestAppendExecutionInput) -> None:
		"""Append an execution entry to an existing ingest row."""
		...
