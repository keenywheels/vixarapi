package models

// Notification represent message layout for notification event
type Notification struct {
	Type             string  `json:"type"`
	UserID           string  `json:"user_id"`
	Email            string  `json:"email"`
	PreviousInterest float64 `json:"previous_interest"`
	CurrentInterest  float64 `json:"current_interest"`
	// TODO: мб надо будет добавить поля; проверить, когда буду делать логику обработки таски в процессоре
}
