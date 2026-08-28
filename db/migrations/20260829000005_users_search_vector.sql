-- migrate:up
ALTER TABLE users ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (to_tsvector('simple', coalesce(name, '') || ' ' || coalesce(email, ''))) STORED;

CREATE INDEX idx_users_search_vector ON users USING GIN (search_vector);

-- migrate:down
DROP INDEX IF EXISTS idx_users_search_vector;
ALTER TABLE users DROP COLUMN IF EXISTS search_vector;
