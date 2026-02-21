-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (address)
VALUES ($1)
RETURNING *;

-- name: GetUserByAddress :one
SELECT * FROM users
WHERE address = $1 LIMIT 1;
