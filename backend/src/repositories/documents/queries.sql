-- name: GetDocumentsByProjectId :many
SELECT
  *
FROM
  documents
WHERE
  project_id = $1;

-- name: AddDocument :one
INSERT INTO
  documents (
    id,
    embeddings_ids,
    filename,
    title,
    project_id,
    created_at
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;
