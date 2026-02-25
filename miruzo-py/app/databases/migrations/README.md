# Migrations and Index Policy

- SQLModel defines logical schema only
- All performance-oriented indexes live in migrations
- Partial indexes are tied to list API semantics
- Indexes are grouped by list API (recently, hall_of_fame, ...)
