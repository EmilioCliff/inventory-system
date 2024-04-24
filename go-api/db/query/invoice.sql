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
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountUserInvoicesByID :one
SELECT COUNT(*) FROM invoices
WHERE user_invoice_id = $1;

-- name: GetUserInvoicesByUsername :many
SELECT * FROM invoices
WHERE user_invoice_username = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountUserInvoicesByUsername :one
SELECT COUNT(*) FROM invoices
WHERE user_invoice_username = $1;

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

-- name: SearchILikeInvoices :many
SELECT invoice_number FROM invoices
WHERE LOWER(invoice_number) LIKE LOWER('%' || $1 || '%');

-- name: SearchUserInvoices :many
SELECT user_invoice_username FROM invoices
WHERE LOWER(user_invoice_username) LIKE LOWER('%' || $1 || '%');

-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices;