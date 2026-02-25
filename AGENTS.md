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
- SQLite 3.35.0+ when using the SQLite backend (`RETURNING` support required).
  Verify Python-linked SQLite via
  `python -c "import sqlite3; print(sqlite3.sqlite_version)"`.
- Install dependencies:
  `cd miruzo-py && pip install -r requirements.txt`
- Copy [`miruzo-py/.env.development`](./miruzo-py/.env.development) to
  `miruzo-py/.env` and adjust paths/DSNs as needed. Core reads all config via
  `pydantic-settings`.
- Start the API locally:
  `cd miruzo-py && python -m scripts.api --dev`
- Rebuild the SQLite database or run importers via
  `cd miruzo-py && python -m scripts.gataku_import --help`
- Run format/lint (if configured):
  `cd miruzo-py && ruff check app tests && ruff format`
- Execute tests: `cd miruzo-py && pytest` (see Testing section for
  Docker-backed suites)

## General Code Style

- Type-hint everything that crosses module boundaries; prefer `Annotated[...]`
  for Pydantic/OpenAPI metadata
- Use tabs for indentation to match the rest of the repository
- Comments must be written in English and focus on intent or tricky behavior
- Keep SQLAlchemy/SQLModel expressions explicit—no implicit string SQL
- Markdown files must be wrapped at 80 characters
- Do not introduce ad-hoc helpers when an equivalent exists in
  `miruzo-py/app/services/images/` or shared test utilities
- Only add `from __future__ import annotations` (or other future imports) when a
  file truly needs it; if a newer language feature would significantly improve
  performance or readability, raise it so we can evaluate bumping the minimum
  Python version together

## API & Services Conventions

- HTTP handlers live under `miruzo-py/app/routers/*` and should delegate to
  services; do not perform DB/variant logic in routers
- Repository protocols and factories live under `miruzo-py/app/persist/*`;
  services should depend on those protocols rather than embedding SQL directly
- Service classes (e.g., `ImageService`) should depend only on repository
  interfaces and pure helpers (`variants.py`, etc.)
- Repositories should inherit from the relevant base repository class and
  store their `Session` via `super().__init__(session)` so shared helpers work
- Always add OpenAPI `title` and `description` via `Field` metadata for new API
  models; omit only when the parent schema already conveys the same meaning
- Enum/string settings should be normalized via validators in
  `miruzo-py/app/config/environments.py`; avoid environment-specific branching
  outside that module
- Importers in `scripts/importers/common/*` must log all destructive actions and
  honor the `force` flag before deleting files
- List APIs must use `ImageListService` + `ImageListRepository` and apply
  `limit + 1` pagination with `paginator.slice_with_cursor_latest` for latest
  and `paginator.slice_with_tuple_cursor` for the other list endpoints

## Commits

- Follow Conventional Commits; keep the subject within 55 characters
- Group related backend changes (API schema + service + repo + tests) in one
  commit; avoid mixing importer work with unrelated API tweaks
- Use the GitHub issue and PR templates provided in `.github/` when filing or
  submitting changes

## Testing

- Default suite: `pytest`
- Run focused suites as needed:
  - `cd miruzo-py && pytest tests/services/images/test_service.py` for service
    logic
  - `cd miruzo-py && pytest tests/services/images/repository/test_sqlite.py`
    for SQLite repo
  - `cd miruzo-py && pytest tests/services/images/repository/test_postgre.py`
    (requires Docker; only needed when touching PostgreSQL code—the rest of the
    suite can run without Docker; this test pulls `postgres:18-alpine`)
  - `cd miruzo-py && pytest tests/services/images/test_variants.py` for
    variant helpers
  - Importer-specific tests live under `miruzo-py/tests/importers`
- Docker-based tests must clean up containers; reuse the helper fixtures
  already committed (do not hand-roll `docker run` invocations)
- Use in-memory SQLite for unit tests unless the code path explicitly requires
  a different backend
- Repository SQL tests live under `miruzo-py/tests/persist/*`; service-level
  list tests should be thin wiring tests that stub the repository and mapper
- Pure helper functions (variant parsing, query splitting, etc.) must be
  covered by unit tests; impure logic should be split into testable pieces
- Avoid network access during tests; mock external calls or use local fixtures
