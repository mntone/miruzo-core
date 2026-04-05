from typing import Annotated, Protocol, final

from pydantic import BaseModel, Field

from app.config import constants as c


@final
class StatsCreateInput(BaseModel):
	ingest_id: Annotated[
		int,
		Field(ge=c.INGEST_ID_MINIMUM, le=c.INGEST_ID_MAXIMUM),
	]
	initial_score: Annotated[
		int,
		Field(ge=c.SCORE_MINIMUM, le=c.SCORE_MAXIMUM),
	]


class StatsRepository(Protocol):
	def create(self, entry: StatsCreateInput) -> None: ...
