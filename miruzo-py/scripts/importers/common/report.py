from dataclasses import dataclass
from typing import IO

from app.models.records import IngestRecord


@dataclass(slots=True)
class ImportStats:
	read: int = 0
	ingested: int = 0
	invalid: int = 0
	missing: int = 0
	fallback: int = 0


def _format_bytes(size: int) -> str:
	"""Convert a size in bytes into a human-friendly string."""

	thresholds = [
		(1024**4, 'TB'),
		(1024**3, 'GB'),
		(1024**2, 'MB'),
		(1024, 'KB'),
	]
	for factor, suffix in thresholds:
		if size >= factor:
			value = size / factor
			if value < 10:
				return f'{value:.2f} {suffix}'
			return f'{value:.1f} {suffix}'
	return f'{size} B'


class ProgressReporter:
	def __init__(
		self,
		*,
		report_variants: bool,
		stream: IO[str] | None = None,
		progress_interval: int = 10,
	) -> None:
		self._report_variants = report_variants
		self._stream = stream
		self._progress_interval = progress_interval

	def _write(self, line: str) -> None:
		if self._stream is None:
			print(line)
			return
		print(line, file=self._stream)

	def maybe_report_variants(self, record: IngestRecord) -> None:
		if not self._report_variants:
			return
		image = record.image
		if image is None:
			return

		self._write(f'[importer] variant report ({record.relative_path}):')
		header = f'{"Label":<10} {"Resolution":<12} {"Size":>10} {"Ratio":<12}'
		self._write(header)
		self._write('-' * len(header))

		original = image.original
		original_width = original['width']
		original_height = original['height']
		original_resolution = f'{original_width}x{original_height}'
		original_size = original['bytes']
		self._write(
			f'{"l0orig":<10} {original_resolution:<12} {_format_bytes(original_size):>10} {"n/a":<12}',
		)

		for variant in image.variants:
			label = 'l' + variant['layer_id'].__str__() + 'w' + variant['width'].__str__()
			width = variant['width']
			height = variant['height']

			size = variant['bytes']
			size_str = _format_bytes(size or 0)

			delta = size - original_size
			delta_str = _format_bytes(abs(delta))

			ratio = (size / original_size) * 100
			ratio_sign = '+' if delta >= 0 else '-'
			ratio_str = f'{ratio:.1f}% ({ratio_sign}{delta_str})'

			resolution = f'{width}x{height}'
			self._write(f'{label:<10} {resolution:<12} {size_str:>10} {ratio_str:<12}')

	def report_progress(self, stats: ImportStats, *, force: bool = False) -> None:
		if not force and stats.read % self._progress_interval != 0:
			return
		line = (
			'[importer] progress: '
			f'read={stats.read}, ingested={stats.ingested}, invalid={stats.invalid}, '
			f'missing={stats.missing}, fallback={stats.fallback}'
		)
		self._write(line)

	def report_summary(self, stats: ImportStats) -> None:
		line = (
			'[importer] summary: '
			f'read={stats.read}, ingested={stats.ingested}, invalid={stats.invalid}, '
			f'missing={stats.missing}, fallback={stats.fallback}'
		)
		self._write(line)
