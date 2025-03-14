-- name: GetProjectsByClientId :many
SELECT * FROM projects WHERE client_id = $1;
