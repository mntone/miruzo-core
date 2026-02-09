[&lsaquo; Back to README](./README.md)

# Contributing to miruzo-core

## 🧭 Overview

miruzo-core is the FastAPI/SQLModel backend that exposes the miruzo photo
archive over a REST API and drives the importer pipeline. This document explains
how to set up a local environment, which conventions to follow, and how to
propose changes. For day-to-day commands (install, dev server, tests) refer to
[README.md](./README.md), and for coding standards see [AGENTS.md](./AGENTS.md).

## 🛠️ Prerequisites

- Python 3.13 (respect [`.python-version`](./.python-version) if present).
  Install dependencies via `pip install -r requirements.txt`.
  All contributions must remain compatible with Python 3.10+, so avoid using
  language/library features newer than 3.10 (e.g., prefer `class Foo(int, Enum)`
  / `class Bar(str, Enum)` instead of newer helpers such as `IntEnum` /
  `StrEnum`).
- Copy [`.env.development`](./.env.development) to `.env` and configure paths/DSNs so API and importer
  flows can run locally.
- Docker is only required for PostgreSQL repository tests
  (`tests/services/images/repository/test_postgre.py`). All other development
  can be done without Docker.
- macOS, Linux, and WSL are supported environments. Ensure `uvicorn` runs via
  `uvicorn app.main:app --reload` and that SQLite 3.35.0+ is available (for
  importer + tests; `RETURNING` support is required).
  Verify Python-linked SQLite with
  `python -c "import sqlite3; print(sqlite3.sqlite_version)"`.
- To exercise importer pipelines, prepare the directories referenced by
  `settings.gataku_root` / `settings.assets_root` (see [README.md](./README.md)).

## 🔁 Workflow

- Follow Conventional Commits (English prefixes such as `feat:`, `fix:`, etc.)
  and keep each commit focused on a single logical change.
- Large or potentially breaking backend work should start as a GitHub issue or
  discussion so scope and migrations can be agreed upon.
- Use the provided GitHub issue/PR templates; fill out the checklists so
  reviewers know what was verified.
- Before requesting review run `pytest` (plus any Docker-dependent suites as
  needed) and `ruff check`. Include OpenAPI/schema changes in your PR summary.

## 🎨 Code style

- `AGENTS.md` is the canonical source of formatting and architectural rules.
  Highlights: tabs for indentation, type-hint all exported surfaces, add OpenAPI
  `title`/`description` metadata for new models, and keep routers thin by
  delegating to services.
- Do not ignore lint errors—fix them or adjust configuration via PR if a rule is
  truly incompatible.
- For shared utilities (variant helpers, repository base classes, importer
  workflows), extend the existing modules in `app/services/images/` or
  `importers/common/` instead of adding ad-hoc versions.

## 🧪 Testing

- We use pytest. Run `pytest` for a full pass or target suites such as
  `pytest tests/services/images/test_service.py` or
  `pytest tests/services/images/repository/test_sqlite.py` while iterating.
- PostgreSQL repository tests (`tests/services/images/repository/test_postgre.py`)
  require Docker and pull `postgres:18-alpine`. Run them only when changing
  Postgres-specific code paths.
- Pure helpers (variant parsing, query models, repository filters) must have
  unit tests unless doing so would add unreasonable complexity. Split impure
  logic into testable helpers when practical.
- Reuse shared fixtures/mocks under `tests/` before writing local ad hoc ones.

## 🐛 Reporting issues

- File bug reports and feature requests through GitHub Issues using the provided
  templates. Include reproduction steps, expected vs. actual behavior, API
  endpoints, and relevant logs.
- Sensitive bugs can be disclosed privately to *mntone* via the contact links in
  `README.md`.

## 📜 License notice

By contributing to miruzo-core you agree to license your work under GPLv3 (same
as the project). Submit only code and assets that you are allowed to relicense
under GPLv3. Verify the compatibility of any third-party dependency before
introducing it.
