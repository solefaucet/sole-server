package models

import "time"

// Session type
const (
	SessionTypeVerifyEmail   = "verify-email"
	SessionTypeResetPassword = "reset-password"
)

// Session model
type Session struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Token     string    `db:"token"`
	Type      string    `db:"type"`
	UpdatedAt time.Time `db:"updated_at"`
}
