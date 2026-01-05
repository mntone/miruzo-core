import argparse
from datetime import timedelta

from app.config.environments import env
from app.database import create_session, init_database
from app.domain.score.calculator import ScoreCalculator
from app.jobs.daily_decay import DailyDecayJob
from app.services.activities.daily_decay import DailyDecayRunner
from app.services.jobs.manager import JobManager
from app.services.jobs.repository.factory import create_job_repository


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
			score_calculator=ScoreCalculator(env.score),
			daily_reset_at=env.time.daily_reset_at,
			base_timezone=env.base_timezone,
		),
		session_factory=create_session,
	)

	job_manager.try_run(job)


if __name__ == '__main__':
	main()
