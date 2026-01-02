from datetime import datetime
from typing import ClassVar, final

from app.jobs.protocol import Job
from app.services.activities.score_decay import ScoreDecayRunner


@final
class ScoreDecayJob(Job):
	"""Scheduled job wrapper for running daily score decay."""

	_NAME: ClassVar[str] = 'score_decay'

	def __init__(self, runner: ScoreDecayRunner) -> None:
		self._runner = runner

	@property
	def name(self) -> str:
		return ScoreDecayJob._NAME

	def run(self, *, evaluated_at: datetime) -> None:
		self._runner.apply_daily_decay(
			evaluated_at=evaluated_at,
		)
