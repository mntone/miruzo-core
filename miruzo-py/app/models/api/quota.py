from datetime import datetime
from typing import Annotated, Literal, final

from pydantic import BaseModel, Field


@final
class QuotaItem(BaseModel):
	"""Quota limits and remaining counts for a single period."""

	period: Annotated[
		Literal['daily'],
		Field(
			title='Quota period',
			description='quota period this limit applies to',
		),
	]
	"""quota period this limit applies to."""

	reset_at: Annotated[
		datetime,
		Field(
			title='Reset timestamp',
			description='timestamp when the quota resets',
		),
	]
	"""timestamp when the quota resets."""

	limit: Annotated[
		int,
		Field(
			title='Quota limit',
			description='maximum allowed actions per period',
			ge=1,
		),
	]
	"""maximum allowed actions per period."""

	remaining: Annotated[
		int,
		Field(
			title='Remaining quota',
			description='remaining actions available in the current period',
			ge=0,
		),
	]
	"""remaining actions available in the current period."""


@final
class QuotaResponse(BaseModel):
	"""Top-level quota response payload."""

	love: QuotaItem
