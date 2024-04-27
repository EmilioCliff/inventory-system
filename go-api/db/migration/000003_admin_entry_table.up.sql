-- Create the entries table
CREATE TABLE "entries" (
    "entry_id" bigserial PRIMARY KEY,
    "product_name" varchar NOT NULL,
    "product_price" integer NOT NULL,
    "quantity_added" integer NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_product_name FOREIGN KEY ("product_name") REFERENCES "products" ("product_name")
);