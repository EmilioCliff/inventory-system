ALTER TABLE stock_value DROP CONSTRAINT fk_stock_value_user_id;
ALTER TABLE invoices DROP COLUMN invoice_date;

DROP TABLE IF EXISTS stock_value; 