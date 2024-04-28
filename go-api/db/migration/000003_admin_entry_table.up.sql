-- Set the time zone to Nairobi, Kenya
SET timezone = 'Africa/Nairobi';

-- Create the entries table
CREATE TABLE "entries" (
    "entry_id" bigserial PRIMARY KEY,
    "product_name" varchar NOT NULL,
    "product_price" integer NOT NULL,
    "quantity_added" integer NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_product_name FOREIGN KEY ("product_name") REFERENCES "products" ("product_name")
);

-- Set the time zone to Nairobi for the current database
ALTER DATABASE inventorydb SET timezone TO 'Africa/Nairobi';

-- ALTER DATABASE inventorydb SET TIMEZONE TO 'Africa/Nairobi';
-- ALTER TABLE products SET TIMEZONE TO 'Africa/Nairobi';
-- ALTER TABLE invoices SET TIMEZONE TO 'Africa/Nairobi';
-- ALTER TABLE receipts SET TIMEZONE TO 'Africa/Nairobi';
-- ALTER TABLE transactions SET TIMEZONE TO 'Africa/Nairobi';
-- ALTER TABLE entries SET TIMEZONE TO 'Africa/Nairobi';