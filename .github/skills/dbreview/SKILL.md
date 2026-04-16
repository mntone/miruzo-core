---
name: dbreview
description: Review DB changes for migration safety, schema/query consistency, and cross-dialect issues (SQLite/PostgreSQL)
---

# DB Review Skill (miruzo)

## Purpose
Review database-related changes and detect risks in:
- migrations
- SQL queries (sqlc)
- repository / adapter changes

Focus on correctness and cross-dialect safety (SQLite / PostgreSQL).

## Input
- migration SQL (up/down)
- sqlc queries
- repository / adapter diffs (if present)

Ignore non-database logic unless it affects DB consistency.

## Output Format

### Result
- OK | NEEDS_REVIEW | RISKY

### Findings
- issue summary
  - scope: migration | query | constraint | index | dialect
  - detail: what is wrong or suspicious
  - impact: what may break and where
  - suggestion: how to fix or what to verify

(repeat as needed)

### Checked
- migration symmetry
- schema/query consistency
- dialect differences (SQLite vs PostgreSQL)
- indexes and constraints

### Required Tests (Must be satisfied)
- Schema changes involving constraints (UNIQUE, CHECK, etc.) must include tests
- Tests must verify constraint behavior, not only schema existence
- Tests should be implemented as contract tests for persistence adapters
  (located under adapter/persistence/contract)
- New persistence repository methods must include or update contract tests
  covering success cases and relevant edge cases

## Review Checklist

### 1. Migration Symmetry
- Is every `up` operation reversible in `down`?
- Are indexes, constraints, triggers also removed in `down`?
- Is rename handled safely (not destructive drop/create)?

### 2. Schema / Query Consistency
- Do queries match updated column names and types?
- Do NULL / NOT NULL changes break scan targets?
- Are default value changes reflected in queries or logic?

### 3. Dialect Differences
Check for differences between SQLite and PostgreSQL:
- `RETURNING` usage
- `ON CONFLICT` behavior
- JSON functions
- datetime representation (TEXT vs TIMESTAMP / INTEGER)
- CHECK constraints behavior
- partial index compatibility

### 4. Index / Constraint
- Are required indexes present for new queries?
- Are UNIQUE / CHECK constraints properly defined?
- Any redundant or unused indexes added?

### 5. Repository Contract Consistency
- Do new or changed repository methods have matching adapter and test updates?
- Are contract tests updated to cover the new repository behavior?

## Naming Conventions
- CHECK (column): ck_`table`_`column`
- CHECK (table):  ck_`table`_`name`
- INDEX:          ix_`table`_`name`
- UNIQUE:         uq_`table`_`name`

Violations must be reported.

## Rules
- Prefer pointing out concrete risks over general advice
- If unsure, mark as NEEDS_REVIEW instead of guessing
- Do not review service/business logic unless it affects DB correctness
- Be concise and technical
- Assume commonly available modern features (e.g., RETURNING, STRICT in SQLite)
- Suggest newer or more efficient SQL features when applicable
