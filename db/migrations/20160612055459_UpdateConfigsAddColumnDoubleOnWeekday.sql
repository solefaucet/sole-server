
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `configs` ADD COLUMN `double_on_weekday` TINYINT(4) NOT NULL DEFAULT -1 COMMENT 'double reward on weekday, starting from sunday 0';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `configs` DROP COLUMN `double_on_weekday`;
