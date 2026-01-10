import logging
import os
from datetime import datetime, timezone
from pathlib import Path

import pytest

from importers.common.ingest_time import resolve_captured_at


def test_resolve_captured_at_prefers_created_at(tmp_path: Path) -> None:
	target = tmp_path / 'sample.jpg'
	target.write_text('x', encoding='utf-8')
	os.utime(target, (1_700_000_000, 1_700_000_000))

	expected = datetime(2024, 1, 2, 3, 4, 5, tzinfo=timezone.utc)
	captured_at, used_fallback, warned_fallback = resolve_captured_at(
		created_at_value=expected.isoformat(),
		src_path=target,
		warned_fallback=False,
	)

	assert captured_at == expected
	assert used_fallback is False
	assert warned_fallback is False


def test_resolve_captured_at_falls_back_on_invalid_created_at(
	tmp_path: Path,
	caplog: pytest.LogCaptureFixture,
) -> None:
	target = tmp_path / 'sample.jpg'
	target.write_text('x', encoding='utf-8')
	os.utime(target, (1_700_000_123, 1_700_000_123))

	caplog.set_level(logging.WARNING)
	captured_at, used_fallback, warned_fallback = resolve_captured_at(
		created_at_value='not-a-date',
		src_path=target,
		warned_fallback=False,
	)

	expected = datetime.fromtimestamp(1_700_000_123, tz=timezone.utc)
	assert captured_at == expected
	assert used_fallback is True
	assert warned_fallback is True
	assert 'invalid created_at detected' in caplog.text


def test_resolve_captured_at_falls_back_on_missing_created_at(
	tmp_path: Path,
	caplog: pytest.LogCaptureFixture,
) -> None:
	target = tmp_path / 'sample.jpg'
	target.write_text('x', encoding='utf-8')
	os.utime(target, (1_700_000_456, 1_700_000_456))

	caplog.set_level(logging.WARNING)
	captured_at, used_fallback, warned_fallback = resolve_captured_at(
		created_at_value=None,
		src_path=target,
		warned_fallback=False,
	)

	expected = datetime.fromtimestamp(1_700_000_456, tz=timezone.utc)
	assert captured_at == expected
	assert used_fallback is True
	assert warned_fallback is True
	assert 'missing created_at detected' in caplog.text
