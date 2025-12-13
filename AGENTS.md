# AGENTS.md

This document captures the non-negotiable guidelines for anyone—human or
automation agents—working on **miruzo-core** (the FastAPI/SQLModel backend).
Follow these instructions when setting up the environment, writing code,
testing, and preparing commits so that tasks stay consistent regardless of
who is executing them.

## Setup Command

- Python 3.13 (respect [`.python-version`](./.python-version) if present). Use
  `pyenv` or your preferred tool to match the pinned version. Code must remain
  compatible with Python 3.10+, so avoid newer-only language features. When
  defining enums, prefer `class Foo(int, Enum)` / `class Bar(str, Enum)` rather
  than Python 3.11+ helpers such as `IntEnum` or `StrEnum`.
- Install dependencies: `pip install -r requirements.txt`
- Copy [`.env.development`](./.env.development) to `.env` and adjust paths/DSNs as needed. Core reads
  all config via `pydantic-settings`.
- Start the API locally: `uvicorn app.main:app --reload`
- Rebuild the SQLite database or run importers via
  `python importers/import.py --help`
- Run format/lint (if configured): `ruff check app tests` and `ruff format`
- Execute tests: `pytest` (see Testing section for Docker-backed suites)

## General Code Style

- Type-hint everything that crosses module boundaries; prefer `Annotated[...]`
  for Pydantic/OpenAPI metadata
- Use tabs for indentation to match the rest of the repository
- Comments must be written in English and focus on intent or tricky behavior
- Keep SQLAlchemy/SQLModel expressions explicit—no implicit string SQL
- Markdown files must be wrapped at 80 characters
- Do not introduce ad-hoc helpers when an equivalent exists in
  `app/services/images/` or shared test utilities

## API & Services Conventions

- HTTP handlers live under `app/routers/*` and should delegate to services; do
  not perform DB/variant logic in routers
- Service classes (e.g., `ImageService`) should depend only on repository
  interfaces and pure helpers (`variants.py`, etc.)
- Repositories must inherit from `ImageRepository` and store their `Session`
  via `super().__init__(session)` so the shared helpers work
- Always add OpenAPI `title` and `description` via `Field` metadata for new API
  models; omit only when the parent schema already conveys the same meaning
- Enum/string settings should be normalized via validators in
  `app/core/settings.py`; avoid environment-specific branching outside that
  module
- Importers in `importers/common/*` must log all destructive actions and honor
  the `force` flag before deleting files

## Commits

- Follow Conventional Commits; keep the subject within 55 characters
- Group related backend changes (API schema + service + repo + tests) in one
  commit; avoid mixing importer work with unrelated API tweaks
- Use the GitHub issue and PR templates provided in `.github/` when filing or
  submitting changes

## Testing

- Default suite: `pytest`
- Run focused suites as needed:
  - `pytest tests/services/images/test_service.py` for service logic
  - `pytest tests/services/images/repository/test_sqlite.py` for SQLite repo
  - `pytest tests/services/images/repository/test_postgre.py` (requires Docker;
    only needed when touching PostgreSQL code—the rest of the suite can run
    without Docker; this test pulls `postgres:18-alpine`)
  - `pytest tests/services/images/test_variants.py` for variant helpers
  - Importer-specific tests live under `tests/importers`
- Docker-based tests must clean up containers; reuse the helper fixtures
  already committed (do not hand-roll `docker run` invocations)
- Use in-memory SQLite for unit tests unless the code path explicitly requires
  a different backend
- Pure helper functions (variant parsing, query splitting, etc.) must be
  covered by unit tests; impure logic should be split into testable pieces
- Avoid network access during tests; mock external calls or use local fixtures
