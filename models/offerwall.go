package models

import "time"

// OfferwowEvent model
type OfferwowEvent struct {
	ID        int64     `db:"id"`
	EventID   string    `db:"event_id"`
	IncomeID  int64     `db:"income_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
