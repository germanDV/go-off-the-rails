-- name: GetMovie :one
SELECT * FROM movies WHERE id = ? AND org_id = ?;

-- name: ListMovies :many
SELECT * FROM movies
WHERE org_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateMovie :exec
INSERT INTO movies (id, title, rating, org_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateMovie :exec
UPDATE movies
SET title = ?, rating = ?, updated_at = ?, version = version + 1
WHERE id = ? AND org_id = ? AND version = ?;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = ? AND org_id = ?;

