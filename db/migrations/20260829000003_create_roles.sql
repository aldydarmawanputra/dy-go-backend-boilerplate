-- migrate:up
CREATE TABLE roles (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (id, name) VALUES (1, 'admin'), (2, 'user') ON CONFLICT DO NOTHING;

-- migrate:down
DROP TABLE IF EXISTS roles;
