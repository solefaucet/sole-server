
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO `reward_rates`(`min`, `max`, `weight`, `type`) VALUES
(1, 10, 90, 'reward-today-less-than-300usd'),
(11, 50, 7, 'reward-today-less-than-300usd'),
(51, 100, 3, 'reward-today-less-than-300usd'),
(1, 10, 95, 'reward-today-more-than-300usd'),
(11, 50, 4, 'reward-today-more-than-300usd'),
(51, 100, 1, 'reward-today-more-than-300usd');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM `reward_rates` WHERE 
`type` = 'reward-today-less-than-300usd' OR
`type` = 'reward-today-more-than-300usd';
