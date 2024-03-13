-- name: GetProduct :one
SELECT * FROM products
WHERE product_id = $1 
LIMIT 1;

-- name: GetProductForUpdate :one
SELECT * FROM products
WHERE product_id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListProduct :many
SELECT * FROM products
ORDER BY product_name
LIMIT $1
OFFSET $2;

-- name: CreateProduct :one
INSERT INTO products (
  product_name, unit_price, packsize
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
  set unit_price = $2,
  product_name = $3,
  packsize = $4
WHERE product_id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_id = $1;

-- name: SearchILikeProducts :many
SELECT * FROM products
WHERE product_name ILIKE $1;

-- name: CountProducts :one
SELECT COUNT(*) FROM products;