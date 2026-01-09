package smtp

import (
	"fmt"
	"net"
	"net/smtp"

	"github.com/keenywheels/backend/pkg/mailer"
)

// check that Mailer implements Mailer interface
var _ mailer.Mailer = (*Mailer)(nil)

// Mailer struct implements Mailer interface using SMTP protocol and contains auth and config for sending email with SMTP
type Mailer struct {
	auth   smtp.Auth
	config *Config
}

// New creates Mailer interface with underlying SMTPMailer struct
func New(config *Config) *Mailer {
	auth := smtp.PlainAuth(
		"",
		config.Username,
		config.Password,
		config.Host,
	)

	return &Mailer{
		auth:   auth,
		config: config,
	}
}

// Send sends email using SMTP protocol
func (m *Mailer) Send(recipients []string, subject, text string) error {
	mailBody := fmt.Sprintf("Subject: %s\n%s", subject, text)

	if err := smtp.SendMail(
		net.JoinHostPort(m.config.Host, m.config.Port),
		m.auth,
		m.config.Username,
		recipients,
		[]byte(mailBody),
	); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
