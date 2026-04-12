# AGENTS.md

This document captures the non-negotiable guidelines for anyone—human or
automation agents—working on **miruzo-core** (the Go API backend and Python
ingest tooling). Follow these instructions when setting up the environment,
writing code, testing, and preparing commits so that tasks stay consistent
regardless of who is executing them.

## Setup Command

- Go 1.26.x (respect [`miruzo/go.mod`](./miruzo/go.mod)).
- Python 3.10+ for ingest tooling (3.13 recommended for local development).
- Database backend version requirements:
  - MySQL `8.0.16+` (`CHECK` support required)
    - Note: MySQL support in Go API is planned but not implemented yet
  - PostgreSQL `14+`
  - SQLite `3.37.0+` (`RETURNING` and `STRICT` support required)
- Python ingest database drivers (`miruzo-py`):
  - MySQL: `mysqlclient` (`MySQLdb`, DSN: `mysql+mysqldb://...`)
  - PostgreSQL: `psycopg3` (`psycopg`, DSN: `postgresql+psycopg://...`)
  - SQLite: `sqlite3` (stdlib, DSN: `sqlite:///...`)
- Install dependencies:
  - `cd miruzo && go mod download`
  - `cd miruzo-py && pip install -r requirements.txt`
- OS-specific setup (required for DB drivers / Go tools):
  - Linux (Debian/Ubuntu):
    - Go tools:
      `cd miruzo && make tools`
    - PostgreSQL Python driver dependencies:
      `sudo apt install -y libpq-dev && pip install psycopg[c,pool]`
    - MySQL Python driver dependencies:
      `sudo apt install -y build-essential default-libmysqlclient-dev pkg-config`
      and `pip install mysqlclient`
  - macOS:
    - Go tools (recommended):
      `brew install gopls delve go-air sqlc goreleaser`
    - Go tools (alternative):
      `cd miruzo && make tools`
    - PostgreSQL Python driver dependencies:
      `brew install libpq && pip install psycopg[c,pool]`
    - MySQL Python driver dependencies:
      `brew install mysql pkg-config && pip install mysqlclient`
- Configure runtime files:
  - API: use `miruzo/config.yaml` (`miruzo/internal/app/config.sample.yaml`
    can be used as the base).
    - Current Go API backends are `sqlite` and `postgres`.
      `database.backend=mysql` is not supported yet.
  - Ingest: copy [`miruzo-py/.env.development`](./miruzo-py/.env.development)
    to `miruzo-py/.env` and set:
    - `ENVIRONMENT` (`development` or `production`)
    - `DATABASE_BACKEND` (`sqlite`, `postgres`, or `mysql`)
    - `DATABASE_URL`
      - SQLite: `sqlite:///...`
      - PostgreSQL: `postgresql+psycopg://...`
      - MySQL: `mysql+mysqldb://...`
    - Path-related variables (`MEDIA_ROOT`, `PUBLIC_MEDIA_ROOT`,
      `GATAKU_ROOT`, `GATAKU_ASSETS_ROOT`, `GATAKU_SYMLINK_DIRNAME`) may use
      defaults on first setup and be customized later if required.
- Optional test database DSN environment variables:
  - Go tests: `MIRUZO_TEST_POSTGRES_URL`
  - Python tests:
    - `MIRUZO_PY_TEST_MYSQL_URL`
    - `MIRUZO_PY_TEST_POSTGRES_URL`
- Run format/lint:
  - Go: `cd miruzo && go test ./...` (run gofmt/goimports as needed before
    commit).
  - Python: `cd miruzo-py && ruff check app scripts tests && ruff format`.
- Execute tests:
  - `cd miruzo && go test ./...`
  - `cd miruzo-py && pytest`

## Operations Commands

- Run CLI commands from `miruzo/`.
- Start API locally:
  - `cd miruzo && air`
  - `cd miruzo && go run ./cmd/miruzo-api`
- Run importer help:
  - `cd miruzo-py && python -m scripts.gataku_import --help`
- Migration commands:
  - `cd miruzo && go run ./cmd/miruzo-cli migrate up`
  - `cd miruzo && go run ./cmd/miruzo-cli migrate down [N]`
  - `cd miruzo && go run ./cmd/miruzo-cli migrate goto V`
  - `cd miruzo && go run ./cmd/miruzo-cli migrate version`
  - `cd miruzo && go run ./cmd/miruzo-cli migrate force V`
- `migrate force V` is for recovery workflows only. Do not use it in normal
  migrations because it changes version state without applying SQL changes.
- Before destructive operations (`migrate down`, `migrate force`, importer
  with `--force`), confirm backup/snapshot availability and impact scope.
- Maintenance jobs:
  - `cd miruzo && go run ./cmd/miruzo-cli job daily-decay`
- If you already built binaries (`cd miruzo && make build`), you may run the
  same commands via `./bin/miruzo-cli ...`.

## General Code Style

- Keep types explicit across module boundaries (Go public APIs and Python
  exported surfaces).
- Use tabs for indentation to match the rest of the repository
- Comments must be written in English and focus on intent or tricky behavior
- Keep SQL explicit:
  - Go: SQL files under `miruzo/internal/database/*/queries/` with sqlc.
  - Python: SQLAlchemy Core expressions (no ORM model magic, no ad-hoc string
    SQL).
- Markdown files must be wrapped at 80 characters
- Do not introduce ad-hoc helpers when an equivalent exists in
  `miruzo/internal/service/*`, `miruzo/internal/api/*`,
  `miruzo-py/app/services/*`, or shared test utilities
- Only add `from __future__ import annotations` (or other future imports) when a
  file truly needs it; if a newer language feature would significantly improve
  performance or readability, raise it so we can evaluate bumping the minimum
  Python version together

## Datetime Naming Conventions

- Keep datetime suffixes consistent across languages:
  - Golang: `...At` (example: `occurredAt`, `loveCanceledAt`)
  - Python: `..._at` (example: `occurred_at`, `love_canceled_at`)
  - SQL: `..._at` (example: `occurred_at`, `love_canceled_at`)
- When introducing new names, follow the mapping above instead of mixing styles
  within a single layer.
- Convert naming styles only at clear boundaries (for example SQL query params
  and generated structs), and keep names internally consistent in each layer.
- For period boundary timestamps, use noun+boundary naming (`periodStartAt`,
  `periodEndAt` / `period_start_at`, `period_end_at`) instead of event-style
  names such as `periodStartedAt` / `period_started_at`.

## API & Services Conventions

- HTTP handlers live under `miruzo/internal/api/*` and should delegate to
  services in `miruzo/internal/service/*`; do not perform DB logic in handlers
- Persistence interfaces live under `miruzo/internal/persist/*`; adapters live
  under `miruzo/internal/adapter/persistence/*`
- Service classes should depend on repository interfaces and pure helpers
- Python ingest persistence lives under `miruzo-py/app/persist/*` and should
  use SQLAlchemy Core table expressions
- Enum/string settings should be normalized via validators in
  `miruzo-py/app/config/environments.py`; avoid environment-specific branching
  outside that module
- Importers in `miruzo-py/scripts/importers/common/*` must log all destructive
  actions and honor the `force` flag before deleting files
- List APIs must preserve `limit + 1` pagination semantics to determine
  `hasNext` and cursor progression
- List API cursors must remain opaque base64url strings at the API boundary.
  Encode/decode through `miruzo/internal/api/image/list/cursor_codec.go`.
- When touching DB queries or schema, run `cd miruzo && make generate` and
  update migrations + generated files in the same commit.

## Commits

- Follow Conventional Commits; keep the subject within 55 characters
- Group related backend changes (API + service + repository + tests) in one
  commit; avoid mixing importer work with unrelated API tweaks
- If backend support/version/driver docs change, update `README.md`,
  `CONTRIBUTING.md`, and `AGENTS.md` in the same PR/commit set.
- Commit code-generated DB artifacts from `make generate` together with source
  SQL/migration edits (`miruzo/internal/database/*/gen`,
  `miruzo/internal/database/*/migrations_min`)
- Do not commit local build/cache outputs such as `miruzo/bin/`, `miruzo/dist/`,
  `**/__pycache__/`, `.pytest_cache/`, or similar machine-local artifacts
- Use the GitHub issue and PR templates provided in `.github/` when filing or
  submitting changes

## Testing

- Default suites:
  - `cd miruzo && go test ./...`
  - `cd miruzo-py && pytest`
- Run focused suites as needed:
  - `cd miruzo && go test ./internal/service/...` for service logic
  - `cd miruzo && go test ./internal/adapter/persistence/contract/...` for
    repository contract behavior (SQLite + PostgreSQL)
  - `cd miruzo-py && pytest tests/importers` for importer flows
  - `cd miruzo-py && pytest tests/persist` for SQLAlchemy Core repositories
- Docker-based tests must clean up containers; reuse the helper fixtures
  already committed (do not hand-roll `docker run` invocations)
- Use in-memory SQLite for unit tests unless the code path explicitly requires
  a different backend
- Pure helper functions must be covered by unit tests; impure logic should be
  split into testable pieces
- Avoid network access during tests; mock external calls or use local fixtures
