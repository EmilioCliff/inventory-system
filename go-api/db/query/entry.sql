-- name: CreateEntry :one
INSERT INTO entries (
    product_id, product_name, product_price, quantity_added
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetEntryByName :many
SELECT
    DATE_TRUNC('day', created_at)::timestamp AS issued_date,
    product_name,
    SUM(product_price) AS total_product_price,
    SUM(quantity_added) AS total_quantity_added
FROM
    entries
GROUP BY
    issued_date, product_name
ORDER BY
    issued_date, product_name DESC;

-- name: ListEntries :many
SELECT * FROM entries
WHERE created_at BETWEEN @from_date AND @to_date;
