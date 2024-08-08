-- name: GetUserStockValue :one
SELECT * FROM stock_value
WHERE user_id = $1 LIMIT 1;

-- name: GetAllStockValue :many
SELECT * FROM stock_value;

-- name: UpdateUserStockValue :one
UPDATE stock_value
  set value = value + $1
WHERE user_id = $2
RETURNING *;

-- name: CreateUserStockValue :one
INSERT INTO stock_value (
  user_id
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteUserStockValue :exec
DELETE FROM stock_value
WHERE user_id = $1;

-- name: TotalStockValue :one
SELECT SUM("value") AS total_value
FROM stock_value;