package models

import (
	"encoding/json"
	"time"
)

// AuthToken model
type AuthToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	AuthToken string    `db:"auth_token"`
	CreatedAt time.Time `db:"created_at"`
}

var _ json.Marshaler = AuthToken{}

// MarshalJSON implements json.Marshaler
func (a AuthToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"auth_token": a.AuthToken,
	})
}
