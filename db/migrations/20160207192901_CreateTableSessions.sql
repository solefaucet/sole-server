
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `sessions` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `token` CHAR(36) NOT NULL COMMENT 'token is v4 uuid',
  `type` VARCHAR(63) NOT NULL DEFAULT '' COMMENT 'type can be reset-password or verify-email',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `sessions`
ADD UNIQUE INDEX (`token`),
ADD UNIQUE INDEX (`user_id`, `type`),
ADD INDEX (`updated_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `sessions`;
