package model

import "time"

// ContainerInfo — структура для хранения информации о контейнере
type ContainerInfo struct {
	ID           string        `json:"id_container"`
	Name         string        `json:"name"`
	IPAddress    string        `json:"ip"`
	Status       bool          `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
	LastSuccess  time.Time     `json:"last_success"`
}
