-- name: CreatePurchaseOrder :one
INSERT INTO purchase_orders (
  id, supplier_name, supplier_address, lpo_data, lpo_pdf
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetPurchaseOrder :one
SELECT * FROM purchase_orders
WHERE id = $1 LIMIT 1;

-- name: ListPurchaseOrder :many
SELECT *  FROM purchase_orders
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: DeletePurchaseOrder :exec
DELETE FROM purchase_orders
WHERE id = $1;

-- name: CountPurchaseOrders :one
SELECT COUNT(*) FROM purchase_orders;