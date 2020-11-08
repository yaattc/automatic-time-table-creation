-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION updated_at_refresh()
    RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TABLE "users" (
     "id" TEXT NOT NULL,
     "email" text NOT NULL,
     "password" text NOT NULL,
     "privileges" jsonb DEFAULT '{}'::jsonb,
     "created_at" timestamptz DEFAULT NOW(),
     "updated_at" timestamptz DEFAULT NOW(),

     CONSTRAINT "users_email_key" UNIQUE ("email"),
     CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE updated_at_refresh();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION updated_at_refresh();
DROP TABLE IF EXISTS "users" CASCADE;
-- +goose StatementEnd