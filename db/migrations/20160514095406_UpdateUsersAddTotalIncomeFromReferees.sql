
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `users` ADD COLUMN `total_income_from_referees` DECIMAL(32, 8) NOT NULL DEFAULT 0 COMMENT 'total income get from referees';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE `users` DROP COLUMN `total_income_from_referees`;
