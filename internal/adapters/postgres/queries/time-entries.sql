-- name: StartTimer :one
INSERT INTO time_entries (user_id, task_id, start_time, description, is_billable)
VALUES ($1, $2, NOW(), $3, $4)
RETURNING *;

-- name: StopTimer :one
UPDATE time_entries
SET end_time = NOW(), updated_at = NOW()
WHERE id = $1 AND end_time IS NULL
RETURNING *;

-- name: StopTimerForUser :one
UPDATE time_entries
SET end_time = NOW(), updated_at = NOW()
WHERE user_id = $1 AND end_time IS NULL
RETURNING *;

-- name: ListTimeEntriesByRange :many
SELECT sqlc.embed(te), sqlc.embed(t), sqlc.embed(p), sqlc.embed(c) FROM time_entries te
JOIN tasks t ON te.task_id = t.id
JOIN projects p ON te.project_id = p.id
LEFT JOIN clients c ON p.client_id = c.id
WHERE te.user_id = $1 AND te.start_time >= $2 AND (te.end_time <= $3 OR te.end_time IS NULL)
ORDER BY te.start_time DESC;

-- name: GetUnbillableTimeEntries :many
SELECT sqlc.embed(te), sqlc.embed(t), sqlc.embed(p) FROM time_entries te
JOIN tasks t ON te.task_id = t.id
JOIN projects p ON te.project_id = p.id
WHERE te.user_id = $1 AND te.is_billable = false AND te.end_time IS NOT NULL
ORDER BY te.start_time DESC;
