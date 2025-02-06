package model

import "time"

// Pings — модель пингера в репо слое
type Pings struct {
	ID           string    `db:"id_container"`
	Name         string    `db:"name"`
	IPAddress    string    `db:"ip"`
	Status       bool      `db:"status"`
	ResponseTime float64   `db:"response_time"`
	LastSuccess  time.Time `db:"last_success"`
}
