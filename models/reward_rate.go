package models

import "time"

// RewardRateType
const (
	RewardRateTypeLess = "reward-today-less"
	RewardRateTypeMore = "reward-today-more"
)

// RewardRate model
type RewardRate struct {
	ID        int64     `db:"id"`
	Min       int64     `db:"min"`
	Max       int64     `db:"max"`
	Weight    int64     `db:"weight"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
