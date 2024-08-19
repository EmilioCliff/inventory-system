CREATE TABLE "purchase_orders" (
    "id" varchar NOT NULL,
    "supplier_name" varchar NOT NULL,
    "supplier_address" varchar NOT NULL,
    "lpo_data" json NOT NULL,
    "lpo_pdf" bytea NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE transactions ADD payment_method varchar DEFAULT 'MPESA';

ALTER TABLE transactions ALTER COLUMN payment_method SET NOT NULL;

ALTER TABLE receipts ADD payment_method varchar DEFAULT 'MPESA';

ALTER TABLE receipts ALTER COLUMN payment_method SET NOT NULL;