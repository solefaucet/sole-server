
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `configs` ADD COLUMN `min_withdrawal_amount` DECIMAL(19, 8) NOT NULL DEFAULT 99999999999.999999 COMMENT 'minimum withdrawal amount';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `configs` DROP COLUMN `min_withdraw_amount`;
