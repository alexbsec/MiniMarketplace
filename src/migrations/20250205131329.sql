-- Create "products" table
CREATE TABLE "public"."products" ("id" bigserial NOT NULL, "name" text NULL, "description" text NULL, "price" numeric NULL, "points" bigint NULL, "category" text NULL, PRIMARY KEY ("id"));
