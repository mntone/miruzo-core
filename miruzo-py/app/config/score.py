from collections.abc import Sequence
from dataclasses import dataclass
from typing import final


@final
@dataclass(frozen=True, slots=True)
class ScoreConfig:
	"""
	Configuration for score calculation.
	All values are interpreted by ScoreCalculator.
	"""

	initial_score: int = 100
	minimum_score: int = 0
	maximum_score: int = 200
	engaged_score_threshold: int = 160

	# --- decay (daily) ---
	decay_high_score_threshold: int = 180  # == hall_of_fame_threshold for now

	daily_decay_penalty: int = -2
	daily_decay_interval10d_penalty: int = -3
	daily_decay_high_score_penalty: int = -3
	daily_no_access_adjustment: int = +1

	# --- view ---
	view_bonus_at_first: int = 10
	view_bonus_by_days: Sequence[tuple[int, int]] = (
		# (days, score)
		(1, 10),
		(3, 7),
		(7, 5),
		(30, 3),
	)
	view_bonus_fallback: int = 2

	# --- memo ---
	memo_bonus: int = 1
	memo_penalty: int = -1

	# --- love ---
	love_bonus: int = 20
	love_penalty: int = -18

	# --- hall of fame ---
	hall_of_fame_threshold: int = 180
