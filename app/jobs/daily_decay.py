from datetime import datetime
from typing import Callable, ClassVar, final

from sqlmodel import Session

from app.jobs.protocol import Job
from app.services.activities.daily_decay import DailyDecayRunner


@final
class DailyDecayJob(Job):
	"""Scheduled job wrapper for running daily decay."""

	_NAME: ClassVar[str] = 'daily_decay'

	def __init__(
		self,
		runner: DailyDecayRunner,
		*,
		session_factory: Callable[[], Session],
	) -> None:
		self._runner = runner
		self._session_factory = session_factory

	@property
	def name(self) -> str:
		return DailyDecayJob._NAME

	def run(self, *, evaluated_at: datetime) -> None:
		with self._session_factory() as session:
			with session.begin():
				self._runner.apply_daily_decay(
					session,
					evaluated_at=evaluated_at,
				)
