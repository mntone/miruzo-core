# pyright: reportMissingTypeStubs=false

from datetime import time, timedelta

from apscheduler.job import Job as APJob
from apscheduler.schedulers.background import BackgroundScheduler

from app.jobs.protocol import Job
from app.services.jobs.manager import JobManager


def _to_time(offset: timedelta) -> time:
	total_seconds = int(offset.total_seconds()) % (24 * 3600)
	hours, remainder = divmod(total_seconds, 3600)
	minutes, seconds = divmod(remainder, 60)

	return time(hour=hours, minute=minutes, second=seconds)


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
	trigger_offset: timedelta,
) -> APJob:
	trigger_time = _to_time(trigger_offset)
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
