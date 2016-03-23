package models

// WebsocketMessage model
type WebsocketMessage struct {
	UsersOnline   int           `json:"users_online,omitempty"`
	LatestIncomes []interface{} `json:"latest_incomes,omitempty"`
	DeltaIncome   interface{}   `json:"delta_income,omitempty"`
}
