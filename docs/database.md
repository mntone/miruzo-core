# Database Notes

## ingests UNIQUE policy

`ingests` has two UNIQUE constraints with different roles.

- `UNIQUE(fingerprint)` is the primary deduplication rule.
  It rejects duplicate image content.
- `UNIQUE(relative_path)` is a secondary safety guard.
  It prevents accidental duplicate inserts for the same path.

For MySQL, `relative_path` comparison assumes a binary collation
(e.g. `utf8mb4_0900_bin`) so case is treated distinctly.
