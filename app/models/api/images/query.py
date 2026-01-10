from typing import Generic, final

from pydantic import ConfigDict

from app.models.api.common.query import PaginationQuery, TCursor
from app.models.api.variants.query import VariantQuery


@final
class ListQuery(PaginationQuery[TCursor], VariantQuery, Generic[TCursor]):
	"""Query parameters for image list endpoints."""

	model_config = ConfigDict(title='Image list query', strict=True)
