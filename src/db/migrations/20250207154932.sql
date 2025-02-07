-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "name" text NOT NULL, "email" text NULL, "password" text NOT NULL, PRIMARY KEY ("id"));
