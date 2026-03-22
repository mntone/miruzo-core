from typing import final

from app.domain.decay_score.context import DecayScoreContext


@final
class StubDecayScoreCalculator:
	def __init__(self) -> None:
		self.apply_called_with: list[tuple[int, DecayScoreContext]] = []

	def apply(
		self,
		*,
		score: int,
		context: DecayScoreContext,
	) -> int:
		self.apply_called_with.append((score, context))

		return score - 2
