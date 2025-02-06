package model

import "time"

// Pings — модель пингера в сервисном слое
type Pings struct {
	ID           string    `json:"id_container"`
	Name         string    `json:"name"`
	IPAddress    string    `json:"ip"`
	Status       bool      `json:"status"`
	ResponseTime float64   `json:"response_time"`
	LastSuccess  time.Time `json:"last_success"`
}

// GetPings — модель пингера в сервисном слое
type GetPings struct {
	ID           string  `json:"id_container"`
	Name         string  `json:"name"`
	IPAddress    string  `json:"ip"`
	Status       bool    `json:"status"`
	ResponseTime float64 `json:"response_time"`
	LastSuccess  string  `json:"last_success"`
}
