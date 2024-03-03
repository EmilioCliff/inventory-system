ALTER TABLE invoices DROP CONSTRAINT fk_user_invoice_id;
ALTER TABLE invoices DROP CONSTRAINT fk_user_invoice_username;
ALTER TABLE receipts DROP CONSTRAINT fk_user_receipt_id;
ALTER TABLE receipts DROP CONSTRAINT fk_user_receipt_username;


DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS receipts;