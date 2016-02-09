package models

import "time"

// TotalReward model
type TotalReward struct {
	ID        int64     `db:"id"`
	Total     int64     `db:"total"`
	CreatedAt time.Time `db:"created_at"`
}
