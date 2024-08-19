-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1 
LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountTransactions :one
SELECT COUNT(*) FROM transactions;

-- name: CreateTransaction :one
INSERT INTO transactions (
    transaction_id, amount, phone_number, transaction_user_id, result_description, data_sold
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ChangeStatus :exec
UPDATE transactions
    set status = $2
WHERE transaction_id = $1;

-- name: ChangePaymentMethod :exec
UPDATE transactions 
    SET payment_method = $1
WHERE transaction_id = $2;

-- name: UpdateResultDescription :one
UPDATE transactions
    set result_description = $2
WHERE transaction_id = $1
RETURNING *;

-- name: UpdateTransaction :one
UPDATE transactions
  set mpesa_receipt_number = $3,
  phone_number = $2,
  result_description = $4
WHERE transaction_id = $1
RETURNING *;

-- name: AllUserTransactions :many
SELECT * FROM transactions
WHERE transaction_user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: AllUserTransactionsNoLimit :many
SELECT * FROM transactions
WHERE transaction_user_id = $1
AND status = TRUE
ORDER BY created_at DESC;

-- name: CountAllUserTransactions :one
SELECT COUNT(*) FROM transactions
WHERE transaction_user_id = $1;

-- name: SuccessUserTransactions :many
SELECT * FROM transactions
WHERE transaction_user_id = $1
AND status = true
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountSuccessfulUserTransactions :one
SELECT COUNT(*) FROM transactions
WHERE transaction_user_id = $1
AND status = true;

-- name: FailedUserTransactions :many
SELECT * FROM transactions
WHERE transaction_user_id = $1
AND status = false
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountFailedUserTransactions :one
SELECT COUNT(*) FROM transactions
WHERE transaction_user_id = $1
AND status = false;

-- name: SuccessTransactions :many
SELECT * FROM transactions
WHERE status = true
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountSuccessfulTransactions :one
SELECT COUNT(*) FROM transactions
WHERE status = true;

-- name: FailedTransactions :many
SELECT * FROM transactions
WHERE status = false
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountFailedTransactions :one
SELECT COUNT(*) FROM transactions
WHERE status = false;

-- name: GetUserTransaction :one
SELECT * FROM transactions
WHERE transaction_user_id = $1 
LIMIT 1;

-- name: SearchILikeTransactions :many
SELECT transaction_id FROM transactions
WHERE LOWER(transaction_id) LIKE LOWER('%' || $1 || '%');

-- name: SearchUserTransactions :many
SELECT transaction_user_id FROM transactions
WHERE LOWER(transaction_user_id) LIKE LOWER('%' || $1 || '%');
