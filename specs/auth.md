# Authentication Implementation Plan

## 1. DB Schema Changes

### Add `verified` column to `users` table

Update `db/migrations/001_orgs_and_users.sql` to include a `verified` boolean column (default `false`):

```sql
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  org_id TEXT NOT NULL REFERENCES orgs(id),
  email TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL,
  verified BOOLEAN NOT NULL DEFAULT false,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Also update the `tempMigrationRun` in `main.go` to match.

### New `verification_tokens` table

Add `db/migrations/003_verification_tokens.sql`:

```sql
CREATE TABLE IF NOT EXISTS verification_tokens (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  token TEXT NOT NULL UNIQUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Also add this to `tempMigrationRun` in `main.go`.

## 2. sqlc Queries

### User queries (`db/queries/users.sql`)

```sql
-- name: CreateUser :exec
INSERT INTO users (id, org_id, email, password_hash, role, verified)
VALUES (?, ?, ?, ?, ?, false);

-- name: GetUserByEmail :one
SELECT id, org_id, email, password_hash, role, verified, created_at, updated_at
FROM users
WHERE email = ?;

-- name: VerifyUser :exec
UPDATE users SET verified = true, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
```

### Verification token queries (`db/queries/verification_tokens.sql`)

```sql
-- name: CreateVerificationToken :exec
INSERT INTO verification_tokens (id, user_id, token)
VALUES (?, ?, ?);

-- name: GetVerificationToken :one
SELECT id, user_id, token, created_at
FROM verification_tokens
WHERE token = ?;

-- name: DeleteVerificationToken :exec
DELETE FROM verification_tokens
WHERE token = ?;
```

Run `sqlc generate` after adding these queries.

## 3. Users Repository (`db/users_repository.go`)

Create `db/users_repository.go` with a `UsersRepository` struct wrapping the generated queries:

- `Create(ctx, user) error` — inserts a new user.
- `GetByEmail(ctx, email) (User, error)` — fetches user by email.
- `Verify(ctx, userID) error` — sets `verified = true`.
- `CreateVerificationToken(ctx, userID, token) error` — stores a verification token.
- `GetVerificationToken(ctx, token) (VerificationToken, error)` — retrieves token record.
- `DeleteVerificationToken(ctx, token) error` — removes token after use.

## 4. Domain Helpers

### Password hashing (`domain/password.go`)

```go
func HashPassword(password string) (string, error)   // bcrypt.GenerateFromPassword
func CheckPassword(hash, password string) error       // bcrypt.CompareHashAndPassword
```

### Verification token generation (`domain/verification_token.go`)

```go
func GenerateVerificationToken() (string, error)  // crypto/rand hex string
```

## 5. Auth Controller (`controllers/auth_controller.go`)

Create `AuthController` with dependencies: `UsersRepository`, `Tokenizer`, `MiddlewareChain`.

### Endpoints

#### `POST /signup`
- Parse `email` and `password` from form body.
- Hash password with `domain.HashPassword`.
- Create user with `verified: false`.
- Generate a verification token and store it.
- Return plain text: `"Signup successful. Please verify your email. Token: <token>"` (token shown in response since no email sending yet).

#### `POST /login`
- Parse `email` and `password` from form body.
- Look up user by email. If not found, return `404 "not found"`.
- Check password with `domain.CheckPassword`. If mismatch, return `404 "not found"` (don't leak info).
- If user has `verified: false`, return `403 "need to verify account"`.
- Generate JWT via `Tokenizer.Generate`.
- Set JWT in an `HttpOnly`, `Secure`, `SameSite=Strict` cookie named `auth_token`.
- Return plain text: `"Login successful"`.

#### `POST /signout`
- Clear the `auth_token` cookie by setting `MaxAge: -1`.
- Return plain text: `"Signed out"`.

#### `GET /verify?token=<token>`
- Parse `token` from query string.
- Look up verification token. If not found, return `400 "invalid token"`.
- Mark user as `verified: true`.
- Delete the verification token.
- Redirect to "/login".

## 6. Route Registration

In `controllers/server.go`'s `routes()` method, add:

```go
usersRepo := db.NewUsersRepository(generated.New(s.dbClient))
authController := NewAuthController(s.mdw, usersRepo, s.mdw.tokenizer)
authController.RegisterRoutes(s.mux)
```

The auth routes (`/signup`, `/login`, `/signout`, `/verify`) are public, use the RBACNone middleware.

## 7. Cookie Handling

Cookie settings for JWT on login:

```go
http.Cookie{
    Name:     AuthCookieName,  // "auth_token"
    Value:    jwtString,
    Path:     "/",
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode,
    MaxAge:   int(domain.TokenExpiration.Seconds()),
}
```

On signout, set the same cookie with `MaxAge: -1` and empty `Value` to destroy it.
