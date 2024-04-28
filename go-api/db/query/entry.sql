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
    product_name,
    SUM(product_price) AS total_product_price,
    SUM(quantity_added) AS total_quantity_added
FROM
    entries
GROUP BY
    issued_date, product_name
ORDER BY
    issued_date, product_name;
