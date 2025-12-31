from datetime import datetime
from typing import Annotated

from pydantic import BaseModel, ConfigDict, Field

from app.config.constants import DEFAULT_LIMIT, LIMIT_MAXIMUM, LIMIT_MINIMUM


class PaginationQuery(BaseModel):
	model_config = ConfigDict(
		title='Pagination query',
		extra='forbid',
		strict=True,
	)

	cursor: Annotated[
		datetime | None,
		Field(
			title='Cursor',
			description='opaque pagination cursor representing the captured_at of the last item returned; `null` for the first page',
		),
	] = None
	"""opaque pagination cursor representing the captured_at of the last item returned; `null` for the first page"""

	limit: Annotated[
		int,
		Field(
			title='Limit',
			description='maximum number of items to return for this request',
			ge=LIMIT_MINIMUM,
			le=LIMIT_MAXIMUM,
		),
	] = DEFAULT_LIMIT
	"""maximum number of items to return for this request"""
