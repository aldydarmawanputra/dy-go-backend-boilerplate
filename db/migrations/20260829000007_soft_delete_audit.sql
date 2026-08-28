-- migrate:up
ALTER TABLE users
    ADD COLUMN created_by UUID,
    ADD COLUMN updated_by UUID,
    ADD COLUMN deleted_at TIMESTAMPTZ;
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

ALTER TABLE user_details
    ADD COLUMN created_by UUID,
    ADD COLUMN updated_by UUID,
    ADD COLUMN deleted_at TIMESTAMPTZ;
CREATE INDEX idx_user_details_deleted_at ON user_details (deleted_at);

ALTER TABLE payments
    ADD COLUMN created_by UUID,
    ADD COLUMN updated_by UUID,
    ADD COLUMN deleted_at TIMESTAMPTZ;
CREATE INDEX idx_payments_deleted_at ON payments (deleted_at);

-- migrate:down
ALTER TABLE payments DROP COLUMN created_by, DROP COLUMN updated_by, DROP COLUMN deleted_at;
ALTER TABLE user_details DROP COLUMN created_by, DROP COLUMN updated_by, DROP COLUMN deleted_at;
ALTER TABLE users DROP COLUMN created_by, DROP COLUMN updated_by, DROP COLUMN deleted_at;
