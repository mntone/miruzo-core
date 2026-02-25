from abc import ABC, abstractmethod
from datetime import datetime
from typing import TypeVar

from sqlalchemy.exc import IntegrityError
from sqlmodel import Session, SQLModel

from app.models.records import JobRecord

TModel = TypeVar('TModel', bound=SQLModel)


class BaseJobRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	@abstractmethod
	def _is_unique_violation(self, error: IntegrityError) -> bool: ...

	def get_or_create(self, job_name: str) -> JobRecord:
		job = self._session.get(JobRecord, job_name)
		if job is not None:
			return job

		job = JobRecord(name=job_name)
		self._session.add(job)

		try:
			self._session.flush()
		except IntegrityError as exc:
			self._session.rollback()
			if not self._is_unique_violation(exc):
				raise
			job = self._session.get_one(JobRecord, job_name)

		return job

	def mark_started(self, job: JobRecord, *, started_at: datetime) -> None:
		job.started_at = started_at
		job.finished_at = None

		self._session.add(job)
		self._session.flush()

	def mark_finished(self, job_name: str, *, finished_at: datetime) -> None:
		job = self._session.get_one(JobRecord, job_name)
		job.finished_at = finished_at

		self._session.add(job)
		self._session.flush()
