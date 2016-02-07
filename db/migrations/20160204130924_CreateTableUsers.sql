
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `users` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL COMMENT 'unique email',
  `bitcoin_address` VARCHAR(63) NOT NULL COMMENT 'unique bitcoin address',
  `verified` BOOL NOT NULL DEFAULT false COMMENT 'indicate if user verifies his email',
  `last_email_sent_at` DATETIME NOT NULL DEFAULT '1000-01-01 00:00:00',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `users` 
ADD UNIQUE INDEX (`email`), 
ADD UNIQUE INDEX (`bitcoin_address`),
ADD INDEX (`verified`),
ADD INDEX (`last_email_sent_at`),
ADD INDEX (`created_at`),
ADD INDEX (`updated_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `users`;
