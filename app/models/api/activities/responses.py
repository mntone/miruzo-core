from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.api.activities.stats import LoveStatsModel


@final
class LoveStatsResponse(BaseModel):
	"""Response payload for love actions."""

	model_config = ConfigDict(
		title='Love stats response',
		extra='forbid',
		frozen=True,
	)

	stats: Annotated[
		LoveStatsModel,
		Field(
			title='Image statistics',
			description='latest statistics for the image',
		),
	]
	"""latest statistics for the image"""
