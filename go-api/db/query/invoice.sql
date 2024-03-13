-- name: GetInvoice :one
SELECT * FROM invoices
WHERE invoice_number = $1 
LIMIT 1;

-- name: GetInvoiceByID :one
SELECT * FROM invoices
WHERE invoice_id = $1 
LIMIT 1;

-- name: GetUserInvoicesByID :many
SELECT * FROM invoices
WHERE user_invoice_id = $1
ORDER BY created_at DESC;

-- name: GetUserInvoicesByUsername :many
SELECT * FROM invoices
WHERE user_invoice_username = $1
ORDER BY created_at DESC;

-- name: ListInvoices :many
SELECT * FROM invoices
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_number, user_invoice_id, invoice_data, user_invoice_username, invoice_pdf
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices;