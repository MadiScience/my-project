package domain

import "time"

// CheckResult описывает модель данных результата проверки сайта.
type CheckResult struct {
	URL          string        `json:"url"`
	StatusCode   int           `json:"status_code"`
	ResponseTime time.Duration `json:"response_time"`
	IsAlive      bool          `json:"is_alive"`
	CheckedAt    string        `json:"checked_at"`
}