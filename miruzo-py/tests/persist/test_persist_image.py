from datetime import datetime, timezone
from typing import Any, Generator

import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from tests.persist.utils import add_ingest_row, get_image_row
from tests.services.images.utils import build_variant

from app.databases.metadata import metadata
from app.models.enums import ImageKind
from app.models.image import Image
from app.models.types import VariantEntry
from app.persist.images.implementation import create_image_repository


@pytest.fixture()
def session() -> Generator[Session, Any, None]:
	engine = create_engine(
		'sqlite+pysqlite:///:memory:',
		connect_args={'check_same_thread': False},
	)
	metadata.create_all(engine)
	with Session(engine) as session:
		yield session


@pytest.mark.parametrize(
	('kind', 'fallback'),
	[
		(ImageKind.PHOTO, None),
		(ImageKind.ILLUST, build_variant('jpeg', 1024, layer_id=9, label='fallback')),
	],
)
def test_create_persists_image_row(
	session: Session,
	kind: ImageKind,
	fallback: VariantEntry | None,
) -> None:
	now = datetime(2026, 1, 1, tzinfo=timezone.utc)
	ingest_id = add_ingest_row(session, ingested_at=now)

	widths = [320, 480, 640, 960]
	original = build_variant('webp', 1024)
	variants = [
		*[build_variant('webp', width, layer_id=1) for width in widths],
		build_variant('jpeg', 320, layer_id=9, label='fallback'),
	]
	create_image_repository(session).create(
		Image(
			ingest_id=ingest_id,
			ingested_at=now,
			kind=kind,
			original=original,
			fallback=fallback,
			variants=variants,
		),
	)

	row = get_image_row(session, ingest_id=ingest_id)
	assert row['ingest_id'] == ingest_id
	assert row['ingested_at'] == now.replace(tzinfo=None)
	assert row['kind'] == kind
	assert row['original'] == original
	assert row['fallback'] == fallback
	assert row['variants'] == variants
