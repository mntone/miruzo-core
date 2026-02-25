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
	PHOTO = 0
	ILLUST = 1
	GRAPHIC = 2


class ActionKind(int, Enum):
	UNKNOWN = 0
	DECAY = 1
	VIEW = 11
	MEMO = 12
	LOVE = 13
	LOVE_CANCELED = 14
	HALL_OF_FAME_ADDED = 15
	HALL_OF_FAME_REMOVED = 16
