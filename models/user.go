package models

import "time"

// User status
const (
	UserStatusBanned     = "banned"
	UserStatusUnverified = "unverified"
	UserStatusVerified   = "verified"
)

// User model
type User struct {
	ID             int64     `db:"id" json:"id,omitempty"`
	Email          string    `db:"email" json:"email,omitempty"`
	BitcoinAddress string    `db:"bitcoin_address" json:"bitcoin_address,omitempty"`
	Status         string    `db:"status" json:"status,omitempty"`
	Balance        int64     `db:"balance" json:"balance"`
	RewardInterval int64     `db:"reward_interval" json:"reward_interval"`
	RewardedAt     time.Time `db:"rewarded_at" json:"rewarded_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"-"`
	CreatedAt      time.Time `db:"created_at" json:"-"`
}
