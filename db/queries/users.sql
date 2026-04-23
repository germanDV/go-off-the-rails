-- name: CreateUser :exec
INSERT INTO users (id, org_id, email, password_hash, role, verified, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT id, org_id, email, password_hash, role, verified, created_at, updated_at
FROM users
WHERE email = ?;

-- name: VerifyUser :exec
UPDATE users SET verified = true, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
