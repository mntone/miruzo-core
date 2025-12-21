from enum import Enum


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
	SUCCESS = 0
	UNKNOWN_ERROR = 1
	IO_ERROR = 2


class ImageKind(int, Enum):
	PHOTO = 0
	ILLUST = 1
	GRAPHIC = 2
