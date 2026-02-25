from typing import final

from app.config.score import ScoreConfig
from app.domain.score.context import ScoreContext
from app.models.enums import ActionKind
from app.models.records import ActionRecord


def _clamp(val: int, lo: int, hi: int) -> int:
	if lo > hi:
		raise ValueError('lo must be <= hi')
	return lo if val < lo else hi if val > hi else val


@final
class ScoreCalculator:
	def __init__(self, config: ScoreConfig) -> None:
		self._config = config

	@property
	def config(self) -> ScoreConfig:
		return self._config

	def _view_bonus(self, context: ScoreContext) -> int:
		if context.last_viewed_at is None:
			return 0

		days = context.days_since_last_view
		for limit, bonus in self._config.view_bonus_by_days:
			if days <= limit:
				return bonus
		return self._config.view_bonus_fallback

	def apply(
		self,
		*,
		action: ActionRecord,
		score: int,
		context: ScoreContext,
	) -> int:
		match action.kind:
			case ActionKind.VIEW:
				if context.last_viewed_at is None:
					score += self._config.view_bonus_at_first
				elif not context.has_view_today:
					score += self._view_bonus(context)

			case ActionKind.DECAY:
				if score >= self._config.decay_high_score_threshold:
					score += self._config.daily_decay_high_score_penalty
				elif context.days_since_last_view != 0 and context.days_since_last_view % 10 == 0:
					score += self._config.daily_decay_interval10d_penalty
				else:
					score += self._config.daily_decay_penalty

				if not context.has_view_today:
					score += self._config.daily_no_access_adjustment

			case ActionKind.MEMO:
				score += self._config.memo_bonus

			case ActionKind.LOVE:
				score += self._config.love_bonus

			case ActionKind.LOVE_CANCELED:
				score += self._config.love_penalty

			case _:
				pass

		return _clamp(
			score,
			self._config.minimum_score,
			self._config.maximum_score,
		)
