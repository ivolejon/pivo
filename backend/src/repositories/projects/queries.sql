-- name: GetProjectsByClientId :many
SELECT DISTINCT
  ON (id) *
FROM
  projects
WHERE
  client_id = $1
ORDER BY
  id,
  created_at DESC;

-- name: AddProject :one
INSERT INTO
  projects (id, client_id, title, created_at)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;

-- name: GetProjectById :one
SELECT DISTINCT
  ON (id) *
FROM
  projects
WHERE
  id = $1
ORDER BY
  id,
  created_at DESC;
