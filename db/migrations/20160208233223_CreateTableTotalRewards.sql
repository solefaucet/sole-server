
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `total_rewards` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `total` INT(11) NOT NULL DEFAULT 0,
  `created_at` DATE NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `total_rewards`
ADD INDEX (`total`),
ADD UNIQUE INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `total_rewards`;
