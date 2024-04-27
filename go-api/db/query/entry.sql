-- name: CreateEntry :one
INSERT INTO entries (
    product_name, product_price, quantity_added
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetEntryByName :many
SELECT
    DATE_TRUNC('day', created_at)::timestamp AS issued_date,
    COUNT(*) AS num_entries,
    product_name AS product_name,
    product_price AS product_price,
    quantity_added AS quantity_added
FROM
    entries
GROUP BY
    issued_date
ORDER BY
    issued_date;