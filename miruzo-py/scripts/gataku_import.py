import argparse

from scripts.importers.common.importer import import_jsonl

from app.models.enums import IngestMode

_MODE_MAP = {
	'copy': IngestMode.COPY,
	'symlink': IngestMode.SYMLINK,
}


def parse_ingest_mode(value: str) -> IngestMode:
	try:
		return _MODE_MAP[value]
	except KeyError as exc:
		raise argparse.ArgumentTypeError(f'Invalid mode: {value}') from exc


def parse_args() -> argparse.Namespace:
	parser = argparse.ArgumentParser(description='Import miruzo images from gataku JSONL outputs.')
	parser.add_argument(
		'--jsonl-path',
		default='../gataku/out/hashdb.jsonl',
		help='Path to the gataku hashdb jsonl file.',
	)
	parser.add_argument('--limit', type=int, default=100, help='Maximum number of records to import.')
	parser.add_argument(
		'--mode',
		type=parse_ingest_mode,
		default=IngestMode.SYMLINK,
		help='How to place images into the media directory. (copy|symlink)',
	)
	parser.add_argument('--force', action='store_true', help='Skip confirmation prompts during import.')
	parser.add_argument(
		'--report-variants',
		action='store_true',
		help='Show thumbnail generation report during import.',
	)
	return parser.parse_args()


def main() -> None:
	args = parse_args()
	import_jsonl(
		jsonl_path=args.jsonl_path,
		limit=args.limit,
		mode=args.mode,
		force=args.force,
		report_variants=args.report_variants,
	)


if __name__ == '__main__':
	main()
