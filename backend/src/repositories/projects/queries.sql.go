// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.sql

package projects

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const AddProject = `-- name: AddProject :one
INSERT INTO
  projects (id, client_id, title, created_at)
VALUES
  ($1, $2, $3, $4)
RETURNING
  id, client_id, title, created_at
`

type AddProjectParams struct {
	ID        uuid.UUID
	ClientID  uuid.UUID
	Title     *string
	CreatedAt time.Time
}

func (q *Queries) AddProject(ctx context.Context, db DBTX, arg AddProjectParams) (Project, error) {
	row := db.QueryRow(ctx, AddProject,
		arg.ID,
		arg.ClientID,
		arg.Title,
		arg.CreatedAt,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.Title,
		&i.CreatedAt,
	)
	return i, err
}

const GetProjectById = `-- name: GetProjectById :one
SELECT DISTINCT
  ON (id) id, client_id, title, created_at
FROM
  projects
WHERE
  id = $1
ORDER BY
  id,
  created_at DESC
`

func (q *Queries) GetProjectById(ctx context.Context, db DBTX, id uuid.UUID) (Project, error) {
	row := db.QueryRow(ctx, GetProjectById, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.Title,
		&i.CreatedAt,
	)
	return i, err
}

const GetProjectsByClientId = `-- name: GetProjectsByClientId :many
SELECT DISTINCT
  ON (id) id, client_id, title, created_at
FROM
  projects
WHERE
  client_id = $1
ORDER BY
  id,
  created_at DESC
`

func (q *Queries) GetProjectsByClientId(ctx context.Context, db DBTX, clientID uuid.UUID) ([]Project, error) {
	rows, err := db.Query(ctx, GetProjectsByClientId, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.ClientID,
			&i.Title,
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
