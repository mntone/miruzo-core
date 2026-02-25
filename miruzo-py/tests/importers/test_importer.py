import argparse

import pytest

from scripts.gataku_import import parse_ingest_mode

from app.models.enums import IngestMode


def test_parse_ingest_mode_accepts_known_values() -> None:
	assert parse_ingest_mode('copy') == IngestMode.COPY
	assert parse_ingest_mode('symlink') == IngestMode.SYMLINK


def test_parse_ingest_mode_rejects_invalid_value() -> None:
	with pytest.raises(argparse.ArgumentTypeError, match='Invalid mode'):
		parse_ingest_mode('bogus')
