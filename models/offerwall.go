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

// SuperrewardsOffer model
type SuperrewardsOffer struct {
	ID            int64     `db:"id"`
	IncomeID      int64     `db:"income_id"`
	UserID        int64     `db:"user_id"`
	TransactionID string    `db:"transaction_id"`
	Amount        float64   `db:"amount"`
	OfferID       string    `db:"offer_id"`
	CreatedAt     time.Time `db:"created_at"`
}
