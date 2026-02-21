-- name: GetAppData :one
SELECT * FROM app_data
WHERE id = ? LIMIT 1;

-- name: CreateAppData :exec
INSERT INTO app_data (
  id, data
)
VALUES (
  ?, ?
);