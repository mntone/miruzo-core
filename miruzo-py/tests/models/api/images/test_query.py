from typing import TypeAlias

import pytest
from pydantic import ValidationError

from app.models.api.images.query import ListQuery

ListQueryNone: TypeAlias = ListQuery[None]


def test_image_list_query_defaults() -> None:
	query = ListQueryNone()
	assert query.exclude_formats == ()


def test_image_list_query_single_format() -> None:
	query = ListQueryNone(exclude_formats=('webp',))
	assert query.exclude_formats == ('webp',)


def test_image_list_query_plus_delimited_formats() -> None:
	query = ListQueryNone(exclude_formats=('webp+jxl',))
	assert query.exclude_formats == ('webp', 'jxl')


def test_image_list_query_deduplicates_formats() -> None:
	query = ListQueryNone(exclude_formats=('webp+webp+jxl',))
	assert query.exclude_formats == ('webp', 'jxl')


def test_image_list_query_supports_iterable_input() -> None:
	query = ListQueryNone(exclude_formats=('webp+jxl', 'avif'))
	assert query.exclude_formats == ('webp', 'jxl', 'avif')


def test_image_list_query_accepts_minimum_limit() -> None:
	query = ListQueryNone(limit=1)
	assert query.limit == 1


def test_image_list_query_rejects_limit_below_minimum() -> None:
	with pytest.raises(ValidationError):
		ListQueryNone(limit=0)


@pytest.mark.parametrize('value', ['webp,jxl', 'webp jxl', 'webp|jxl', 'WEBP+JXL'])
def test_image_list_query_rejects_invalid_separators(value: str) -> None:
	with pytest.raises(ValidationError):
		ListQueryNone(exclude_formats=(value,))


@pytest.mark.parametrize('value', ['webp+$jxl', 'webp+av-if', 'webp+av_if'])
def test_image_list_query_rejects_invalid_characters(value: str) -> None:
	with pytest.raises(ValidationError):
		ListQueryNone(exclude_formats=(value,))
