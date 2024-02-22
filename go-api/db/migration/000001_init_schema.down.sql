ALTER TABLE invoices DROP CONSTRAINT invoices_user_invoice_id_fkey;
ALTER TABLE receipts DROP CONSTRAINT receipts_user_receipt_id_fkey;


DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS receipts;