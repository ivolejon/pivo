// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package documents

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const AddDocument = `-- name: AddDocument :one
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
  id, embeddings_ids, filename, title, project_id, created_at
`

type AddDocumentParams struct {
	ID            uuid.UUID   `json:"id"`
	EmbeddingsIds []uuid.UUID `json:"embeddingsIds"`
	Filename      string      `json:"filename"`
	Title         *string     `json:"title"`
	ProjectID     uuid.UUID   `json:"projectId"`
	CreatedAt     time.Time   `json:"createdAt"`
}

func (q *Queries) AddDocument(ctx context.Context, db DBTX, arg AddDocumentParams) (Document, error) {
	row := db.QueryRow(ctx, AddDocument,
		arg.ID,
		arg.EmbeddingsIds,
		arg.Filename,
		arg.Title,
		arg.ProjectID,
		arg.CreatedAt,
	)
	var i Document
	err := row.Scan(
		&i.ID,
		&i.EmbeddingsIds,
		&i.Filename,
		&i.Title,
		&i.ProjectID,
		&i.CreatedAt,
	)
	return i, err
}

const GetDocumentsByProjectId = `-- name: GetDocumentsByProjectId :many
SELECT
  id, embeddings_ids, filename, title, project_id, created_at
FROM
  documents
WHERE
  project_id = $1
`

func (q *Queries) GetDocumentsByProjectId(ctx context.Context, db DBTX, projectID uuid.UUID) ([]Document, error) {
	rows, err := db.Query(ctx, GetDocumentsByProjectId, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Document
	for rows.Next() {
		var i Document
		if err := rows.Scan(
			&i.ID,
			&i.EmbeddingsIds,
			&i.Filename,
			&i.Title,
			&i.ProjectID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
