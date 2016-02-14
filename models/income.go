package models

import (
	"encoding/json"
	"fmt"
	"time"
)

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

// MarshalJSON implements json.Marshaler interface
func (i Income) MarshalJSON() ([]byte, error) {
	t := ""
	switch i.Type {
	case IncomeTypeReward:
		t = "reward"
	case IncomeTypeOfferwall:
		t = "offerwall"
	default:
		panic(fmt.Sprintf("Invalid income type %v", i.Type))
	}

	return json.Marshal(map[string]interface{}{
		"type":           t,
		"income":         i.Income,
		"referer_income": i.RefererIncome,
		"created_at":     i.CreatedAt,
	})
}
