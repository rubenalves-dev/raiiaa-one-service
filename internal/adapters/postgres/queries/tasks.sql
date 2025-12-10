-- name: CreateTask :one
INSERT INTO tasks (user_id, project_id, title, description, status, priority, due_date, estimated_minutes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListTasksByProjectID :many
SELECT * FROM tasks
WHERE user_id = $1 AND project_id = $2
ORDER BY created_at DESC;

-- name: ListTasksByUserID :many
SELECT * FROM tasks
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = $2, updated_at = NOW()
WHERE id = $1;
