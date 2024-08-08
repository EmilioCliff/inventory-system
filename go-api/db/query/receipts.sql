-- name: GetReceipt :one
SELECT * FROM receipts
WHERE receipt_number = $1 
LIMIT 1;

-- name: GetReceiptByID :one
SELECT * FROM receipts
WHERE receipt_id = $1 
LIMIT 1;

-- name: GetUserReceiptsByID :many
SELECT * FROM receipts
WHERE user_receipt_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetAllUserReceiptsByID :many
SELECT * FROM receipts
WHERE user_receipt_id = $1
ORDER BY created_at DESC;

-- name: CountUserReceiptsByID :one
SELECT COUNT(*) FROM receipts
WHERE user_receipt_id = $1;

-- name: GetUserReceiptsByUsername :many
SELECT * FROM receipts
WHERE user_receipt_username = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountUserReceiptsByUsername :one
SELECT COUNT(*) FROM receipts
WHERE user_receipt_username= $1;

-- name: ListReceipts :many
SELECT * FROM receipts
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CreateReceipt :one
INSERT INTO receipts(
    receipt_number, user_receipt_id,  user_receipt_username, receipt_pdf
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: SearchILikeReceipts :many
SELECT receipt_number FROM receipts
WHERE LOWER(receipt_number) LIKE LOWER('%' || $1 || '%');

-- name: SearchUserReceipts :many
SELECT user_receipt_username FROM receipts
WHERE LOWER(user_receipt_username) LIKE LOWER('%' || $1 || '%');

-- name: CountReceipts :one
SELECT COUNT(*) FROM receipts;
