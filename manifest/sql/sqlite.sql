CREATE TABLE hd_system_config(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name`        varchar(128)    NOT NULL ,
    `key`         varchar(128)    NOT NULL,
    `value`       varchar(128)    NOT NULL ,
    `description` varchar(128)    NOT NULL DEFAULT '' ,
    `created_at`  datetime ,
    `updated_at`  datetime 
);
CREATE UNIQUE  INDEX uk_key
on hd_system_config (key);

CREATE TABLE hd_contract_event(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `contract_name`    varchar(64)  NOT NULL ,
    `contract_address` varchar(42)  NOT NULL ,
    `tx_hash`          varchar(66)  NOT NULL ,
    `event_hash`       varchar(66)  NOT NULL ,
    `event_id`         BIGINT(64)   NOT NULL ,
    `block_number`     BIGINT(64)   NOT NULL ,
    `event_topics`     BLOB    NOT NULL ,
    `event_data`       BLOB     NOT NULL ,
    `state`            TINYINT NOT NULL DEFAULT 0 ,
    `created_at`       datetime,
    `updated_at`       datetime
);
CREATE UNIQUE  INDEX uk_tx_hash
on hd_contract_event (tx_hash, event_id);