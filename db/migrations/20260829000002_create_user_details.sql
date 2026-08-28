-- migrate:up
CREATE TABLE user_details (
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    phone      VARCHAR(30)  NOT NULL DEFAULT '',
    address    VARCHAR(255) NOT NULL DEFAULT '',
    city       VARCHAR(100) NOT NULL DEFAULT '',
    country    VARCHAR(100) NOT NULL DEFAULT '',
    bio        VARCHAR(500) NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS user_details;
