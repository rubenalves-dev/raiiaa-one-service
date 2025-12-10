-- name: CreateClient :one
INSERT INTO clients (user_id, company_name, contact_name, contact_email, contact_phone, address, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListClientByUserID :many
SELECT * FROM clients
WHERE user_id = $1 and is_archived = false
ORDER BY created_at DESC;

-- name: GetClient :one
SELECT * FROM clients
WHERE user_id = $1 and id = $2;
