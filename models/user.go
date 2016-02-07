package models

import (
	"encoding/json"
	"time"
)

// User model
type User struct {
	ID              int       `db:"id"`
	Email           string    `db:"email"`
	BitcoinAddress  string    `db:"bitcoin_address"`
	Verified        bool      `db:"verified"`
	LastEmailSentAt time.Time `db:"last_email_sent_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	CreatedAt       time.Time `db:"created_at"`
}

var _ json.Marshaler = User{}

// MarshalJSON implements json.Marshaler
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"email":              u.Email,
		"bitcoin_address":    u.BitcoinAddress,
		"verified":           u.Verified,
		"last_email_sent_at": u.LastEmailSentAt,
	})
}
