
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `configs` MODIFY `referer_reward_rate` DECIMAL(5, 4);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `configs` MODIFY `referer_reward_rate` DECIMAL(4, 4);
