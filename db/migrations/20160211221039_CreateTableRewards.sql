
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `rewards` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `reward` INT(11) NOT NULL COMMENT 'reward in 1 Satonish',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `rewards`
ADD INDEX (`user_id`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `rewards`;
