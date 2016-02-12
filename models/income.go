package models

import "time"

// Income Type
const (
	IncomeTypeReward    = 0
	IncomeTypeOfferwall = 1
)

// Income model
type Income struct {
	ID            int64     `db:"id"`
	UserID        int64     `db:"user_id"`
	RefererID     int64     `db:"referer_id"`
	Type          int64     `db:"type"`
	Income        int64     `db:"income"`
	RefererIncome int64     `db:"referer_income"`
	CreatedAt     time.Time `db:"created_at"`
}
