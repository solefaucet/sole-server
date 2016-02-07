package models

import "time"

// Session model
type Session struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	Type      string    `db:"type"`
	UpdatedAt time.Time `db:"updated_at"`
}
