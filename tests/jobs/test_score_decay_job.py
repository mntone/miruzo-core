from datetime import datetime, timezone
from typing import final

from app.jobs.score_decay import ScoreDecayJob


@final
class _StubScoreDecayRunner:
	def __init__(self) -> None:
		self.called_with: list[datetime] = []

	def apply_daily_decay(self, *, evaluated_at: datetime) -> None:
		self.called_with.append(evaluated_at)


def test_job_name_is_stable() -> None:
	job = ScoreDecayJob(_StubScoreDecayRunner())  # pyright: ignore[reportArgumentType]

	assert job.name == 'score_decay'


def test_run_delegates_to_runner() -> None:
	evaluated_at = datetime(2026, 1, 2, 6, 0, tzinfo=timezone.utc)
	runner = _StubScoreDecayRunner()
	job = ScoreDecayJob(runner)  # pyright: ignore[reportArgumentType]

	job.run(evaluated_at=evaluated_at)

	assert runner.called_with == [evaluated_at]
