package models

import "time"

// TotalReward model
type TotalReward struct {
	ID        int64     `db:"id"`
	Total     float64   `db:"total"`
	CreatedAt time.Time `db:"created_at"`
}

// IsSameDay checks if created_at and now are in the same day
func (t *TotalReward) IsSameDay(now time.Time) bool {
	return t.CreatedAt.YearDay() == now.UTC().YearDay()
}
