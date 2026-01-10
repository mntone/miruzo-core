from app.persist.jobs.base import BaseJobRepository
from app.persist.mixins.postgre import PostgreSQLUniqueViolationMixin


class PostgreSQLJobRepository(PostgreSQLUniqueViolationMixin, BaseJobRepository): ...
