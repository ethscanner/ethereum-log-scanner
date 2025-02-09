DROP TABLE IF EXISTS `hd_system_config`;
CREATE TABLE `hd_system_config`
(
    `id`          BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'id',
    `name`        varchar(128)    NOT NULL COMMENT '配置名称',
    `key`         varchar(128)    NOT NULL COMMENT '配置key',
    `value`       varchar(128)    NOT NULL COMMENT '配置值',
    `description` varchar(128)    NOT NULL DEFAULT '' COMMENT '描述',
    `created_at`  datetime COMMENT '创建时间',
    `updated_at`  datetime COMMENT '更新时间',
    unique uk_key (`key`),
    index idx_created_at (created_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='系统参数配置';

DROP TABLE IF EXISTS `hd_contract_event`;
CREATE TABLE `hd_contract_event`
(
    `id`               BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'id',
    `contract_name`    varchar(64)     NOT NULL COMMENT '合约名',
    `contract_address` varchar(42)     NOT NULL COMMENT '合约地址',
    `tx_hash`          varchar(66)     NOT NULL COMMENT '交易哈希',
    `event_hash`       varchar(66)     NOT NULL COMMENT '事件名',
    `event_id`         BIGINT(64)     NOT NULL COMMENT '事件id',
    `block_number`     BIGINT(64)     NOT NULL COMMENT '区块编号',
    `block_hash`          varchar(66)     NOT NULL COMMENT '交易哈希',
    `event_topics`     BLOB    NOT NULL COMMENT '事件头',
    `event_data`       BLOB     NOT NULL COMMENT 'event数据',
    `state`            TINYINT         NOT NULL DEFAULT 0 COMMENT '0:待处理 10:已处理',
    `check_state`  TINYINT         NOT NULL DEFAULT 0 COMMENT '链上状态: 0:待处理 10:已确认 20:确认异常',
    `checked_block` BIGINT  UNSIGNED   NOT NULL DEFAULT 0 COMMENT '已确认区块',
    `created_at`       datetime COMMENT '创建时间',
    `updated_at`       datetime COMMENT '更新时间',
    unique uk_tx_hash (contract_name,tx_hash,block_hash,event_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='合约事件汇总';
