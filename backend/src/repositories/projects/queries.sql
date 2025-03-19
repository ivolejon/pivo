-- name: GetProjectsByClientId :many
SELECT
  *
FROM
  projects
WHERE
  client_id = $1;

-- name: AddProject :one
INSERT INTO
  projects (id, client_id, title, created_at)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;

-- name: GetProjectById :one
SELECT
  *
FROM
  projects
WHERE
  id = $1;
