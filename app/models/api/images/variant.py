from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.config.environments import env
from app.models.api.utils.units import bytes_to_manbytes
from app.models.types import VariantEntry


@final
class VariantModel(BaseModel):
	"""Normalized metadata for a single rendered asset."""

	model_config = ConfigDict(
		title='Image variant',
		extra='forbid',
		frozen=True,
	)

	src: Annotated[
		str,
		Field(
			title='Variant source URL',
			description=f'server-relative path to the asset (e.g. `{env.public_media_root}/foo.webp`) so clients can fetch it directly',
		),
	]
	"""server-relative path to the asset, so clients can fetch it directly"""

	format: Annotated[
		str,
		Field(
			title='Variant format',
			description='container format string (e.g. `webp`) that tells browsers how to decode the file',
		),
	]
	"""container format string (e.g. `webp`) that tells browsers how to decode the file"""

	codecs: Annotated[
		str | None,
		Field(
			title='Variant codec hint',
			description="optional codec hint (`vp8`, `vp8l`, etc.) for cases where the format alone isn't specific enough",
		),
	]
	"""optional codec hint (`vp8`, `vp8l`, etc.) for cases where the format alone isn't specific enough"""

	manbytes: Annotated[
		int,
		Field(
			title='Variant manbytes',
			description='file size expressed in manbytes (see docs/unit.md); typically ≥1 but 0 indicates an unexpected/invalid asset',
			ge=0,
		),
	]
	"""file size expressed in manbytes (see docs/unit.md); typically ≥1 but 0 indicates an unexpected/invalid asset"""

	w: Annotated[
		int,
		Field(
			title='Variant width',
			description='width of this rendition in pixels; guaranteed to be a positive integer',
			gt=0,
		),
	]
	"""width of this rendition in pixels; guaranteed to be a positive integer"""

	h: Annotated[
		int,
		Field(
			title='Variant height',
			description='height of this rendition in pixels; guaranteed to be a positive integer',
			gt=0,
		),
	]
	"""height of this rendition in pixels; guaranteed to be a positive integer"""

	@classmethod
	def from_record(cls, variant: VariantEntry) -> 'VariantModel':
		return cls(
			src=env.public_media_root + variant['rel'],
			format=variant['format'],
			codecs=variant['codecs'],
			manbytes=bytes_to_manbytes(variant['bytes']),
			w=variant['width'],
			h=variant['height'],
		)
