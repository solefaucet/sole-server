
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `reward_rates` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `min` INT(11) NOT NULL COMMENT 'min value as 1 Penny',
  `max` INT(11) NOT NULL COMMENT 'max value as 1 Penny',
  `weight` INT(11) NOT NULL COMMENT 'weight of rate of this type',
  `type` VARCHAR(63) NOT NULL DEFAULT '' COMMENT 'type can be reward-today-less-than or reward-today-more-than',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `reward_rates`
ADD INDEX (`type`),
ADD INDEX (`updated_at`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `reward_rates`;
