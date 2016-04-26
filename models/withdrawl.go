package models

import "time"

// WithdrawalStatus
const (
	WithdrawalStatusPending    = 0
	WithdrawalStatusProcessing = 1
	WithdrawalStatusProcessed  = 2
)

// Withdrawal model
type Withdrawal struct {
	ID            int64     `db:"id" json:"id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	Address       string    `db:"address" json:"address"`
	Amount        float64   `db:"amount" json:"amount"`
	Status        int64     `db:"status" json:"status"`
	TransactionID string    `db:"transaction_id" json:"tx_id"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
