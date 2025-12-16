from datetime import datetime, timezone
from typing import Annotated, Optional, TypedDict, final

from pydantic import Field
from sqlalchemy import JSON, Column, Integer
from sqlmodel import Field as SQLField
from sqlmodel import Relationship, SQLModel

from app.config.constants import DEFAULT_SCORE, SCORE_MAXIMUM, SCORE_MINIMUM
from app.models.enums import ImageStatus


@final
class VariantRecord(TypedDict):
	filepath: str
	format: Annotated[str, Field(ge=3, le=8)]
	codecs: str | None = None
	size: Annotated[int, Field(ge=1)]
	width: Annotated[int, Field(ge=1, le=10240)]
	height: Annotated[int, Field(ge=1, le=10240)]
	quality: Annotated[int | None, Field(ge=1, le=100)] = None


@final
class ImageRecord(SQLModel, table=True):
	__tablename__ = 'images'

	id: int | None = SQLField(primary_key=True, default=None)
	fingerprint: str = SQLField(min_length=64, max_length=64, unique=True)
	captured_at: datetime | None = SQLField(default=None)
	ingested_at: datetime = SQLField(default=datetime.now(timezone.utc))
	status: ImageStatus = SQLField(default=ImageStatus.ACTIVE, sa_column=Column(Integer))

	original: VariantRecord = SQLField(sa_column=Column(JSON))
	fallback: VariantRecord | None = SQLField(default=None, sa_column=Column(JSON))
	variants: list[list[VariantRecord]] = SQLField(sa_column=Column(JSON))

	stats: Optional['StatsRecord'] = Relationship(back_populates='image')


@final
class StatsRecord(SQLModel, table=True):
	__tablename__ = 'stats'

	image_id: int = SQLField(primary_key=True, foreign_key='images.id', nullable=False)
	favorite: bool = SQLField(default=False)
	score: int = SQLField(default=DEFAULT_SCORE, ge=SCORE_MINIMUM, le=SCORE_MAXIMUM, index=True)
	view_count: int = SQLField(default=0, ge=0)
	last_viewed_at: datetime | None = SQLField(default=None, index=True)

	image: ImageRecord | None = Relationship(back_populates='stats')
