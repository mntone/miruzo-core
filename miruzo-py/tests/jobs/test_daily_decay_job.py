from datetime import datetime, timezone
from typing import final

from tests.stubs.session import StubSession, create_stub_session

from app.jobs.daily_decay import DailyDecayJob


@final
class _StubScoreDecayRunner:
	def __init__(self) -> None:
		self.called_with: list[datetime] = []

	def apply_daily_decay(self, _: StubSession, *, evaluated_at: datetime) -> None:
		self.called_with.append(evaluated_at)


def test_job_name_is_stable() -> None:
	runner = _StubScoreDecayRunner()
	job = DailyDecayJob(
		runner,  # pyright: ignore[reportArgumentType]
		session_factory=create_stub_session,  # pyright: ignore[reportArgumentType]
	)

	assert job.name == 'daily_decay'


def test_run_delegates_to_runner() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	runner = _StubScoreDecayRunner()
	job = DailyDecayJob(
		runner,  # pyright: ignore[reportArgumentType]
		session_factory=create_stub_session,  # pyright: ignore[reportArgumentType]
	)

	job.run(evaluated_at=evaluated_at)

	assert runner.called_with == [evaluated_at]
