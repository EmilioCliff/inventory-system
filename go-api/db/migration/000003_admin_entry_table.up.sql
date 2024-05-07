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
ALTER DATABASE railway SET timezone TO 'Africa/Nairobi';