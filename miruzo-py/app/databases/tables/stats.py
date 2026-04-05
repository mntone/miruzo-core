from sqlalchemy import CheckConstraint, Column, ForeignKey, SmallInteger, Table

from app.databases.metadata import metadata
from app.databases.tables.ingests import _ingest_table

_stats_table = Table(
	'stats',
	metadata,
	Column(
		'ingest_id',
		ForeignKey(_ingest_table.c.id, ondelete='CASCADE'),
		primary_key=True,
	),
	Column('score', SmallInteger, nullable=False),
	Column('score_evaluated', SmallInteger, nullable=False),
)

# SQLite constraints
_stats_table.append_constraint(
	CheckConstraint(
		'score BETWEEN -32768 AND 32767',
		'ck_stats_score',
	).ddl_if(dialect='sqlite'),
)
_stats_table.append_constraint(
	CheckConstraint(
		'score_evaluated BETWEEN -32768 AND 32767',
		'ck_stats_score_evaluated',
	).ddl_if(dialect='sqlite'),
)
