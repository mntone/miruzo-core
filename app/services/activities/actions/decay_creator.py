from datetime import datetime
from typing import final

from app.domain.activities.daily_period import DailyPeriodResolver
from app.models.enums import ActionKind
from app.models.records import ActionRecord
from app.persist.actions.protocol import ActionRepository


@final
class DecayActionCreator:
	"""Create a daily decay action unless one already exists for the period."""

	def __init__(
		self,
		*,
		repository: ActionRepository,
		period_resolver: DailyPeriodResolver,
	) -> None:
		self._repository = repository
		self._period_resolver = period_resolver

	def create(
		self,
		ingest_id: int,
		*,
		occurred_at: datetime,
	) -> ActionRecord | None:
		# --- get period_start ---
		period_start, period_end = self._period_resolver.resolve_period_range(occurred_at)

		# --- get decay action ---
		action = self._repository.select_latest_one(
			ingest_id,
			kind=ActionKind.DECAY,
			since_occurred_at=period_start,
			until_occurred_at=period_end,
			require_unique=True,
		)
		if action is not None:
			return None

		# --- insert new action ---
		new_action = self._repository.insert(
			ingest_id,
			kind=ActionKind.DECAY,
			occurred_at=occurred_at,
		)

		return new_action
