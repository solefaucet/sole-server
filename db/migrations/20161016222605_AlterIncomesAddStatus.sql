
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `incomes` ADD COLUMN `status` VARCHAR(127) NOT NULL DEFAULT 'Charged' COMMENT 'Pending Charged Chargeback';
ALTER TABLE `incomes` ADD INDEX (`status`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `incomes` DROP COLUMN `status`;
ALTER TABLE `incomes` DROP INDEX `status`;
