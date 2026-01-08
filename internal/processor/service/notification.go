package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/keenywheels/backend/internal/vixarapi/models"
)

// NotifyUser notify user using
func (s *Service) NotifyUser(ctx context.Context, message string) error {
	op := "Service.NotifyUser"

	var event models.Notification
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		return fmt.Errorf("[%s] failed to unmarshal: %w", op, err)
	}

	fmt.Println("TESTING", event)

	return nil
}
