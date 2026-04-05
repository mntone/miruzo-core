from sqlalchemy import (
	CheckConstraint,
	Column,
	DateTime,
	ForeignKey,
	SmallInteger,
	Table,
	text,
)

from app.databases.metadata import metadata
from app.databases.tables.ingests import _ingest_table
from app.databases.types import JSON_VALUE
from app.models.enums import ImageKind

_image_table = Table(
	'images',
	metadata,
	Column(
		'ingest_id',
		ForeignKey(_ingest_table.c.id),
		primary_key=True,
	),
	Column('ingested_at', DateTime, nullable=False),
	Column(
		'kind',
		SmallInteger,
		CheckConstraint('kind IN(0, 1, 2, 3)', 'ck_images_kind'),
		nullable=False,
		server_default=text(str(int(ImageKind.UNSPECIFIED))),
	),
	Column('original', JSON_VALUE, nullable=False),
	Column('fallback', JSON_VALUE),
	Column('variants', JSON_VALUE, nullable=False),
)

# PostgreSQL constraints
_image_table.append_constraint(
	CheckConstraint(
		"jsonb_typeof(original) = 'object'",
		'ck_images_original',
	).ddl_if(dialect='postgresql'),
)
_image_table.append_constraint(
	CheckConstraint(
		"jsonb_typeof(fallback) = 'object'",
		'ck_images_fallback',
	).ddl_if(dialect='postgresql'),
)
_image_table.append_constraint(
	CheckConstraint(
		"jsonb_typeof(variants) = 'array'",
		'ck_images_variants',
	).ddl_if(dialect='postgresql'),
)

# SQLite constraints
_image_table.append_constraint(
	CheckConstraint(
		"json_type(original) IS 'object'",
		'ck_images_original',
	).ddl_if(dialect='sqlite'),
)
_image_table.append_constraint(
	CheckConstraint(
		"fallback IS NULL OR json_type(fallback) IS 'object'",
		'ck_images_fallback',
	).ddl_if(dialect='sqlite'),
)
_image_table.append_constraint(
	CheckConstraint(
		"json_type(variants) IS 'array'",
		'ck_images_variants',
	).ddl_if(dialect='sqlite'),
)
