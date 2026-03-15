import argparse
from datetime import timedelta

from app.config.environments import env
from app.databases import create_session, init_database
from app.domain.score.calculator import ScoreCalculator
from app.jobs.daily_decay import DailyDecayJob
from app.persist.jobs.factory import create_job_repository
from app.services.activities.daily_decay import DailyDecayRunner
from app.services.jobs.manager import JobManager
from app.services.settings.factory import build_daily_period_resolver


def parse_args() -> argparse.Namespace:
	parser = argparse.ArgumentParser(description='Run daily score decay job once.')
	parser.add_argument(
		'--force',
		action='store_true',
		help='Ignore the minimum interval guard and run immediately.',
	)
	return parser.parse_args()


def main() -> None:
	args = parse_args()
	init_database()

	min_interval = timedelta(minutes=3)
	if args.force:
		min_interval = timedelta()

	job_manager = JobManager(
		session_factory=create_session,
		job_repo_factory=create_job_repository,
		min_interval=min_interval,
	)

	job = DailyDecayJob(
		DailyDecayRunner(
			period_resolver=build_daily_period_resolver(
				day_start_offset=env.period.day_start_offset,
				initial_location=env.period.initial_location,
			),
			score_calculator=ScoreCalculator(env.score),
		),
		session_factory=create_session,
	)

	job_manager.try_run(job)


if __name__ == '__main__':
	main()
