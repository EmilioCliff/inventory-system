-- name: GetInvoicesByDate :many
SELECT 
    -- DATE_TRUNC('day', created_at)::timestamp AS issued_date,
    DATE_TRUNC('minute', created_at) +
    INTERVAL '10 minutes' *
    FLOOR(EXTRACT(MINUTE FROM created_at) / 10) AS issued_date,  
    COUNT(*) AS num_invoices, 
    JSON_AGG(invoice_data) AS invoice_data
FROM 
    invoices
GROUP BY 
    issued_date
ORDER BY 
    issued_date;

-- name: GetUserInvoicesByDate :many
SELECT 
    -- DATE_TRUNC('day', created_at)::timestamp AS issued_date, 
    DATE_TRUNC('minute', created_at) +
    INTERVAL '10 minutes' *
    FLOOR(EXTRACT(MINUTE FROM created_at) / 10) AS issued_date, 
    COUNT(*) AS num_invoices, 
    JSON_AGG(invoice_data) AS invoice_data
FROM 
    invoices
WHERE
    user_invoice_id = $1
GROUP BY 
    issued_date
ORDER BY 
    issued_date;

-- name: GetReceiptsByDate :many
SELECT
    -- DATE_TRUNC('day', created_at)::timestamp AS issued_date,
    DATE_TRUNC('minute', created_at) +
    INTERVAL '10 minutes' *
    FLOOR(EXTRACT(MINUTE FROM created_at) / 10) AS issued_date, 
    COUNT(*) AS num_receipts,
    JSON_AGG(receipt_data) AS receipt_data
FROM 
    receipts
GROUP BY
    issued_date
ORDER BY
    issued_date;

-- name: GetUserReceiptsByDate :many
SELECT
    -- DATE_TRUNC('day', created_at)::timestamp AS issued_date,
    DATE_TRUNC('minute', created_at) +
    INTERVAL '10 minutes' *
    FLOOR(EXTRACT(MINUTE FROM created_at) / 10) AS issued_date, 
    COUNT(*) AS num_receipts,
    JSON_AGG(receipt_data) AS receipt_data
FROM 
    receipts
WHERE
    user_receipt_id = $1
GROUP BY
    issued_date
ORDER BY
    issued_date;
