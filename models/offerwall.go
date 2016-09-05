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

// ClixwallOffer model
type ClixwallOffer struct {
	ID        int64     `db:"id"`
	IncomeID  int64     `db:"income_id"`
	UserID    int64     `db:"user_id"`
	OfferID   string    `db:"offer_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

// PtcwallOffer model
type PtcwallOffer struct {
	ID        int64     `db:"id"`
	IncomeID  int64     `db:"income_id"`
	UserID    int64     `db:"user_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

// PersonalyOffer model
type PersonalyOffer struct {
	ID        int64     `db:"id"`
	IncomeID  int64     `db:"income_id"`
	UserID    int64     `db:"user_id"`
	OfferID   string    `db:"offer_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

// TrialpayOffer model
type TrialpayOffer struct {
	ID        int64     `db:"id"`
	IncomeID  int64     `db:"income_id"`
	UserID    int64     `db:"user_id"`
	OfferID   string    `db:"offer_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

// KiwiwallOffer model
type KiwiwallOffer struct {
	ID            int64     `db:"id"`
	IncomeID      int64     `db:"income_id"`
	UserID        int64     `db:"user_id"`
	TransactionID string    `db:"transaction_id"`
	Amount        float64   `db:"amount"`
	OfferID       string    `db:"offer_id"`
	CreatedAt     time.Time `db:"created_at"`
}

// AdscendMedia model
type AdscendMedia struct {
	ID            int64     `db:"id"`
	IncomeID      int64     `db:"income_id"`
	UserID        int64     `db:"user_id"`
	TransactionID string    `db:"transaction_id"`
	Amount        float64   `db:"amount"`
	OfferID       string    `db:"offer_id"`
	CreatedAt     time.Time `db:"created_at"`
}
