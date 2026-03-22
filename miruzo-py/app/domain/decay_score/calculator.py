from typing import final

from app.config.score import ScoreConfig
from app.domain.decay_score.context import DecayScoreContext


def _clamp(val: int, lo: int, hi: int) -> int:
	if lo > hi:
		raise ValueError('lo must be <= hi')
	return lo if val < lo else hi if val > hi else val


@final
class DecayScoreCalculator:
	def __init__(self, config: ScoreConfig) -> None:
		self._config = config

	@property
	def config(self) -> ScoreConfig:
		return self._config

	def apply(
		self,
		*,
		score: int,
		context: DecayScoreContext,
	) -> int:
		if score >= self._config.decay_high_score_threshold:
			score += self._config.daily_decay_high_score_penalty
		elif context.days_since_last_view != 0 and context.days_since_last_view % 10 == 0:
			score += self._config.daily_decay_interval10d_penalty
		else:
			score += self._config.daily_decay_penalty

		if not context.has_view_today:
			score += self._config.daily_no_access_adjustment

		return _clamp(
			score,
			self._config.minimum_score,
			self._config.maximum_score,
		)
