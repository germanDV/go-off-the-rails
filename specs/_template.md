# Feature Spec Template

Use this template for every new feature. Copy it to `specs/{feature-name}.md` and fill it in during the discussion phase.

## Overview

One-paragraph description of what this feature does and why it exists.

## Goals

- Specific, measurable goals
- Acceptance criteria

## Out of Scope

- Things we are explicitly NOT doing in this iteration

## DB Schema Changes

List all migration files and SQL changes:

```sql
-- db/migrations/XXX_description.sql
```

## sqlc Queries

List all query files and example queries:

```sql
-- db/queries/thing.sql
-- name: GetThing :one
SELECT * FROM things WHERE id = ? AND org_id = ?;
```

## Domain Model

Go structs and business logic:

```go
// domain/thing.go
type Thing struct { ... }

func NewThing(...) (Thing, error) { ... }
```

## Repository Layer

```go
// db/things_repository.go
type ThingsRepository struct { ... }
```

List methods and their signatures.

## Controller Endpoints

| Method | Route | Handler | Description |
|--------|-------|---------|-------------|
| GET    | /things | ThingsIndex | List all things |
| POST   | /things | ThingsCreate | Create a thing |

## Views / Templates

- `views/things_index.templ`
- `views/things_new.templ`

## Middleware

- Which middleware applies (auth, rbac, etc.)

## Route Registration

Where and how routes are wired in `controllers/server.go`.

## Testing Notes

- What should be manually tested
- Edge cases to verify

## Dependencies

- Any new Go packages or tools needed
