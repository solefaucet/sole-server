package models

// WebsocketMessage model
type WebsocketMessage struct {
	BitcoinPrice  int64         `json:"bitcoin_price,omitempty"`
	UsersOnline   int           `json:"users_online,omitempty"`
	LatestIncomes []interface{} `json:"latest_incomes,omitempty"`
	DeltaIncome   interface{}   `json:"delta_income,omitempty"`
}
