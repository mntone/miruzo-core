from dataclasses import dataclass
from typing import final


@final
@dataclass(frozen=True, slots=True)
class ScoreConfig:
	"""
	Configuration for score calculation.
	All values are interpreted by DecayScoreCalculator.
	"""

	initial_score: int = 100
	minimum_score: int = -20000
	public_minimum_score: int = 0
	maximum_score: int = 200
	engaged_score_threshold: int = 160

	# --- decay (daily) ---
	decay_high_score_threshold: int = 180  # == hall_of_fame_threshold for now

	daily_decay_penalty: int = -2
	daily_decay_interval10d_penalty: int = -3
	daily_decay_high_score_penalty: int = -3
	daily_no_access_adjustment: int = +1
