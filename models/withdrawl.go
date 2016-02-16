package models

import "time"

// WithdrawlStatus
const (
	WithdrawlStatusPending    = 0
	WithdrawlStatusProcessing = 1
	WithdrawlStatusProcessed  = 2
	WithdrawlStatusRejected   = 3
)

// Withdrawl model
type Withdrawl struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"user_id"`
	BitcoinAddress  string    `db:"bitcoin_address" json:"bitcoin_address"`
	Amount          int64     `db:"amount" json:"amount"`
	Status          int64     `db:"status" json:"status"`
	TransactionHash string    `db:"transaction_hash" json:"tx_hash"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
