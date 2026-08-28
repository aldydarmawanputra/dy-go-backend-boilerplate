-- migrate:up
CREATE TABLE payments (
    id          UUID PRIMARY KEY,
    order_id    VARCHAR(100) NOT NULL UNIQUE,
    provider    VARCHAR(50)  NOT NULL,
    external_id VARCHAR(255) NOT NULL DEFAULT '',
    amount      BIGINT       NOT NULL,
    currency    VARCHAR(10)  NOT NULL,
    status      VARCHAR(20)  NOT NULL,
    payment_url TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_status ON payments (status);

-- migrate:down
DROP TABLE IF EXISTS payments;
