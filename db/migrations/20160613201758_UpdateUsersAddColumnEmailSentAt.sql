
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `users` ADD COLUMN `email_sent_at` DATETIME NOT NULL DEFAULT '1970-01-01 00:00:01' COMMENT 'last email sent time';

ALTER TABLE `users` ADD INDEX (`email_sent_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `users` DROP COLUMN `email_sent_at`;
