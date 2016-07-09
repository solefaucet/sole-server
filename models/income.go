package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Income Type
const (
	IncomeTypeReward       = 0
	IncomeTypeOfferwow     = 1
	IncomeTypeSuperrewards = 2
	IncomeTypeClixwall     = 3
)

var incomeTypes = map[int64]string{
	IncomeTypeReward:       "reward",
	IncomeTypeOfferwow:     "offerwow",
	IncomeTypeSuperrewards: "superrewards",
	IncomeTypeClixwall:     "clixwall",
}

// Income model
type Income struct {
	ID            int64     `db:"id"`
	UserID        int64     `db:"user_id"`
	RefererID     int64     `db:"referer_id"`
	Type          int64     `db:"type"`
	Income        float64   `db:"income"`
	RefererIncome float64   `db:"referer_income"`
	CreatedAt     time.Time `db:"created_at"`
}

// MarshalJSON implements json.Marshaler interface
func (i Income) MarshalJSON() ([]byte, error) {
	t, ok := incomeTypes[i.Type]
	if !ok {
		panic(fmt.Sprintf("Invalid income type %v", i.Type))
	}

	return json.Marshal(map[string]interface{}{
		"type":           t,
		"income":         i.Income,
		"referer_income": i.RefererIncome,
		"created_at":     i.CreatedAt,
	})
}
