from app.databases.mixins.sqlite import SQLiteUniqueViolationMixin
from app.persist.jobs.base import BaseJobRepository


class SQLiteJobRepository(SQLiteUniqueViolationMixin, BaseJobRepository): ...
