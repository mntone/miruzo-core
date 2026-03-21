# Go Style

## Constructor Return Type (`New`)

- Default: constructors should return a value (`Service`), not a pointer.
- Return a pointer (`*Service`) only when the service should not be copied.
- Typical pointer cases:
  - the struct contains non-copy-safe fields (for example `sync.Mutex`)
  - the call path is likely to create extra value copies
- Size alone is not a hard rule.
- Even for a large service, returning a value is acceptable if it is stored once
  and used through a pointer receiver handler path, so extra copies are
  effectively avoided.
