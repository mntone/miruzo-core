[&lsaquo; Back to README](./README.md)

# Contributing to miruzo-core

## 🧭 Overview

miruzo-core is composed of the Go API backend (`miruzo/`) and the Python ingest
tooling (`miruzo-py/`). This document describes contribution workflow,
environment setup, and review expectations.

For canonical coding rules and non-negotiable conventions, always follow
[AGENTS.md](./AGENTS.md).

<a id="prerequisites"></a>

## 🛠️ Prerequisites

- Go 1.26.x (see [`miruzo/go.mod`](./miruzo/go.mod))
- Python 3.10+ (3.14.x recommended for local development; see
  [`miruzo-py/.python-version`](./miruzo-py/.python-version))
- `uv` (Python package/dependency manager)
- Database minimum versions:
  - MySQL 8.0.16+ (`CHECK` support required)
  - PostgreSQL 14+
  - SQLite 3.37.0+ (`RETURNING` and `STRICT` support required)
- Python ingest drivers:
  - MySQL: `mysqlclient` (`MySQLdb`, DSN `mysql+mysqldb://...`)
  - PostgreSQL: `psycopg3` (`psycopg`, DSN `postgresql+psycopg://...`)
  - SQLite: stdlib `sqlite3` (DSN `sqlite:///...`)

Install dependencies:

```bash
cd miruzo && go mod download
cd ../miruzo-py && uv sync --extra dev
```

OS-specific setup (required for DB drivers / Go tools):

- Linux (Debian/Ubuntu):
  - Go tools:
    ```bash
    cd miruzo && make tools
    ```
  - PostgreSQL Python driver dependencies:
    ```bash
    sudo apt install -y libpq-dev
    uv sync --extra dev --extra postgres
    ```
  - MySQL Python driver dependencies:
    ```bash
    sudo apt install -y build-essential default-libmysqlclient-dev pkg-config
    uv sync --extra dev --extra mysql
    ```
- macOS:
  - Go tools (recommended):
    ```bash
    brew install gopls delve go-air sqlc goreleaser
    ```
  - Go tools (alternative):
    ```bash
    cd miruzo && make tools
    ```
  - PostgreSQL Python driver dependencies:
    ```bash
    brew install libpq
    uv sync --extra dev --extra postgres
    ```
  - MySQL Python driver dependencies:
    ```bash
    brew install mysql pkg-config
    uv sync --extra dev --extra mysql
    ```

Configure runtime files:

- API: use `miruzo/config.yaml` (base: `miruzo/internal/app/config.sample.yaml`)
- Ingest: copy `miruzo-py/.env.development` to `miruzo-py/.env` and set:
  - `ENVIRONMENT` (`development` or `production`)
  - `DATABASE_BACKEND` (`mysql`, `postgres`, or `sqlite`)
  - `DATABASE_URL` (`mysql+mysqldb://...`, `postgresql+psycopg://...`,
    `sqlite:///...`)
  - Path-related variables (`MEDIA_ROOT`, `PUBLIC_MEDIA_ROOT`, `GATAKU_ROOT`,
    `GATAKU_ASSETS_ROOT`, `GATAKU_SYMLINK_DIRNAME`) are optional for initial
    setup and can stay on defaults.
- Optional test database DSN environment variables:
  - Go tests:
    - `MIRUZO_TEST_MYSQL_URL`
    - `MIRUZO_TEST_POSTGRES_URL`
  - Python tests:
    - `MIRUZO_PY_TEST_MYSQL_URL`
    - `MIRUZO_PY_TEST_POSTGRES_URL`

Current Go API database backends are `mysql`, `postgresql`, and `sqlite`.

## 🔁 Workflow

- Follow Conventional Commits (`feat:`, `fix:`, `refactor:`, etc.)
- Keep each commit focused on one logical change
- Start larger or risky changes from a GitHub issue/discussion
- Use issue/PR templates under `.github/`
- Before requesting review, run relevant tests and linters
- When backend support/version/driver documentation changes, update
  `README.md`, `CONTRIBUTING.md`, and `AGENTS.md` together.

If you modify DB schema/query sources (`miruzo/internal/database/*/queries`,
migrations), run:

```bash
cd miruzo && make generate
```

Commit generated artifacts together with source SQL/migration changes:

- `miruzo/internal/database/*/gen`
- `miruzo/internal/database/*/migrations_min`

Do not commit local build/cache outputs such as:
`miruzo/bin/`, `miruzo/dist/`, `**/__pycache__/`, `.pytest_cache/`.

## 🧰 Operations

Run API and CLI commands from `miruzo/`.

Start API locally:

```bash
cd miruzo && air
# or
cd miruzo && go run ./cmd/miruzo-api
```

Run importer help:

```bash
cd miruzo-py && uv run python -m scripts.gataku_import --help
```

Common CLI operations:

```bash
cd miruzo && go run ./cmd/miruzo-cli migrate up
cd miruzo && go run ./cmd/miruzo-cli migrate down [N]
cd miruzo && go run ./cmd/miruzo-cli migrate goto V
cd miruzo && go run ./cmd/miruzo-cli migrate version
cd miruzo && go run ./cmd/miruzo-cli migrate force V
cd miruzo && go run ./cmd/miruzo-cli job daily-decay
```

`migrate force V` is for recovery workflows only. Before destructive operations
(`migrate down`, `migrate force`, importer `--force`), confirm backup/snapshot
availability and impact scope.

## 🎨 Code Style

`AGENTS.md` is the source of truth. Important highlights:

- Use tabs for indentation
- Keep interfaces and module boundaries strongly typed
- Keep SQL explicit:
  - Go: SQL files + sqlc under `miruzo/internal/database/*`
  - Python: SQLAlchemy Core expressions (no ORM model magic)
- Keep HTTP handlers thin and delegate to services/repositories
- Keep markdown wrapped at 80 columns

## 🧪 Testing

Default suites:

```bash
cd miruzo && go test ./...
cd miruzo-py && uv run pytest
```

Focused suites:

- `cd miruzo && go test ./internal/service/...`
- `cd miruzo && go test ./internal/adapter/persistence/contract/...`
- `cd miruzo-py && uv run pytest tests/importers`
- `cd miruzo-py && uv run pytest tests/persist`

Notes:

- Reuse committed fixtures/helpers; do not hand-roll container lifecycle scripts
- Keep unit tests backend-light unless backend-specific behavior is under test
- Avoid external network dependencies in tests

## 🐛 Reporting Issues

- File bugs/feature requests through GitHub Issues using provided templates
- Include reproduction steps, expected/actual behavior, and relevant logs
- For sensitive issues, contact *mntone* via links in [README.md](./README.md)

## 📜 License Notice

By contributing, you agree to license your contributions under GPLv3 (same as
the project). Submit only code/assets you are allowed to relicense under GPLv3,
and verify compatibility of third-party dependencies before introducing them.
