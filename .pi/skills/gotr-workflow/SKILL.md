---
name: gotr-workflow
description: GOTR project conventions and the three-phase workflow (discuss -> spec -> implement). Use when working on the Go Off The Rails (GOTR) framework project.
---

# GOTR Workflow

This skill defines how to work on the GOTR project. Always follow this three-phase workflow.

## Architecture

GOTR is a layered web framework:

```
views (templ + HTMX)
  ↑
controllers (orchestration, HTTP handlers)
  ↑
db (repositories, sqlc generated code, migrations)
  ↑
domain (business logic, validation, pure functions)
```

- **domain**: Pure Go structs and logic. No DB or HTTP dependencies. Returns errors for invalid state.
- **db**: Owns all data access. Repositories wrap sqlc-generated queries and map sqlc structs to domain structs. Handles migrations.
- **controllers**: Orchestration layer. Reads from repositories, calls domain logic, persists changes, renders views.
- **views**: Templ templates that render domain data. Uses HTMX for interactivity.

## Database

- SQLite only. IDs are UUIDv7 (TEXT PRIMARY KEY).
- SQL is first-class. No ORMs. Write migrations in `db/migrations/`.
- Write queries in `db/queries/*.sql`, then run `sqlc generate`.
- Every table must have: `id`, `version` (for optimistic concurrency), `created_at`, `updated_at`.
- Tenant-scoped tables must have `org_id TEXT NOT NULL REFERENCES orgs(id)`.

## Migrations & Queries

- Migration files: `db/migrations/{VERSION}_{description}.sql`
- Query files: `db/queries/{resource}_queries.sql`
- After adding queries, run `sqlc generate`.
- Optimistic concurrency: `version INTEGER NOT NULL DEFAULT 1`, increment on UPDATE.

## Auth & Multi-Tenancy

- JWT in HttpOnly Secure SameSite=Strict cookie named `auth_token`.
- Roles: `ADMIN`, `USER`, `SUPERADMIN`.
- Every user belongs to an org. Signup creates a new org and assigns `ADMIN`.
- Invite flow allows joining existing orgs as `USER`.
- Controllers read `org_id` and `role` from JWT claims via context.

## Three-Phase Workflow

### Phase 1: Discuss (/think)
- Use read-only exploration. Only read, bash, grep, find, ls, questionnaire.
- Explore existing code to understand patterns.
- Ask clarifying questions.
- Evaluate approaches.
- Do NOT modify files.

### Phase 2: Spec (/spec)
- Write a detailed implementation plan to `specs/{feature}.md`.
- Follow `specs/_template.md`.
- Include: DB schema, sqlc queries, domain model, repository, controller endpoints, views, routes, middleware.
- Number all implementation steps.
- This is the contract. Do NOT implement yet.

### Phase 3: Build (/build)
- Read the spec.
- Implement step by step.
- Mark completed steps with `[DONE:n]` in responses.
- If the spec needs to change during implementation, update the spec file first.
- After DB query changes, run `sqlc generate`.
- After Go changes, run `go build ./...`.

## Verification (/verify)
- After implementation, cross-check the code against the spec.
- Report any gaps.

## Common Patterns

### Domain Validation
```go
func NewThing(orgID uuid.UUID, title string) (Thing, error) {
    title = strings.TrimSpace(title)
    if utf8.RuneCountInString(title) == 0 {
        return Thing{}, errors.New("title is required")
    }
    return Thing{ID: uuid.Must(uuid.NewV7()), OrgID: orgID, Title: title}, nil
}
```

### Repository Method
```go
func (r *ThingsRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Thing, error) {
    row, err := r.queries.GetThing(ctx, id.String())
    if err != nil { return domain.Thing{}, err }
    return domain.Thing{ID: uuid.MustParse(row.ID), Title: row.Title}, nil
}
```

### Controller Endpoint
```go
func (c *ThingsController) Index(w http.ResponseWriter, r *http.Request) {
    orgID := domain.OrgIDFromContext(r.Context())
    things, err := c.repo.List(r.Context(), orgID)
    if err != nil { /* handle */ }
    views.ThingsIndex(things).Render(r.Context(), w)
}
```

### HTMX Form
```html
// Use data attributes for verbs other than POST
// e.g., PUT for updates, DELETE for deletion
```

## Form Values (Not JSON)

All HTML forms submit as form values, never JSON. Parse with `r.ParseForm()` and `r.FormValue("key")`.
