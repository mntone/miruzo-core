# pyright: reportMissingTypeStubs=false

from datetime import time

from apscheduler.job import Job as APJob
from apscheduler.schedulers.background import BackgroundScheduler

from app.jobs.protocol import Job
from app.services.jobs.manager import JobManager


def create_scheduler() -> BackgroundScheduler:
	scheduler = BackgroundScheduler(timezone='UTC')
	return scheduler


def _dispatch_job(job_manager: JobManager, job: Job) -> None:
	job_manager.try_run(job)


def register_daily_job(
	*,
	scheduler: BackgroundScheduler,
	job_manager: JobManager,
	job: Job,
	trigger_time: time,
) -> APJob:
	ap_job = scheduler.add_job(  # pyright: ignore[reportUnknownMemberType]
		_dispatch_job,
		trigger='cron',
		hour=trigger_time.hour,
		minute=trigger_time.minute,
		second=trigger_time.second,
		args=[job_manager, job],
		id=job.name,
		replace_existing=True,
	)

	return ap_job
