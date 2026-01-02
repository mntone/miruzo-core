from typing import final

from app.domain.score.context import ScoreContext
from app.models.enums import ActionKind
from app.models.records import ActionRecord


@final
class StubScoreCalculator:
	def __init__(self) -> None:
		self.apply_called_with: list[tuple[ActionRecord, int, ScoreContext]] = []

	def apply(
		self,
		*,
		action: ActionRecord,
		score: int,
		context: ScoreContext,
	) -> int:
		self.apply_called_with.append((action, score, context))

		match action.kind:
			case ActionKind.DECAY:
				return score - 2
			case _:
				return score + 2
