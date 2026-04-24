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
