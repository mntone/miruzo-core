# pyright: reportAssignmentType=false
# pyright: reportUnknownVariableType=false

from collections.abc import Sequence
from datetime import datetime
from typing import Optional, final

from sqlalchemy import (
	JSON,
	BigInteger,
	CheckConstraint,
	Column,
	SmallInteger,
	String,
)
from sqlmodel import Field as SQLField
from sqlmodel import Relationship, SQLModel

from app.config.constants import (
	DAILY_LOVE_USED_MAXIMUM,
	EXECUTION_MAXIMUM,
	INGEST_ID_MAXIMUM,
	INGEST_ID_MINIMUM,
)
from app.config.environments import env
from app.models.enums import ActionKind, ImageKind, ProcessStatus, VisibilityStatus
from app.models.types import ExecutionEntry, ExecutionsJSON, UTCDateTime, VariantEntry


@final
class SettingsRecord(SQLModel, table=True):
	__tablename__ = 'settings'

	key: str = SQLField(
		min_length=4,
		max_length=8,
		sa_column=Column(
			String,
			CheckConstraint('length(key) BETWEEN 4 AND 8', 'ck_settings_key'),
			primary_key=True,
		),
	)
	value: str | None = SQLField(
		min_length=1,
		sa_column=Column(
			String,
			CheckConstraint('length(value) >= 1', 'ck_settings_value'),
		),
	)


@final
class IngestRecord(SQLModel, table=True):
	__tablename__ = 'ingests'

	id: int = SQLField(
		default=None,
		ge=INGEST_ID_MINIMUM,
		le=INGEST_ID_MAXIMUM,
		primary_key=True,
		nullable=False,
	)
	process: ProcessStatus = SQLField(
		default=ProcessStatus.PROCESSING,
		sa_column=Column(
			SmallInteger,
			CheckConstraint('process IN (0, 1)', 'ck_ingests_process'),
			autoincrement=False,
			default=ProcessStatus.PROCESSING,
			nullable=False,
		),
	)
	visibility: VisibilityStatus = SQLField(
		default=VisibilityStatus.PRIVATE,
		sa_column=Column(
			SmallInteger,
			CheckConstraint('visibility IN (0, 1)', 'ck_ingests_visibility'),
			autoincrement=False,
			default=VisibilityStatus.PRIVATE,
			nullable=False,
		),
	)
	relative_path: str = SQLField(
		min_length=4,
		sa_column=Column(
			String,
			CheckConstraint('length(relative_path) >= 4', 'ck_ingests_relative_path'),
		),
	)
	fingerprint: str = SQLField(min_length=64, max_length=64, unique=True)
	ingested_at: datetime = SQLField(sa_column=Column(UTCDateTime(), nullable=False))
	captured_at: datetime = SQLField(
		sa_column=Column(
			UTCDateTime(),
			CheckConstraint('captured_at <= ingested_at'),
			nullable=False,
		),
	)
	updated_at: datetime = SQLField(
		sa_column=Column(
			UTCDateTime(),
			CheckConstraint('updated_at >= ingested_at'),
			nullable=False,
		),
	)
	executions: Sequence[ExecutionEntry] | None = SQLField(
		default=None,
		min_length=1,
		max_length=EXECUTION_MAXIMUM,
		sa_column=Column(ExecutionsJSON),
	)

	image: Optional['ImageRecord'] = Relationship(back_populates='ingest')
	actions: Optional['ActionRecord'] = Relationship(back_populates='ingest', cascade_delete=True)
	stats: Optional['StatsRecord'] = Relationship(back_populates='ingest', cascade_delete=True)


@final
class ImageRecord(SQLModel, table=True):
	__tablename__ = 'images'

	ingest_id: int = SQLField(primary_key=True, foreign_key='ingests.id', nullable=False)
	ingested_at: datetime = SQLField(default=datetime.min, sa_column=Column(UTCDateTime(), nullable=False))
	kind: ImageKind = SQLField(
		default=ImageKind.UNSPECIFIED,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			default=ImageKind.UNSPECIFIED,
			nullable=False,
		),
	)

	original: VariantEntry = SQLField(sa_column=Column(JSON))
	fallback: VariantEntry | None = SQLField(default=None, sa_column=Column(JSON(none_as_null=True)))
	variants: Sequence[VariantEntry] = SQLField(sa_column=Column(JSON))

	ingest: IngestRecord = Relationship(back_populates='image')


@final
class ActionRecord(SQLModel, table=True):
	__tablename__ = 'actions'

	id: int = SQLField(default=None, primary_key=True, nullable=False)
	ingest_id: int = SQLField(foreign_key='ingests.id', ondelete='CASCADE', nullable=False)
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

	ingest: IngestRecord = Relationship(back_populates='actions')


@final
class StatsRecord(SQLModel, table=True):
	__tablename__ = 'stats'

	ingest_id: int = SQLField(
		primary_key=True,
		foreign_key='ingests.id',
		ondelete='CASCADE',
		nullable=False,
	)
	score: int = SQLField(
		ge=env.score.minimum_score,
		le=env.score.maximum_score,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			nullable=False,
		),
	)
	score_evaluated: int = SQLField(
		ge=env.score.minimum_score,
		le=env.score.maximum_score,
		sa_column=Column(
			SmallInteger,
			autoincrement=False,
			nullable=False,
		),
	)
	score_evaluated_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	first_loved_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	last_loved_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	hall_of_fame_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))
	last_viewed_at: datetime | None = SQLField(default=None, sa_column=Column(UTCDateTime()))

	view_count: int = SQLField(
		default=0,
		ge=0,
		sa_column=Column(
			BigInteger,
			CheckConstraint('view_count >= 0'),
			autoincrement=False,
			default=0,
			nullable=False,
		),
	)
	view_milestone_count: int = SQLField(
		default=0,
		ge=0,
		sa_column=Column(
			BigInteger,
			CheckConstraint('view_milestone_count >= 0 AND view_milestone_count <= view_count'),
			autoincrement=False,
			default=0,
			nullable=False,
		),
	)
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
		le=DAILY_LOVE_USED_MAXIMUM,
		sa_column=Column(
			SmallInteger,
			CheckConstraint('daily_love_used >= 0'),
			autoincrement=False,
			default=0,
			nullable=False,
		),
	)
