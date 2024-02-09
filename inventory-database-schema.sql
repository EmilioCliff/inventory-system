SET timezone = 'Africa/Nairobi';

CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "address" varchar NOT NULL,
  "stock" jsonb,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "products" (
  "product_id" bigserial PRIMARY KEY,
  "product_name" varchar UNIQUE NOT NULL,
  "unit_price" int NOT NULL,
  "packsize" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "invoices" (
  "invoice_id" bigserial PRIMARY KEY,
  "invoice_number" string UNIQUE NOT NULL,
  "user_invoice_id" integer,
  "invoce_data" jsonb NOT NULL,
  "created_at" timestamptz DEFAULT 'now()'
);

CREATE TABLE "receipts" (
  "receipt_id" bigserial PRIMARY KEY,
  "receipt_number" string UNIQUE NOT NULL,
  "user_receipt_id" integer,
  "receipt_data" jsonb,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "products" ("product_name");

CREATE INDEX ON "invoices" ("user_invoice_id");

CREATE INDEX ON "invoices" ("invoice_number");

CREATE INDEX ON "receipts" ("user_receipt_id");

CREATE INDEX ON "receipts" ("receipt_number");

ALTER TABLE "invoices" ADD FOREIGN KEY ("user_invoice_id") REFERENCES "users" ("user_id");

ALTER TABLE "receipts" ADD FOREIGN KEY ("user_receipt_id") REFERENCES "users" ("user_id");
