from typing import Annotated, Literal, final

from pydantic import ConfigDict, Field

from app.models.api.variants.query import VariantQuery

_ContextImageLevel = Literal['default', 'rich']


@final
class ContextQuery(VariantQuery):
	"""Query parameters for the context API."""

	model_config = ConfigDict(
		title='Context query',
		extra='forbid',
		strict=True,
	)

	level: Annotated[
		_ContextImageLevel,
		Field(
			title='Image level',
			description='controls whether the response includes variant layers (`rich`) or summary-only (`default`).',
		),
	] = 'default'
	"""controls whether the response includes variant layers (`rich`) or summary-only (`default`)."""
