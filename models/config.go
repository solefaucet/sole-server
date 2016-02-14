package models

import "time"

// Config model
type Config struct {
	ID                   int64     `db:"id"`
	TotalRewardThreshold int64     `db:"total_reward_threshold"`
	RefererRewardRate    float64   `db:"referer_reward_rate"`
	BitcoinPrice         int64     `db:"bitcoin_price"`
	UpdatedAt            time.Time `db:"updated_at"`
	CreatedAt            time.Time `db:"created_at"`
}
