from datetime import datetime, timedelta, timezone
from types import SimpleNamespace
from typing import Callable

import pytest

from tests.jobs.stub import StubJob
from tests.services.jobs.stubs import StubJobRepository
from tests.stubs.ref import RefInt
from tests.stubs.session import create_stub_session_factory

import app.services.jobs.manager as manager
from app.models.records import JobRecord
from app.services.jobs.manager import JobManager


def test_try_run_executes_job_and_marks_records(monkeypatch: pytest.MonkeyPatch) -> None:
	started_at = datetime(2026, 1, 1, 0, 0, tzinfo=timezone.utc)
	finished_at = datetime(2026, 1, 1, 0, 1, tzinfo=timezone.utc)
	times = [started_at, finished_at]

	def now(_: object | None = None) -> datetime:
		return times.pop(0)

	monkeypatch.setattr(manager, 'datetime', SimpleNamespace(now=now))

	job = StubJob()
	record = JobRecord(name=job._NAME)
	repo = StubJobRepository()
	repo.jobs = {job.name: record}
	session_calls = RefInt(0)

	manager_instance = JobManager(
		session_factory=create_stub_session_factory(session_calls),  # pyright: ignore[reportArgumentType]
		job_repo_factory=lambda _: repo,
		min_interval=timedelta(minutes=10),
	)

	ran = manager_instance.try_run(job)

	assert ran is True
	assert session_calls.value == 2
	assert repo.get_called_with == [job.name]
	assert repo.started_called_with == [started_at]
	assert repo.finished_called_with == [(job.name, finished_at)]
	assert job.run_called_with == [started_at]
	assert record.started_at == started_at
	assert record.finished_at == finished_at


def test_try_run_skips_recent_run(monkeypatch: pytest.MonkeyPatch) -> None:
	current = datetime(2026, 1, 1, 0, 0, tzinfo=timezone.utc)
	monkeypatch.setattr(manager, 'datetime', SimpleNamespace(now=lambda _: current))

	job = StubJob()
	record = JobRecord(
		name=job.name,
		started_at=current - timedelta(minutes=1),
	)
	repo = StubJobRepository()
	repo.jobs = {job.name: record}
	session_calls = RefInt(0)

	manager_instance = JobManager(
		session_factory=create_stub_session_factory(session_calls),  # pyright: ignore[reportArgumentType]
		job_repo_factory=lambda _: repo,
		min_interval=timedelta(minutes=10),
	)

	ran = manager_instance.try_run(job)

	assert ran is False
	assert session_calls.value == 1
	assert repo.get_called_with == [job.name]
	assert repo.started_called_with == []
	assert repo.finished_called_with == []
	assert job.run_called_with == []


def test_try_run_runs_when_interval_equals_minimum(monkeypatch: pytest.MonkeyPatch) -> None:
	current = datetime(2026, 1, 1, 0, 0, tzinfo=timezone.utc)
	finished_at = datetime(2026, 1, 1, 0, 1, tzinfo=timezone.utc)
	times = [current, finished_at]

	def now(_: object | None = None) -> datetime:
		return times.pop(0)

	monkeypatch.setattr(manager, 'datetime', SimpleNamespace(now=now))

	job = StubJob()
	record = JobRecord(
		name=job.name,
		started_at=current - timedelta(minutes=10),
	)
	repo = StubJobRepository()
	repo.jobs = {job.name: record}
	session_calls = RefInt(0)

	manager_instance = JobManager(
		session_factory=create_stub_session_factory(session_calls),  # pyright: ignore[reportArgumentType]
		job_repo_factory=lambda _: repo,
		min_interval=timedelta(minutes=10),
	)

	ran = manager_instance.try_run(job)

	assert ran is True
	assert session_calls.value == 2
	assert repo.started_called_with == [current]
	assert repo.finished_called_with == [(job.name, finished_at)]
	assert job.run_called_with == [current]


def test_try_run_does_not_mark_finished_on_run_error(
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	current = datetime(2026, 1, 1, 0, 0, tzinfo=timezone.utc)
	monkeypatch.setattr(manager, 'datetime', SimpleNamespace(now=lambda _: current))

	class _FailingJob(StubJob):
		def run(self, *, evaluated_at: datetime) -> None:
			raise RuntimeError('boom')

	job = _FailingJob()
	record = JobRecord(name=job.name)
	repo = StubJobRepository()
	repo.jobs = {job.name: record}
	session_calls = RefInt(0)

	manager_instance = JobManager(
		session_factory=create_stub_session_factory(session_calls),  # pyright: ignore[reportArgumentType]
		job_repo_factory=lambda _: repo,
		min_interval=timedelta(minutes=10),
	)

	with pytest.raises(RuntimeError, match='boom'):
		manager_instance.try_run(job)

	assert session_calls.value == 1
	assert repo.started_called_with == [current]
	assert repo.finished_called_with == []
