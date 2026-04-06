import io

import pytest

from scripts.importers.common.report import ImportStats, ProgressReporter, _format_bytes
from tests.fixtures.image import make_image_fixture
from tests.fixtures.ingest import make_ingest_fixture


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
	ingest = make_ingest_fixture(6)
	image = make_image_fixture(6)

	reporter.maybe_report_variants((ingest, image))

	output = stream.getvalue().splitlines()
	assert output[0].startswith('[importer] variant report (')
	assert output[1].startswith('Label      Resolution')
	assert output[2].startswith('---------------------------')
	assert output[3].startswith('l0orig     1024x768')
	assert output[4].startswith('l1w320     320x240')
	assert output[5].startswith('l1w480     480x360')
	assert output[6].startswith('l1w640     640x480')
	assert output[7].startswith('l1w960     960x720')
	assert output[8].startswith('l9w320     320x240')
