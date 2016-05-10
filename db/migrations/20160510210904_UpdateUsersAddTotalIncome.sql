
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `users` ADD COLUMN `total_income` DECIMAL(32, 8) NOT NULL DEFAULT 0 COMMENT 'user total income';
ALTER TABLE `users` ADD COLUMN `referer_total_income` DECIMAL(32, 8) NOT NULL DEFAULT 0 COMMENT 'referer total income';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE `users` DROP COLUMN `total_income`;
ALTER TABLE `users` DROP COLUMN `referer_total_income`;
