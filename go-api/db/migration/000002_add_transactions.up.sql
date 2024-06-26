-- Create the transaction table
CREATE TABLE "transactions" (
  "transaction_id" varchar PRIMARY KEY,
  "transaction_user_id" integer NOT NULL,
  "amount" integer NOT NULL,
  "status" boolean NOT NULL DEFAULT false,
  "data_sold" json NOT NULL,
  "phone_number" varchar NOT NULL DEFAULT '00',
  "mpesa_receipt_number" varchar NOT NULL DEFAULT 'No Receipt Number',
  "result_description" varchar NOT NULL DEFAULT 'Mpesa Not Called Back Description',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_transaction_user_id FOREIGN KEY ("transaction_user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE ON UPDATE CASCADE
);

ALTER TABLE "receipts" ADD FOREIGN KEY ("receipt_number") REFERENCES "transactions" ("transaction_id");
