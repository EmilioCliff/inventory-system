-- Set the time zone to Nairobi, Kenya
SET timezone = 'Africa/Nairobi';

-- Create the users table
CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "address" varchar NOT NULL,
  "stock" json NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

-- Create the products table
CREATE TABLE "products" (
  "product_id" bigserial PRIMARY KEY,
  "product_name" varchar NOT NULL UNIQUE,
  "unit_price" integer NOT NULL,
  "packsize" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

-- Create the invoices table with ON DELETE CASCADE and ON UPDATE CASCADE
CREATE TABLE "invoices" (
  "invoice_id" bigserial PRIMARY KEY,
  "invoice_number" varchar UNIQUE NOT NULL,
  "user_invoice_id" integer NOT NULL,
  "user_invoice_username" varchar NOT NULL,
  "invoice_data" jsonb NOT NULL,
  "invoice_pdf" bytea NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_user_invoice_id FOREIGN KEY ("user_invoice_id") REFERENCES "users" ("user_id") ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_user_invoice_username FOREIGN KEY ("user_invoice_username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create the receipts table with ON DELETE CASCADE and ON UPDATE CASCADE
CREATE TABLE "receipts" (
  "receipt_id" bigserial PRIMARY KEY,
  "receipt_number" varchar UNIQUE NOT NULL,
  "user_receipt_id" integer NOT NULL,
  "user_receipt_username" varchar NOT NULL,
  "receipt_data" json NOT NULL,
  "receipt_pdf" bytea NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_user_receipt_id FOREIGN KEY ("user_receipt_id") REFERENCES "users" ("user_id") ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_user_receipt_username FOREIGN KEY ("user_receipt_username") REFERENCES "users" ("username") ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create indexes
CREATE INDEX ON "users" ("username");
CREATE INDEX ON "products" ("product_name");
CREATE INDEX ON "invoices" ("user_invoice_id");
CREATE INDEX ON "invoices" ("invoice_number");
CREATE INDEX ON "receipts" ("user_receipt_id");
CREATE INDEX ON "receipts" ("receipt_number");
