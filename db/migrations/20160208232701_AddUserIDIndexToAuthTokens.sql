
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `auth_tokens`
ADD INDEX (`user_id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX `user_id` ON `auth_tokens`;
