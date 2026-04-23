-- name: CreateOrg :exec
INSERT INTO orgs (id, name, created_at, updated_at)
VALUES (?, ?, ?, ?);
