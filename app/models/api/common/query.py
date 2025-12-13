from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.core.constants import DEFAULT_LIMIT, LIMIT_MAXIMUM, LIMIT_MINIMUM


@final
class PaginationQuery(BaseModel):
	model_config = ConfigDict(title='Pagination query', strict=True)

	cursor: Annotated[
		datetime | None,
		Field(
			default=None,
			title='Cursor',
			description='opaque pagination cursor representing the captured_at of the last item returned; `null` for the first page',
		),
	]
	limit: Annotated[
		int,
		Field(
			title='Limit',
			description='maximum number of items to return for this request',
			gt=LIMIT_MINIMUM,
			le=LIMIT_MAXIMUM,
		),
	] = DEFAULT_LIMIT
