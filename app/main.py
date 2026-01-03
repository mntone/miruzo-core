from contextlib import asynccontextmanager
from datetime import timedelta
from logging import getLogger
from typing import AsyncGenerator

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from fastapi.staticfiles import StaticFiles

from app.config.environments import env
from app.database import create_session, init_database
from app.domain.score.calculator import ScoreCalculator
from app.infrastructures.scheduler import create_scheduler, register_daily_job
from app.jobs.score_decay import ScoreDecayJob
from app.routers.health import router as health
from app.routers.images import router as images
from app.services.activities.score_decay import ScoreDecayRunner
from app.services.ingests.bootstrap import ensure_ingest_layout
from app.services.jobs.manager import JobManager
from app.services.jobs.repository.factory import create_job_repository

log = getLogger('uvicorn.error')


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None]:
	ensure_ingest_layout(env)
	init_database()
	log.info(f'Starting miruzo API in {env.environment.value} mode')

	scheduler = create_scheduler()

	job_manager = JobManager(
		session_factory=create_session,
		job_repo_factory=create_job_repository,
		min_interval=timedelta(minutes=3),
	)

	score_decay = ScoreDecayJob(
		ScoreDecayRunner(
			score_calculator=ScoreCalculator(env.score),
			daily_reset_at=env.time.daily_reset_at,
			base_timezone=env.base_timezone,
		),
		session_factory=create_session,
	)

	register_daily_job(
		scheduler=scheduler,
		job_manager=job_manager,
		job=score_decay,
		trigger_time=env.time.daily_reset_at,
	)

	scheduler.start()  # pyright: ignore[reportUnknownMemberType]
	yield
	scheduler.shutdown(wait=False)  # pyright: ignore[reportUnknownMemberType]


app = FastAPI(title='miruzo API', lifespan=lifespan)
app.add_middleware(GZipMiddleware, minimum_size=1000, compresslevel=9)
app.include_router(images, prefix='/api')
app.include_router(health, prefix='/api')

app.add_middleware(
	CORSMiddleware,
	allow_origins=['*'],
	allow_credentials=True,
	allow_methods=['*'],
	allow_headers=['*'],
)
app.mount(
	env.public_media_root,
	StaticFiles(directory=env.media_root, follow_symlink=True),
	name='media',
)
