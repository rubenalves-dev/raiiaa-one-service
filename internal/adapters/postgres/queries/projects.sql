-- name: CreateProject :one
INSERT INTO projects (user_id, client_id, name, description, status, color, deadline, hourly_rate, budget, repo_url, live_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: ListProjectsByClientID :many
SELECT * FROM projects
WHERE user_id = $1 AND client_id = $2
ORDER BY deadline ASC;

-- name: ListProjectsByUserID :many
SELECT * FROM projects
WHERE user_id = $1
ORDER BY deadline DESC;
