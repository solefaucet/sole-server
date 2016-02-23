
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `users` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL COMMENT 'unique email',
  `bitcoin_address` VARCHAR(63) NOT NULL COMMENT 'unique bitcoin address',
  `status` VARCHAR(15) NOT NULL DEFAULT 'unverified' COMMENT 'indicate account status, can be unverified|verified|banned',
  `balance` MEDIUMINT(8) NOT NULL DEFAULT 0 COMMENT 'user account balance count in bitcoin satonish',
  `min_withdrawal_amount` MEDIUMINT(8) NOT NULL DEFAULT 100000 COMMENT 'minimum withdrawal amount',
  `reward_interval` SMALLINT(6) NOT NULL DEFAULT 900 COMMENT 'users can get reward every $reward_interval seconds',
  `rewarded_at` DATETIME NOT NULL DEFAULT '1970-01-01 00:00:01',
  `referer_id` INT(11) NOT NULL DEFAULT 0 COMMENT 'default no referer',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `users` 
ADD UNIQUE INDEX (`email`), 
ADD UNIQUE INDEX (`bitcoin_address`),
ADD INDEX (`status`),
ADD INDEX (`balance`),
ADD INDEX (`reward_interval`),
ADD INDEX (`rewarded_at`),
ADD INDEX (`referer_id`),
ADD INDEX (`created_at`),
ADD INDEX (`updated_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `users`;
