"""
Add love action index

Revision ID: 919723f05419
Revises: a46b415897c4
Create Date: 2026-02-07 12:04:19.000000
"""

from collections.abc import Sequence
from typing import Union

from alembic import op

# revision identifiers, used by Alembic.
revision: str = '919723f05419'
down_revision: Union[str, Sequence[str], None] = 'a46b415897c4'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
	"""Add love action index for effective love queries."""

	dialect_name = op.get_bind().dialect.name

	if dialect_name == 'postgresql':
		op.execute("""
			CREATE INDEX ix_actions_love
			ON actions (ingest_id, occurred_at DESC, id DESC)
			INCLUDE (kind)
			WHERE kind IN (13, 14);
		""")
	elif dialect_name == 'sqlite':
		op.execute("""
			CREATE INDEX ix_actions_love
			ON actions (ingest_id, occurred_at DESC, id DESC)
			WHERE kind IN (13, 14);
		""")
	else:
		raise RuntimeError(f'Unsupported dialect for migration {revision}: {dialect_name}')


def downgrade() -> None:
	dialect_name = op.get_bind().dialect.name

	if dialect_name in {'postgresql', 'sqlite'}:
		op.drop_index('ix_actions_love', table_name='actions')
	else:
		raise RuntimeError(f'Unsupported dialect for migration {revision}: {dialect_name}')
