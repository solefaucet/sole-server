
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `incomes` ADD COLUMN `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
ALTER TABLE `incomes` ADD INDEX (`updated_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `incomes` DROP COLUMN `updated_at`;
ALTER TABLE `incomes` DROP INDEX `updated_at`;
