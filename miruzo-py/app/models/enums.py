from enum import Enum
from typing import final


@final
class IngestMode(int, Enum):
	COPY = 0
	SYMLINK = 1


@final
class ProcessStatus(int, Enum):
	PROCESSING = 0
	FINISHED = 1


@final
class VisibilityStatus(int, Enum):
	PRIVATE = 0
	PUBLIC = 1


@final
class ExecutionStatus(int, Enum):
	SUCCESS = 0
	UNKNOWN_ERROR = 1
	DB_ERROR = 2
	IO_ERROR = 3
	IMAGE_ERROR = 4


@final
class ImageKind(int, Enum):
	UNSPECIFIED = 0
	PHOTO = 1
	ILLUST = 2
	GRAPHIC = 3
