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

-- Acquiring stack: merchants, funding lines, and screening lists (phase 6).
CREATE TABLE IF NOT EXISTS merchants (
    id                 VARCHAR(64) PRIMARY KEY,
    name               VARCHAR(128) NOT NULL,
    dba                VARCHAR(128) NOT NULL,
    tax_id             VARCHAR(32),
    principals         TEXT,
    mccs               TEXT,
    status             VARCHAR(16)  NOT NULL,
    risk_tier          VARCHAR(8)   NOT NULL,
    reserve_rate_bps   BIGINT       NOT NULL DEFAULT 0,
    funding_delay_days INTEGER      NOT NULL DEFAULT 0,
    transaction_limit  BIGINT       NOT NULL DEFAULT 0,
    reserve_balance    BIGINT       NOT NULL DEFAULT 0,
    volume             BIGINT       NOT NULL DEFAULT 0,
    decline_reason     VARCHAR(256),
    approved_at        TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS funding_lines (
    id           BIGSERIAL PRIMARY KEY,
    batch_id     VARCHAR(64) NOT NULL,
    merchant_id  VARCHAR(64) NOT NULL,
    gross        BIGINT      NOT NULL,
    fees         BIGINT      NOT NULL DEFAULT 0,
    reserve_hold BIGINT      NOT NULL DEFAULT 0,
    net          BIGINT      NOT NULL,
    date         TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_funding_lines_merchant
    ON funding_lines (merchant_id, date);

CREATE TABLE IF NOT EXISTS screening_lists (
    id     BIGSERIAL PRIMARY KEY,
    list   VARCHAR(8)  NOT NULL,
    name   VARCHAR(128) NOT NULL,
    tax_id VARCHAR(32),
    detail VARCHAR(128),
    UNIQUE (list, name)
);

CREATE INDEX IF NOT EXISTS idx_screening_lists_list ON screening_lists (list, name);

-- Disputes engine: cases and monitored transactions (phase 7).
CREATE TABLE IF NOT EXISTS disputes (
    id              VARCHAR(64) PRIMARY KEY,
    ref_id          VARCHAR(64) NOT NULL,
    merchant_id     VARCHAR(64) NOT NULL,
    cardholder      VARCHAR(64) NOT NULL,
    amount_minor    BIGINT      NOT NULL,
    currency        VARCHAR(3)  NOT NULL,
    reason_code     VARCHAR(8)  NOT NULL,
    category        VARCHAR(16) NOT NULL,
    stage           VARCHAR(16) NOT NULL,
    status          VARCHAR(16) NOT NULL,
    filed_at        TIMESTAMPTZ NOT NULL,
    response_due    TIMESTAMPTZ NOT NULL,
    responded_at    TIMESTAMPTZ,
    escalated_at    TIMESTAMPTZ,
    evidence        TEXT,
    decision        VARCHAR(16),
    winner          VARCHAR(16),
    decision_at     TIMESTAMPTZ,
    dispute_fee     BIGINT      NOT NULL DEFAULT 0,
    arbitration_fee BIGINT      NOT NULL DEFAULT 0,
    note            VARCHAR(256)
);

CREATE INDEX IF NOT EXISTS idx_disputes_merchant ON disputes (merchant_id, filed_at);
CREATE INDEX IF NOT EXISTS idx_disputes_stage ON disputes (stage);

CREATE TABLE IF NOT EXISTS dispute_transactions (
    ref_id       VARCHAR(64) PRIMARY KEY,
    merchant_id  VARCHAR(64) NOT NULL,
    amount_minor BIGINT      NOT NULL,
    currency     VARCHAR(3)  NOT NULL,
    is_credit    BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_dispute_transactions_merchant
    ON dispute_transactions (merchant_id);


