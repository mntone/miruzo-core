from dataclasses import dataclass
from datetime import datetime
from typing import final


@dataclass(frozen=True, slots=True)
@final
class ScoreContext:
	"""Snapshot used to compute score adjustments for a single evaluation."""

	evaluated_at: datetime

	last_viewed_at: datetime | None

	days_since_last_view: int

	has_view_today: bool
	"""Derived state for scoring, not a persisted fact"""
