-- name: CreateUser :one
INSERT INTO users (
  username, password, email, phone_number, address, stock, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsename :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUser :many
SELECT * FROM users
ORDER BY username;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: UpdateUserStock :one
UPDATE users
  set stock = $2
WHERE user_id = $1
RETURNING *;

-- name: UpdateUserCredentials :one
UPDATE users
  set password = $3,
  email = $2,
  address = $4,
  phone_number = $5,
  username = $6
WHERE user_id = $1
RETURNING *;

-- name: UpdateUserPasswordFisrtLogin :one
UPDATE users
  set password = $2
WHERE user_id = $1
RETURNING *;

-- name: SearchILikeUsers :many
SELECT username
FROM users
WHERE username ILIKE $1;