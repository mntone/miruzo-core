import argparse

import pytest

from importers.common.importer import _format_bytes
from importers.gataku_import import parse_ingest_mode

from app.models.enums import IngestMode


@pytest.mark.parametrize(
	('size', 'expected'),
	[
		(0, '0 B'),
		(999, '999 B'),
		(1024, '1.00 KB'),
		(10 * 1024, '10.0 KB'),
		(5 * 1024**2, '5.00 MB'),
		(10 * 1024**2, '10.0 MB'),
		(3 * 1024**3, '3.00 GB'),
		(12 * 1024**3, '12.0 GB'),
		(4.53 * 1024**4, '4.53 TB'),
		(16.49 * 1024**4, '16.5 TB'),
	],
)
def test_format_bytes_formats_sizes(size: float, expected: str) -> None:
	assert _format_bytes(int(round(size))) == expected


def test_parse_ingest_mode_accepts_known_values() -> None:
	assert parse_ingest_mode('copy') == IngestMode.COPY
	assert parse_ingest_mode('symlink') == IngestMode.SYMLINK


def test_parse_ingest_mode_rejects_invalid_value() -> None:
	with pytest.raises(argparse.ArgumentTypeError, match='Invalid mode'):
		parse_ingest_mode('bogus')
