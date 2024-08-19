-- Set the time zone to Nairobi, Kenya
SET timezone = 'Africa/Nairobi';

-- Create the entries table
CREATE TABLE "entries" (
    "entry_id" bigserial PRIMARY KEY,
    "product_id" integer NOT NULL,
    "product_name" varchar NOT NULL,
    "product_price" integer NOT NULL,
    "quantity_added" integer NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_product_id FOREIGN KEY ("product_id") REFERENCES "products" ("product_id") ON UPDATE CASCADE ON DELETE CASCADE
);

-- Set the time zone to Nairobi for the current database
-- ALTER DATABASE railway SET timezone TO 'Africa/Nairobi';