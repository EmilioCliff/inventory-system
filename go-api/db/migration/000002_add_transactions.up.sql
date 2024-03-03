-- Create the transaction table
CREATE TABLE "transactions" (
  "transaction_id" varchar PRIMARY KEY,
  "amount" integer NOT NULL,
  "status" boolean NOT NULL DEFAULT false,
  "data_sold" json NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE "receipts" ADD FOREIGN KEY ("receipt_number") REFERENCES "transactions" ("transaction_id"); 
