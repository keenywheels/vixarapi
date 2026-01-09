package models

import "time"

// Notification represent message layout for notification event
type Notification struct {
	NotifyWith       string    `json:"notify_with"`
	Type             string    `json:"type"`
	UserID           string    `json:"user_id"`
	Email            string    `json:"email"`
	Username         string    `json:"username"`
	Token            string    `json:"token"`
	Category         string    `json:"category"`
	Threshold        float64   `json:"threshold"`
	PreviousInterest float64   `json:"previous_interest"`
	CurrentInterest  float64   `json:"current_interest"`
	ScanDate         time.Time `json:"scan_date"`
	// TODO: мб надо будет добавить поля; проверить, когда буду делать логику обработки таски в процессоре
}
