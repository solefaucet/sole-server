
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `reward_rates` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `min` DECIMAL(19, 8) NOT NULL COMMENT 'min value',
  `max` DECIMAL(19, 8) NOT NULL COMMENT 'max value',
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

INSERT INTO `reward_rates`(`min`, `max`, `weight`, `type`) VALUES
(0.00001, 0.0001, 90, 'reward-today-less'),
(0.00011, 0.0005, 7, 'reward-today-less'),
(0.00051, 0.001, 3, 'reward-today-less'),
(0.00001, 0.0001, 95, 'reward-today-more'),
(0.00011, 0.0005, 4, 'reward-today-more'),
(0.00051, 0.001, 1, 'reward-today-more');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `reward_rates`;
