
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `users` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL COMMENT 'unique email',
  `bitcoin_address` VARCHAR(63) NOT NULL COMMENT 'unique bitcoin address',
  `udpated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `users` ADD UNIQUE INDEX (`email`), ADD UNIQUE INDEX (`bitcoin_address`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `users` DROP INDEX `email`, DROP INDEX `bitcoin_address`;
DROP TABLE `users`;
