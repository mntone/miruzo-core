from datetime import datetime

from app.models.records import JobRecord


class StubJobRepository:
	def __init__(self) -> None:
		self.jobs: dict[str, JobRecord] = {}
		self.get_called_with: list[str] = []
		self.started_called_with: list[datetime] = []
		self.finished_called_with: list[tuple[str, datetime]] = []

	def get_or_create(self, job_name: str) -> JobRecord:
		self.get_called_with.append(job_name)

		job_record = self.jobs.get(job_name)
		if job_record is not None:
			return job_record

		job_record = JobRecord(name=job_name)
		self.jobs[job_name] = job_record
		return job_record

	def mark_started(self, job: JobRecord, *, started_at: datetime) -> None:
		self.started_called_with.append(started_at)
		job.started_at = started_at

	def mark_finished(self, job_name: str, *, finished_at: datetime) -> None:
		self.finished_called_with.append((job_name, finished_at))

		job_record = self.jobs.get(job_name)
		if job_record is not None:
			job_record.finished_at = finished_at
		else:
			raise RuntimeError('job not found')
