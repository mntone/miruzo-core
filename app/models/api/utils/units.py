from math import ceil


def bytes_to_manbytes(size_bytes: int) -> int:
	return ceil(size_bytes / 10_000)
