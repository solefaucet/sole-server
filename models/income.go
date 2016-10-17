package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Income Status
const (
	IncomeStatusPending    = "Pending"
	IncomeStatusCharged    = "Charged"
	IncomeStatusChargeback = "Chargeback"
)

// Income Type
const (
	IncomeTypeReward       = 0
	IncomeTypeOfferwow     = 1
	IncomeTypeSuperrewards = 2
	IncomeTypeClixwall     = 3
	IncomeTypePtcwall      = 4
	IncomeTypePersonaly    = 5
	IncomeTypeKiwiwall     = 7
	IncomeTypeAdscendMedia = 8
	IncomeTypeAdgateMedia  = 9
	IncomeTypeOffertoro    = 10
)

var incomeTypes = map[int64]string{
	IncomeTypeReward:       "reward",
	IncomeTypeOfferwow:     "offerwow",
	IncomeTypeSuperrewards: "superrewards",
	IncomeTypeClixwall:     "clixwall",
	IncomeTypePtcwall:      "ptcwall",
	IncomeTypePersonaly:    "personaly",
	IncomeTypeKiwiwall:     "kiwiwall",
	IncomeTypeAdscendMedia: "adscend media",
	IncomeTypeAdgateMedia:  "adgate media",
	IncomeTypeOffertoro:    "offertoro",
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
	Status        string    `db:"status"`
}

// MarshalJSON implements json.Marshaler interface
func (i Income) MarshalJSON() ([]byte, error) {
	t, ok := incomeTypes[i.Type]
	if !ok {
		panic(fmt.Sprintf("Invalid income type %v", i.Type))
	}

	// FIXME: it's silly to put income status in type, FUCK MY CODE
	if i.Type != IncomeTypeReward {
		t = fmt.Sprintf("%s.%s", t, i.Status)
	}

	return json.Marshal(map[string]interface{}{
		"type":           t,
		"income":         i.Income,
		"referer_income": i.RefererIncome,
		"created_at":     i.CreatedAt,
	})
}
