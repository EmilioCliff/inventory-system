-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1 
LIMIT 1;

-- name: CreateTransaction :one
INSERT INTO transactions (
    transaction_id, amount, data_sold
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ChangeStatus :one
UPDATE transactions
    set status = $2
WHERE transaction_id = $1
RETURNING *;
