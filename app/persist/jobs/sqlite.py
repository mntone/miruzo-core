from app.persist.jobs.base import BaseJobRepository
from app.persist.mixins.sqlite import SQLiteUniqueViolationMixin


class SQLiteJobRepository(SQLiteUniqueViolationMixin, BaseJobRepository): ...
