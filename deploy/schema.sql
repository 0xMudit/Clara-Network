-- Clara Network phase 1 schema.
CREATE TABLE IF NOT EXISTS switch_transactions (
    id            BIGSERIAL PRIMARY KEY,
    stan          VARCHAR(6)   NOT NULL,
    mti           VARCHAR(4)   NOT NULL,
    pan_masked    VARCHAR(20),
    amount        VARCHAR(12),
    response_code VARCHAR(3),
    destination   VARCHAR(11),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_switch_transactions_stan
    ON switch_transactions (stan, created_at);
