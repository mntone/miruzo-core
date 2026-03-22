from dataclasses import dataclass
from typing import final


@dataclass(frozen=True, slots=True)
@final
class DecayScoreContext:
	"""Snapshot used to compute score adjustments for a single evaluation."""

	days_since_last_view: int

	has_view_today: bool
	"""Derived state for scoring, not a persisted fact"""
