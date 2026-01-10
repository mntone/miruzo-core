from pathlib import Path

from importers.common.readers.jsonl import JsonlReader


def _write_jsonl(path: Path, lines: list[str]) -> None:
	content = '\n'.join(lines) + '\n'
	path.write_text(content, encoding='utf-8')


def test_jsonl_reader_parses_rows_and_marks_invalid_lines(tmp_path: Path) -> None:
	path = tmp_path / 'data.jsonl'
	sha_primary = 'a' * 64
	sha_secondary = 'b' * 64
	lines = [
		(
			'{"filepath": "foo/bar.jpg", "sha256": "'
			+ sha_primary
			+ '", "created_at": "2024-01-02T03:04:05+00:00"}'
		),
		'not-json',
		'{"filepath": "/abs/path.jpg", "sha256": "' + sha_secondary + '"}',
	]
	_write_jsonl(path, lines)

	reader = JsonlReader(path)
	rows = list(reader.read())

	assert len(rows) == 3
	assert rows[1] is None

	first = rows[0]
	assert first is not None
	assert first.filepath == Path('foo/bar.jpg')
	assert first.sha256 == sha_primary
	assert first.created_at == '2024-01-02T03:04:05+00:00'

	third = rows[2]
	assert third is not None
	assert third.filepath == Path('/abs/path.jpg')
	assert third.sha256 == sha_secondary
	assert third.created_at is None


def test_jsonl_reader_respects_limit(tmp_path: Path) -> None:
	path = tmp_path / 'data.jsonl'
	sha_primary = 'a' * 64
	sha_secondary = 'b' * 64
	lines = [
		'{"filepath": "one.jpg", "sha256": "' + sha_primary + '"}',
		'{"filepath": "two.jpg", "sha256": "' + sha_secondary + '"}',
	]
	_write_jsonl(path, lines)

	reader = JsonlReader(path)
	rows = list(reader.read(limit=1))

	assert len(rows) == 1
	assert rows[0] is not None
	assert rows[0].filepath == Path('one.jpg')
