
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `offerwow` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(255) NOT NULL,
  `income_id` INT(11) NOT NULL,
  `amount` DECIMAL(19, 8) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `offerwow`
ADD UNIQUE INDEX (`event_id`),
ADD UNIQUE INDEX (`income_id`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `offerwow`;
