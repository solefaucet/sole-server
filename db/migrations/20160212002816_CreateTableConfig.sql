
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `configs` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `total_reward_threshold` BIGINT(19) NOT NULL COMMENT 'threshold count in Satonish that determines reward_rate_type',
  `referer_reward_rate` TINYINT(4) NOT NULL COMMENT 'referer reward rate means the percentage get from reward of referee',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `configs` (total_reward_threshold, referer_reward_rate) VALUES (100000000, 10);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `configs`;
