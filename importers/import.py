import argparse

from importers.common.importer import import_jsonl


def parse_args() -> argparse.Namespace:
	parser = argparse.ArgumentParser(description='Import miruzo images from gataku outputs.')
	parser.add_argument(
		'--jsonl-path',
		default='../gataku/out/hashdb.jsonl',
		help='Path to the gataku hashdb jsonl file.',
	)
	parser.add_argument('--static-dir', default='./static', help='Directory to populate with image files.')
	parser.add_argument('--limit', type=int, default=100, help='Maximum number of records to import.')
	parser.add_argument(
		'--mode',
		choices=['copy', 'symlink'],
		default='symlink',
		help='How to place images into the static directory.',
	)
	parser.add_argument(
		'--orig-dir',
		default='gataku',
		help='Name of the directory (under static/) that stores original assets.',
	)
	parser.add_argument('--force', action='store_true', help='Skip confirmation prompts during import.')
	parser.add_argument(
		'--report-variants',
		action='store_true',
		help='Show thumbnail generation report during import.',
	)
	parser.add_argument(
		'--repair',
		action='store_true',
		help='Repair database entries without regenerating thumbnails.',
	)
	return parser.parse_args()


def main() -> None:
	args = parse_args()
	import_jsonl(
		jsonl_path=args.jsonl_path,
		static_dir=args.static_dir,
		limit=args.limit,
		mode=args.mode,
		original_subdir=args.orig_dir,
		force=args.force,
		report_variants=args.report_variants,
		repair=args.repair,
	)


if __name__ == '__main__':
	main()
