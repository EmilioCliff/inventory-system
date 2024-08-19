-- name: GetUserStockDistributed :one
SELECT COALESCE(SUM((elem->>'totalBill')::integer), 0) AS invoice_total
FROM invoices i, LATERAL (
    SELECT jsonb_array_elements(i.invoice_data) AS elem
    OFFSET 1
) AS elems
WHERE i.user_invoice_id = @user_id AND i.invoice_date BETWEEN @from_date
AND @to_date;

-- name: GetUserInvoiceSummaryBtwnPeriod :many
SELECT invoice_number, invoice_data, 
	(SELECT COALESCE(SUM((elem->>'totalBill')::integer), 0)
     	FROM LATERAL (
            SELECT jsonb_array_elements(invoice_data) AS elem 
            OFFSET 1
        ) AS elem) AS total, invoice_date
FROM invoices WHERE user_invoice_id = @user_id AND invoice_date BETWEEN @from_date AND @to_date;

-- name: GetAllInvoiceSummaryBtwnPeriod :many
SELECT i.invoice_number, i.invoice_data, 
	(SELECT COALESCE(SUM((elem->>'totalBill')::integer), 0)
     	FROM LATERAL (
            SELECT jsonb_array_elements(i.invoice_data) AS elem 
            OFFSET 1
        ) AS elem) AS total, i.invoice_date, u.username
FROM invoices i JOIN users u ON u.user_id = i.user_invoice_id
WHERE invoice_date BETWEEN @from_date AND @to_date;

-- name: GetUserReceiptPaidTotal :one
SELECT COALESCE(SUM(t.amount), 0) AS receipt_total
FROM receipts r JOIN transactions t ON r.receipt_number = t.transaction_id 
WHERE r.user_receipt_id = @user_id AND r.created_at BETWEEN @from_date AND @to_date;

-- name: GetUserReceiptSummaryBtwnPeriod :many
SELECT r.receipt_number, r.receipt_data, t.amount, r.payment_method, r.created_at,  t.mpesa_receipt_number, t.phone_number
FROM receipts r JOIN transactions t ON r.receipt_number = t.transaction_id
WHERE user_receipt_id = @user_id AND r.created_at BETWEEN @from_date AND @to_date;

-- name: GetAllReceiptSummaryBtwnPeriod :many
SELECT r.receipt_number, r.receipt_data, t.amount, r.payment_method, r.created_at, t.mpesa_receipt_number, t.phone_number, u.username
FROM receipts r JOIN transactions t ON r.receipt_number = t.transaction_id
JOIN users u ON u.user_id = r.user_receipt_id
WHERE r.created_at BETWEEN @from_date AND @to_date;

-- name: GetAdminPurchaseOrders :many
SELECT (SELECT SUM((elem->>'quantity')::integer * (elem->>'unit_price')::integer)
FROM LATERAL (SELECT jsonb_array_elements(lpo_data) AS elem) AS elem) AS total, id, supplier_name, created_at, lpo_data
FROM purchase_orders WHERE created_at BETWEEN @from_date AND @to_date;