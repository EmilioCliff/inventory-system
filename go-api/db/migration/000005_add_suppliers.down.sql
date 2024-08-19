DROP TABLE IF EXISTS purchase_orders;

ALTER TABLE transactions DROP COLUMN payment_method;

ALTER TABLE receipts DROP COLUMN payment_method;