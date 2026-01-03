from datetime import datetime, timedelta, timezone
from typing import Callable

from sqlmodel import Session

from app.jobs.protocol import Job
from app.models.records import JobRecord
from app.services.jobs.repository.protocol import JobRepository


class JobManager:
	def __init__(
		self,
		*,
		session_factory: Callable[[], Session],
		job_repo_factory: Callable[[Session], JobRepository],
		min_interval: timedelta,
	) -> None:
		self._session_factory = session_factory
		self._job_repo_factory = job_repo_factory
		self._min_interval = min_interval

	def _should_skip_run(self, job_record: JobRecord, current: datetime) -> bool:
		if job_record.started_at is None:
			return False

		interval = current - job_record.started_at
		return interval < self._min_interval

	def try_run(self, job: Job) -> bool:
		current = datetime.now(timezone.utc)

		with self._session_factory() as session:
			with session.begin():
				job_repo = self._job_repo_factory(session)

				job_record = job_repo.get_or_create(job.name)

				skip_run = self._should_skip_run(
					job_record=job_record,
					current=current,
				)
				if skip_run:
					return False

				job_repo.mark_started(job_record, started_at=current)

		job.run(evaluated_at=current)

		with self._session_factory() as session:
			with session.begin():
				job_repo = self._job_repo_factory(session)

				job_repo.mark_finished(job.name, finished_at=datetime.now(timezone.utc))

		return True
