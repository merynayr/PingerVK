package model

import "time"

// Pings модель пингера в сервисном слое
type Pings struct {
	ID           int64
	IP           string
	Status       bool
	ResponseTime time.Time
	LastSuccess  time.Time
}
