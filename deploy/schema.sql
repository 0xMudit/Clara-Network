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

-- Clearing and net settlement (phase 3).
CREATE TABLE IF NOT EXISTS clearing_records (
    id            BIGSERIAL PRIMARY KEY,
    cycle_id      VARCHAR(16)   NOT NULL,
    stan          VARCHAR(6)    NOT NULL,
    mti           VARCHAR(4)    NOT NULL,
    sender        VARCHAR(20)   NOT NULL,
    receiver      VARCHAR(20)   NOT NULL,
    amount_minor  BIGINT        NOT NULL,
    interchange   BIGINT        NOT NULL DEFAULT 0,
    currency      VARCHAR(3)    NOT NULL,
    ref_id        VARCHAR(32),
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_clearing_records_cycle
    ON clearing_records (cycle_id);

CREATE TABLE IF NOT EXISTS net_positions (
    id         BIGSERIAL PRIMARY KEY,
    cycle_id   VARCHAR(16) NOT NULL,
    member     VARCHAR(20) NOT NULL,
    net_minor  BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_net_positions_cycle
    ON net_positions (cycle_id);

CREATE TABLE IF NOT EXISTS settlement_instructions (
    id               BIGSERIAL PRIMARY KEY,
    cycle_id         VARCHAR(16) NOT NULL,
    msg_id           VARCHAR(32) NOT NULL,
    member           VARCHAR(20) NOT NULL,
    amount_minor     BIGINT      NOT NULL,
    direction        VARCHAR(8)  NOT NULL,
    currency         VARCHAR(3)  NOT NULL,
    instruction_time TIMESTAMPTZ NOT NULL,
    final            BOOLEAN     NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_settlement_instructions_cycle
    ON settlement_instructions (cycle_id);

CREATE TABLE IF NOT EXISTS prefund_accounts (
    member   VARCHAR(20) PRIMARY KEY,
    balance  BIGINT NOT NULL DEFAULT 0,
    cap      BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS default_fund (
    id      SMALLINT PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0
);

-- Append-only double-entry ledger (phase 4).
CREATE TABLE IF NOT EXISTS ledger_accounts (
    id       VARCHAR(40) PRIMARY KEY,
    type     VARCHAR(16) NOT NULL,
    currency VARCHAR(3)  NOT NULL
);

CREATE TABLE IF NOT EXISTS ledger_entries (
    id         BIGSERIAL PRIMARY KEY,
    journal_id VARCHAR(64) NOT NULL,
    account_id VARCHAR(40) NOT NULL,
    direction  VARCHAR(8)  NOT NULL,
    amount     BIGINT      NOT NULL,
    currency   VARCHAR(3)  NOT NULL,
    reference  VARCHAR(64) NOT NULL,
    posted_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_account
    ON ledger_entries (account_id, posted_at);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_journal
    ON ledger_entries (journal_id);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_reference
    ON ledger_entries (reference);

-- Issuing stack: BIN ranges, cards, and payment tokens (phase 5).
CREATE TABLE IF NOT EXISTS bin_ranges (
    bin      VARCHAR(6) PRIMARY KEY,
    low      BIGINT      NOT NULL,
    high     BIGINT      NOT NULL,
    currency VARCHAR(3)  NOT NULL,
    product  VARCHAR(16) NOT NULL
);

CREATE TABLE IF NOT EXISTS cards (
    ref        VARCHAR(64) PRIMARY KEY,
    pan_hash   BYTEA       NOT NULL UNIQUE,
    pan_masked VARCHAR(24) NOT NULL,
    bin        VARCHAR(6)  NOT NULL,
    expiry     VARCHAR(4)  NOT NULL,
    status     VARCHAR(16) NOT NULL,
    product    VARCHAR(16) NOT NULL,
    udk        BYTEA       NOT NULL,
    last_atc   INTEGER     NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tokens (
    token      VARCHAR(16) PRIMARY KEY,
    pan_hash   BYTEA       NOT NULL,
    par        VARCHAR(29) NOT NULL,
    status     VARCHAR(16) NOT NULL,
    bin        VARCHAR(6)  NOT NULL,
    trid       VARCHAR(16) NOT NULL DEFAULT '',
    device_id  VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tokens_par ON tokens (par);
CREATE INDEX IF NOT EXISTS idx_tokens_pan_hash ON tokens (pan_hash);
