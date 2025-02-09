-- Create "wallets" table
CREATE TABLE "public"."wallets" ("id" bigserial NOT NULL, "amount" numeric NOT NULL, "points" numeric NOT NULL, "user_id" bigint NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_wallets_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
