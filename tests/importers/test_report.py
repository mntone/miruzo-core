import io
from datetime import datetime, timezone

import pytest

from importers.common.report import ImportStats, ProgressReporter, _format_bytes

from app.models.records import ImageRecord, IngestRecord
from app.models.types import VariantEntry


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


def _build_record() -> IngestRecord:
	original: VariantEntry = {
		'rel': 'l0orig/foo.jpg',
		'layer_id': 0,
		'format': 'jpeg',
		'codecs': None,
		'bytes': 100,
		'width': 10,
		'height': 10,
		'quality': None,
	}
	variant: VariantEntry = {
		'rel': 'l1w5/foo.jpg',
		'layer_id': 1,
		'format': 'webp',
		'codecs': 'vp8',
		'bytes': 50,
		'width': 5,
		'height': 5,
		'quality': 80,
	}

	image = ImageRecord(
		ingest_id=1,
		ingested_at=datetime.now(timezone.utc),
		original=original,
		fallback=None,
		variants=[variant],
	)

	record = IngestRecord(
		relative_path='l0orig/foo.jpg',
		fingerprint='a' * 64,
		ingested_at=datetime.now(timezone.utc),
		captured_at=datetime.now(timezone.utc),
	)
	record.image = image

	return record


def test_progress_reporter_writes_progress_and_summary() -> None:
	stream = io.StringIO()
	reporter = ProgressReporter(report_variants=False, stream=stream, progress_interval=2)
	stats = ImportStats(read=1)

	reporter.report_progress(stats)
	assert stream.getvalue() == ''

	stats.read = 2
	stats.ingested = 1
	reporter.report_progress(stats)
	reporter.report_summary(stats)

	lines = stream.getvalue().strip().splitlines()
	assert len(lines) == 2
	assert lines[0].startswith('[importer] progress:')
	assert 'read=2' in lines[0]
	assert lines[1].startswith('[importer] summary:')
	assert 'ingested=1' in lines[1]


def test_progress_reporter_writes_variant_report() -> None:
	stream = io.StringIO()
	reporter = ProgressReporter(report_variants=True, stream=stream)
	record = _build_record()

	reporter.maybe_report_variants(record)

	output = stream.getvalue()
	assert '[importer] variant report' in output
	assert 'l0orig' in output
	assert 'l1w5' in output
