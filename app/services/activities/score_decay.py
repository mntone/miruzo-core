from datetime import datetime, time
from typing import Callable, final

from sqlmodel import Session

from app.domain.score.calculator import ScoreCalculator
from app.services.activities.actions.decay_creator import DecayActionCreator
from app.services.activities.actions.repository import ActionRepository
from app.services.activities.stats.repository.factory import create_stats_repository
from app.services.activities.stats.score_factory import make_score_context


@final
class ScoreDecayRunner:
	"""Run daily decay for all stats using the provided dependencies."""

	def __init__(
		self,
		*,
		score_calculator: ScoreCalculator,
		session_factory: Callable[[], Session],
		daily_reset_at: time,
	) -> None:
		self._score_calculator = score_calculator
		self._session_factory = session_factory
		self._daily_reset_at = daily_reset_at

	def apply_daily_decay(self, *, evaluated_at: datetime) -> None:
		with self._session_factory() as session:
			stats_repo = create_stats_repository(session)
			decay_creator = DecayActionCreator(
				ActionRepository(session),
				daily_reset_at=self._daily_reset_at,
			)

			stats_list = stats_repo.iterable()
			for stats in stats_list:
				new_action = decay_creator.create(
					stats.ingest_id,
					occurred_at=evaluated_at,
				)

				# None means this decay has already been applied for the period
				if new_action is None:
					continue

				context = make_score_context(
					stats=stats,
					evaluated_at=evaluated_at,
					daily_reset_at=self._daily_reset_at,
				)

				new_score = self._score_calculator.apply(
					action=new_action,
					score=stats.score,
					context=context,
				)

				stats.score = new_score
