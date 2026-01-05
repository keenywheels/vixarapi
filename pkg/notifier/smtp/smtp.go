package smtp

import (
	"fmt"
	"net"
	"net/smtp"
)

// Mailer struct implements Mailer interface using SMTP protocol and contains auth and config for sending email with SMTP
type Mailer struct {
	auth   smtp.Auth
	config *Config
}

// New creates Mailer interface with underlying SMTPMailer struct
func New(config *Config) *Mailer {
	auth := smtp.PlainAuth(
		"",
		config.username,
		config.password,
		config.host,
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
		net.JoinHostPort(m.config.host, m.config.port),
		m.auth,
		m.config.username,
		recipients,
		[]byte(mailBody),
	); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
