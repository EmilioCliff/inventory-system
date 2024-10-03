CREATE TABLE "c2b_transactions" (
    "id" bigserial PRIMARY KEY,
    "fullname" varchar NOT NULL,
    "phone" varchar NOT NULL,
    "amount" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "org_account_balance" varchar NOT NULL,
    "transaction_time" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);