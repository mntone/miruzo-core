from enum import Enum, auto


class IngestMode(int, Enum):
	COPY = 0
	SYMLINK = 1


class ProcessStatus(int, Enum):
	PROCESSING = 0
	FINISHED = 1


class VisibilityStatus(int, Enum):
	PRIVATE = 0
	PUBLIC = 1


class ExecutionStatus(int, Enum):
	SUCCESS = auto()
	UNKNOWN_ERROR = auto()
	DB_ERROR = auto()
	IO_ERROR = auto()
	IMAGE_ERROR = auto()


class ImageKind(int, Enum):
	UNSPECIFIED = 0
	PHOTO = 1
	ILLUST = 2
	GRAPHIC = 3
