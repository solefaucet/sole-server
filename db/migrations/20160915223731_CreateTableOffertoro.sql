
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `offertoro` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `income_id` INT(11) NOT NULL,
  `user_id` INT(11) NOT NULL,
  `transaction_id` VARCHAR(255) NOT NULL,
  `offer_id` VARCHAR(127) NOT NULL,
  `amount` DECIMAL(19, 8) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `offertoro`
ADD UNIQUE INDEX (`income_id`),
ADD UNIQUE INDEX (`user_id`, `transaction_id`),
ADD INDEX (`offer_id`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `offertoro`;

