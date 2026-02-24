import argparse
import os

import uvicorn


def parse_args() -> argparse.Namespace:
	parser = argparse.ArgumentParser(description='Run miruzo API server.')
	parser.add_argument(
		'--dev',
		action='store_true',
		help='Run in development mode with auto reload.',
	)
	parser.add_argument(
		'--host',
		default=os.getenv('MIRUZO_HOST', '0.0.0.0'),
		help='Host address to bind.',
	)
	parser.add_argument(
		'--port',
		type=int,
		default=int(os.getenv('MIRUZO_PORT', '1024')),
		help='Port number to bind.',
	)
	parser.add_argument(
		'--env-file',
		default=None,
		help='Path to env file passed to uvicorn.',
	)

	return parser.parse_args()


def _run(args: argparse.Namespace) -> None:
	dev_mode = args.dev
	env_file = args.env_file or ('.env.development' if dev_mode else None)

	uvicorn.run(
		'app.main:app',
		host=args.host,
		port=args.port,
		env_file=env_file,
		reload=dev_mode,
		reload_dirs=['./app'] if dev_mode else None,
	)


def main() -> None:
	args = parse_args()
	_run(args)


def main_dev() -> None:
	args = parse_args()
	args.dev = True
	_run(args)


if __name__ == '__main__':
	main()
