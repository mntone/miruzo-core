from app.databases.mixins.postgre import PostgreSQLUniqueViolationMixin
from app.persist.jobs.base import BaseJobRepository


class PostgreSQLJobRepository(PostgreSQLUniqueViolationMixin, BaseJobRepository): ...
