from datetime import datetime
from typing import Protocol

from app.models.records import JobRecord


class JobRepository(Protocol):
	def get_or_create(self, job_name: str) -> JobRecord: ...
	def mark_started(self, job: JobRecord, *, started_at: datetime) -> None: ...
	def mark_finished(self, job_name: str, *, finished_at: datetime) -> None: ...
