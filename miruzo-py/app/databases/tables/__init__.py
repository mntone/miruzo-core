from app.databases.tables.images import _image_table as image_table
from app.databases.tables.ingests import _ingest_table as ingest_table
from app.databases.tables.stats import _stats_table as stats_table

__all__ = [
	'image_table',
	'ingest_table',
	'stats_table',
]
