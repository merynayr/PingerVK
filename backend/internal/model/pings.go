package model

import "time"

// Pings — модель пингера в сервисном слое
type Pings struct {
	IP           string    `json:"ip" binding:"required"`
	Status       bool      `json:"status"`
	ResponseTime float64   `json:"response_time"`
	LastSuccess  time.Time `json:"last_success"`
}
