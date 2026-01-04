# pyright: reportAssignmentType=false
# pyright: reportUnknownVariableType=false

from collections.abc import Sequence
from datetime import datetime
from typing import Optional, final

from sqlalchemy import JSON, Column, Integer, SmallInteger
from sqlmodel import Field as SQLField
from sqlmodel import Relationship, SQLModel

from app.config.constants import EXECUTION_MAXIMUM
from app.config.environments import env
from app.models.enums import ActionKind, ImageKind, ProcessStatus, VisibilityStatus
from app.models.types import ExecutionEntry, ExecutionsJSON, UTCDateTime, VariantEntry


@final
class IngestRecord(SQLModel, table=True):
	__tablename__ = 'ingests'

	id: int = SQLField(default=None, primary_key=True, nullable=False)
	process: ProcessStatus = SQLField(
		default=ProcessStatus.PROCESSING,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			default=ProcessStatus.PROCESSING,
			nullable=False,
		),
	)
	visibility: VisibilityStatus = SQLField(
		default=VisibilityStatus.PRIVATE,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			default=VisibilityStatus.PRIVATE,
			nullable=False,
		),
	)
	relative_path: str
	fingerprint: str = SQLField(min_length=64, max_length=64, unique=True)
	ingested_at: datetime = SQLField(default=datetime.min, sa_column=Column(UTCDateTime(), nullable=False))
	captured_at: datetime = SQLField(default=datetime.min, sa_column=Column(UTCDateTime(), nullable=False))
	updated_at: datetime = SQLField(default=datetime.min, sa_column=Column(UTCDateTime(), nullable=False))
	executions: Sequence[ExecutionEntry] | None = SQLField(
		default=None,
		min_length=1,
		max_length=EXECUTION_MAXIMUM,
		sa_column=Column(ExecutionsJSON),
	)

	image: Optional['ImageRecord'] = Relationship(back_populates='ingest')
	stats: Optional['StatsRecord'] = Relationship(back_populates='ingest')


@final
class ImageRecord(SQLModel, table=True):
	__tablename__ = 'images'

	ingest_id: int = SQLField(primary_key=True, foreign_key='ingests.id', nullable=False)
	ingested_at: datetime = SQLField(default=datetime.min, sa_column=Column(UTCDateTime(), nullable=False))
	kind: ImageKind = SQLField(default=ImageKind.PHOTO, sa_column=Column(Integer))

	original: VariantEntry = SQLField(sa_column=Column(JSON))
	fallback: VariantEntry | None = SQLField(default=None, sa_column=Column(JSON))
	variants: Sequence[VariantEntry] = SQLField(sa_column=Column(JSON))

	ingest: IngestRecord = Relationship(back_populates='image')


@final
class ActionRecord(SQLModel, table=True):
	__tablename__ = 'actions'

	id: int = SQLField(default=None, primary_key=True, nullable=False)
	ingest_id: int = SQLField(foreign_key='ingests.id', nullable=False)
	kind: ActionKind = SQLField(
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			nullable=False,
		),
	)
	occurred_at: datetime = SQLField(
		default=datetime.min,
		sa_column=Column(UTCDateTime(), nullable=False),
	)


@final
class StatsRecord(SQLModel, table=True):
	__tablename__ = 'stats'

	ingest_id: int = SQLField(primary_key=True, foreign_key='ingests.id', nullable=False)
	score: int = SQLField(
		le=env.score.maximum_score,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			default=env.score.initial_score,
			index=True,
			nullable=False,
		),
	)
	view_count: int = SQLField(default=0, ge=0, nullable=False)
	last_viewed_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	first_loved_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	hall_of_fame_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))

	view_milestone_count: int = SQLField(default=0, ge=0, nullable=False)
	view_milestone_archived_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))

	ingest: IngestRecord = Relationship(back_populates='stats')


@final
class JobRecord(SQLModel, table=True):
	__tablename__ = 'jobs'

	name: str = SQLField(min_length=8, max_length=16, primary_key=True)
	started_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	finished_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))


@final
class UserRecord(SQLModel, table=True):
	__tablename__ = 'users'

	id: int = SQLField(default=None, primary_key=True, nullable=False)
	daily_love_used: int = SQLField(
		default=0,
		ge=0,
		le=10,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			default=0,
			nullable=False,
		),
	)
