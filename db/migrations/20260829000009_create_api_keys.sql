-- migrate:up
CREATE TABLE api_keys (
    id           UUID PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    key_hash     VARCHAR(64)  NOT NULL UNIQUE,
    prefix       VARCHAR(20)  NOT NULL,
    revoked      BOOLEAN      NOT NULL DEFAULT FALSE,
    last_used_at TIMESTAMPTZ,
    created_by   UUID,
    updated_by   UUID,
    deleted_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_api_keys_deleted_at ON api_keys (deleted_at);

-- migrate:down
DROP TABLE IF EXISTS api_keys;
