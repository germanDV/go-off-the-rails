# Invite Flow Specification

## 1. DB Schema

### New `invites` table

Add to `db/migrations/001_orgs_and_users.sql`:

```sql
CREATE TABLE IF NOT EXISTS invites (
  id TEXT PRIMARY KEY,
  org_id TEXT NOT NULL REFERENCES orgs(id),
  email TEXT NOT NULL,
  token TEXT NOT NULL UNIQUE,
  created_at DATETIME NOT NULL,
  expires_at DATETIME NOT NULL
);
```

Also add this migration to `tempMigrationRun` in `main.go`.

## 2. sqlc Queries

### Invite queries (`db/queries/invites.sql`)

```sql
-- name: CreateInvite :exec
INSERT INTO invites (id, org_id, email, token, created_at, expires_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetInviteByToken :one
SELECT id, org_id, email, token, created_at, expires_at
FROM invites
WHERE token = ?;

-- name: DeleteInvite :exec
DELETE FROM invites
WHERE token = ?;
```

Run `sqlc generate` after adding these queries.

## 3. Domain Model (`domain/invite.go`)

```go
type Invite struct {
    ID        uuid.UUID
    OrgID     uuid.UUID
    Email     string
    Token     string
    CreatedAt time.Time
    ExpiresAt time.Time
}
```

Include a `GenerateInviteToken()` function (similar to verification tokens) that produces a crypto/rand hex string.

## 4. Repository (`db/invites_repository.go`)

Create `InvitesRepository` wrapping the generated queries:

- `Create(ctx, invite) error` — inserts a new invite.
- `GetByToken(ctx, token) (Invite, error)` — retrieves an invite by token.
- `Delete(ctx, token) error` — removes an invite after use.

## 5. Endpoints

### `POST /invites` (admin-only)

- Protected by `RBACAdmin` middleware — only authenticated admins can create invites.
- Parse `email` from form body.
- Read the admin's `org_id` from the JWT claims (via `domain.Actor` in context).
- Generate an invite token and store it with the org ID and email.
- Return plain text: `"Invite created. Token: <token>"` (token shown in response since no email sending yet).
- Add to the existing AuthController

## 6. Modified Signup Behavior

### `POST /signup`

The signup handler accepts an optional `invite_token` form parameter.

#### With `invite_token` (joining an existing org)

1. Validate the invite token — look it up in the `invites` table. If not found, return `400 "invalid invite token"`.
2. Parse `email` and `password` from form body.
3. Verify the email matches the invite's email. If mismatch, return `400 "email does not match invite"`. Ensure the token is not expired.
4. Hash the password.
5. Create the user with the invite's `org_id`, `role = USER`, and `verified = true` (no verification token is generated — the admin already validated the email via the invite).
6. Delete the invite token (cleanup).
7. All DB operations (create user, delete invite) run in a single transaction.
8. Return plain text: `"Signup successful. You can now log in."`

#### Without `invite_token` (creating a new org)

Keep the current behavior unchanged:

1. Create a new org named `"<email>'s org"`.
2. Create the user with `role = ADMIN`.
3. Generate a verification token.
4. All in a single transaction.

## 7. Route Registration

In `controllers/server.go`'s `routes()` method, add:

```go
invitesRepo := db.NewInvitesRepository(generated.New(s.dbClient))
```

Pass `invitesRepo` to `AuthController` (for signup invite handling) and register the `POST /invites` route with `RBACAdmin` middleware.

## 8. Cleanup

Used invite tokens are deleted immediately upon successful signup. No background cleanup job is needed for now.
