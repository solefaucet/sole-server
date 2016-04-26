
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `incomes` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `referer_id` INT(11) NOT NULL,
  `type` TINYINT(4) NOT NULL COMMENT '0: reward, 1: offerwall',
  `income` DECIMAL(19, 8) NOT NULL,
  `referer_income` DECIMAL(19, 8) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `incomes`
ADD INDEX (`user_id`),
ADD INDEX (`referer_id`),
ADD INDEX (`type`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `incomes`;
