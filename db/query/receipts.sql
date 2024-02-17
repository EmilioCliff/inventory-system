-- name: GetReceipt :one
SELECT * FROM receipts
WHERE receipt_number = $1 
LIMIT 1;

-- name: GetUserReceiptsByID :many
SELECT * FROM receipts
WHERE user_receipt_id = $1;

-- name: GetUserReceiptsByUsername :many
SELECT * FROM receipts
WHERE user_receipt_username = $1;

-- name: ListReceipts :many
SELECT * FROM receipts
ORDER BY receipt_id;

-- name: CreateReceipt :one
INSERT INTO receipts(
    receipt_number, user_receipt_id, receipt_data, user_receipt_username, receipt_pdf
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;
