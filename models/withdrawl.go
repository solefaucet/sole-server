package models

import "time"

// WithdrawalStatus
const (
	WithdrawalStatusPending    = 0
	WithdrawalStatusProcessing = 1
	WithdrawalStatusProcessed  = 2
	WithdrawalStatusRejected   = 3
)

// Withdrawal model
type Withdrawal struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"user_id"`
	BitcoinAddress  string    `db:"bitcoin_address" json:"bitcoin_address"`
	Amount          int64     `db:"amount" json:"amount"`
	Status          int64     `db:"status" json:"status"`
	TransactionHash string    `db:"transaction_hash" json:"tx_hash"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
