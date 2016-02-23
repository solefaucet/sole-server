
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `withdrawals` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `bitcoin_address` VARCHAR(63) NOT NULL COMMENT 'withdraw to bitcoin address',
  `amount` MEDIUMINT(8) NOT NULL,
  `status` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '0: pending, 1: processing 2: processed',
  `transaction_hash` VARCHAR(127) NOT NULL DEFAULT '' COMMENT 'transaction_hash identifies the transaction in bitcoin blockchain',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `withdrawals`
ADD INDEX (`user_id`),
ADD INDEX (`bitcoin_address`),
ADD INDEX (`status`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `withdrawals`;
