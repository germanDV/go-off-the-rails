-- name: CreateVerificationToken :exec
INSERT INTO verification_tokens (token, user_id, created_at)
VALUES (?, ?, ?);

-- name: GetVerificationToken :one
SELECT token, user_id, created_at
FROM verification_tokens
WHERE token = ?;

-- name: DeleteVerificationToken :exec
DELETE FROM verification_tokens
WHERE token = ?;
