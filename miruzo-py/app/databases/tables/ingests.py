from sqlalchemy import (
	BigInteger,
	CheckConstraint,
	Column,
	DateTime,
	Integer,
	String,
	Table,
	text,
)

from app.databases.metadata import metadata
from app.databases.types import JSON_VALUE, LEAST8_INT

_ingest_table = Table(
	'ingests',
	metadata,
	Column(
		'id',
		BigInteger().with_variant(Integer, 'sqlite'),
		CheckConstraint('id BETWEEN 1 AND 9007199254740991', 'ck_ingests_id'),
		primary_key=True,
	),
	Column(
		'process',
		LEAST8_INT,
		CheckConstraint('process IN(0, 1)', 'ck_ingests_process'),
		nullable=False,
		server_default=text('0'),
	),
	Column(
		'visibility',
		LEAST8_INT,
		CheckConstraint('visibility IN(0, 1)', 'ck_ingests_visibility'),
		nullable=False,
		server_default=text('0'),
	),
	Column(
		'relative_path',
		String,
		nullable=False,
	),
	Column('fingerprint', String(length=64), nullable=False, unique=True),
	Column('ingested_at', DateTime, nullable=False),
	Column(
		'captured_at',
		DateTime,
		CheckConstraint('captured_at <= ingested_at'),
		nullable=False,
	),
	Column(
		'updated_at',
		DateTime,
		CheckConstraint('updated_at >= ingested_at'),
		nullable=False,
	),
	Column(
		'executions',
		JSON_VALUE,
		nullable=False,
		server_default=text("'[]'"),
	),
	sqlite_autoincrement=True,
)

# PostgreSQL constraints
_ingest_table.append_constraint(
	CheckConstraint(
		"length(relative_path) >= 4 AND relative_path !~ '^/'",
		'ck_ingests_relative_path',
	).ddl_if(dialect='postgresql'),
)
_ingest_table.append_constraint(
	CheckConstraint(
		"jsonb_typeof(executions) = 'array' AND jsonb_array_length(executions) <= 5",
		'ck_ingests_executions',
	).ddl_if(dialect='postgresql'),
)

# SQLite constraints
_ingest_table.append_constraint(
	CheckConstraint(
		"length(relative_path) >= 4 AND relative_path NOT LIKE '/%'",
		'ck_ingests_relative_path',
	).ddl_if(dialect='sqlite'),
)
_ingest_table.append_constraint(
	CheckConstraint(
		"json_type(executions) IS 'array' AND json_array_length(executions) <= 5",
		'ck_ingests_executions',
	).ddl_if(dialect='sqlite'),
)
