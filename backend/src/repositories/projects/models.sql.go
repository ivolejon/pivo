// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package projects

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `json:"id"`
	ClientID  uuid.UUID `json:"clientId"`
	Title     *string   `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

type SchemaMigration struct {
	Version string `json:"version"`
}
