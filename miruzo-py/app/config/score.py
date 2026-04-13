from dataclasses import dataclass
from typing import final


@dataclass(frozen=True, slots=True)
@final
class ScoreConfig:
	"""
	Configuration for score calculation.
	All values are interpreted by DecayScoreCalculator.
	"""

	initial_score: int = 100
	minimum_score: int = -20000
	maximum_score: int = 200
