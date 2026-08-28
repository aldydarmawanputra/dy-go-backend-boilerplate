-- migrate:up
ALTER TABLE users ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT FALSE;

-- migrate:down
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
