---
name: General change
about: Default template for features, fixes, and refactors
---

# Summary

<!--
Add one or two bullet points. Describe what changed and why.
Mention related issues (e.g., `Fixes #123`) when applicable.
-->

- 

# Checklist

- [ ] Run `pytest` (plus targeted suites such as repository/service tests) and
  `ruff check` if the change affects logic or schemas.
- [ ] Commits follow Conventional Commits (e.g., `feat:`, `fix:`) and do not mix
  unrelated changes.
- [ ] New API models/fields include OpenAPI `title` and `description`
  metadata.
- [ ] Update configuration or documentation (`.env*`, importer README, etc.) when
  the change requires it.
- [ ] No sensitive or secret information is in the commits.
