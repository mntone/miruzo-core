"""
Add image list indexes

Revision ID: a46b415897c4
Revises:
Create Date: 2026-01-01 09:58:54.621175
"""

from collections.abc import Sequence
from typing import Union

from alembic import op

# revision identifiers, used by Alembic.
revision: str = 'a46b415897c4'
down_revision: Union[str, Sequence[str], None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
	"""Add list indexes for ImageListQueryBuilder performance."""

	op.execute("""
		CREATE INDEX ix_images_latest
		ON images (ingested_at DESC, ingest_id DESC);
	""")
	op.execute("""
		CREATE INDEX ix_ingests_chronological
		ON ingests (captured_at DESC, id DESC);
	""")
	op.execute("""
		CREATE INDEX ix_stats_recently
		ON stats (last_viewed_at DESC, ingest_id DESC)
		WHERE last_viewed_at IS NOT NULL;
	""")
	op.execute("""
		CREATE INDEX ix_stats_first_love
		ON stats (first_loved_at DESC, ingest_id DESC)
		WHERE first_loved_at IS NOT NULL;
	""")
	op.execute("""
		CREATE INDEX ix_stats_hall_of_fame
		ON stats (hall_of_fame_at DESC, ingest_id DESC)
		WHERE hall_of_fame_at IS NOT NULL;
	""")


def downgrade() -> None:
	op.drop_index('ix_stats_hall_of_fame', table_name='stats')
	op.drop_index('ix_stats_first_love', table_name='stats')
	op.drop_index('ix_stats_recently', table_name='stats')
	op.drop_index('ix_ingests_chronological', table_name='ingests')
	op.drop_index('ix_images_latest', table_name='images')
