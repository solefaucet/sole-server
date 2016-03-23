
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `withdrawals` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `address` VARCHAR(63) NOT NULL COMMENT 'withdraw to address',
  `amount` INT(11) NOT NULL,
  `status` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '0: pending, 1: processing 2: processed',
  `transaction_id` VARCHAR(127) NOT NULL DEFAULT '' COMMENT 'transaction_id identifies the transaction',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `withdrawals`
ADD INDEX (`user_id`),
ADD INDEX (`address`),
ADD INDEX (`status`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `withdrawals`;
