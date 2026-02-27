from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Any

import pytest

import app.services.images.list as list_service
from app.config.environments import env
from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursorMode,
	UInt8ImageListCursor,
)
from app.services.images.list import ImageListService


@dataclass(frozen=True)
class _ImageStub:
	ingest_id: int
	ingested_at: datetime


def _build_repository(select_name: str, rows: Any) -> Any:
	class _Repository:
		def __init__(self, rows: Any) -> None:
			self.rows = rows
			self.cursor = None
			self.limit = None

		def _select(self, *, cursor: Any, limit: int) -> Any:
			self.cursor = cursor
			self.limit = limit
			return self.rows

	repository = _Repository(rows)
	setattr(repository, select_name, repository._select)
	return repository


def test_get_latest_wires_paginator_and_mapper(monkeypatch: pytest.MonkeyPatch) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	images = [_ImageStub(1, now), _ImageStub(2, now)]
	cursor = DatetimeImageListCursor(
		mode=ImageListCursorMode.LATEST,
		value=now,
		ingest_id=images[0].ingest_id,
	)
	repository = _build_repository('select_latest', images)

	calls: dict[str, Any] = {}

	def fake_fetch(repo: Any, cursor: Any, limit: int) -> Any:
		return repo.select_latest(cursor=cursor, limit=limit)

	def fake_slice(rows: Any, limit: int) -> tuple[Any, Any]:
		calls['slice'] = (rows, limit)
		return rows[:limit], cursor

	def fake_map(items: Any, *, next_cursor: Any, exclude_formats: Any, variant_layers: Any) -> str:
		calls['map'] = (items, next_cursor, exclude_formats, variant_layers)
		return 'response'

	monkeypatch.setattr(
		list_service.spec,
		'LATEST_SPEC',
		list_service.ImageListSpec(fetch=fake_fetch, slice=fake_slice),
	)
	monkeypatch.setattr(list_service, 'map_image_records_to_list_response', fake_map)

	service = ImageListService(repository=repository, variant_layers=env.variant_layers)

	response = service.get_latest(cursor=None, limit=1, exclude_formats=('avif',))

	assert repository.cursor is None
	assert repository.limit == 2
	assert calls['slice'] == (images, 1)
	assert calls['map'] == ([images[0]], cursor, ('avif',), env.variant_layers)
	assert response == 'response'


@pytest.mark.parametrize(
	'method_name, spec_name, select_name, cursor_value, row_cursor',
	[
		(
			'get_chronological',
			'CHRONOLOGICAL_SPEC',
			'select_chronological',
			DatetimeImageListCursor(
				mode=ImageListCursorMode.CHRONOLOGICAL,
				value=datetime(2024, 1, 2, tzinfo=timezone.utc),
				ingest_id=1,
			),
			datetime(2024, 1, 2, tzinfo=timezone.utc),
		),
		(
			'get_recently',
			'RECENTLY_SPEC',
			'select_recently',
			DatetimeImageListCursor(
				mode=ImageListCursorMode.RECENTLY,
				value=datetime(2024, 1, 3, tzinfo=timezone.utc),
				ingest_id=2,
			),
			datetime(2024, 1, 3, tzinfo=timezone.utc),
		),
		(
			'get_first_love',
			'FIRST_LOVE_SPEC',
			'select_first_love',
			DatetimeImageListCursor(
				mode=ImageListCursorMode.FIRST_LOVE,
				value=datetime(2024, 1, 4, tzinfo=timezone.utc),
				ingest_id=3,
			),
			datetime(2024, 1, 4, tzinfo=timezone.utc),
		),
		(
			'get_hall_of_fame',
			'HALL_OF_FAME_SPEC',
			'select_hall_of_fame',
			DatetimeImageListCursor(
				mode=ImageListCursorMode.HALL_OF_FAME,
				value=datetime(2024, 1, 5, tzinfo=timezone.utc),
				ingest_id=4,
			),
			datetime(2024, 1, 5, tzinfo=timezone.utc),
		),
		(
			'get_engaged',
			'ENGAGED_SPEC',
			'select_engaged',
			UInt8ImageListCursor(
				mode=ImageListCursorMode.ENGAGED,
				value=180,
				ingest_id=5,
			),
			200,
		),
	],
)
def test_list_methods_use_tuple_paginator(
	monkeypatch: pytest.MonkeyPatch,
	method_name: str,
	spec_name: str,
	select_name: str,
	cursor_value: DatetimeImageListCursor | UInt8ImageListCursor,
	row_cursor: datetime | int,
) -> None:
	now = datetime(2024, 1, 2, tzinfo=timezone.utc)
	images = [_ImageStub(1, now)]
	rows = [(images[0], row_cursor)]
	repository = _build_repository(select_name, rows)

	calls: dict[str, Any] = {}

	def fake_fetch(repo: Any, cursor: Any, limit: int) -> Any:
		return getattr(repo, select_name)(cursor=cursor, limit=limit)

	def fake_slice(rows_arg: Any, limit_arg: int) -> tuple[Any, Any]:
		calls['slice'] = (rows_arg, limit_arg)
		return [rows_arg[0][0]], 'next'

	def fake_map(items: Any, *, next_cursor: Any, exclude_formats: Any, variant_layers: Any) -> str:
		calls['map'] = (items, next_cursor, exclude_formats, variant_layers)
		return 'response'

	monkeypatch.setattr(
		list_service.spec,
		spec_name,
		list_service.ImageListSpec(fetch=fake_fetch, slice=fake_slice),
	)
	monkeypatch.setattr(list_service, 'map_image_records_to_list_response', fake_map)

	service = ImageListService(repository=repository, variant_layers=env.variant_layers)
	method = getattr(service, method_name)
	response = method(cursor=cursor_value, limit=3, exclude_formats=())

	assert repository.cursor == cursor_value
	assert repository.limit == 4
	assert calls['slice'] == (rows, 3)
	assert calls['map'] == ([images[0]], 'next', (), env.variant_layers)
	assert response == 'response'
