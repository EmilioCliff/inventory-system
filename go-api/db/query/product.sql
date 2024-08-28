-- name: GetProduct :one
SELECT * FROM products
WHERE product_id = $1 
LIMIT 1;

-- name: GetProductByProductName :one
SELECT * FROM products
WHERE product_name = $1 
LIMIT 1;

-- name: GetProductPrice :one
SELECT unit_price FROM products
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

-- name: ListAllProduct :many
SELECT * FROM products
ORDER BY product_name;

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
SELECT product_name FROM products
WHERE LOWER(product_name) LIKE LOWER('%' || $1 || '%');

-- name: CountProducts :one
SELECT COUNT(*) FROM products;