-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, default_hourly_rate, currency)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
