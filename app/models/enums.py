from enum import Enum


class ImageStatus(int, Enum):
	"""Lifecycle state for an image record."""

	ACTIVE = 0  #: Image is available and can be served.
	DELETED = 1  #: Image was removed intentionally.
	MISSING = 2  #: Image metadata exists but the asset is missing.
