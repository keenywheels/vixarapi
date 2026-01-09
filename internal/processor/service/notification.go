package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/keenywheels/backend/internal/vixarapi/models"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.tmpl"))

const (
	notifyWithEmail                   = "email"
	notificationTypeInterestIncreased = "interest_increased"
	headerTypeInterestIncreased       = "Увеличился интерес к отслеживаемой теме"
)

var (
	ErrUnknownNotificationType = errors.New("unknown notification type")
)

// NotifyUser notify user using
func (s *Service) NotifyUser(ctx context.Context, message string) error {
	op := "Service.NotifyUser"

	var event models.Notification
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		return fmt.Errorf("[%s] failed to unmarshal: %w", op, err)
	}

	var err error

	// check notification type
	switch event.NotifyWith {
	case notifyWithEmail:
		err = s.notifyWithEmail(ctx, &event)
	default:
		err = ErrUnknownNotificationType
	}

	if err != nil {
		return fmt.Errorf("[%s] failed to notify user: %w", op, err)
	}

	return nil
}

// notifyWithEmail send notification email to the user
func (s *Service) notifyWithEmail(ctx context.Context, event *models.Notification) error {
	var (
		body   bytes.Buffer
		header string
		err    error
	)

	// generate message body based on notification type
	switch event.Type {
	case notificationTypeInterestIncreased:
		header = headerTypeInterestIncreased

		// TODO: подумать, что можно сделать с тем, что в сообщение токен в денормализованном виде
		err = templates.ExecuteTemplate(&body, fmt.Sprintf("%s.tmpl", event.Type), struct {
			Name             string
			Token            string
			Percentage       int
			CurrentInterest  float64
			PreviousInterest float64
			ScanTime         string
		}{
			Name:             event.Username,
			Token:            event.Token,
			Percentage:       int(event.CurrentInterest / event.PreviousInterest * 100),
			CurrentInterest:  event.CurrentInterest,
			PreviousInterest: event.PreviousInterest,
			ScanTime:         event.ScanDate.Format(time.DateTime),
		})
	}
	if err != nil {
		return fmt.Errorf("failed to create message template for type=%s: %w", event.Type, err)
	}

	// send email
	if err := s.mailer.Send([]string{event.Email}, header, body.String()); err != nil {
		return fmt.Errorf("failed to send email for user %s: %w", event.Email, err)
	}

	return nil
}
