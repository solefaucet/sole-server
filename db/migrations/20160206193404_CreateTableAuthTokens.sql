
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `auth_tokens` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `auth_token` CHAR(36) NOT NULL COMMENT 'auth token is v4 uuid',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `auth_tokens`
ADD UNIQUE INDEX (`auth_token`),
ADD INDEX (`user_id`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `auth_tokens`;
