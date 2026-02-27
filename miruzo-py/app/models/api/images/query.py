from typing import final

from pydantic import ConfigDict

from app.models.api.common.query import PaginationQuery
from app.models.api.variants.query import VariantQuery


@final
class ListQuery(PaginationQuery[str], VariantQuery):
	"""Query parameters for image list endpoints."""

	model_config = ConfigDict(title='Image list query', strict=True)
