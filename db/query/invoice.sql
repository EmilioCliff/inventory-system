-- name: GetInvoice :one
SELECT * FROM invoices
WHERE invoice_number = $1 
LIMIT 1;

-- name: GetUserInvoicesByID :many
SELECT * FROM invoices
WHERE user_invoice_id = $1;

-- name: GetUserInvoicesByUsername :many
SELECT * FROM invoices
WHERE user_invoice_username = $1;

-- name: ListInvoices :many
SELECT * FROM invoices
ORDER BY invoice_id;

-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_number, user_invoice_id, invoice_data, user_invoice_username, invoice_pdf
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;