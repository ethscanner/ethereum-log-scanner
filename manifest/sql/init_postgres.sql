DROP TABLE IF EXISTS hd_contract_event;
CREATE TABLE hd_contract_event
(
    id               BIGSERIAL PRIMARY KEY,
    contract_name    VARCHAR(64)     NOT NULL,
    contract_address VARCHAR(42)     NOT NULL,
    tx_hash          VARCHAR(66)     NOT NULL,
    event_hash       VARCHAR(66)     NOT NULL,
    event_id         BIGINT          NOT NULL,
    block_number     BIGINT          NOT NULL,
    block_hash       VARCHAR(66)     NOT NULL,
    event_topics     BYTEA           NOT NULL,
    event_data       BYTEA           NOT NULL,
    state            SMALLINT        NOT NULL DEFAULT 0,
    check_state      SMALLINT        NOT NULL DEFAULT 0,
    checked_block    BIGINT          NOT NULL DEFAULT 0,
    created_at       TIMESTAMP,
    updated_at       TIMESTAMP,
    UNIQUE (block_hash, event_id,tx_hash)
); 