from sqlalchemy.exc import IntegrityError


class PostgreSQLUniqueViolationMixin:
	def _is_unique_violation(self, error: IntegrityError) -> bool:
		orig = error.orig
		if orig is None:
			return False

		pgcode = getattr(orig, 'pgcode', None)
		if pgcode == '23505':
			return True

		try:
			from psycopg2.errors import UniqueViolation
		except ImportError:
			return False

		return isinstance(orig, UniqueViolation)
