-- name: CreateC2BTransaction :one
INSERT INTO c2b_transactions (fullname, phone, amount, transaction_id, org_account_balance, transaction_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetC2BTransaction :one
SELECT * FROM c2b_transactions
WHERE transaction_id = $1
LIMIT 1;

-- name: ListC2BTransactions :many
SELECT * FROM c2b_transactions
ORDER BY created_at DESC;

