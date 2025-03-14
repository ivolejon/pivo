// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package jobs

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const GetProjectsByClientId = `-- name: GetProjectsByClientId :many
SELECT id, client_id, title, created_at FROM projects WHERE client_id = $1
`

func (q *Queries) GetProjectsByClientId(ctx context.Context, db DBTX, clientID pgtype.UUID) ([]Project, error) {
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
