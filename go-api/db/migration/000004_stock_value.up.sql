CREATE TABLE "stock_value" (
	"user_id" integer NOT NULL,
	"value" bigint NOT NULL DEFAULT 0,
	CONSTRAINT fk_stock_value_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);

ALTER TABLE invoices ADD invoice_date timestamptz DEFAULT NOW();

ALTER TABLE invoices ALTER COLUMN invoice_date SET NOT NULL;
