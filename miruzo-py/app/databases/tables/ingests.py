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
from sqlalchemy.dialects import mysql

from app.databases.defaults import empty_json_array_default
from app.databases.metadata import metadata
from app.databases.types import JSON_VALUE, LEAST8_INT

INGEST_ID_TYPE = (
	BigInteger().with_variant(mysql.BIGINT(unsigned=True), 'mysql').with_variant(Integer, 'sqlite')
)

_ingest_table = Table(
	'ingests',
	metadata,
	Column('id', INGEST_ID_TYPE, primary_key=True),
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
		String(length=255),
		nullable=False,
	),
	Column('fingerprint', String(length=64), nullable=False, unique=True),
	Column('ingested_at', DateTime, nullable=False),
	Column('captured_at', DateTime, nullable=False),
	Column('updated_at', DateTime, nullable=False),
	Column(
		'executions',
		JSON_VALUE,
		nullable=False,
		server_default=empty_json_array_default(),
	),
	sqlite_autoincrement=True,
)

# Default constraints
_ingest_table.append_constraint(CheckConstraint('captured_at <= ingested_at'))
_ingest_table.append_constraint(CheckConstraint('updated_at >= ingested_at'))

# MySQL constraints
_ingest_table.append_constraint(
	CheckConstraint(
		"char_length(relative_path) >= 5 AND relative_path NOT LIKE '/%' AND relative_path NOT LIKE '..%'",
		'ck_ingests_relative_path',
	).ddl_if(dialect='mysql'),
)
_ingest_table.append_constraint(
	CheckConstraint(
		"JSON_TYPE(executions) = 'ARRAY' AND JSON_LENGTH(executions) <= 5",
		'ck_ingests_executions',
	).ddl_if(dialect='mysql'),
)

# PostgreSQL constraints
_ingest_table.append_constraint(
	CheckConstraint(
		'id BETWEEN 1 AND 9007199254740991',
		'ck_ingests_id',
	).ddl_if(dialect='postgresql'),
)
_ingest_table.append_constraint(
	CheckConstraint(
		"length(relative_path) >= 5 AND relative_path NOT LIKE '/%' AND relative_path NOT LIKE '..%'",
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
		'id BETWEEN 1 AND 9007199254740991',
		'ck_ingests_id',
	).ddl_if(dialect='sqlite'),
)
_ingest_table.append_constraint(
	CheckConstraint(
		"length(relative_path) BETWEEN 5 AND 255 AND relative_path NOT LIKE '/%' AND relative_path NOT LIKE '..%'",
		'ck_ingests_relative_path',
	).ddl_if(dialect='sqlite'),
)
_ingest_table.append_constraint(
	CheckConstraint(
		"json_type(executions) IS 'array' AND json_array_length(executions) <= 5",
		'ck_ingests_executions',
	).ddl_if(dialect='sqlite'),
)
