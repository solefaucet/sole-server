
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `incomes` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`, `user_id`) PARTITION BY HASH(`user_id`) PARTITIONS 4096;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `incomes` REMOVE PARTITIONING;
