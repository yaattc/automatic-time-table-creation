-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;
CREATE TABLE "users" (
     "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
     "email" text NOT NULL,
     "password" text NOT NULL,
     "privileges" jsonb DEFAULT '{}'::jsonb,
     "created_at" timestamptz DEFAULT NOW(),
     "updated_at" timestamptz DEFAULT NOW(),

     CONSTRAINT "users_email_key" UNIQUE ("email"),
     CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users" CASCADE;
DROP SEQUENCE IF EXISTS users_id_seq;
-- +goose StatementEnd
