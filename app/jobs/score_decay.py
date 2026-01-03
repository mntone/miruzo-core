from datetime import datetime
from typing import Callable, ClassVar, final

from sqlmodel import Session

from app.jobs.protocol import Job
from app.services.activities.score_decay import ScoreDecayRunner


@final
class ScoreDecayJob(Job):
	"""Scheduled job wrapper for running daily score decay."""

	_NAME: ClassVar[str] = 'score_decay'

	def __init__(
		self,
		runner: ScoreDecayRunner,
		*,
		session_factory: Callable[[], Session],
	) -> None:
		self._runner = runner
		self._session_factory = session_factory

	@property
	def name(self) -> str:
		return ScoreDecayJob._NAME

	def run(self, *, evaluated_at: datetime) -> None:
		with self._session_factory() as session:
			with session.begin():
				self._runner.apply_daily_decay(
					session,
					evaluated_at=evaluated_at,
				)
