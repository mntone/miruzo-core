"""
Ensure singleton user exists

Revision ID: d4dfd8798fc0
Revises: 919723f05419
Create Date: 2026-02-08 22:00:45.000000
"""

from collections.abc import Sequence
from typing import Union

from alembic import op

# revision identifiers, used by Alembic.
revision: str = 'd4dfd8798fc0'
down_revision: Union[str, Sequence[str], None] = '919723f05419'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
	"""Insert users(id=1) if missing."""

	op.execute("""
		INSERT INTO users (id, daily_love_used)
		VALUES (1, 0)
		ON CONFLICT (id) DO NOTHING;
	""")


def downgrade() -> None:
	# This migration only guarantees singleton presence; avoid destructive downgrade.
	pass
