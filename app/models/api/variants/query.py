import re
from collections.abc import Iterable
from typing import Annotated, cast

from pydantic import BaseModel, Field, field_validator

_FORMAT_TOKEN_PATTERN = re.compile(r'^[a-z0-9]+$')


class VariantQuery(BaseModel):
	"""Variant-related query parameters shared across API responses."""

	exclude_formats: Annotated[
		tuple[str, ...],
		Field(
			default=(),
			title='Excluded formats',
			description='list of formats to exclude from the response (e.g. `exclude_formats=webp+jxl`); empty means allow everything',
		),
	] = ()
	"""list of formats to exclude from the response (e.g. `exclude_formats=webp+jxl`); empty means allow everything"""

	@staticmethod
	def _validate_token(token: str) -> str:
		if not token or not _FORMAT_TOKEN_PATTERN.fullmatch(token):
			raise ValueError('format must consist of lowercase letters and digits and use "+" as separator')
		return token

	@staticmethod
	def _split_single(raw: str) -> list[str]:
		raw_value = raw.strip()
		if raw_value == '':
			return []

		chunks = raw_value.split('+')
		return [VariantQuery._validate_token(chunk) for chunk in chunks]

	@field_validator('exclude_formats', mode='before')
	@classmethod
	def _split_formats(cls, value: object) -> tuple[str, ...]:
		if value is None or value == '':
			return ()

		if isinstance(value, str):
			parts = cls._split_single(value)
		elif isinstance(value, (list, tuple, set)):
			iterable_value = cast(Iterable[object], value)
			parts: list[str] = []
			for item in iterable_value:
				parts.extend(cls._split_single(str(item)))
		else:
			raise TypeError('format must be supplied as a string or list of strings')

		seen: set[str] = set()
		normalized: list[str] = []
		for token in parts:
			if token not in seen:
				seen.add(token)
				normalized.append(token)

		return tuple(normalized)
