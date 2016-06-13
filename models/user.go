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
	ID                      int64     `db:"id" json:"id,omitempty"`
	Email                   string    `db:"email" json:"email,omitempty"`
	Address                 string    `db:"address" json:"address,omitempty"`
	Status                  string    `db:"status" json:"status,omitempty"`
	Balance                 float64   `db:"balance" json:"balance"`
	TotalIncome             float64   `db:"total_income" json:"total_income"`
	TotalIncomeFromReferees float64   `db:"total_income_from_referees" json:"total_income_from_referees"`
	RefererTotalIncome      float64   `db:"referer_total_income" json:"referer_total_income"`
	RewardInterval          int64     `db:"reward_interval" json:"reward_interval"`
	RewardedAt              time.Time `db:"rewarded_at" json:"rewarded_at"`
	RefererID               int64     `db:"referer_id" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
}

// HasReferer indicates if the user is referred by another user
func (u User) HasReferer() bool {
	return u.RefererID > 0
}
